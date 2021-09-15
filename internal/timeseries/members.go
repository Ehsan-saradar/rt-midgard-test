package timeseries

import (
	"context"
	"fmt"

	"gitlab.com/thorchain/midgard/internal/db"
	"gitlab.com/thorchain/midgard/internal/fetch/record"
	"gitlab.com/thorchain/midgard/internal/util"
	"gitlab.com/thorchain/midgard/openapi/generated/oapigen"
)

// Represents membership in a pool
type membership struct {
	runeAddress    string
	assetAddress   string
	liquidityUnits int64
}

type addrIndex map[string](map[string]*membership)

func (index addrIndex) getMembership(address, pool string) (*membership, bool) {
	_, ok := index[address]
	if ok {
		ret, ok := index[address][pool]
		return ret, ok
	} else {
		return nil, false
	}
}

func (index addrIndex) setMembership(address, pool string, newMembership *membership) {
	if index[address] == nil {
		index[address] = make(map[string]*membership)
	}
	index[address][pool] = newMembership
}

// MemberAddrs gets all member known addresses.
// When there's a rune/asset address pair or a rune addres for the member,
// the rune asset is shown.
// Else the asset address is shown.
// If an address participates in multiple pools it will be shown only once
func GetMemberAddrs(ctx context.Context, pool *string) (addrs []string, err error) {
	// Build indexes: nested maps -> address and pools for each address as keys
	// Needed to access each member from any address and also to identify unique addresses

	// runeAddrIndex: all memberships with a rune address
	// using the rune address as key
	runeAddrIndex := make(addrIndex)

	// asymAddrIndex: all memberships with only an asset address
	// none of the pointes here should be stored in runeAddrIndex
	// A single asset address can stake in different pools
	// (E.g.: ETH address in mutiple ERC20 tokens)
	asymAssetAddrIndex := make(addrIndex)

	poolFilter := ""
	qargs := []interface{}{}
	if pool != nil {
		poolFilter = "pool = $1"
		qargs = append(qargs, pool)
	}

	// Rune asset queryies. If a liquidity provider has a rune address then it is identified
	// by its rune address.
	// NOTE: Assumes only a single asset address per pool can be paired with a single rune
	// address
	runeALQ := `
		SELECT
			rune_addr,
			COALESCE(MAX(asset_addr), ''),
			pool,
			SUM(stake_units) as liquidity_units
		FROM stake_events
		` + db.Where(poolFilter, "rune_addr IS NOT NULL") + `
		GROUP BY rune_addr, pool
	`
	runeALRows, err := db.Query(ctx, runeALQ, qargs...)
	if err != nil {
		return nil, err
	}
	defer runeALRows.Close()

	for runeALRows.Next() {
		var newMembership membership
		var pool string
		err := runeALRows.Scan(
			&newMembership.runeAddress,
			&newMembership.assetAddress,
			&pool,
			&newMembership.liquidityUnits)
		if err != nil {
			return nil, err
		}
		runeAddrIndex.setMembership(newMembership.runeAddress, pool, &newMembership)
	}

	// Asymmetrical addLiquidity with asset only
	// part of asym membership (as if there was a rune address present, the liquidity provider
	// would be matched using the rune address)
	asymAssetALQ := `
		SELECT
			asset_addr,
			pool,
			SUM(stake_units) as liquidity_units
		FROM stake_events
		` + db.Where(poolFilter, "asset_addr IS NOT NULL AND rune_addr IS NULL") + `
		GROUP BY asset_addr, pool
	`

	asymAssetALRows, err := db.Query(ctx, asymAssetALQ, qargs...)
	if err != nil {
		return nil, err
	}
	defer asymAssetALRows.Close()
	for asymAssetALRows.Next() {
		var assetAddress, pool string
		var liquidityUnits int64
		err := asymAssetALRows.Scan(&assetAddress, &pool, &liquidityUnits)
		if err != nil {
			return nil, err
		}
		newMembership := membership{
			assetAddress:   assetAddress,
			liquidityUnits: liquidityUnits,
		}
		asymAssetAddrIndex.setMembership(assetAddress, pool, &newMembership)
	}

	// Withdraws: try matching from address to a membreship from
	// the index and subtract addLiquidityUnits.
	// If there's no match either there's an error with the
	// implementation or the Thorchain events.
	withdrawQ := `
		SELECT
			from_addr,
			pool,
			SUM(stake_units) as liquidity_units
		FROM unstake_events
		` + db.Where(poolFilter) + `
		GROUP BY from_addr, pool
	`
	withdrawRows, err := db.Query(ctx, withdrawQ, qargs...)
	if err != nil {
		return nil, err
	}
	defer withdrawRows.Close()

	for withdrawRows.Next() {
		var fromAddr, pool string
		var liquidityUnits int64
		err := withdrawRows.Scan(&fromAddr, &pool, &liquidityUnits)
		if err != nil {
			return nil, err
		}

		existingMembership, ok := runeAddrIndex.getMembership(fromAddr, pool)
		if ok && (existingMembership.runeAddress == fromAddr) {
			existingMembership.liquidityUnits -= liquidityUnits
			continue
		}

		existingMembership, ok = asymAssetAddrIndex.getMembership(fromAddr, pool)
		if ok && (existingMembership.assetAddress == fromAddr) {
			existingMembership.liquidityUnits -= liquidityUnits
			continue
		}

		return nil, fmt.Errorf("Address %s, pool %s, found in withdraw events should have a matching membership", fromAddr, pool)
	}

	// Lookup membership addresses:
	// Either in runeIndex or asymIndex with at least one pool
	// with positive liquidityUnits balance
	addrs = make([]string, 0, len(runeAddrIndex)+len(asymAssetAddrIndex))

	for address, poolMemberships := range runeAddrIndex {
		// if it has at least a non zero balance, add it to the result
		isMember := false
		for _, memb := range poolMemberships {
			if memb.liquidityUnits > 0 {
				isMember = true
				break
			}
		}

		if isMember {
			addrs = append(addrs, address)
		}
	}

	for address, poolMemberships := range asymAssetAddrIndex {
		// if it has at least a non zero balance, add it to the result
		isMember := false
		for _, memb := range poolMemberships {
			if memb.liquidityUnits > 0 {
				isMember = true
				break
			}
		}

		if isMember {
			addrs = append(addrs, address)
		}
	}

	return addrs, nil
}

// Info of a member in a specific pool.
type MemberPool struct {
	Pool           string
	RuneAddress    string
	AssetAddress   string
	LiquidityUnits int64
	RuneAdded      int64
	AssetAdded     int64
	RunePending    int64
	AssetPending   int64
	DateFirstAdded int64
	DateLastAdded  int64
	RuneWithdrawn  int64
	AssetWithdrawn int64
}

func (memberPool MemberPool) toOapigen() oapigen.MemberPool {
	return oapigen.MemberPool{
		Pool:           memberPool.Pool,
		RuneAddress:    memberPool.RuneAddress,
		AssetAddress:   memberPool.AssetAddress,
		LiquidityUnits: util.IntStr(memberPool.LiquidityUnits),
		RuneAdded:      util.IntStr(memberPool.RuneAdded),
		AssetAdded:     util.IntStr(memberPool.AssetAdded),
		RuneWithdrawn:  util.IntStr(memberPool.RuneWithdrawn),
		AssetWithdrawn: util.IntStr(memberPool.AssetWithdrawn),
		RunePending:    util.IntStr(memberPool.RunePending),
		AssetPending:   util.IntStr(memberPool.AssetPending),
		DateFirstAdded: util.IntStr(memberPool.DateFirstAdded),
		DateLastAdded:  util.IntStr(memberPool.DateLastAdded),
	}
}

// Pools data associated with a single member
type MemberPools []MemberPool

func (memberPools MemberPools) ToOapigen() []oapigen.MemberPool {
	ret := make([]oapigen.MemberPool, len(memberPools))
	for i, memberPool := range memberPools {
		ret[i] = memberPool.toOapigen()
	}

	return ret
}

func GetMemberPools(ctx context.Context, address string) (MemberPools, error) {
	if record.AddressIsRune(address) {
		return memberDetailsRune(ctx, address)
	} else {
		return memberDetailsAsset(ctx, address)
	}
}

const mpAddLiquidityQFields = `
		COALESCE(SUM(asset_E8), 0),
		COALESCE(SUM(rune_E8), 0),
		COALESCE(SUM(stake_units), 0),
		COALESCE(MIN(block_timestamp) / 1000000000, 0),
		COALESCE(MAX(block_timestamp) / 1000000000, 0)
`

const mpWithdrawQFields = `
		COALESCE(SUM(emit_asset_e8), 0),
		COALESCE(SUM(emit_rune_e8), 0),
		COALESCE(SUM(stake_units), 0)
`

const mpPendingQFields = `
		COALESCE(SUM(asset_e8), 0),
		COALESCE(SUM(rune_e8), 0)
`

// RUNE addresses
func memberDetailsRune(ctx context.Context, runeAddress string) (MemberPools, error) {
	// Aggregate the add- and withdraw- liquidity events. Conceptually we need to
	// union the stake_events and unstake_events tables and aggregate the add
	// and withdrawal amounts grouping by pool and member id. In practice the
	// query gets a bit complicated because it needs to account for situations
	// like the following:
	//   1. liquidity is added symmetrically
	//   2. all of the asset is withdrawn
	// In this case, the asset address should be forgotten. To achieve this,
	// the events are assigned a partition number which is incremented each time
	// all assets are withdrawn (i.e. basis_points=10000 for the withdrawal event).
	// Then, when the events are aggregated, they are grouped by pool, rune_address
	// and partition number, with only the rows with the highest partition number
	// for each pool/rune_address group returned.
	rolledUp := `
select distinct on(pool, rune_addr)
	pool,
	coalesce(last_value(asset_addr) over wnd, ''),
	coalesce(last_value(added_asset_e8) over wnd, 0),
	coalesce(last_value(added_rune_e8) over wnd, 0),
	coalesce(last_value(withdrawn_asset_e8) over wnd, 0),
    coalesce(last_value(withdrawn_rune_e8) over wnd, 0),
	coalesce(last_value(added_stake) over wnd, 0) -
	    coalesce(last_value(withdrawn_stake) over wnd, 0),
	coalesce(min_add_timestamp / 1000000000, 0),
	coalesce(max_add_timestamp / 1000000000, 0)
from (
	select
		pool,
		rune_addr,
		min(asset_addr) as asset_addr,
		asset_addr_partition,
		sum(added_asset_e8) as added_asset_e8,
        sum(added_rune_e8) as added_rune_e8,
		sum(added_stake) as added_stake,
		sum(withdrawn_asset_e8) as withdrawn_asset_e8,
		sum(withdrawn_rune_e8) as withdrawn_rune_e8,
		sum(withdrawn_stake) as withdrawn_stake,
		min(add_timestamp) as min_add_timestamp,
		max(add_timestamp) as max_add_timestamp
	from (
		select
			coalesce(stake.pool, unstake.pool) as pool,
			stake.block_timestamp as add_timestamp,
			coalesce(stake.rune_addr, unstake.from_addr) as rune_addr,
			asset_addr,
			stake.rune_e8 as added_rune_e8,
			stake.asset_e8 as added_asset_e8,
			stake.stake_units as added_stake,
			unstake.emit_rune_e8 as withdrawn_rune_e8,
			unstake.emit_asset_e8 as withdrawn_asset_e8,
			unstake.stake_units as withdrawn_stake,
			coalesce(
				sum(case when unstake.basis_points = 10000 then 1 else 0 end)
				over (partition by coalesce(stake.pool, unstake.pool),
					               coalesce(stake.rune_addr, unstake.from_addr)
					order by coalesce(stake.block_timestamp, unstake.block_timestamp)
					rows between unbounded preceding and 1 preceding), 0) as asset_addr_partition
		from midgard.stake_events as stake full outer join
		   (select * from midgard.unstake_events) as unstake
		on stake.block_timestamp = unstake.block_timestamp
		order by pool, coalesce(stake.block_timestamp, unstake.block_timestamp) asc) as timeseries
	group by pool, rune_addr, asset_addr_partition
	order by pool, asset_addr_partition) as rolled_up
where rune_addr = $1
window wnd as (partition by pool, rune_addr order by asset_addr_partition
				rows between unbounded preceding and unbounded following)`

	rows, err := db.Query(ctx, rolledUp, runeAddress)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	memberPoolsMap := make(map[string]MemberPool)

	for rows.Next() {
		memberPool := MemberPool{}
		err := rows.Scan(
			&memberPool.Pool,
			&memberPool.AssetAddress,
			&memberPool.AssetAdded,
			&memberPool.RuneAdded,
			&memberPool.AssetWithdrawn,
			&memberPool.RuneWithdrawn,
			&memberPool.LiquidityUnits,
			&memberPool.DateFirstAdded,
			&memberPool.DateLastAdded,
		)
		if err != nil {
			return nil, err
		}

		memberPool.RuneAddress = runeAddress
		memberPoolsMap[memberPool.Pool] = memberPool
	}

	pendingLiquidityQ := `SELECT
		pool,
	` + mpPendingQFields + `
	FROM midgard_agg.pending_adds
	WHERE rune_addr = $1
	GROUP BY pool`

	pendingLiquidityRows, err := db.Query(ctx, pendingLiquidityQ, runeAddress)
	if err != nil {
		return nil, err
	}
	defer pendingLiquidityRows.Close()

	for pendingLiquidityRows.Next() {
		var pool string
		var assetE8, runeE8 int64

		err := pendingLiquidityRows.Scan(
			&pool,
			&assetE8,
			&runeE8,
		)
		if err != nil {
			return nil, err
		}

		memberPool, ok := memberPoolsMap[pool]
		if !ok {
			memberPool.Pool = pool
			memberPool.RuneAddress = runeAddress
		}

		memberPool.AssetPending = assetE8
		memberPool.RunePending = runeE8
		memberPoolsMap[memberPool.Pool] = memberPool
	}

	ret := make(MemberPools, 0, len(memberPoolsMap))
	for _, memberPool := range memberPoolsMap {
		if memberPool.LiquidityUnits > 0 ||
			0 < memberPool.AssetPending || 0 < memberPool.RunePending {
			ret = append(ret, memberPool)
		}
	}

	return ret, nil
}

func memberDetailsAsset(ctx context.Context, assetAddress string) (MemberPools, error) {
	// Get all the rune addresses the asset address is paired with
	addressesQ := `SELECT
		se.pool,
		COALESCE(se.rune_addr, '') as pair_rune_addr
	FROM stake_events AS se
	WHERE se.asset_addr = $1
	GROUP BY pool, pair_rune_addr
	`

	addressesRows, err := db.Query(ctx, addressesQ, assetAddress)
	if err != nil {
		return nil, err
	}
	defer addressesRows.Close()

	var memberPools MemberPools
	for addressesRows.Next() {
		memberPool := MemberPool{AssetAddress: assetAddress}
		err := addressesRows.Scan(&memberPool.Pool, &memberPool.RuneAddress)

		var whereAddLiquidityAddresses, queryAddress string
		if memberPool.RuneAddress == "" {
			// asym liquidity provider, asset address is used to identify it
			// (if there is a rune_addr it will always be used to get the lp so it has to be NULL)
			whereAddLiquidityAddresses = "WHERE asset_addr = $1 AND rune_addr IS NULL"
			queryAddress = memberPool.AssetAddress
		} else {
			// sym liquidity provider, rune address is used to identify it
			whereAddLiquidityAddresses = "WHERE rune_addr = $1"
			queryAddress = memberPool.RuneAddress
		}

		addLiquidityQ := `SELECT ` + mpAddLiquidityQFields + `FROM stake_events ` + whereAddLiquidityAddresses + ` AND pool = $2`

		addLiquidityRow, err := db.Query(ctx, addLiquidityQ, queryAddress, memberPool.Pool)
		if err != nil {
			return nil, err
		}
		defer addLiquidityRow.Close()
		if addLiquidityRow.Next() {
			err := addLiquidityRow.Scan(&memberPool.AssetAdded, &memberPool.RuneAdded, &memberPool.LiquidityUnits, &memberPool.DateFirstAdded, &memberPool.DateLastAdded)
			if err != nil {
				return nil, err
			}
		}

		pendingLiquidityQ := `SELECT ` + mpPendingQFields + `FROM midgard_agg.pending_adds ` + whereAddLiquidityAddresses + ` AND pool = $2`

		pendingLiquidityRow, err := db.Query(ctx, pendingLiquidityQ, queryAddress, memberPool.Pool)
		if err != nil {
			return nil, err
		}
		defer pendingLiquidityRow.Close()
		if pendingLiquidityRow.Next() {
			err := pendingLiquidityRow.Scan(&memberPool.AssetPending, &memberPool.RunePending)
			if err != nil {
				return nil, err
			}
		}

		withdrawQ := `SELECT ` + mpWithdrawQFields + ` FROM unstake_events WHERE from_addr=$1 AND pool=$2`
		withdrawRow, err := db.Query(ctx, withdrawQ, queryAddress, memberPool.Pool)
		if err != nil {
			return nil, err
		}
		defer withdrawRow.Close()
		if withdrawRow.Next() {
			var unitsWithdrawn int64
			err = withdrawRow.Scan(&memberPool.AssetWithdrawn, &memberPool.RuneWithdrawn, &unitsWithdrawn)
			if err != nil {
				return nil, err
			}
			memberPool.LiquidityUnits -= unitsWithdrawn
		}

		if memberPool.LiquidityUnits > 0 {
			memberPools = append(memberPools, memberPool)
		}
	}

	return memberPools, nil
}
