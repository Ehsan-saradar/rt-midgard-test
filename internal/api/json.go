package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"gitlab.com/thorchain/midgard/internal/db"
	"gitlab.com/thorchain/midgard/internal/graphql/model"

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

func jsonHealth(w http.ResponseWriter, r *http.Request) {
	height, _, _ := timeseries.LastBlock()
	synced := InSync()
	respJSON(w, oapigen.HealthResponse{
		InSync:        synced,
		Database:      true,
		ScannerHeight: intStr(height + 1),
	})
}

func jsonEarningsHistory(w http.ResponseWriter, r *http.Request) {
	buckets, err := db.BucketsFromQuery(r.Context(), r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := stat.GetEarningsHistory(r.Context(), buckets)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respJSON(w, res)
}

func jsonLiquidityHistory(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	buckets, err := db.BucketsFromQuery(r.Context(), query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pool := query.Get("pool")
	if pool == "" {
		pool = "*"
	}

	res, err := stat.GetLiquidityHistory(r.Context(), buckets, pool)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respJSON(w, res)
}

func jsonDepths(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	buckets, err := db.BucketsFromQuery(r.Context(), query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO(acsaba): check if pool exists.
	pool := path.Base(r.URL.Path)

	res, err := stat.PoolDepthHistory(r.Context(), buckets, pool)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respJSON(w, toOapiDepthResponse(res))
}

func toOapiDepthResponse(buckets []stat.PoolDepthBucket) (result oapigen.DepthHistoryResponse) {
	result.Intervals = make(oapigen.DepthHistoryIntervals, 0, len(buckets))
	for _, bucket := range buckets {
		result.Intervals = append(result.Intervals, oapigen.DepthHistoryItem{
			StartTime:  intStr(bucket.StartTime.ToI()),
			EndTime:    intStr(bucket.EndTime.ToI()),
			AssetDepth: intStr(bucket.AssetDepth),
			RuneDepth:  intStr(bucket.RuneDepth),
		})
	}
	result.Meta.StartTime = intStr(buckets[0].StartTime.ToI())
	result.Meta.EndTime = intStr(buckets[len(buckets)-1].EndTime.ToI())
	return
}

func jsonSwapHistory(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	buckets, err := db.BucketsFromQuery(r.Context(), query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pool := query.Get("pool")
	if pool == "" {
		pool = "*"
	}

	res, err := swapHistory(r.Context(), buckets, pool)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respJSON(w, res)
}

func swapHistory(
	ctx context.Context,
	buckets db.Buckets,
	pool string) (oapigen.SwapHistoryResponse, error) {

	mergedPoolSwaps, err := stat.GetPoolSwaps(ctx, pool, buckets)
	if err != nil {
		return oapigen.SwapHistoryResponse{}, err
	}

	return createVolumeIntervals(mergedPoolSwaps), nil
}

func toSwapHistoryItem(bucket stat.SwapBucket) oapigen.SwapHistoryItem {
	var averageSlip float64 = 0
	if 0 < bucket.TotalCount {
		averageSlip = float64(bucket.TotalSlip) / float64(bucket.TotalCount)
	}
	return oapigen.SwapHistoryItem{
		StartTime:     intStr(bucket.StartTime.ToI()),
		EndTime:       intStr(bucket.EndTime.ToI()),
		ToRuneVolume:  intStr(bucket.ToRuneVolume),
		ToAssetVolume: intStr(bucket.ToAssetVolume),
		TotalVolume:   intStr(bucket.TotalVolume),
		ToAssetCount:  intStr(bucket.ToAssetCount),
		ToRuneCount:   intStr(bucket.ToRuneCount),
		TotalCount:    intStr(bucket.TotalCount),
		TotalFees:     intStr(bucket.TotalFees),
		AverageSlip:   floatStr(averageSlip),
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

func jsonNetwork(w http.ResponseWriter, r *http.Request) {
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

func jsonNodes(w http.ResponseWriter, r *http.Request) {
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

	array := make([]oapigen.NodeKey, 0, len(m))
	for _, e := range m {
		array = append(array, oapigen.NodeKey{
			Secp256k1: e.Secp,
			Ed25519:   e.Ed,
		})
	}
	respJSON(w, array)
}

func filteredPoolsByStatus(r *http.Request, statusMap map[string]string) ([]string, error) {
	pools, err := timeseries.Pools(r.Context())
	if err != nil {
		return nil, err
	}
	ret := pools
	statusParams := r.URL.Query()["status"]
	if len(statusParams) != 0 {
		const errormsg = "Max one status parameter, accepted values: available, staged, suspended"
		if 1 < len(statusParams) {
			return nil, fmt.Errorf(errormsg)
		}
		status := statusParams[0]
		status = strings.ToLower(status)
		// Allowed statuses in
		// https://gitlab.com/thorchain/thornode/-/blob/master/x/thorchain/types/type_pool.go
		if status != "available" && status != "staged" && status != "suspended" {
			return nil, fmt.Errorf(errormsg)
		}
		ret = []string{}
		for _, pool := range pools {
			poolStatus := poolStatusFromMap(pool, statusMap)
			if poolStatus == status {
				ret = append(ret, pool)
			}
		}
	}
	return ret, nil
}

type poolAggregates struct {
	dailyVolumes        map[string]int64
	poolUnits           map[string]int64
	poolWeeklyRewards   map[string]int64
	assetE8DepthPerPool map[string]int64
	runeE8DepthPerPool  map[string]int64
}

func getPoolAggregates(ctx context.Context, pools []string) (*poolAggregates, error) {
	assetE8DepthPerPool, runeE8DepthPerPool, timestamp := timeseries.AssetAndRuneDepths()

	dailyVolumes, err := stat.PoolsTotalVolume(ctx, pools, timestamp.Add(-24*time.Hour), timestamp)
	if err != nil {
		return nil, err
	}

	poolUnits, err := timeseries.PoolsUnits(ctx, pools)
	if err != nil {
		return nil, err
	}

	poolWeeklyRewards, err := timeseries.PoolsTotalIncome(ctx, pools, timestamp.Add(-1*time.Hour*24*7), timestamp)
	if err != nil {
		return nil, err
	}

	aggregates := poolAggregates{
		dailyVolumes:        dailyVolumes,
		poolUnits:           poolUnits,
		poolWeeklyRewards:   poolWeeklyRewards,
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

func buildPoolDetail(pool, status string, aggregates poolAggregates) oapigen.PoolDetail {
	assetDepth := aggregates.assetE8DepthPerPool[pool]
	runeDepth := aggregates.runeE8DepthPerPool[pool]
	dailyVolume := aggregates.dailyVolumes[pool]
	poolUnits := aggregates.poolUnits[pool]
	rewards := aggregates.poolWeeklyRewards[pool]
	poolAPY := timeseries.GetPoolAPY(runeDepth, rewards)
	var price float64
	if assetDepth != 0 {
		price = float64(runeDepth) / float64(assetDepth)
	}
	return oapigen.PoolDetail{
		Asset:      pool,
		AssetDepth: intStr(assetDepth),
		RuneDepth:  intStr(runeDepth),
		PoolAPY:    floatStr(poolAPY),
		AssetPrice: floatStr(price),
		Status:     status,
		Units:      intStr(poolUnits),
		Volume24h:  intStr(dailyVolume),
	}
}

func jsonPools(w http.ResponseWriter, r *http.Request) {
	_, lastTime, _ := timeseries.LastBlock()
	statusMap, err := timeseries.GetPoolsStatuses(r.Context(), db.Nano(lastTime.UnixNano()))
	if err != nil {
		respError(w, r, err)
		return
	}
	pools, err := filteredPoolsByStatus(r, statusMap)
	if err != nil {
		respError(w, r, err)
		return
	}

	aggregates, err := getPoolAggregates(r.Context(), pools)
	if err != nil {
		respError(w, r, err)
		return
	}

	poolsResponse := make(oapigen.PoolsResponse, len(pools))
	for i, pool := range pools {
		status := poolStatusFromMap(pool, statusMap)
		poolsResponse[i] = buildPoolDetail(pool, status, *aggregates)
	}

	respJSON(w, poolsResponse)
}

func jsonPool(w http.ResponseWriter, r *http.Request) {
	pool := path.Base(r.URL.Path)

	assetE8DepthPerPool, runeE8DepthPerPool, timestamp := timeseries.AssetAndRuneDepths()
	_, assetOk := assetE8DepthPerPool[pool]
	_, runeOk := runeE8DepthPerPool[pool]

	// Return not found if there's no track of the pool
	if !assetOk && !runeOk {
		http.NotFound(w, r)
		return
	}

	status, err := timeseries.PoolStatus(r.Context(), pool, timestamp)
	if err != nil {
		respError(w, r, err)
		return
	}

	aggregates, err := getPoolAggregates(r.Context(), []string{pool})
	if err != nil {
		respError(w, r, err)
		return
	}

	poolData := buildPoolDetail(pool, status, *aggregates)

	respJSON(w, poolData)
}

// returns string array
func jsonMembers(w http.ResponseWriter, r *http.Request) {
	addrs, err := timeseries.MemberAddrs(r.Context())
	if err != nil {
		respError(w, r, err)
		return
	}
	result := oapigen.MembersResponse(addrs)
	respJSON(w, result)
}

func jsonMemberDetails(w http.ResponseWriter, r *http.Request) {
	addr := path.Base(r.URL.Path)
	// TODO(elfedy): validate that the address is from the same chain as
	// the RUNE asset and return 400 if not

	poolsDeposits, err := stat.AddressPoolDepositsLookup(r.Context(), addr)
	if err != nil {
		respError(w, r, err)
		return
	}

	if len(poolsDeposits) == 0 {
		http.NotFound(w, r)
		return
	}

	poolsWithdrawals, err := stat.AddressPoolWithdrawalsLookup(r.Context(), addr)
	if err != nil {
		respError(w, r, err)
		return
	}

	var pools []oapigen.MemberPoolDetails
	for pool, poolDeposits := range poolsDeposits {
		poolWithdrawals := poolsWithdrawals[pool]

		detail := oapigen.MemberPoolDetails{
			AssetAdded:     intStr(poolDeposits.AssetE8Total),
			AssetWithdrawn: intStr(poolWithdrawals.AssetE8Total),
			DateFirstAdded: intStr(poolDeposits.DateFirstAdded),
			DateLastAdded:  intStr(poolDeposits.DateLastAdded),
			LiquidityUnits: intStr(poolDeposits.UnitsTotal - poolWithdrawals.UnitsTotal),
			Pool:           pool,
			RuneAdded:      intStr(poolDeposits.RuneE8Total),
			RuneWithdrawn:  intStr(poolWithdrawals.RuneE8Total),
		}

		pools = append(pools, detail)
	}

	respJSON(w, oapigen.MemberDetailsResponse{
		Pools: pools,
	})
}

func jsonStats(w http.ResponseWriter, r *http.Request) {
	_, runeE8DepthPerPool, timestamp := timeseries.AssetAndRuneDepths()
	window := db.Window{From: 0, Until: db.TimeToSecond(timestamp)}

	stakes, err := stat.StakesLookup(r.Context(), window)
	if err != nil {
		respError(w, r, err)
		return
	}
	unstakes, err := stat.UnstakesLookup(r.Context(), window)
	if err != nil {
		respError(w, r, err)
		return
	}
	swapsFromRune, err := stat.SwapsFromRuneLookup(r.Context(), window)
	if err != nil {
		respError(w, r, err)
		return
	}
	swapsToRune, err := stat.SwapsToRuneLookup(r.Context(), window)
	if err != nil {
		respError(w, r, err)
		return
	}
	tSec := db.TimeToSecond(timestamp)
	dailySwapsFromRune, err := stat.SwapsFromRuneLookup(r.Context(), db.Window{From: tSec.Add(-24 * time.Hour), Until: tSec})
	if err != nil {
		respError(w, r, err)
		return
	}
	dailySwapsToRune, err := stat.SwapsToRuneLookup(r.Context(), db.Window{From: tSec.Add(-24 * time.Hour), Until: tSec})
	if err != nil {
		respError(w, r, err)
		return
	}
	monthlySwapsFromRune, err := stat.SwapsFromRuneLookup(r.Context(), db.Window{From: tSec.Add(-30 * 24 * time.Hour), Until: tSec})
	if err != nil {
		respError(w, r, err)
		return
	}
	monthlySwapsToRune, err := stat.SwapsToRuneLookup(r.Context(), db.Window{From: tSec.Add(-30 * 24 * time.Hour), Until: tSec})
	if err != nil {
		respError(w, r, err)
		return
	}

	var runeDepth int64
	for _, depth := range runeE8DepthPerPool {
		runeDepth += depth
	}

	respJSON(w, oapigen.StatsResponse{
		DailyActiveUsers:   intStr(dailySwapsFromRune.RuneAddrCount + dailySwapsToRune.RuneAddrCount),
		DailyTx:            intStr(dailySwapsFromRune.TxCount + dailySwapsToRune.TxCount),
		MonthlyActiveUsers: intStr(monthlySwapsFromRune.RuneAddrCount + monthlySwapsToRune.RuneAddrCount),
		MonthlyTx:          intStr(monthlySwapsFromRune.TxCount + monthlySwapsToRune.TxCount),
		TotalAssetBuys:     intStr(swapsFromRune.TxCount),
		TotalAssetSells:    intStr(swapsToRune.TxCount),
		TotalDepth:         intStr(runeDepth),
		TotalUsers:         intStr(swapsFromRune.RuneAddrCount + swapsToRune.RuneAddrCount),
		TotalStakeTx:       intStr(stakes.TxCount + unstakes.TxCount),
		TotalStaked:        intStr(stakes.RuneE8Total - unstakes.RuneE8Total),
		TotalTx:            intStr(swapsFromRune.TxCount + swapsToRune.TxCount + stakes.TxCount + unstakes.TxCount),
		TotalVolume:        intStr(swapsFromRune.RuneE8Total + swapsToRune.RuneE8Total),
		TotalWithdrawTx:    intStr(unstakes.RuneE8Total),
	})
	/* TODO(pascaldekloe)
	   "poolCount":"20",
	   "totalEarned":"1827445688454",
	   "totalVolume24hr":"37756279870656",
	*/
}

func jsonTx(w http.ResponseWriter, r *http.Request) {
	// Parse params
	urlParams := r.URL.Query()
	lookupParamKeys := []string{"limit", "offset", "type", "address", "txid", "asset"}
	lookupParams := make(map[string]string)
	for _, key := range lookupParamKeys {
		val := ""
		if urlParams[key] != nil {
			val = urlParams[key][0]
		}
		lookupParams[key] = val
	}

	// Get results
	txs, err := timeseries.TxList(r.Context(), time.Time{}, lookupParams)

	// Send response
	if err != nil {
		respError(w, r, err)
		return
	}
	respJSON(w, txs)
}

func jsonSwagger(w http.ResponseWriter, r *http.Request) {
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
	// TODO(acsaba): logging HTTP errors somewhere else then stdout.
	// log.Printf("HTTP %q %q: %s", r.Method, r.URL.Path, err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

// IntStr returns the value as a decimal string.
// JSON numbers are double-precision floating-points.
// We don't want any unexpected rounding due to the 57-bit limit.
func intStr(v int64) string {
	return strconv.FormatInt(v, 10)
}

func intArrayStrs(a []*int64) []string {
	b := make([]string, len(a))
	for i, v := range a {
		b[i] = intStr(*v)
	}
	return b
}

func floatStr(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}