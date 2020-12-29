package timeseries

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"gitlab.com/thorchain/midgard/event"
	"gitlab.com/thorchain/midgard/internal/db"
	"gitlab.com/thorchain/midgard/openapi/generated/oapigen"
)

func intStr(v int64) string {
	return strconv.FormatInt(v, 10)
}
func floatStr(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

type TxTransactions struct {
	Count uint            `json:"count,string"`
	Txs   []TxTransaction `json:"txs"`
}

type TxTransaction struct {
	Pools     []string         `json:"pools"`
	EventType string           `json:"type"`
	Status    string           `json:"status"`
	In        []oapigen.Tx     `json:"in"`
	Out       []oapigen.Tx     `json:"out"`
	Date      int64            `json:"date,string"`
	Height    int64            `json:"height,string"`
	Metadata  oapigen.Metadata `json:"metadata"`
}

const blankTxId = "0000000000000000000000000000000000000000000000000000000000000000"

// Gets a list of operations generated by external transactions and return its associated data
func TxList(ctx context.Context, moment time.Time, params map[string]string) (oapigen.TxResponse, error) {
	// CHECK PARAMS
	// give latest value if zero moment
	_, timestamp, _ := LastBlock()
	if moment.IsZero() {
		moment = timestamp
	} else if timestamp.Before(moment) {
		return oapigen.TxResponse{}, errBeyondLast
	}

	// check limit param
	if params["limit"] == "" {
		return oapigen.TxResponse{}, errors.New("Query parameter limit is required")
	}
	limit, err := strconv.ParseUint(params["limit"], 10, 64)
	if err != nil || limit < 1 || limit > 50 {
		return oapigen.TxResponse{}, errors.New("limit must be an integer between 1 and 50")
	}

	// check offset param
	if params["offset"] == "" {
		return oapigen.TxResponse{}, errors.New("Query parameter offset is required")
	}
	offset, err := strconv.ParseUint(params["offset"], 10, 64)
	if err != nil || offset < 0 {
		return oapigen.TxResponse{}, errors.New("offset must be a positive integer")
	}

	// build types from type param
	types := make([]string, 0)
	for k := range txInSelectQueries {
		types = append(types, k)
	}
	if params["type"] != "" {
		types = strings.Split(params["type"], ",")
	}

	// EXECUTE QUERIES
	countPS, resultsPS, err := txPreparedStatements(
		moment,
		params["txid"],
		params["address"],
		params["asset"],
		types,
		limit,
		offset)
	if err != nil {
		return oapigen.TxResponse{}, fmt.Errorf("tx prepared statements error: %w", err)
	}

	// Get count
	countRows, err := db.Query(ctx, countPS.Query, countPS.Values...)

	if err != nil {
		return oapigen.TxResponse{}, fmt.Errorf("tx count lookup: %w", err)
	}
	defer countRows.Close()
	var txCount uint
	countRows.Next()
	err = countRows.Scan(&txCount)
	if err != nil {
		return oapigen.TxResponse{}, fmt.Errorf("tx count resolve: %w", err)
	}

	// Get results subset
	rows, err := db.Query(ctx, resultsPS.Query, resultsPS.Values...)
	if err != nil {
		return oapigen.TxResponse{}, fmt.Errorf("tx lookup: %w", err)
	}
	defer rows.Close()

	// PROCESS RESULTS
	transactions := []TxTransaction{}
	// TODO(elfedy): This is a hack to get block heights in a semi-performant way,
	// where we get min and max timestamp and get all the relevant heights
	// If we want to make this operation faster we should consider indexing
	// the block_log table by timestamp or making it an hypertable
	var minTimestamp, maxTimestamp int64
	minTimestamp = math.MaxInt64

	for rows.Next() {
		var result txQueryResult
		err := rows.Scan(
			&result.tx,
			&result.fromAddr,
			&result.tx_2nd,
			&result.fromAddr_2nd,
			&result.toAddr,
			&result.asset,
			&result.assetE8,
			&result.asset_2nd,
			&result.asset_2nd_E8,
			&result.memo,
			&result.pool,
			&result.pool_2nd,
			&result.liquidityFee,
			&result.liquidityUnits,
			&result.tradeSlip,
			&result.tradeTarget,
			&result.asymmetry,
			&result.basisPoints,
			&result.reason,
			&result.eventType,
			&result.blockTimestamp)
		if err != nil {
			return oapigen.TxResponse{}, fmt.Errorf("tx resolve: %w", err)
		}
		var transaction TxTransaction

		transaction, err = txProcessQueryResult(ctx, result)
		if err != nil {
			return oapigen.TxResponse{}, fmt.Errorf("tx resolve: %w", err)
		}

		// compute min/max timestamp to get heights later
		if transaction.Date < minTimestamp {
			minTimestamp = transaction.Date
		}
		if transaction.Date > maxTimestamp {
			maxTimestamp = transaction.Date
		}

		transactions = append(transactions, transaction)
	}

	// get heights and store them in a map
	heights := make(map[int64]int64)
	heightsQuery := "SELECT timestamp, height FROM block_log WHERE TIMESTAMP >= $1 AND TIMESTAMP <= $2"
	heightRows, err := db.Query(ctx, heightsQuery, minTimestamp, maxTimestamp)
	if err != nil {
		return oapigen.TxResponse{}, fmt.Errorf("tx height lookup: %w", err)
	}
	defer heightRows.Close()

	for heightRows.Next() {
		var timestamp, height int64
		err = heightRows.Scan(&timestamp, &height)
		if err != nil {
			return oapigen.TxResponse{}, fmt.Errorf("tx height resolve: %w", err)
		}
		heights[timestamp] = height
	}

	// Add height to each result set
	for i := range transactions {
		transactions[i].Height = heights[transactions[i].Date]
	}

	// TODO(acsaba): copying everything over use these structs above.
	//   this will need a little refactoring because these structs hold strings
	//   instead of ints, and those values are used for calculations.
	txs := make([]oapigen.TxDetails, len(transactions))
	for i, v := range transactions {
		txs[i] = oapigen.TxDetails{
			Pools:    v.Pools,
			Type:     v.EventType,
			Status:   v.Status,
			In:       v.In,
			Out:      v.Out,
			Metadata: v.Metadata,
			Date:     intStr(v.Date),
			Height:   intStr(v.Height),
		}
	}
	return oapigen.TxResponse{Count: intStr(int64(txCount)), Txs: txs}, rows.Err()
}

// Helper structs to build needed queries
// Query key is used in the query to then be replaced when parsed
// This way arguments can be dynamically inserted in query strings
type namedSqlValue struct {
	QueryKey string
	Value    interface{}
}

type preparedSqlStatement struct {
	Query  string
	Values []interface{}
}

// Builds prepared statements for Tx lookup. Two queries are needed, one to get the count
// of the total entries for the query, and one to get the subset that will actually be
// returned to the caller.
// The two queries are built form a base query with the structure:
// SELECT * FROM (inTxType1Query UNION_ALL inTxType2Query...inTxTypeNQuery) WHERE <<conditions>>
func txPreparedStatements(moment time.Time,
	txid,
	address,
	asset string,
	types []string,
	limit,
	offset uint64) (preparedSqlStatement, preparedSqlStatement, error) {

	var countPS, resultsPS preparedSqlStatement
	// Initialize query param slices (to dynamically insert query params)
	baseValues := make([]namedSqlValue, 0)
	subsetValues := make([]namedSqlValue, 0)

	baseValues = append(baseValues, namedSqlValue{"#MOMENT#", moment.UnixNano()})
	subsetValues = append(subsetValues, namedSqlValue{"#LIMIT#", limit}, namedSqlValue{"#OFFSET#", offset})

	// Build select part of the query by taking the tx in queries from the selected types
	// and joining them using UNION ALL
	usedSelectQueries := make([]string, 0)
	for _, eventType := range types {
		q := txInSelectQueries[eventType]
		if q == nil {
			return countPS, resultsPS, fmt.Errorf("invalid type %q", eventType)
		}

		usedSelectQueries = append(usedSelectQueries, q...)
	}
	selectQuery := "SELECT * FROM (" + strings.Join(usedSelectQueries, " UNION ALL ") + ") union_results"

	// TODO(elfedy): this is a temporary hack as for some reason the count query that has
	// a single select query is much slower when no UNIONS happen, and making a union into
	// itself makes it faster. Profiling and optimizing should be done for this at a later stage
	countSelectQuery := selectQuery
	if len(usedSelectQueries) == 1 {
		countSelectQuery = "SELECT * FROM (" + usedSelectQueries[0] + " UNION " + usedSelectQueries[0] + ") union_results"
	}

	// Replace all #RUNE# values with actual asset
	selectQuery = strings.ReplaceAll(selectQuery, "#RUNE#", `'`+event.RuneAsset()+`'`)
	countSelectQuery = strings.ReplaceAll(countSelectQuery, "#RUNE#", `'`+event.RuneAsset()+`'`)

	// build WHERE clause applied to the union_all result, based on filter arguments
	// (txid, address, asset)
	whereQuery := `
	WHERE union_results.block_timestamp <= #MOMENT#`

	if txid != "" {
		baseValues = append(baseValues, namedSqlValue{"#TXID#", txid})
		whereQuery += ` AND (
			union_results.tx = #TXID# OR
			union_results.tx IN (
				SELECT in_tx FROM outbound_events WHERE
					outbound_events.tx = #TXID#
			)
		)`
	}

	if address != "" {
		baseValues = append(baseValues, namedSqlValue{"#ADDRESS#", address})
		whereQuery += ` AND (
			union_results.to_addr = #ADDRESS# OR
			union_results.from_addr = #ADDRESS# OR
			union_results.tx IN (
				SELECT in_tx FROM outbound_events WHERE
					outbound_events.to_addr = #ADDRESS# OR
					outbound_events.from_addr = #ADDRESS#
			)
		)`
	}

	if asset != "" {
		baseValues = append(baseValues, namedSqlValue{"#ASSET#", asset})
		whereQuery += ` AND (
			union_results.asset = #ASSET# OR
			union_results.asset_2nd = #ASSET# OR 
			union_results.tx IN (
				SELECT in_tx FROM outbound_events WHERE
					outbound_events.asset = #ASSET#
			)
		)`
	}

	// build subset query for the results being shown (based on limit and offset)
	subsetQuery := `
	ORDER BY union_results.block_timestamp DESC
	LIMIT #LIMIT# 
	OFFSET #OFFSET# 
	`
	// build and return final queries
	countTxQuery := countSelectQuery + " " + whereQuery
	countQuery := "SELECT count(*) FROM (" + countTxQuery + ") AS count"
	countQueryValues := make([]interface{}, 0)
	for i, queryValue := range baseValues {
		position := i + 1
		positionLabel := fmt.Sprintf("$%d", position)
		countQuery = strings.ReplaceAll(countQuery, queryValue.QueryKey, positionLabel)
		countQueryValues = append(countQueryValues, queryValue.Value)
	}
	countPS = preparedSqlStatement{countQuery, countQueryValues}

	txQuery := selectQuery + " " + whereQuery
	resultsQuery := txQuery + subsetQuery
	resultsQueryValues := make([]interface{}, 0)
	for i, queryValue := range append(baseValues, subsetValues...) {
		position := i + 1
		positionLabel := fmt.Sprintf("$%d", position)
		resultsQuery = strings.ReplaceAll(resultsQuery, queryValue.QueryKey, positionLabel)
		resultsQueryValues = append(resultsQueryValues, queryValue.Value)
	}
	resultsPS = preparedSqlStatement{resultsQuery, resultsQueryValues}

	return countPS, resultsPS, nil
}

type txQueryResult struct {
	tx             string
	fromAddr       string
	tx_2nd         string
	fromAddr_2nd   string
	toAddr         string
	asset          sql.NullString
	assetE8        int64
	asset_2nd      sql.NullString
	asset_2nd_E8   int64
	memo           string
	pool           sql.NullString
	pool_2nd       sql.NullString
	liquidityFee   int64
	liquidityUnits int64
	tradeSlip      int64
	tradeTarget    int64
	asymmetry      float64
	basisPoints    int64
	reason         string
	eventType      string
	blockTimestamp int64
}

func txProcessQueryResult(ctx context.Context, result txQueryResult) (TxTransaction, error) {

	// build incoming related transactions
	var inTxs []oapigen.Tx

	if result.eventType != "addLiquidity" || result.tx == result.tx_2nd {
		inTx := oapigen.Tx{
			Address: result.fromAddr,
			Memo:    result.memo,
			TxID:    result.tx,
		}
		if result.asset.Valid && result.assetE8 > 0 {
			inTx.Coins = append(inTx.Coins, oapigen.Coin{Amount: intStr(result.assetE8), Asset: result.asset.String})
		}
		if result.asset_2nd.Valid && result.asset_2nd_E8 > 0 {
			inTx.Coins = append(inTx.Coins, oapigen.Coin{Amount: intStr(result.asset_2nd_E8), Asset: result.asset_2nd.String})
		}
		inTxs = []oapigen.Tx{inTx}
	} else {
		// Deposit with two separate transactions
		inTx1 := oapigen.Tx{
			Address: result.fromAddr,
			Memo:    result.memo,
			TxID:    result.tx,
			Coins:   []oapigen.Coin{{Amount: intStr(result.assetE8), Asset: result.asset.String}},
		}
		inTx2 := oapigen.Tx{
			Address: result.fromAddr_2nd,
			Memo:    result.memo,
			TxID:    result.tx_2nd,
			Coins:   []oapigen.Coin{{Amount: intStr(result.asset_2nd_E8), Asset: result.asset_2nd.String}},
		}
		inTxs = []oapigen.Tx{inTx1, inTx2}
	}

	// get outbounds and network fees
	outTxs := []oapigen.Tx{}
	var networkFees []oapigen.Coin
	switch result.eventType {
	case "swap", "refund", "withdraw":
		var err error
		outTxs, networkFees, err = getOutboundsAndNetworkFees(ctx, result)
		if err != nil {
			return TxTransaction{}, err
		}
	}

	// process status
	status := "pending"
	switch result.eventType {
	case "swap":
		if len(outTxs) == 1 {
			status = "success"
		}
	case "refund":
		if len(outTxs) == len(inTxs[0].Coins) {
			status = "success"
		}
	case "withdraw":
		if len(outTxs) == 2 {
			status = "success"
		}
	case "donate", "addLiquidity":
		status = "success"
	}

	pools := []string{}
	if result.pool.Valid {
		pools = append(pools, result.pool.String)
	}
	if result.pool_2nd.Valid {
		pools = append(pools, result.pool_2nd.String)
	}

	// Build action metadata
	metadata := oapigen.Metadata{}

	switch result.eventType {
	case "swap":
		metadata.Swap = &oapigen.SwapMetadata{
			LiquidityFee: intStr(result.liquidityFee),
			TradeSlip:    floatStr(float64(result.tradeSlip) / 10000),
			TradeTarget:  intStr(result.tradeTarget),
			NetworkFees:  networkFees,
		}
	case "addLiquidity":
		metadata.AddLiquidity = &oapigen.AddLiquidityMetadata{
			LiquidityUnits: intStr(result.liquidityUnits),
		}
	case "withdraw":
		metadata.Withdraw = &oapigen.WithdrawMetadata{
			LiquidityUnits: intStr(result.liquidityUnits),
			Asymmetry:      floatStr(result.asymmetry),
			NetworkFees:    networkFees,
		}
	case "refund":
		metadata.Refund = &oapigen.RefundMetadata{
			NetworkFees: networkFees,
			Reason:      result.reason,
		}
	}

	transaction := TxTransaction{
		EventType: result.eventType,
		Date:      result.blockTimestamp,
		Metadata:  metadata,
		In:        inTxs,
		Out:       outTxs,
		Pools:     pools,
		Status:    status,
	}

	return transaction, nil
}

func getOutboundsAndNetworkFees(ctx context.Context, result txQueryResult) ([]oapigen.Tx, []oapigen.Coin, error) {
	blockTime := time.Unix(0, result.blockTimestamp)
	outboundTimeLower := blockTime.Add(-OutboundTimeout).UnixNano()
	outboundTimeUpper := blockTime.Add(OutboundTimeout).UnixNano()

	// Get and process outbound transactions (from vault address to external address)
	outboundsQuery := `
	SELECT 
	COALESCE(tx, '` + blankTxId + `'),
	from_addr,
	memo,
	asset,
	asset_E8
	FROM outbound_events
	WHERE in_tx = $1 AND block_timestamp > $2 AND block_timestamp < $3
	`

	networkFeesQuery := `
	SELECT 
	asset,
	asset_E8
	FROM fee_events
	WHERE tx = $1 AND block_timestamp > $2 AND block_timestamp < $3
	`

	outboundRows, err := db.Query(ctx, outboundsQuery, result.tx, outboundTimeLower, outboundTimeUpper)
	if err != nil {
		return nil, nil, fmt.Errorf("outbound tx lookup: %w", err)
	}
	defer outboundRows.Close()

	networkFeeRows, err := db.Query(ctx, networkFeesQuery, result.tx, outboundTimeLower, outboundTimeUpper)
	if err != nil {
		return nil, nil, fmt.Errorf("network fee lookup: %w", err)
	}
	defer networkFeeRows.Close()

	outTxs := []oapigen.Tx{}

	for outboundRows.Next() {
		var tx,
			address,
			memo,
			asset string
		var assetE8 int64

		err := outboundRows.Scan(&tx, &address, &memo, &asset, &assetE8)
		if err != nil {
			return nil, nil, fmt.Errorf("outbound tx lookup: %w", err)
		}
		outTx := oapigen.Tx{
			Address: address,
			Coins:   []oapigen.Coin{{Amount: intStr(assetE8), Asset: asset}},
			Memo:    memo,
			TxID:    tx,
		}
		outTxs = append(outTxs, outTx)
	}

	networkFees := []oapigen.Coin{}

	for networkFeeRows.Next() {
		var asset string
		var assetE8 int64

		err := networkFeeRows.Scan(&asset, &assetE8)
		if err != nil {
			return nil, nil, fmt.Errorf("network fee lookup: %w", err)
		}
		networkFee := oapigen.Coin{
			Amount: intStr(assetE8),
			Asset:  asset,
		}
		networkFees = append(networkFees, networkFee)
	}

	return outTxs, networkFees, nil
}

// txIn select queries: list of queries that have inbound
// transactions as rows. They are given a type based on the operation they relate to.
// These queries are built using data from events sent by Thorchain
var txInSelectQueries = map[string][]string{
	"swap": {`SELECT 
				tx,
				from_addr,
				'' as tx_2nd,
				'' as from_addr_2nd,
				to_addr,
				from_asset as asset,
				from_E8 as asset_E8,
				'' as asset_2nd,
				0 as asset_2nd_E8,
				memo,
				pool,
				NULL as pool_2nd,
				liq_fee_in_rune_E8,
				0 as stake_units,
				trade_slip_BP,
				to_E8_min as trade_target,
				0 as asymmetry,
				0 as basis_points,
				'' as reason,
				'swap' as type,
				block_timestamp
			FROM swap_events AS single_swaps
			WHERE NOT EXISTS (
				SELECT tx FROM swap_events WHERE block_timestamp = single_swaps.block_timestamp AND tx = single_swaps.tx AND from_asset <> single_swaps.from_asset
			)`,
		`SELECT
				swap_in.tx as tx,
				swap_in.from_addr as from_addr,
				'' as tx_2nd,
				'' as from_addr_2nd,
				swap_in.to_addr as to_addr,
				swap_in.from_asset as asset,
				swap_in.from_E8 as asset_E8,
				NULL as asset_2nd,
				0 as asset_2nd_E8,
				swap_in.memo as memo,
				swap_in.pool as pool,
				swap_out.pool as pool_2nd,
				(swap_in.liq_fee_in_rune_E8 + swap_out.liq_fee_in_rune_E8) as liq_fee_E8,
				swap_out.to_E8_min as trade_target,
				0 as stake_units,
				(swap_in.trade_slip_BP + swap_out.trade_slip_BP) as trade_slip_BP,
				0 as asymmetry,
				0 as basis_points,
				'' as reason,
				'swap' as type,
				swap_in.block_timestamp as block_timestamp
			FROM
			swap_events AS swap_in
			INNER JOIN
			swap_events AS swap_out
			ON swap_in.tx = swap_out.tx
			WHERE swap_in.from_asset = swap_in.pool AND swap_out.from_asset <> swap_out.pool AND swap_in.block_timestamp = swap_out.block_timestamp`},
	"addLiquidity": {
		// TODO(elfedy): v1 queries thorchain to get some tx details when it parses the events
		// (i.e: the memo, non rune address) those are currently missing in this implementation.
		`SELECT 
					rune_tx as tx,
					rune_addr as from_addr,
					asset_tx as tx_2nd,
					asset_addr as from_addr_2nd,
					'' as to_addr,
					#RUNE# as asset,
					rune_E8 as asset_E8,
					pool as asset_2nd,
					asset_E8 as asset_2nd_E8,
					'' as memo,
					pool,
					NULL as pool_2nd,
					0 as liq_fee_E8,
					stake_units,
					0 as trade_slip_BP,
					0 as trade_target,
					0 as asymmetry,
					0 as basis_points,
					'' as reason,
					'addLiquidity' as type,
					block_timestamp
				FROM stake_events`},
	"withdraw": {`
			SELECT 
				tx,
				from_addr,
				'' as tx_2nd,
				'' as from_addr_2nd,
				to_addr,
				asset,
				asset_E8,
				'' as asset_2nd,
				0 as asset_2nd_E8,
				memo,
				pool,
				NULL as pool_2nd,
				0 as liq_fee_E8,
				(stake_units * -1) as stake_units,
				0 as trade_slip_BP,
				0 as trade_target,
				asymmetry,
				basis_points,
				'' as reason,
				'withdraw' as type,
				block_timestamp
			FROM unstake_events`},
	"donate": {`
			SELECT 
				tx,
				from_addr,
				'' as tx_2nd,
				'' as from_addr_2nd,
				to_addr,
				asset,
				asset_E8,
				#RUNE# as asset_2nd,
				rune_E8 as asset_2nd_E8,
				memo,
				pool,
				NULL as pool_2nd,
				0 as liq_fee_E8,
				0 as stake_units,
				0 as trade_slip_BP,
				0 as trade_target,
				0 as asymmetry,
				0 as basis_points,
				'' as reason,
				'add' as type,
				block_timestamp
			FROM add_events`},
	"refund": {`SELECT 
				tx,
				from_addr,
				'' as tx_2nd,
				'' as from_addr_2nd,
				to_addr,
				asset,
				asset_E8,
				asset_2nd,
				asset_2nd_E8,
				memo,
				NULL as pool,
				NULL as pool_2nd,
				0 as liq_fee_E8,
				0 as stake_units,
				0 as trade_slip_BP,
				0 as trade_target,
				0 as asymmetry,
				0 as basis_points,
				reason,
				'refund' as type,
				block_timestamp
			FROM refund_events`},
}
