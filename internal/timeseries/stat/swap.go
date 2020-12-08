package stat

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"gitlab.com/thorchain/midgard/internal/db"
)

// Swaps are generic swap statistics.
type Swaps struct {
	TxCount       int64
	RuneAddrCount int64 // Number of unique addresses involved.
	RuneE8Total   int64
}

func SwapsFromRuneLookup(ctx context.Context, w db.Window) (*Swaps, error) {
	const q = `SELECT COALESCE(COUNT(*), 0), COALESCE(COUNT(DISTINCT(from_addr)), 0), COALESCE(SUM(from_E8), 0)
        FROM swap_events
        WHERE pool = from_asset AND block_timestamp >= $1 AND block_timestamp <= $2`

	return querySwaps(ctx, q, w.From.ToNano(), w.Until.ToNano())
}

// TODO(acsaba): change graphql to use the same as json and probably delete this.
func SwapsToRuneLookup(ctx context.Context, w db.Window) (*Swaps, error) {
	const q = `SELECT COALESCE(COUNT(*), 0), COALESCE(COUNT(DISTINCT(swap.from_addr)), 0), COALESCE(SUM(out.asset_E8), 0)
        FROM swap_events swap
	JOIN outbound_events out ON
		/* limit comparison set—no indinces */
		swap.block_timestamp <= out.block_timestamp AND
		swap.block_timestamp + 36000000000000 >= out.block_timestamp AND
		swap.tx = out.in_tx
        WHERE swap.block_timestamp >= $1 AND swap.block_timestamp <= $2 AND swap.pool <> swap.from_asset`

	return querySwaps(ctx, q, w.From.ToNano(), w.Until.ToNano())
}

func querySwaps(ctx context.Context, q string, args ...interface{}) (*Swaps, error) {
	rows, err := db.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var swaps Swaps
	if rows.Next() {
		err := rows.Scan(&swaps.TxCount, &swaps.RuneAddrCount, &swaps.RuneE8Total)
		if err != nil {
			return nil, err
		}
	}
	return &swaps, rows.Err()
}

type SwapBucket struct {
	StartTime     db.Second
	EndTime       db.Second
	ToAssetCount  int64
	ToRuneCount   int64
	TotalCount    int64
	ToAssetVolume int64
	ToRuneVolume  int64
	TotalVolume   int64
	TotalFees     int64
	TotalSlip     int64
}

func (meta *SwapBucket) AddBucket(bucket SwapBucket) {
	meta.ToAssetCount += bucket.ToAssetCount
	meta.ToRuneCount += bucket.ToRuneCount
	meta.TotalCount += bucket.TotalCount
	meta.ToAssetVolume += bucket.ToAssetVolume
	meta.ToRuneVolume += bucket.ToRuneVolume
	meta.TotalVolume += bucket.TotalVolume
	meta.TotalFees += bucket.TotalFees
	meta.TotalSlip += bucket.TotalSlip
}

type oneDirectionSwapBucket struct {
	Time         db.Second
	Count        int64
	VolumeInRune int64
	TotalFees    int64
	TotalSlip    int64
}

// Returns sparse buckets, when there are no swaps in the bucket, the bucket is missing.
func getSwapBuckets(ctx context.Context, pool string, interval db.Interval, w db.Window, swapToAsset bool) ([]oneDirectionSwapBucket, error) {
	queryArguments := []interface{}{w.From.ToNano(), w.Until.ToNano(), db.DBIntervalName[interval]}

	var poolFilter string
	if pool != "*" {
		poolFilter = `swap.pool = $4 AND`
		queryArguments = append(queryArguments, pool)
	}

	var directionFilter, volume string
	if swapToAsset {
		// from rune to asset
		volume = `COALESCE(SUM(from_E8), 0)`
		directionFilter = ` from_asset <> pool`
	} else {
		// from asset to Rune
		volume = `COALESCE(SUM(to_e8), 0) + COALESCE(SUM(liq_fee_in_rune_e8), 0)`
		directionFilter = ` from_asset = pool`
	}

	q := fmt.Sprintf(`
		SELECT
			` + db.QuerySelectTimestampInSecondsForInterval("swap.block_timestamp", "$3") + ` AS time,
			COALESCE(COUNT(*), 0) AS count,
			` + volume + ` AS volume,
			COALESCE(SUM(liq_fee_in_rune_E8), 0) AS fee,
			COALESCE(SUM(trade_slip_bp), 0) AS slip
		FROM swap_events AS swap
		WHERE ` + poolFilter + directionFilter + `
		    AND block_timestamp >= $1 AND block_timestamp < $2
		GROUP BY time
		ORDER BY time ASC`,
	)

	rows, err := db.Query(ctx, q, queryArguments...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ret := []oneDirectionSwapBucket{}
	for rows.Next() {
		var bucket oneDirectionSwapBucket
		err := rows.Scan(&bucket.Time, &bucket.Count, &bucket.VolumeInRune, &bucket.TotalFees, &bucket.TotalSlip)
		if err != nil {
			return []oneDirectionSwapBucket{}, err
		}
		ret = append(ret, bucket)
	}
	return ret, rows.Err()
}

// Returns gapfilled PoolSwaps for given pool, window and interval
func GetPoolSwaps(ctx context.Context, pool string, buckets db.Buckets) ([]SwapBucket, error) {

	toAsset, err := getSwapBuckets(ctx, pool, buckets.Interval, buckets.Window(), true)
	if err != nil {
		return nil, err
	}

	toRune, err := getSwapBuckets(ctx, pool, buckets.Interval, buckets.Window(), false)
	if err != nil {
		return nil, err
	}

	return mergeSwapsGapfill(buckets.Timestamps, toAsset, toRune), nil
}

func intStr(v int64) string {
	return strconv.FormatInt(v, 10)
}

func mergeSwapsGapfill(timestamps []db.Second, toAsset, toRune []oneDirectionSwapBucket) []SwapBucket {
	ret := make([]SwapBucket, len(timestamps)-1)

	timeAfterLast := timestamps[len(timestamps)-1] + 1
	toAsset = append(toAsset, oneDirectionSwapBucket{Time: timeAfterLast})
	toRune = append(toRune, oneDirectionSwapBucket{Time: timeAfterLast})

	for i, trIdx, taIdx := 0, 0, 0; i < len(timestamps)-1; i++ {
		current := &ret[i]
		current.StartTime = timestamps[i]
		current.EndTime = timestamps[i+1]
		ta := toAsset[taIdx]
		tr := toRune[trIdx]

		if current.StartTime == ta.Time {
			// We have swap to Asset in this bucket
			current.ToAssetCount = ta.Count
			current.ToAssetVolume = ta.VolumeInRune
			current.TotalFees += ta.TotalFees
			current.TotalSlip += ta.TotalSlip
			taIdx++
		}
		if current.StartTime == tr.Time {
			// We have swap to Rune in this bucket
			current.ToRuneCount = tr.Count
			current.ToRuneVolume = tr.VolumeInRune
			current.TotalFees += tr.TotalFees
			current.TotalSlip += tr.TotalSlip
			trIdx++
		}
		current.TotalCount = current.ToAssetCount + current.ToRuneCount
		current.TotalVolume = current.ToAssetVolume + current.ToRuneVolume
	}

	return ret
}

// PoolsTotalVolume computes total volume amount for given timestamps (from/to) and pools
// TODO(acsaba): replace this with event based volume. Maybe call previous with interval=NONE.
// TODO(acsaba): check that this result is consistent with interval search.
func PoolsTotalVolume(ctx context.Context, pools []string, from, to time.Time) (map[string]int64, error) {
	toRuneVolumeQ := `SELECT pool,
		COALESCE(CAST(SUM(CAST(rune_e8 as NUMERIC) / CAST(asset_e8 as NUMERIC) * swap.from_e8) as bigint), 0)
		FROM swap_events AS swap
			LEFT JOIN LATERAL (
				SELECT depths.asset_e8, depths.rune_e8
					FROM block_pool_depths as depths
				WHERE
				depths.block_timestamp <= swap.block_timestamp AND swap.pool = depths.pool
				ORDER BY depths.block_timestamp DESC
				LIMIT 1
			) AS joined on TRUE
		WHERE swap.from_asset = swap.pool AND swap.pool = ANY($1) AND swap.block_timestamp >= $2 AND swap.block_timestamp <= $3
		GROUP BY pool
	`
	toRuneRows, err := db.Query(ctx, toRuneVolumeQ, pools, from.UnixNano(), to.UnixNano())
	if err != nil {
		return nil, err
	}
	defer toRuneRows.Close()

	poolVolumes := make(map[string]int64)
	for toRuneRows.Next() {
		var pool string
		var volume int64
		err := toRuneRows.Scan(&pool, &volume)
		if err != nil {
			return nil, err
		}
		poolVolumes[pool] = volume
	}

	fromRuneVolumeQ := `SELECT pool, COALESCE(SUM(from_e8), 0)
	FROM swap_events
	WHERE from_asset <> pool AND pool = ANY($1) AND block_timestamp >= $2 AND block_timestamp <= $3
	GROUP BY pool
	`
	fromRuneRows, err := db.Query(ctx, fromRuneVolumeQ, pools, from.UnixNano(), to.UnixNano())
	if err != nil {
		return nil, err
	}
	defer fromRuneRows.Close()
	for fromRuneRows.Next() {
		var pool string
		var volume int64
		err := fromRuneRows.Scan(&pool, &volume)
		if err != nil {
			return nil, err
		}
		poolVolumes[pool] = poolVolumes[pool] + volume
	}

	return poolVolumes, nil
}
