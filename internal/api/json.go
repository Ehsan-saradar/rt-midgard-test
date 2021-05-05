package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/thorchain/midgard/internal/db"
	"gitlab.com/thorchain/midgard/internal/graphql/model"
	"gitlab.com/thorchain/midgard/internal/util/miderr"

	"gitlab.com/thorchain/midgard/internal/timeseries"
	"gitlab.com/thorchain/midgard/internal/timeseries/stat"
	"gitlab.com/thorchain/midgard/openapi/generated/oapigen"
)

// Version 1 compatibility is a minimal effort attempt to provide smooth migration.

// InSync returns whether the entire blockchain is processed.
var InSync func() bool

type Health struct {
	CatchingUp    bool  `json:"catching_up"`
	Database      bool  `json:"database"`
	ScannerHeight int64 `json:"scannerHeight,string"`
}

func jsonHealth(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	height, _, _ := timeseries.LastBlock()
	synced := InSync()
	respJSON(w, oapigen.HealthResponse{
		InSync:        synced,
		Database:      true,
		ScannerHeight: intStr(height + 1),
	})
}

func jsonEarningsHistory(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	buckets, merr := db.BucketsFromQuery(r.Context(), r.URL.Query())
	if merr != nil {
		merr.ReportHTTP(w)
		return
	}

	var res oapigen.EarningsHistoryResponse
	res, err := stat.GetEarningsHistory(r.Context(), buckets)
	if err != nil {
		miderr.InternalErrE(err).ReportHTTP(w)
		return
	}
	if buckets.OneInterval() {
		res.Intervals = oapigen.EarningsHistoryIntervals{}
	}
	respJSON(w, res)
}

func jsonLiquidityHistory(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	query := r.URL.Query()

	buckets, merr := db.BucketsFromQuery(r.Context(), query)
	if merr != nil {
		merr.ReportHTTP(w)
		return
	}

	pool := query.Get("pool")
	if pool == "" {
		pool = "*"
	}
	var res oapigen.LiquidityHistoryResponse
	res, err := stat.GetLiquidityHistory(r.Context(), buckets, pool)
	if err != nil {
		miderr.InternalErrE(err).ReportHTTP(w)
		return
	}
	if buckets.OneInterval() {
		res.Intervals = oapigen.LiquidityHistoryIntervals{}
	}
	respJSON(w, res)
}

func jsonDepths(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pool := ps[0].Value

	if !timeseries.PoolExists(pool) {
		miderr.BadRequestF("Unknown pool: %s", pool).ReportHTTP(w)
		return
	}

	query := r.URL.Query()

	buckets, merr := db.BucketsFromQuery(r.Context(), query)
	if merr != nil {
		merr.ReportHTTP(w)
		return
	}

	depths, err := stat.PoolDepthHistory(r.Context(), buckets, pool)
	if err != nil {
		miderr.InternalErrE(err).ReportHTTP(w)
		return
	}
	units, err := stat.PoolLiquidityUnitsHistory(r.Context(), buckets, pool)
	if err != nil {
		miderr.InternalErrE(err).ReportHTTP(w)
		return
	}
	if len(depths) != len(units) || depths[0].Window != units[0].Window {
		miderr.InternalErr("Buckets misalligned").ReportHTTP(w)
		return
	}
	var result oapigen.DepthHistoryResponse = toOapiDepthResponse(depths, units)
	respJSON(w, result)
}

func toOapiDepthResponse(
	depths []stat.PoolDepthBucket,
	units []stat.UnitsBucket) (
	result oapigen.DepthHistoryResponse) {
	result.Intervals = make(oapigen.DepthHistoryIntervals, 0, len(depths))
	for i, bucket := range depths {
		result.Intervals = append(result.Intervals, oapigen.DepthHistoryItem{
			StartTime:      intStr(bucket.Window.From.ToI()),
			EndTime:        intStr(bucket.Window.Until.ToI()),
			AssetDepth:     intStr(bucket.Depths.AssetDepth),
			RuneDepth:      intStr(bucket.Depths.RuneDepth),
			AssetPrice:     floatStr(bucket.Depths.AssetPrice()),
			AssetPriceUSD:  floatStr(bucket.AssetPriceUSD),
			LiquidityUnits: intStr(units[i].Units),
		})
	}
	result.Meta.StartTime = intStr(depths[0].Window.From.ToI())
	result.Meta.EndTime = intStr(depths[len(depths)-1].Window.Until.ToI())
	return
}

func jsonSwapHistory(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	query := r.URL.Query()

	buckets, merr := db.BucketsFromQuery(r.Context(), query)
	if merr != nil {
		merr.ReportHTTP(w)
		return
	}

	var pool *string
	poolParam := query.Get("pool")
	if poolParam != "" {
		pool = &poolParam
	}

	mergedPoolSwaps, err := stat.GetPoolSwaps(r.Context(), pool, buckets)
	if err != nil {
		miderr.InternalErr(err.Error()).ReportHTTP(w)
		return
	}
	var result oapigen.SwapHistoryResponse = createVolumeIntervals(mergedPoolSwaps)
	if buckets.OneInterval() {
		result.Intervals = oapigen.SwapHistoryIntervals{}
	}
	respJSON(w, result)
}

func toSwapHistoryItem(bucket stat.SwapBucket) oapigen.SwapHistoryItem {
	return oapigen.SwapHistoryItem{
		StartTime:          intStr(bucket.StartTime.ToI()),
		EndTime:            intStr(bucket.EndTime.ToI()),
		ToRuneVolume:       intStr(bucket.ToRuneVolume),
		ToAssetVolume:      intStr(bucket.ToAssetVolume),
		TotalVolume:        intStr(bucket.TotalVolume),
		ToAssetCount:       intStr(bucket.ToAssetCount),
		ToRuneCount:        intStr(bucket.ToRuneCount),
		TotalCount:         intStr(bucket.TotalCount),
		ToAssetFees:        intStr(bucket.ToAssetFees),
		ToRuneFees:         intStr(bucket.ToRuneFees),
		TotalFees:          intStr(bucket.TotalFees),
		ToAssetAverageSlip: ratioStr(bucket.ToAssetSlip, bucket.ToAssetCount),
		ToRuneAverageSlip:  ratioStr(bucket.ToRuneSlip, bucket.ToRuneCount),
		AverageSlip:        ratioStr(bucket.TotalSlip, bucket.TotalCount),
		RunePriceUSD:       floatStr(bucket.RunePriceUSD),
	}
}

func createVolumeIntervals(buckets []stat.SwapBucket) (result oapigen.SwapHistoryResponse) {
	metaBucket := stat.SwapBucket{}

	for _, bucket := range buckets {
		metaBucket.AddBucket(bucket)

		result.Intervals = append(result.Intervals, toSwapHistoryItem(bucket))
	}

	result.Meta = toSwapHistoryItem(metaBucket)
	result.Meta.StartTime = result.Intervals[0].StartTime
	result.Meta.EndTime = result.Intervals[len(result.Intervals)-1].EndTime
	result.Meta.RunePriceUSD = result.Intervals[len(result.Intervals)-1].RunePriceUSD
	return
}

// TODO(huginn): remove when bonds are fixed
var ShowBonds bool = false

func jsonTVLHistory(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	query := r.URL.Query()

	buckets, merr := db.BucketsFromQuery(r.Context(), query)
	if merr != nil {
		merr.ReportHTTP(w)
		return
	}

	depths, err := stat.TVLDepthHistory(r.Context(), buckets)
	if err != nil {
		miderr.InternalErrE(err).ReportHTTP(w)
		return
	}
	bonds, err := stat.BondsHistory(r.Context(), buckets)
	if err != nil {
		miderr.InternalErrE(err).ReportHTTP(w)
		return
	}
	if len(depths) != len(bonds) || depths[0].Window != bonds[0].Window {
		miderr.InternalErr("Buckets misalligned").ReportHTTP(w)
		return
	}
	var result oapigen.TVLHistoryResponse = toTVLHistoryResponse(depths, bonds)
	respJSON(w, result)
}

func toTVLHistoryResponse(depths []stat.TVLDepthBucket, bonds []stat.BondBucket) (result oapigen.TVLHistoryResponse) {
	showBonds := func(value string) *string {
		if !ShowBonds {
			return nil
		}
		return &value
	}

	result.Intervals = make(oapigen.TVLHistoryIntervals, 0, len(depths))
	for i, bucket := range depths {
		pools := 2 * bucket.TotalPoolDepth
		bonds := bonds[i].Bonds
		result.Intervals = append(result.Intervals, oapigen.TVLHistoryItem{
			StartTime:        intStr(bucket.Window.From.ToI()),
			EndTime:          intStr(bucket.Window.Until.ToI()),
			TotalValuePooled: intStr(pools),
			TotalValueBonded: showBonds(intStr(bonds)),
			TotalValueLocked: showBonds(intStr(pools + bonds)),
			RunePriceUSD:     floatStr(bucket.RunePriceUSD),
		})
	}
	result.Meta = result.Intervals[len(depths)-1]
	result.Meta.StartTime = result.Intervals[0].StartTime
	return
}

type Network struct {
	ActiveBonds     []string `json:"activeBonds,string"`
	ActiveNodeCount int      `json:"activeNodeCount,string"`
	BlockRewards    struct {
		BlockReward int64 `json:"blockReward,string"`
		BondReward  int64 `json:"bondReward,string"`
		PoolReward  int64 `json:"poolReward,string"`
	} `json:"blockRewards"`
	BondMetrics struct {
		TotalActiveBond    int64 `json:"totalActiveBond,string"`
		AverageActiveBond  int64 `json:"averageActiveBond,string"`
		MedianActiveBond   int64 `json:"medianActiveBond,string"`
		MinimumActiveBond  int64 `json:"minimumActiveBond,string"`
		MaximumActiveBond  int64 `json:"maximumActiveBond,string"`
		TotalStandbyBond   int64 `json:"totalStandbyBond,string"`
		MinimumStandbyBond int64 `json:"minimumStandbyBond,string"`
		MaximumStandbyBond int64 `json:"maximumStandbyBond,string"`
		AverageStandbyBond int64 `json:"averageStandbyBond,string"`
		MedianStandbyBond  int64 `json:"medianStandbyBond,string"`
	} `json:"bondMetrics"`
	StandbyBonds            []string `json:"standbyBonds,string"`
	StandbyNodeCount        int      `json:"standbyNodeCount,string"`
	TotalPooledRune         int64    `json:"totalPooledRune,string"`
	TotalReserve            int64    `json:"totalReserve,string"`
	NextChurnHeight         int64    `json:"nextChurnHeight,string"`
	PoolActivationCountdown int64    `json:"poolActivationCountdown,string"`
	PoolShareFactor         float64  `json:"poolShareFactor,string"`
	BondingAPY              float64  `json:"bondingAPY,string"`
	LiquidityAPY            float64  `json:"liquidityAPY,string"`
}

func jsonNetwork(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	network, err := timeseries.GetNetworkData(r.Context())
	if err != nil {
		respError(w, r, err)
		return
	}

	respJSON(w, convertNetwork(network))
}

func convertNetwork(network model.Network) oapigen.Network {
	return oapigen.Network{
		ActiveBonds:     intArrayStrs(network.ActiveBonds),
		ActiveNodeCount: intStr(network.ActiveNodeCount),
		BlockRewards: oapigen.BlockRewards{
			BlockReward: intStr(network.BlockRewards.BlockReward),
			BondReward:  intStr(network.BlockRewards.BondReward),
			PoolReward:  intStr(network.BlockRewards.PoolReward),
		},
		// TODO(acsaba): create bondmetrics right away with this type.
		BondMetrics: oapigen.BondMetrics{
			TotalActiveBond:    intStr(network.BondMetrics.Active.TotalBond),
			AverageActiveBond:  intStr(network.BondMetrics.Active.AverageBond),
			MedianActiveBond:   intStr(network.BondMetrics.Active.MedianBond),
			MinimumActiveBond:  intStr(network.BondMetrics.Active.MinimumBond),
			MaximumActiveBond:  intStr(network.BondMetrics.Active.MaximumBond),
			TotalStandbyBond:   intStr(network.BondMetrics.Standby.TotalBond),
			AverageStandbyBond: intStr(network.BondMetrics.Standby.AverageBond),
			MedianStandbyBond:  intStr(network.BondMetrics.Standby.MedianBond),
			MinimumStandbyBond: intStr(network.BondMetrics.Standby.MinimumBond),
			MaximumStandbyBond: intStr(network.BondMetrics.Standby.MaximumBond),
		},
		BondingAPY:              floatStr(network.BondingApy),
		LiquidityAPY:            floatStr(network.LiquidityApy),
		NextChurnHeight:         intStr(network.NextChurnHeight),
		PoolActivationCountdown: intStr(network.PoolActivationCountdown),
		PoolShareFactor:         floatStr(network.PoolShareFactor),
		StandbyBonds:            intArrayStrs(network.StandbyBonds),
		StandbyNodeCount:        intStr(network.StandbyNodeCount),
		TotalReserve:            intStr(network.TotalReserve),
		TotalPooledRune:         intStr(network.TotalPooledRune),
	}
}

type Node struct {
	Secp256K1 string `json:"secp256k1"`
	Ed25519   string `json:"ed25519"`
}

func jsonNodes(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	secpAddrs, edAddrs, err := timeseries.NodesSecpAndEd(r.Context(), time.Now())
	if err != nil {
		respError(w, r, err)
		return
	}

	m := make(map[string]struct {
		Secp string
		Ed   string
	}, len(secpAddrs))
	for key, addr := range secpAddrs {
		e := m[addr]
		e.Secp = key
		m[addr] = e
	}
	for key, addr := range edAddrs {
		e := m[addr]
		e.Ed = key
		m[addr] = e
	}

	array := make([]oapigen.Node, 0, len(m))
	for key, e := range m {
		array = append(array, oapigen.Node{
			Secp256k1:   e.Secp,
			Ed25519:     e.Ed,
			NodeAddress: key,
		})
	}
	respJSON(w, array)
}

// Filters out Suspended pools.
// If there is a status url parameter then returns pools with that status only.
func poolsWithRequestedStatus(r *http.Request, statusMap map[string]string) ([]string, error) {
	pools, err := timeseries.PoolsWithDeposit(r.Context())
	if err != nil {
		return nil, err
	}
	statusParams := r.URL.Query()["status"]
	requestedStatus := ""
	if len(statusParams) != 0 {
		const errormsg = "Max one status parameter, accepted values: available, staged, suspended"
		if 1 < len(statusParams) {
			return nil, fmt.Errorf(errormsg)
		}
		requestedStatus = statusParams[0]
		requestedStatus = strings.ToLower(requestedStatus)
		// Allowed statuses in
		// https://gitlab.com/thorchain/thornode/-/blob/master/x/thorchain/types/type_pool.go
		if requestedStatus != "available" && requestedStatus != "staged" && requestedStatus != "suspended" {
			return nil, fmt.Errorf(errormsg)
		}
	}
	ret := []string{}
	for _, pool := range pools {
		poolStatus := poolStatusFromMap(pool, statusMap)
		if poolStatus != "suspended" && (requestedStatus == "" || poolStatus == requestedStatus) {
			ret = append(ret, pool)
		}
	}
	return ret, nil
}

type poolAggregates struct {
	dailyVolumes        map[string]int64
	poolUnits           map[string]int64
	poolAPYs            map[string]float64
	assetE8DepthPerPool map[string]int64
	runeE8DepthPerPool  map[string]int64
}

func getPoolAggregates(ctx context.Context, pools []string) (*poolAggregates, error) {
	assetE8DepthPerPool, runeE8DepthPerPool, timestamp := timeseries.AssetAndRuneDepths()
	now := db.TimeToSecond(timestamp)
	dayAgo := now - 24*60*60

	dailyVolumes, err := stat.PoolsTotalVolume(ctx, pools, dayAgo.ToNano(), now.ToNano())
	if err != nil {
		return nil, err
	}

	poolUnits, err := stat.CurrentPoolsLiquidityUnits(ctx, pools)
	if err != nil {
		return nil, err
	}

	week := db.Window{From: now - 7*24*60*60, Until: now}
	poolAPYs, err := timeseries.GetPoolAPY(ctx, runeE8DepthPerPool, pools, week)

	aggregates := poolAggregates{
		dailyVolumes:        dailyVolumes,
		poolUnits:           poolUnits,
		poolAPYs:            poolAPYs,
		assetE8DepthPerPool: assetE8DepthPerPool,
		runeE8DepthPerPool:  runeE8DepthPerPool,
	}

	return &aggregates, nil
}

func poolStatusFromMap(pool string, statusMap map[string]string) string {
	status, ok := statusMap[pool]
	if !ok {
		status = timeseries.DefaultPoolStatus
	}
	return status
}

func buildPoolDetail(
	pool, status string, aggregates poolAggregates, runePriceUsd float64) oapigen.PoolDetail {
	assetDepth := aggregates.assetE8DepthPerPool[pool]
	runeDepth := aggregates.runeE8DepthPerPool[pool]
	dailyVolume := aggregates.dailyVolumes[pool]
	poolUnits := aggregates.poolUnits[pool]
	poolAPY := aggregates.poolAPYs[pool]
	price := timeseries.AssetPrice(assetDepth, runeDepth)
	priceUSD := price * runePriceUsd

	return oapigen.PoolDetail{
		Asset:         pool,
		AssetDepth:    intStr(assetDepth),
		RuneDepth:     intStr(runeDepth),
		PoolAPY:       floatStr(poolAPY),
		AssetPrice:    floatStr(price),
		AssetPriceUSD: floatStr(priceUSD),
		Status:        status,
		Units:         intStr(poolUnits),
		Volume24h:     intStr(dailyVolume),
	}
}

func jsonPools(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_, lastTime, _ := timeseries.LastBlock()
	statusMap, err := timeseries.GetPoolsStatuses(r.Context(), db.Nano(lastTime.UnixNano()))
	if err != nil {
		respError(w, r, err)
		return
	}
	pools, err := poolsWithRequestedStatus(r, statusMap)
	if err != nil {
		respError(w, r, err)
		return
	}

	aggregates, err := getPoolAggregates(r.Context(), pools)
	if err != nil {
		respError(w, r, err)
		return
	}

	runePriceUsd := stat.RunePriceUSD()

	poolsResponse := make(oapigen.PoolsResponse, len(pools))
	for i, pool := range pools {
		status := poolStatusFromMap(pool, statusMap)
		poolsResponse[i] = buildPoolDetail(pool, status, *aggregates, runePriceUsd)
	}

	respJSON(w, poolsResponse)
}

func jsonPool(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pool := ps[0].Value

	if !timeseries.PoolExistsNow(pool) {
		miderr.BadRequestF("Unknown pool: %s", pool).ReportHTTP(w)
		return
	}

	status, err := timeseries.PoolStatus(r.Context(), pool)
	if err != nil {
		miderr.InternalErrE(err).ReportHTTP(w)
		return
	}

	aggregates, err := getPoolAggregates(r.Context(), []string{pool})
	if err != nil {
		miderr.InternalErrE(err).ReportHTTP(w)
		return
	}

	runePriceUsd := stat.RunePriceUSD()

	var poolResponse oapigen.PoolResponse
	poolResponse = oapigen.PoolResponse(
		buildPoolDetail(pool, status, *aggregates, runePriceUsd))
	respJSON(w, poolResponse)
}

// returns string array
func jsonMembers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	query := r.URL.Query()

	var pool *string
	poolParam := query.Get("pool")
	if poolParam != "" {
		pool = &poolParam
		if !timeseries.PoolExists(*pool) {
			miderr.BadRequestF("Unknown pool: %s", *pool).ReportHTTP(w)
			return
		}

	}

	addrs, err := timeseries.GetMemberAddrs(r.Context(), pool)
	if err != nil {
		respError(w, r, err)
		return
	}
	result := oapigen.MembersResponse(addrs)
	respJSON(w, result)
}

func jsonMemberDetails(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	addr := ps[0].Value

	pools, err := timeseries.GetMemberPools(r.Context(), addr)
	if err != nil {
		respError(w, r, err)
		return
	}
	if len(pools) == 0 {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	respJSON(w, oapigen.MemberDetailsResponse{
		Pools: pools.ToOapigen(),
	})
}

func jsonStats(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()

	state := timeseries.Latest.GetState()
	now := db.NowSecond()
	window := db.Window{From: 0, Until: now}

	stakes, err := stat.StakesLookup(ctx, window)
	if err != nil {
		respError(w, r, err)
		return
	}
	unstakes, err := stat.UnstakesLookup(ctx, window)
	if err != nil {
		respError(w, r, err)
		return
	}
	swapsFromRune, err := stat.SwapsFromRuneLookup(ctx, window)
	if err != nil {
		respError(w, r, err)
		return
	}
	swapsToRune, err := stat.SwapsToRuneLookup(ctx, window)
	if err != nil {
		respError(w, r, err)
		return
	}

	window24h := db.Window{From: now - 24*60*60, Until: now}
	window30d := db.Window{From: now - 30*24*60*60, Until: now}

	dailySwapsFromRune, err := stat.SwapsFromRuneLookup(ctx, window24h)
	if err != nil {
		respError(w, r, err)
		return
	}
	dailySwapsToRune, err := stat.SwapsToRuneLookup(ctx, window24h)
	if err != nil {
		respError(w, r, err)
		return
	}
	monthlySwapsFromRune, err := stat.SwapsFromRuneLookup(ctx, window30d)
	if err != nil {
		respError(w, r, err)
		return
	}
	monthlySwapsToRune, err := stat.SwapsToRuneLookup(ctx, window30d)
	if err != nil {
		respError(w, r, err)
		return
	}

	var runeDepth int64
	for _, poolInfo := range state.Pools {
		runeDepth += poolInfo.RuneDepth
	}

	switchedRune, err := stat.SwitchedRune(ctx)
	if err != nil {
		respError(w, r, err)
		return
	}

	runePrice := stat.RunePriceUSD()

	// TODO(acsaba): validate/correct calculations:
	//   - UniqueSwapperCount is it correct to do fromRune+toRune with multichain? (Now overlap?)
	//   - Swap count with doubleswaps are counted twice?
	//   - Predecessor to AddLiquidityVolume was totalStaked, which was stakes-withdraws.
	//       Is the new one ok?
	//   - AddLiquidityVolume looks only on rune, doesn't work with assymetric.
	//   - consider adding 24h 30d and total for everything.
	respJSON(w, oapigen.StatsResponse{
		RuneDepth:                     intStr(runeDepth),
		SwitchedRune:                  intStr(switchedRune),
		RunePriceUSD:                  floatStr(runePrice),
		SwapVolume:                    intStr(swapsFromRune.RuneE8Total + swapsToRune.RuneE8Total),
		SwapCount24h:                  intStr(dailySwapsFromRune.TxCount + dailySwapsToRune.TxCount),
		SwapCount30d:                  intStr(monthlySwapsFromRune.TxCount + monthlySwapsToRune.TxCount),
		SwapCount:                     intStr(swapsFromRune.TxCount + swapsToRune.TxCount),
		ToAssetCount:                  intStr(swapsFromRune.TxCount),
		ToRuneCount:                   intStr(swapsToRune.TxCount),
		DailyActiveUsers:              intStr(dailySwapsFromRune.RuneAddrCount + dailySwapsToRune.RuneAddrCount),
		MonthlyActiveUsers:            intStr(monthlySwapsFromRune.RuneAddrCount + monthlySwapsToRune.RuneAddrCount),
		UniqueSwapperCount:            intStr(swapsFromRune.RuneAddrCount + swapsToRune.RuneAddrCount),
		AddLiquidityVolume:            intStr(stakes.TotalVolume),
		WithdrawVolume:                intStr(unstakes.TotalVolume),
		ImpermanentLossProtectionPaid: intStr(unstakes.ImpermanentLossProtection),
		AddLiquidityCount:             intStr(stakes.Count),
		WithdrawCount:                 intStr(unstakes.Count),
	})
	/* TODO(pascaldekloe)
	   "poolCount":"20",
	   "totalEarned":"1827445688454",
	   "totalVolume24hr":"37756279870656",
	*/
}

func jsonActions(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	urlParams := r.URL.Query()
	params := timeseries.ActionsParams{
		Limit:      urlParams.Get("limit"),
		Offset:     urlParams.Get("offset"),
		ActionType: urlParams.Get("type"),
		Address:    urlParams.Get("address"),
		TXId:       urlParams.Get("txid"),
		Asset:      urlParams.Get("asset"),
	}

	// Get results
	actions, err := timeseries.GetActions(r.Context(), time.Time{}, params)
	// Send response
	if err != nil {
		respError(w, r, err)
		return
	}
	respJSON(w, actions)
}

func jsonSwagger(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	swagger, err := oapigen.GetSwagger()
	if err != nil {
		respError(w, r, err)
		return
	}
	respJSON(w, swagger)
}

func respJSON(w http.ResponseWriter, body interface{}) {
	w.Header().Set("Content-Type", "application/json")

	e := json.NewEncoder(w)
	e.SetIndent("", "\t")
	// Error discarded
	_ = e.Encode(body)
}

func respError(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

// IntStr returns the value as a decimal string.
// JSON numbers are double-precision floating-points.
// We don't want any unexpected rounding due to the 57-bit limit.
func intStr(v int64) string {
	return strconv.FormatInt(v, 10)
}

func intArrayStrs(a []int64) []string {
	b := make([]string, len(a))
	for i, v := range a {
		b[i] = intStr(v)
	}
	return b
}

func ratioStr(a, b int64) string {
	if b == 0 {
		return "0"
	} else {
		return strconv.FormatFloat(float64(a)/float64(b), 'f', -1, 64)
	}
}

func floatStr(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
