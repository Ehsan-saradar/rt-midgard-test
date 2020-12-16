// Package oapigen provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package oapigen

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// BlockRewards defines model for BlockRewards.
type BlockRewards struct {
	BlockReward string `json:"blockReward"`
	BondReward  string `json:"bondReward"`
	PoolReward  string `json:"poolReward"`
}

// BondMetrics defines model for BondMetrics.
type BondMetrics struct {

	// Int64, Average bond of active nodes.
	AverageActiveBond string `json:"averageActiveBond"`

	// Int64, Average bond of standby nodes
	AverageStandbyBond string `json:"averageStandbyBond"`

	// Int64, Maxinum bond of active nodes.
	MaximumActiveBond string `json:"maximumActiveBond"`

	// Int64, Maximum bond of standby nodes
	MaximumStandbyBond string `json:"maximumStandbyBond"`

	// Int64, Median bond of active nodes.
	MedianActiveBond string `json:"medianActiveBond"`

	// Int64, Median bond of standby nodes
	MedianStandbyBond string `json:"medianStandbyBond"`

	// Int64, Minumum bond of active nodes.
	MinimumActiveBond string `json:"minimumActiveBond"`

	// Int64, Minumum bond of standby nodes
	MinimumStandbyBond string `json:"minimumStandbyBond"`

	// Int64, Total bond of active nodes.
	TotalActiveBond string `json:"totalActiveBond"`

	// Int64, Total bond of standby nodes
	TotalStandbyBond string `json:"totalStandbyBond"`
}

// BoolConstants defines model for BoolConstants.
type BoolConstants struct {
	StrictBondLiquidityRatio bool `json:"StrictBondLiquidityRatio"`
}

// Constants defines model for Constants.
type Constants struct {
	BoolValues   BoolConstants   `json:"bool_values"`
	Int64Values  Int64Constants  `json:"int_64_values"`
	StringValues StringConstants `json:"string_values"`
}

// DepthHistory defines model for DepthHistory.
type DepthHistory struct {
	Intervals DepthHistoryIntervals `json:"intervals"`
	Meta      DepthHistoryMeta      `json:"meta"`
}

// DepthHistoryIntervals defines model for DepthHistoryIntervals.
type DepthHistoryIntervals []DepthHistoryItem

// DepthHistoryItem defines model for DepthHistoryItem.
type DepthHistoryItem struct {

	// Int64, the amount of Asset in the pool.
	AssetDepth string `json:"assetDepth"`

	// Float, price of asset in rune. I.e. rune amount / asset amount.
	AssetPrice string `json:"assetPrice"`

	// Int64, The end time of bucket in unix timestamp
	EndTime string `json:"endTime"`

	// Int64, the amount of Rune in the pool.
	RuneDepth string `json:"runeDepth"`

	// Int64, The beginning time of bucket in unix timestamp
	StartTime string `json:"startTime"`
}

// DepthHistoryMeta defines model for DepthHistoryMeta.
type DepthHistoryMeta struct {

	// Int64, The end time of bucket in unix timestamp
	EndTime string `json:"endTime"`

	// Int64, The beginning time of bucket in unix timestamp
	StartTime string `json:"startTime"`
}

// EarningsHistory defines model for EarningsHistory.
type EarningsHistory struct {
	Intervals EarningsHistoryIntervals `json:"intervals"`
	Meta      EarningsHistoryItem      `json:"meta"`
}

// EarningsHistoryIntervals defines model for EarningsHistoryIntervals.
type EarningsHistoryIntervals []EarningsHistoryItem

// EarningsHistoryItem defines model for EarningsHistoryItem.
type EarningsHistoryItem struct {

	// float64, Average amount of active nodes during the time interval
	AvgNodeCount string `json:"avgNodeCount"`

	// Int64, Total block rewards emitted during the time interval
	BlockRewards string `json:"blockRewards"`

	// Int64, Share of earnings sent to nodes during the time interval
	BondingEarnings string `json:"bondingEarnings"`

	// Int64, System income generated during the time interval. It is the sum of liquidity fees and block rewards.
	Earnings string `json:"earnings"`

	// Int64, The end time of interval in unix timestamp
	EndTime string `json:"endTime"`

	// Int64, Share of earnings sent to pools during the time interval
	LiquidityEarnings string `json:"liquidityEarnings"`

	// Int64, Total liquidity fees, converted to RUNE, collected during the time interval.
	LiquidityFees string `json:"liquidityFees"`

	// Earnings data for each pool for the time interval
	Pools []EarningsHistoryItemPool `json:"pools"`

	// Int64, The beginning time of interval in unix timestamp
	StartTime string `json:"startTime"`
}

// EarningsHistoryItemPool defines model for EarningsHistoryItemPool.
type EarningsHistoryItemPool struct {

	// Int64, Share of earnings sent to the pool during the time interval
	Earnings string `json:"earnings"`

	// asset for the given pool
	Pool string `json:"pool"`
}

// Health defines model for Health.
type Health struct {

	// True means healthy, connected to database
	Database bool `json:"database"`

	// True means healthy. False means Midgard is still catching up to the chain
	InSync bool `json:"inSync"`

	// Int64, the current block count
	ScannerHeight string `json:"scannerHeight"`
}

// InboundAddresses defines model for InboundAddresses.
type InboundAddresses struct {
	Current []InboundAddressesItem `json:"current"`
}

// InboundAddressesItem defines model for InboundAddressesItem.
type InboundAddressesItem struct {
	Address string `json:"address"`
	Chain   string `json:"chain"`

	// indicate whether this chain has halted
	Halted bool   `json:"halted"`
	PubKey string `json:"pub_key"`
}

// Int64Constants defines model for Int64Constants.
type Int64Constants struct {
	AsgardSize                  int64 `json:"AsgardSize"`
	BadValidatorRate            int64 `json:"BadValidatorRate"`
	BlocksPerYear               int64 `json:"BlocksPerYear"`
	ChurnInterval               int64 `json:"ChurnInterval"`
	ChurnRetryInterval          int64 `json:"ChurnRetryInterval"`
	CliTxCost                   int64 `json:"CliTxCost"`
	DesiredValidatorSet         int64 `json:"DesiredValidatorSet"`
	DoubleSignMaxAge            int64 `json:"DoubleSignMaxAge"`
	EmissionCurve               int64 `json:"EmissionCurve"`
	FailKeygenSlashPoints       int64 `json:"FailKeygenSlashPoints"`
	FailKeysignSlashPoints      int64 `json:"FailKeysignSlashPoints"`
	FundMigrationInterval       int64 `json:"FundMigrationInterval"`
	JailTimeKeygen              int64 `json:"JailTimeKeygen"`
	JailTimeKeysign             int64 `json:"JailTimeKeysign"`
	LackOfObservationPenalty    int64 `json:"LackOfObservationPenalty"`
	LiquidityLockUpBlocks       int64 `json:"LiquidityLockUpBlocks"`
	MinimumBondInRune           int64 `json:"MinimumBondInRune"`
	MinimumNodesForBFT          int64 `json:"MinimumNodesForBFT"`
	MinimumNodesForYggdrasil    int64 `json:"MinimumNodesForYggdrasil"`
	NativeChainGasFee           int64 `json:"NativeChainGasFee"`
	NewPoolCycle                int64 `json:"NewPoolCycle"`
	ObservationDelayFlexibility int64 `json:"ObservationDelayFlexibility"`
	ObserveSlashPoints          int64 `json:"ObserveSlashPoints"`
	OldValidatorRate            int64 `json:"OldValidatorRate"`
	OutboundTransactionFee      int64 `json:"OutboundTransactionFee"`
	SigningTransactionPeriod    int64 `json:"SigningTransactionPeriod"`
	YggFundLimit                int64 `json:"YggFundLimit"`
}

// Lastblock defines model for Lastblock.
type Lastblock struct {
	Current []LastblockItem `json:"current"`
}

// LastblockItem defines model for LastblockItem.
type LastblockItem struct {
	Chain          string `json:"chain"`
	LastObservedIn string `json:"last_observed_in"`
	LastSignedOut  string `json:"last_signed_out"`
	Thorchain      string `json:"thorchain"`
}

// LiquidityHistory defines model for LiquidityHistory.
type LiquidityHistory struct {
	Intervals LiquidityHistoryIntervals `json:"intervals"`
	Meta      LiquidityHistoryItem      `json:"meta"`
}

// LiquidityHistoryIntervals defines model for LiquidityHistoryIntervals.
type LiquidityHistoryIntervals []LiquidityHistoryItem

// LiquidityHistoryItem defines model for LiquidityHistoryItem.
type LiquidityHistoryItem struct {

	// Int64, total deposits (liquidity additions) during the time interval
	Deposits string `json:"deposits"`

	// Int64, The end time of bucket in unix timestamp
	EndTime string `json:"endTime"`

	// Int64, net liquidity changes (withdrawals - deposits) during the time interval
	Net string `json:"net"`

	// Int64, The beginning time of bucket in unix timestamp
	StartTime string `json:"startTime"`

	// Int64, total withdrawals during the time interval
	Withdrawals string `json:"withdrawals"`
}

// MemberDetails defines model for MemberDetails.
type MemberDetails struct {

	// Liquidity provider data for all the pools of a given member
	Pools []MemberPoolDetails `json:"pools"`
}

// MemberPoolDetails defines model for MemberPoolDetails.
type MemberPoolDetails struct {

	// Int64, total asset added to the pool by member
	AssetAdded string `json:"assetAdded"`

	// Int64, total asset withdrawn from the pool by member
	AssetWithdrawn string `json:"assetWithdrawn"`

	// Int64, Unix timestamp for the first time member deposited into the pool
	DateFirstAdded string `json:"dateFirstAdded"`

	// Int64, Unix timestamp for the last time member deposited into the pool
	DateLastAdded string `json:"dateLastAdded"`

	// Int64, pool liquidity units that belong the the member
	LiquidityUnits string `json:"liquidityUnits"`

	// Pool rest of the data refers to
	Pool string `json:"pool"`

	// Int64, total RUNE added to the pool by member
	RuneAdded string `json:"runeAdded"`

	// Int64, total RUNE withdrawn from the pool by member
	RuneWithdrawn string `json:"runeWithdrawn"`
}

// Members defines model for Members.
type Members []string

// Network defines model for Network.
type Network struct {
	ActiveBonds []string `json:"activeBonds"`

	// Int64, Number of Active Nodes
	ActiveNodeCount string       `json:"activeNodeCount"`
	BlockRewards    BlockRewards `json:"blockRewards"`
	BondMetrics     BondMetrics  `json:"bondMetrics"`

	// Float, (1 + (bondReward * blocksPerMonth/totalActiveBond)) ^ 12 -1
	BondingAPY string `json:"bondingAPY"`
	BondingROI string `json:"bondingROI"`

	// Float, (1 + (stakeReward * blocksPerMonth/totalDepth of active pools)) ^ 12 -1
	LiquidityAPY string `json:"liquidityAPY"`

	// Int64, next height of blocks.
	NextChurnHeight string `json:"nextChurnHeight"`

	// Int64, the remaining time of pool activation (in blocks).
	PoolActivationCountdown string  `json:"poolActivationCountdown"`
	PoolShareFactor         string  `json:"poolShareFactor"`
	StakingROI              *string `json:"stakingROI,omitempty"`

	// Array of Standby Bonds
	StandbyBonds []string `json:"standbyBonds"`

	// Int64, Number of Standby Nodes
	StandbyNodeCount string `json:"standbyNodeCount"`

	// Int64, Total Rune pooled in all pools.
	TotalPooledRune string `json:"totalPooledRune"`

	// Int64, Total left in Reserve
	TotalReserve string `json:"totalReserve"`
}

// NodeKey defines model for NodeKey.
type NodeKey struct {

	// ed25519 public key
	Ed25519 string `json:"ed25519"`

	// secp256k1 public key
	Secp256k1 string `json:"secp256k1"`
}

// NodeKeys defines model for NodeKeys.
type NodeKeys []NodeKey

// PoolDetail defines model for PoolDetail.
type PoolDetail struct {
	Asset string `json:"asset"`

	// Int64, the amount of Asset in the pool.
	AssetDepth string `json:"assetDepth"`

	// Float, price of asset in rune. I.e. rune amount / asset amount.
	AssetPrice string `json:"assetPrice"`

	// Float, Average Percentage Yield: annual return estimated using last weeks income, taking compound interest into account.
	PoolAPY string `json:"poolAPY"`

	// Int64, the amount of Rune in the pool.
	RuneDepth string `json:"runeDepth"`

	// The state of the pool, e.g. Available, Staged.
	Status string `json:"status"`

	// Int64, Liquidity Units in the pool.
	Units string `json:"units"`

	// Int64, the total volume of swaps in the last 24h to and from Rune denoted in Rune.
	Volume24h string `json:"volume24h"`
}

// PoolDetails defines model for PoolDetails.
type PoolDetails []PoolDetail

// Queue defines model for Queue.
type Queue struct {
	Outbound string `json:"outbound"`
	Swap     string `json:"swap"`
}

// StatsData defines model for StatsData.
type StatsData struct {

	// Daily active users (unique addresses interacting)
	DailyActiveUsers string `json:"dailyActiveUsers"`

	// Daily transactions
	DailyTx string `json:"dailyTx"`

	// Monthly active users
	MonthlyActiveUsers string `json:"monthlyActiveUsers"`

	// Monthly transactions
	MonthlyTx string `json:"monthlyTx"`

	// Total buying transactions
	TotalAssetBuys string `json:"totalAssetBuys"`

	// Total selling transactions
	TotalAssetSells string `json:"totalAssetSells"`

	// Total RUNE balances
	TotalDepth string `json:"totalDepth"`

	// Total staking transactions
	TotalStakeTx string `json:"totalStakeTx"`

	// Total staked (in RUNE Value).
	TotalStaked string `json:"totalStaked"`

	// Total transactions
	TotalTx string `json:"totalTx"`

	// Total unique swappers & members
	TotalUsers string `json:"totalUsers"`

	// Total (in RUNE Value) of all assets swapped since start.
	TotalVolume string `json:"totalVolume"`

	// Total withdrawing transactions
	TotalWithdrawTx string `json:"totalWithdrawTx"`
}

// StringConstants defines model for StringConstants.
type StringConstants struct {
	DefaultPoolStatus string `json:"DefaultPoolStatus"`
}

// SwapHistory defines model for SwapHistory.
type SwapHistory struct {
	Intervals SwapHistoryIntervals `json:"intervals"`
	Meta      SwapHistoryItem      `json:"meta"`
}

// SwapHistoryIntervals defines model for SwapHistoryIntervals.
type SwapHistoryIntervals []SwapHistoryItem

// SwapHistoryItem defines model for SwapHistoryItem.
type SwapHistoryItem struct {

	// Int64, the average slip by swap. Big swaps have the same weight as small swaps
	AverageSlip string `json:"averageSlip"`

	// Int64, The end time of bucket in unix timestamp
	EndTime string `json:"endTime"`

	// Int64, The beginning time of bucket in unix timestamp
	StartTime string `json:"startTime"`

	// Int64, count of swaps from rune to asset
	ToAssetCount string `json:"toAssetCount"`

	// Int64, volume of swaps from rune to asset denoted in rune
	ToAssetVolume string `json:"toAssetVolume"`

	// Int64, count of swaps from asset to rune
	ToRuneCount string `json:"toRuneCount"`

	// Int64, volume of swaps from asset to rune denoted in rune
	ToRuneVolume string `json:"toRuneVolume"`

	// Int64, toAssetCount + toRuneCount
	TotalCount string `json:"totalCount"`

	// Int64, the sum of all fees collected denoted in rune
	TotalFees string `json:"totalFees"`

	// Int64, toAssetVolume + toRuneVolume (denoted in rune)
	TotalVolume string `json:"totalVolume"`
}

// TxDetails defines model for TxDetails.
type TxDetails struct {

	// Int64, Unix timestamp.
	Date   string `json:"date"`
	Events Event  `json:"events"`
	Height string `json:"height"`
	In     Tx     `json:"in"`
	Out    []Tx   `json:"out"`
	Pool   string `json:"pool"`
	Status string `json:"status"`
	Type   string `json:"type"`
}

// Coin defines model for coin.
type Coin struct {
	Amount string `json:"amount"`
	Asset  string `json:"asset"`
}

// Coins defines model for coins.
type Coins []Coin

// Event defines model for event.
type Event struct {
	Fee        string `json:"fee"`
	Slip       string `json:"slip"`
	StakeUnits string `json:"stakeUnits"`
}

// Option defines model for option.
type Option struct {
	Asymmetry           string `json:"asymmetry"`
	PriceTarget         string `json:"priceTarget"`
	WithdrawBasisPoints string `json:"withdrawBasisPoints"`
}

// Tx defines model for tx.
type Tx struct {
	Address string `json:"address"`
	Coins   Coins  `json:"coins"`
	Memo    string `json:"memo"`
	Options Option `json:"options"`
	TxID    string `json:"txID"`
}

// ConstantsResponse defines model for ConstantsResponse.
type ConstantsResponse Constants

// DepthHistoryResponse defines model for DepthHistoryResponse.
type DepthHistoryResponse DepthHistory

// EarningsHistoryResponse defines model for EarningsHistoryResponse.
type EarningsHistoryResponse EarningsHistory

// HealthResponse defines model for HealthResponse.
type HealthResponse Health

// InboundAddressesResponse defines model for InboundAddressesResponse.
type InboundAddressesResponse InboundAddresses

// LastblockResponse defines model for LastblockResponse.
type LastblockResponse Lastblock

// LiquidityHistoryResponse defines model for LiquidityHistoryResponse.
type LiquidityHistoryResponse LiquidityHistory

// MemberDetailsResponse defines model for MemberDetailsResponse.
type MemberDetailsResponse MemberDetails

// MembersResponse defines model for MembersResponse.
type MembersResponse Members

// NetworkResponse defines model for NetworkResponse.
type NetworkResponse Network

// NodeKeyResponse defines model for NodeKeyResponse.
type NodeKeyResponse NodeKeys

// PoolResponse defines model for PoolResponse.
type PoolResponse PoolDetail

// PoolsResponse defines model for PoolsResponse.
type PoolsResponse PoolDetails

// QueueResponse defines model for QueueResponse.
type QueueResponse Queue

// StatsResponse defines model for StatsResponse.
type StatsResponse StatsData

// SwapHistoryResponse defines model for SwapHistoryResponse.
type SwapHistoryResponse SwapHistory

// TxResponse defines model for TxResponse.
type TxResponse struct {

	// Int64, count of txs matching the filters.
	Count string      `json:"count"`
	Txs   []TxDetails `json:"txs"`
}

// GetDepthHistoryParams defines parameters for GetDepthHistory.
type GetDepthHistoryParams struct {

	// Interval of calculations
	Interval string `json:"interval"`

	// Number of intervals to return. Should be between [1..100].
	Count *int `json:"count,omitempty"`

	// End time of the query as unix timestamp. If only count is given, defaults to now.
	To *int64 `json:"to,omitempty"`

	// Start time of the query as unix timestamp
	From *int64 `json:"from,omitempty"`
}

// GetEarningsHistoryParams defines parameters for GetEarningsHistory.
type GetEarningsHistoryParams struct {

	// Interval of calculations
	Interval string `json:"interval"`

	// Number of intervals to return. Should be between [1..100].
	Count *int `json:"count,omitempty"`

	// End time of the query as unix timestamp. If only count is given, defaults to now.
	To *int64 `json:"to,omitempty"`

	// Start time of the query as unix timestamp
	From *int64 `json:"from,omitempty"`
}

// GetLiquidityHistoryParams defines parameters for GetLiquidityHistory.
type GetLiquidityHistoryParams struct {

	// Return stats for given pool. Returns sum of all pools if missing.
	Pool *string `json:"pool,omitempty"`

	// Interval of calculations
	Interval string `json:"interval"`

	// Number of intervals to return. Should be between [1..100].
	Count *int `json:"count,omitempty"`

	// End time of the query as unix timestamp. If only count is given, defaults to now.
	To *int64 `json:"to,omitempty"`

	// Start time of the query as unix timestamp
	From *int64 `json:"from,omitempty"`
}

// GetSwapHistoryParams defines parameters for GetSwapHistory.
type GetSwapHistoryParams struct {

	// Return history given pool. Returns sum of all pools if missing.
	Pool *string `json:"pool,omitempty"`

	// Interval of calculations
	Interval string `json:"interval"`

	// Number of intervals to return. Should be between [1..100].
	Count *int `json:"count,omitempty"`

	// End time of the query as unix timestamp. If only count is given, defaults to now.
	To *int64 `json:"to,omitempty"`

	// Start time of the query as unix timestamp
	From *int64 `json:"from,omitempty"`
}

// GetPoolsParams defines parameters for GetPools.
type GetPoolsParams struct {

	// Filter for only pools with this status.
	Status *string `json:"status,omitempty"`
}

// GetTxDetailsParams defines parameters for GetTxDetails.
type GetTxDetailsParams struct {

	// Address of sender or recipient of any in/out tx in event
	Address *string `json:"address,omitempty"`

	// ID of any in/out tx in event
	Txid *string `json:"txid,omitempty"`

	// Any asset used in event (CHAIN.SYMBOL)
	Asset *string `json:"asset,omitempty"`

	// One or more comma separated unique types of event
	Type *string `json:"type,omitempty"`

	// pagination limit
	Limit int64 `json:"limit"`

	// pagination offset
	Offset int64 `json:"offset"`
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file.
func GetSwagger() (*openapi3.Swagger, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromData(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error loading Swagger: %s", err)
	}
	return swagger, nil
}
