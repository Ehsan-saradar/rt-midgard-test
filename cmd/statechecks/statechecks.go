package main

// Automated check. To manually check values go here:
// https://testnet.thornode.thorchain.info/thorchain/pools
// https://testnet.midgard.thorchain.info/v2/pools

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"gitlab.com/thorchain/midgard/config"
	"gitlab.com/thorchain/midgard/internal/db"
	"gitlab.com/thorchain/midgard/internal/timeseries"
	"gitlab.com/thorchain/midgard/internal/timeseries/stat"
)

const usageStr = `Checks state at latest height.
Usage:
$ go run ./cmd/statechecks [--onlydepthdiff] [--nonodescheck] [--seachmin] config
`

func init() {
	flag.Usage = func() {
		fmt.Println(usageStr)
		flag.PrintDefaults()
	}
}

var OnlyStructuredDiff = flag.Bool("onlydepthdiff", false,
	"No binary search, only the latest depth differences in structured form.")

var NoNodesCheck = flag.Bool("nonodescheck", false,
	"Skip active node count and bonds check.")

var BinarySearchMin = flag.Int64("searchmin", 1,
	"Base of the binary search, a known good state.")

var (
	CheckUnits bool = true
	CheckBonds bool = false
)

type Pool struct {
	Pool        string `json:"asset"`
	AssetDepth  int64  `json:"balance_asset,string"`
	RuneDepth   int64  `json:"balance_rune,string"`
	SynthSupply int64  `json:"synth_supply,string"`
	LPUnits     int64  `json:"LP_units,string"`
	Status      string `json:"status"`
	Timestamp   db.Nano
}

func (pool Pool) String() string {
	return fmt.Sprintf("%s [Asset: %d, Rune: %d, Synth: %d, Units: %d]",
		pool.Pool, pool.AssetDepth, pool.RuneDepth, pool.SynthSupply, pool.LPUnits)
}

type State struct {
	Pools           map[string]Pool
	ActiveNodeCount int64
	TotalBonded     int64
}

type Node struct {
	Status      string `json:"status"`
	Address     string `json:"node_address"`
	Bond        string `json:"bond"`
	BondAddress string `json:"bond_address"`
}

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true})
	logrus.SetLevel(logrus.InfoLevel)
	// logrus.SetLevel(logrus.DebugLevel)

	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("Missing config argument!", flag.Args())
		flag.Usage()
		return
	}

	config.ReadGlobalFrom(flag.Arg(0))

	ctx := context.Background()

	db.Setup()

	db.LoadFirstBlockFromDB(ctx)

	lastHeight, lastTimestamp := getLastBlockFromDB(ctx)
	logrus.Infof("Latest height: %d, timestamp: %d", lastHeight, lastTimestamp)

	midgardState := getMidgardState(ctx, lastHeight, lastTimestamp)
	logrus.Debug("Pools checked: ", midgardState)

	thorNodeURL := config.Global.ThorChain.ThorNodeURL
	thornodeState := getThornodeState(ctx, thorNodeURL, lastHeight)

	if *OnlyStructuredDiff {
		reportStructuredDiff(midgardState, thornodeState)
	} else {
		problems := compareStates(midgardState, thornodeState)

		for _, pool := range problems.mismatchingPools {
			binarySearchPool(ctx, thorNodeURL, pool, *BinarySearchMin, lastHeight)
		}

		if problems.activeNodeCountError {
			binarySearchNodes(ctx, thorNodeURL, *BinarySearchMin, lastHeight)
		}

		if problems.bondError {
			BondDetails(ctx, thorNodeURL)
		}
	}
}

func getLastBlockFromDB(ctx context.Context) (lastHeight int64, lastTimestamp db.Nano) {
	logrus.Info("Getting latest recorded height...")
	lastHeightRows, err := db.Query(ctx, "SELECT height, timestamp from block_log order by height desc limit 1")
	if err != nil {
		logrus.Fatal(err)
	}
	defer lastHeightRows.Close()

	if lastHeightRows.Next() {
		err := lastHeightRows.Scan(&lastHeight, &lastTimestamp)
		if err != nil {
			logrus.Fatal(err)
		}
	}
	return
}

func getMidgardState(ctx context.Context, height int64, timestamp db.Nano) (state State) {
	logrus.Debug("Getting Midgard data at height: ", height, ", timestamp: ", timestamp)

	poolsQ := `
		SELECT asset FROM pool_events WHERE block_timestamp <= $1 GROUP BY asset ORDER BY asset`
	poolsRows, err := db.Query(ctx, poolsQ, timestamp)
	if err != nil {
		logrus.Fatal(err)
	}
	defer poolsRows.Close()

	poolsWithStatus, err := timeseries.GetPoolsStatuses(ctx, timestamp)
	if err != nil {
		logrus.Fatal(err)
	}

	state.Pools = map[string]Pool{}
	pools := []string{}

	for poolsRows.Next() {
		var poolName, status string

		err := poolsRows.Scan(&poolName)
		if err != nil {
			logrus.Fatal(err)
		}
		logrus.Debug("Fetching Midgard pool: ", poolName)

		status = poolsWithStatus[poolName]
		if status == "" {
			status = timeseries.DefaultPoolStatus
		}
		if status == "suspended" {
			continue
		}

		pool := midgardPoolAtHeight(ctx, poolName, height)
		pool.Status = status
		if 0 < pool.RuneDepth && 0 < pool.AssetDepth {
			state.Pools[pool.Pool] = pool
			pools = append(pools, pool.Pool)
		}
	}

	if !*NoNodesCheck {
		state.ActiveNodeCount, err = timeseries.ActiveNodeCount(ctx, timestamp)
		if err != nil {
			logrus.Fatal(err)
		}
		state.TotalBonded, err = stat.GetTotalBond(ctx, -1)
		if err != nil {
			logrus.Fatal(err)
		}
	}
	return
}

func queryThorNode(thorNodeUrl string, urlPath string, height int64, dest interface{}) {
	url := thorNodeUrl + urlPath
	if 0 < height {
		url += "?height=" + strconv.FormatInt(height, 10)
	}
	logrus.Debug("Querying thornode: ", url)
	resp, err := http.Get(url)
	if err != nil {
		logrus.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, dest)
	if err != nil {
		logrus.Warn("Json unmarshal error for url: ", url)
		logrus.Fatal(err)

	}
}

func getThornodeNodesInfo(ctx context.Context, thorNodeUrl string, height int64) (
	nodeCount int64, totalBonded int64,
) {
	if *NoNodesCheck {
		return 0, 0
	}

	var nodes []Node
	queryThorNode(thorNodeUrl, "/nodes", height, &nodes)
	for _, node := range nodes {
		if strings.ToLower(node.Status) == "active" {
			nodeCount++
		}
		bond, err := strconv.ParseInt(node.Bond, 10, 64)
		if err != nil {
			logrus.Fatal(err)
		}
		totalBonded += bond
	}
	return
}

// true if active
func allThornodeNodes(ctx context.Context, thorNodeUrl string, height int64) map[string]bool {
	if *NoNodesCheck {
		return map[string]bool{}
	}

	var nodes []Node
	queryThorNode(thorNodeUrl, "/nodes", height, &nodes)
	ret := map[string]bool{}
	for _, node := range nodes {
		ret[node.Address] = (strings.ToLower(node.Status) == "active")
	}
	return ret
}

func getThornodeState(ctx context.Context, thorNodeUrl string, height int64) (state State) {
	logrus.Debug("Getting ThorNode data...")

	var pools []Pool

	queryThorNode(thorNodeUrl, "/pools", height, &pools)
	state.Pools = map[string]Pool{}
	for _, pool := range pools {
		state.Pools[pool.Pool] = pool
	}

	state.ActiveNodeCount, state.TotalBonded = getThornodeNodesInfo(ctx, thorNodeUrl, height)
	return
}

func reportStructuredDiff(midgardState, thornodeState State) {
	existenceDiff := strings.Builder{}
	depthDiffs := strings.Builder{}
	for _, thornodePool := range thornodeState.Pools {
		midgardPool, ok := midgardState.Pools[thornodePool.Pool]
		delete(midgardState.Pools, thornodePool.Pool)
		if !ok {
			fmt.Fprintf(&existenceDiff, "%s - did not find pool in Midgard (Exists in Thornode)\n", thornodePool.Pool)
			continue
		}

		runeDiff := thornodePool.RuneDepth - midgardPool.RuneDepth
		assetDiff := thornodePool.AssetDepth - midgardPool.AssetDepth
		if runeDiff != 0 || assetDiff != 0 {
			fmt.Fprintf(
				&depthDiffs, `{"%s", %d, %d},`+"\n",
				thornodePool.Pool, runeDiff, assetDiff)
		}

	}

	for pool := range midgardState.Pools {
		fmt.Fprintf(&existenceDiff, "%s - did not find pool in Thornode (Exists in Midgard)\n", pool)
		continue
	}

	if existenceDiff.Len() != 0 {
		logrus.Warn("Pool existence differences:\n", existenceDiff.String())
	}
	if depthDiffs.Len() != 0 {
		logrus.Warn("Depth differences:\n", depthDiffs.String())
	}
}

type Problems struct {
	mismatchingPools     []string
	activeNodeCountError bool
	bondError            bool
}

func compareStates(midgardState, thornodeState State) (problems Problems) {
	mismatchingPools := map[string]bool{}

	errors := strings.Builder{}

	for _, thornodePool := range thornodeState.Pools {
		midgardPool, ok := midgardState.Pools[thornodePool.Pool]
		prompt := fmt.Sprintf("\t- [Pool:%s]:", thornodePool.Pool)
		delete(midgardState.Pools, thornodePool.Pool)
		if !ok {
			fmt.Fprintf(&errors, "%s Did not find pool in Midgard (Exists in Thornode)\n", prompt)
			continue
		}

		if midgardPool.RuneDepth != thornodePool.RuneDepth {
			mismatchingPools[thornodePool.Pool] = true
			fmt.Fprintf(
				&errors, "%s RUNE Depth mismatch Thornode: %d, Midgard: %d\n",
				prompt, thornodePool.RuneDepth, midgardPool.RuneDepth)

		}

		if midgardPool.AssetDepth != thornodePool.AssetDepth {
			mismatchingPools[thornodePool.Pool] = true
			fmt.Fprintf(
				&errors, "%s Asset Depth mismatch Thornode: %d, Midgard: %d\n",
				prompt, thornodePool.AssetDepth, midgardPool.AssetDepth)
		}

		if midgardPool.SynthSupply != thornodePool.SynthSupply {
			mismatchingPools[thornodePool.Pool] = true
			fmt.Fprintf(
				&errors, "%s Synth Supply mismatch Thornode: %d, Midgard: %d\n",
				prompt, thornodePool.SynthSupply, midgardPool.SynthSupply)
		}

		if CheckUnits && midgardPool.LPUnits != thornodePool.LPUnits {
			mismatchingPools[thornodePool.Pool] = true
			fmt.Fprintf(
				&errors, "%s Pool Units mismatch Thornode: %d, Midgard: %d\n",
				prompt, thornodePool.LPUnits, midgardPool.LPUnits)
		}

		if midgardPool.Status != strings.ToLower(thornodePool.Status) {
			fmt.Fprintf(&errors, "%s Status mismatch Thornode: %s, Midgard: %s\n",
				prompt, strings.ToLower(thornodePool.Status), midgardPool.Status)
		}
	}

	for pool := range midgardState.Pools {
		prompt := fmt.Sprintf("\t- [Pool:%s]:", pool)
		fmt.Fprintf(&errors, "%s Did not find pool in Thornode (Exists in Midgard)\n", prompt)
		continue
	}

	if thornodeState.ActiveNodeCount != midgardState.ActiveNodeCount {
		problems.activeNodeCountError = true
		fmt.Fprintf(
			&errors, "\t- [Nodes]: Active Node Count mismatch Thornode: %d, Midgard %d\n",
			thornodeState.ActiveNodeCount, midgardState.ActiveNodeCount)
	}

	if CheckBonds && thornodeState.TotalBonded != midgardState.TotalBonded {
		problems.bondError = true
		tBonded := thornodeState.TotalBonded
		mBonded := midgardState.TotalBonded
		fmt.Fprintf(
			&errors,
			"\t- [Bonded]: Total Bonded mismatch Thornode: %d Midgard %d MidgardExcess %.2f%%\n",
			tBonded, mBonded, 100*float64(mBonded-tBonded)/float64(tBonded))
	}

	if errors.Len() > 0 {
		logrus.Warnf("ERRORS where found\n%s", errors.String())
	} else {
		logrus.Infof("All state checks OK")
	}

	for pool := range mismatchingPools {
		problems.mismatchingPools = append(problems.mismatchingPools, pool)
	}
	sort.Strings(problems.mismatchingPools)
	return
}

func midgardPoolAtHeight(ctx context.Context, pool string, height int64) Pool {
	logrus.Debug("Getting Midgard data at height: ", height, " pool: ", pool)

	q := `
	SELECT block_log.timestamp, asset_e8, rune_e8, synth_e8
	FROM block_pool_depths
	INNER JOIN block_log
	ON block_pool_depths.block_timestamp <= block_log.timestamp
	WHERE height=$1 AND pool = $2
	ORDER BY block_timestamp DESC
	LIMIT 1
	`

	rows, err := db.Query(ctx, q, height, pool)
	if err != nil {
		logrus.Fatal(err)
	}
	defer rows.Close()

	ret := Pool{Pool: pool}
	if rows.Next() {
		err := rows.Scan(&ret.Timestamp, &ret.AssetDepth, &ret.RuneDepth, &ret.SynthSupply)
		if err != nil {
			logrus.Fatal(err)
		}
	}

	until := ret.Timestamp + 1
	unitsMap, err := stat.PoolsLiquidityUnitsBefore(ctx, []string{pool}, &until)
	if err != nil {
		logrus.Fatal(err)
	}
	ret.LPUnits = unitsMap[pool]

	return ret
}

func findTablesWithColumns(ctx context.Context, columnName string) map[string]bool {
	q := `
	SELECT
		table_name
	FROM information_schema.columns
	WHERE table_schema='midgard' and column_name=$1
	`
	rows, err := db.Query(ctx, q, columnName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer rows.Close()

	ret := map[string]bool{}
	for rows.Next() {
		var table string
		err := rows.Scan(&table)
		if err != nil {
			logrus.Fatal(err)
		}
		ret[table] = true
	}

	return ret
}

type EventTable struct {
	TableName      string
	PoolColumnName string // "pool" or "asset" or ""
}

func findEventTables(ctx context.Context) []EventTable {
	blockTimestampTables := findTablesWithColumns(ctx, "block_timestamp")
	blockTimestampTables["block_pool_depths"] = false

	poolTables := findTablesWithColumns(ctx, "pool")
	assetTables := findTablesWithColumns(ctx, "asset")
	ret := []EventTable{}
	for table := range blockTimestampTables {
		if poolTables[table] {
			ret = append(ret, EventTable{TableName: table, PoolColumnName: "pool"})
		} else if assetTables[table] {
			ret = append(ret, EventTable{TableName: table, PoolColumnName: "asset"})
		} else {
			ret = append(ret, EventTable{TableName: table, PoolColumnName: ""})
		}
	}
	return ret
}

var (
	eventTablesCache []EventTable
	eventTablesOnce  sync.Once
)

func getEventTables(ctx context.Context) []EventTable {
	eventTablesOnce.Do(func() {
		eventTablesCache = findEventTables(ctx)
	})
	return eventTablesCache
}

func logEventsFromTable(ctx context.Context, eventTable EventTable, pool string, timestamp db.Nano) {
	if eventTable.TableName == "message_events" || eventTable.TableName == "block_pool_depths" {
		return
	}
	poolFilters := []string{"block_timestamp = $1"}
	qargs := []interface{}{timestamp}
	if eventTable.PoolColumnName != "" {
		synthPool := strings.Replace(pool, ".", "/", -1)
		poolFilters = append(poolFilters,
			fmt.Sprintf("(%[1]s = $2) OR (%[1]s = $3)",
				eventTable.PoolColumnName))
		qargs = append(qargs, pool, synthPool)
	}

	q := `
	SELECT *
	FROM ` + eventTable.TableName + ` ` + db.Where(poolFilters...)

	rows, err := db.Query(ctx, q, qargs...)
	if err != nil {
		logrus.Fatal(err)
	}
	defer rows.Close()

	colNames, err := rows.Columns()
	if err != nil {
		panic(err)
	}

	eventNum := 0
	for rows.Next() {
		eventNum++
		colsPtr := make([]interface{}, len(colNames))
		for i := range colsPtr {
			var tmp interface{}
			colsPtr[i] = &tmp
		}
		err := rows.Scan(colsPtr...)
		if err != nil {
			logrus.Fatal(err)
		}
		buf := strings.Builder{}

		fmt.Fprintf(&buf, "%s [", eventTable.TableName)
		for i := range colNames {
			if i != 0 {
				fmt.Fprintf(&buf, ", ")
			}
			fmt.Fprintf(&buf, "%s: %v", colNames[i], *(colsPtr[i].(*interface{})))
		}
		fmt.Fprintf(&buf, "]")
		logrus.Infof(buf.String())
	}
}

func logAllEventsAtHeight(ctx context.Context, pool string, timestamp db.Nano) {
	eventTables := getEventTables(ctx)
	for _, eventTable := range eventTables {
		logEventsFromTable(ctx, eventTable, pool, timestamp)
	}
}

// Looks up the first difference in the (min, max) range. May choose max.
func binarySearchPool(ctx context.Context, thorNodeUrl string, pool string, minHeight, maxHeight int64) {
	logrus.Infof("=====  [%s] Binary searching in range [%d, %d)", pool, minHeight, maxHeight)

	for 1 < maxHeight-minHeight {
		middleHeight := (minHeight + maxHeight) / 2
		logrus.Debugf(
			"--- [%s] Binary search step [%d, %d] height: %d",
			pool, minHeight, maxHeight, middleHeight)
		var thorNodePool Pool
		queryThorNode(thorNodeUrl, "/pool/"+pool, middleHeight, &thorNodePool)
		logrus.Debug("Thornode: ", thorNodePool)
		midgardPool := midgardPoolAtHeight(ctx, pool, middleHeight)
		logrus.Debug("Midgard: ", midgardPool)
		ok := (thorNodePool.AssetDepth == midgardPool.AssetDepth &&
			thorNodePool.RuneDepth == midgardPool.RuneDepth &&
			thorNodePool.SynthSupply == midgardPool.SynthSupply &&
			(!CheckUnits || thorNodePool.LPUnits == midgardPool.LPUnits))
		if ok {
			logrus.Debug("Same at height ", middleHeight)
			minHeight = middleHeight
		} else {
			logrus.Debug("Differ at height ", middleHeight)
			maxHeight = middleHeight
		}
	}

	midgardPoolBefore := midgardPoolAtHeight(ctx, pool, maxHeight-1)

	var thorNodePool Pool
	queryThorNode(thorNodeUrl, "/pool/"+pool, maxHeight, &thorNodePool)
	midgardPool := midgardPoolAtHeight(ctx, pool, maxHeight)

	logrus.Infof("[%s] First differenct at height: %d timestamp: %d date: %s",
		pool, maxHeight, midgardPool.Timestamp,
		midgardPool.Timestamp.ToSecond().ToTime().Format("2006-01-02 15:04:05"))
	logrus.Info("Previous state:  ", midgardPoolBefore)
	logrus.Info("Thornode:        ", thorNodePool)
	logrus.Info("Midgard:         ", midgardPool)

	logWithPercent := func(msg string, diffValue int64, base int64) {
		percent := 100 * float64(diffValue) / float64(base)
		if base == 0 && diffValue == 0 {
			percent = 0
		}
		logrus.Infof("%s:  %d (%f%%)", msg, diffValue, percent)
	}
	logWithPercent("Midgard Asset excess",
		midgardPool.AssetDepth-thorNodePool.AssetDepth,
		midgardPoolBefore.AssetDepth)
	logWithPercent("Midgard Rune excess",
		midgardPool.RuneDepth-thorNodePool.RuneDepth,
		midgardPoolBefore.RuneDepth)
	logWithPercent("Midgard Synth excess",
		midgardPool.SynthSupply-thorNodePool.SynthSupply,
		midgardPoolBefore.SynthSupply)
	logWithPercent("Midgard Unit excess",
		midgardPool.LPUnits-thorNodePool.LPUnits,
		midgardPoolBefore.LPUnits)

	logAllEventsAtHeight(ctx, pool, midgardPool.Timestamp)
}

func timestampAtHeight(ctx context.Context, height int64) db.Nano {
	q := `
	SELECT timestamp
	FROM block_log
	WHERE height=$1
	`
	rows, err := db.Query(ctx, q, height)
	if err != nil {
		logrus.Fatal(err)
	}
	defer rows.Close()

	if !rows.Next() {
		logrus.Fatal("No rows selected:", q)
	}
	var ts db.Nano
	err = rows.Scan(&ts)
	if err != nil {
		logrus.Fatal(err)
	}
	return ts
}

func midgardActiveNodeCount(ctx context.Context, height int64) int64 {
	timestamp := timestampAtHeight(ctx, height)
	midgardCount, err := timeseries.ActiveNodeCount(ctx, timestamp)
	if err != nil {
		logrus.Fatal(err)
	}
	return midgardCount
}

func allMidgardNodes(ctx context.Context, height int64) map[string]bool {
	timestamp := timestampAtHeight(ctx, height)
	q := `
	SELECT
		node_addr,
		last(current, block_timestamp) AS status
	FROM update_node_account_status_events
	WHERE block_timestamp <= $1
	GROUP BY node_addr
	`
	rows, err := db.Query(ctx, q, timestamp)
	if err != nil {
		logrus.Fatal(err)
	}
	defer rows.Close()

	ret := map[string]bool{}
	for rows.Next() {
		var addr, status string
		err = rows.Scan(&addr, &status)
		if err != nil {
			logrus.Fatal(err)
		}
		logrus.Debug("Status: ", strings.ToLower(status))
		ret[addr] = (strings.ToLower(status) == "active")
	}
	return ret
}

func excessNodes(str string, a, b map[string]bool) {
	buf := strings.Builder{}
	hasdiff := false
	for node, status := range a {
		if !status {
			continue
		}
		status2b, ok := b[node]
		if !ok {
			fmt.Fprintf(&buf, "present: %s - ", node)
			hasdiff = true
		} else if status2b == false {
			fmt.Fprintf(&buf, "active: %s - ", node)
			hasdiff = true
		}
	}
	if hasdiff {
		logrus.Info(str, " excess: ", buf.String())
	} else {
		logrus.Info(str, " OK")
	}
}

// Looks up the first difference in the (min, max) range. May choose max.
func binarySearchNodes(ctx context.Context, thorNodeUrl string, minHeight, maxHeight int64) {
	logrus.Infof("=====  Binary searching active nodes in range [%d, %d)", minHeight, maxHeight)

	for 1 < maxHeight-minHeight {
		middleHeight := (minHeight + maxHeight) / 2
		logrus.Debugf(
			"--- Binary search step [%d, %d] height: %d",
			minHeight, maxHeight, middleHeight)
		thorNodeCount, _ := getThornodeNodesInfo(ctx, thorNodeUrl, middleHeight)
		logrus.Debug("Thornode: ", thorNodeCount)

		midgardCount := midgardActiveNodeCount(ctx, middleHeight)
		logrus.Debug("Midgard: ", midgardCount)
		ok := midgardCount == thorNodeCount
		if ok {
			logrus.Debug("Same at height ", middleHeight)
			minHeight = middleHeight
		} else {
			logrus.Debug("Differ at height ", middleHeight)
			maxHeight = middleHeight
		}
	}

	countBefore := midgardActiveNodeCount(ctx, maxHeight-1)

	thorNodeCount, _ := getThornodeNodesInfo(ctx, thorNodeUrl, maxHeight)
	midgardCount := midgardActiveNodeCount(ctx, maxHeight)

	logrus.Infof("First node differenct at height: %d timestamp: %d",
		maxHeight, timestampAtHeight(ctx, maxHeight))
	logrus.Info("Previous state:  ", countBefore)
	logrus.Info("Thornode:        ", thorNodeCount)
	logrus.Info("Midgard:         ", midgardCount)

	prevThornodeNodes := allThornodeNodes(ctx, thorNodeUrl, maxHeight-1)
	prevMidgardNodes := allMidgardNodes(ctx, maxHeight-1)
	excessNodes("previous thornode vs midgard", prevThornodeNodes, prevMidgardNodes)
	excessNodes("previous midgard vs thornode", prevMidgardNodes, prevThornodeNodes)

	curentThornodeNodes := allThornodeNodes(ctx, thorNodeUrl, maxHeight)
	curentMidgardNodes := allMidgardNodes(ctx, maxHeight)
	excessNodes("Current thornode vs midgard", curentThornodeNodes, curentMidgardNodes)
	excessNodes("Current midgard vs thornode", curentMidgardNodes, curentThornodeNodes)
}
