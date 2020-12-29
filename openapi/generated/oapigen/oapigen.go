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

	// Int64, Average bond of active nodes
	AverageActiveBond string `json:"averageActiveBond"`

	// Int64, Average bond of standby nodes
	AverageStandbyBond string `json:"averageStandbyBond"`

	// Int64, Maxinum bond of active nodes
	MaximumActiveBond string `json:"maximumActiveBond"`

	// Int64, Maximum bond of standby nodes
	MaximumStandbyBond string `json:"maximumStandbyBond"`

	// Int64, Median bond of active nodes
	MedianActiveBond string `json:"medianActiveBond"`

	// Int64, Median bond of standby nodes
	MedianStandbyBond string `json:"medianStandbyBond"`

	// Int64, Minumum bond of active nodes
	MinimumActiveBond string `json:"minimumActiveBond"`

	// Int64, Minumum bond of standby nodes
	MinimumStandbyBond string `json:"minimumStandbyBond"`

	// Int64, Total bond of active nodes
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

	// Int64, the amount of Asset in the pool
	AssetDepth string `json:"assetDepth"`

	// Float, price of asset in rune. I.e. rune amount / asset amount
	AssetPrice string `json:"assetPrice"`

	// Int64, The end time of bucket in unix timestamp
	EndTime string `json:"endTime"`

	// Int64, the amount of Rune in the pool
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

	// Int64, System income generated during the time interval. It is the sum of liquidity fees and block rewards
	Earnings string `json:"earnings"`

	// Int64, The end time of interval in unix timestamp
	EndTime string `json:"endTime"`

	// Int64, Share of earnings sent to pools during the time interval
	LiquidityEarnings string `json:"liquidityEarnings"`

	// Int64, Total liquidity fees, converted to RUNE, collected during the time interval
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

	// Float, (1 + (stakeReward * blocksPerMonth/totalDepth of active pools)) ^ 12 -1
	LiquidityAPY string `json:"liquidityAPY"`

	// Int64, next height of blocks
	NextChurnHeight string `json:"nextChurnHeight"`

	// Int64, the remaining time of pool activation (in blocks)
	PoolActivationCountdown string `json:"poolActivationCountdown"`
	PoolShareFactor         string `json:"poolShareFactor"`

	// Array of Standby Bonds
	StandbyBonds []string `json:"standbyBonds"`

	// Int64, Number of Standby Nodes
	StandbyNodeCount string `json:"standbyNodeCount"`

	// Int64, Total Rune pooled in all pools
	TotalPooledRune string `json:"totalPooledRune"`

	// Int64, Total left in Reserve
	TotalReserve string `json:"totalReserve"`
}

// NodeKey defines model for NodeKey.
type NodeKey struct {

	// ed25519 public key
	Ed25519 string `json:"ed25519"`

	// node thorchain address
	NodeAddress string `json:"nodeAddress"`

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

// PoolLegacyDetail defines model for PoolLegacyDetail.
type PoolLegacyDetail struct {
	Asset string `json:"asset"`

	// Int64, the amount of Asset in the pool
	AssetDepth string `json:"assetDepth"`

	// Float, price of asset in rune. I.e. rune amount / asset amount
	AssetPrice string `json:"assetPrice"`

	// Int64, same as history/swaps:toAssetCount
	BuyAssetCount string `json:"buyAssetCount"`

	// Float, same as history/swaps:toRuneVolume/toRuneCount
	BuyTxAverage string `json:"buyTxAverage"`

	// Int64, same as history/swaps:toAssetVolume
	BuyVolume string `json:"buyVolume"`

	// Float, Average Percentage Yield: annual return estimated using last weeks income, taking compound interest into account.
	PoolAPY string `json:"poolAPY"`

	// Int64, same as assetDepth + runeDepth
	PoolDepth string `json:"poolDepth"`

	// Float, same as history/swaps:totalFees/totalCount
	PoolFeeAverage string `json:"poolFeeAverage"`

	// Int64, same as history/swaps:totalFees
	PoolFeesTotal string `json:"poolFeesTotal"`

	// Float, same as history/swaps:averageSlip
	PoolSlipAverage string `json:"poolSlipAverage"`

	// Float, same as history/swaps:totalVolume/totalCount
	PoolTxAverage string `json:"poolTxAverage"`

	// Int64, same as buyVolume + sellVolume
	PoolVolume string `json:"poolVolume"`

	// Int64, the amount of Rune in the pool
	RuneDepth string `json:"runeDepth"`

	// Int64, same as history/swaps:toRuneCount
	SellAssetCount string `json:"sellAssetCount"`

	// Float, same as history/swaps:toAssetVolume/toAssetCount
	SellTxAverage string `json:"sellTxAverage"`

	// Int64, same as history/swaps:toRuneVolume
	SellVolume string `json:"sellVolume"`

	// The state of the pool, e.g. Available, Staged
	Status string `json:"status"`

	// Int64, same as history/swaps:totalCount
	SwappingTxCount string `json:"swappingTxCount"`

	// Int64, Liquidity Units in the pool
	Units string `json:"units"`

	// Int64, the total volume of swaps in the last 24h to and from Rune denoted in Rune
	Volume24h string `json:"volume24h"`
}

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

	// Total staked (in RUNE Value)
	TotalStaked string `json:"totalStaked"`

	// Total transactions
	TotalTx string `json:"totalTx"`

	// Total unique swappers & members
	TotalUsers string `json:"totalUsers"`

	// Total (in RUNE Value) of all assets swapped since start
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

	// Float, the average slip by swap. Big swaps have the same weight as small swaps
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

	// Int64, Unix timestamp
	Date     string   `json:"date"`
	Height   string   `json:"height"`
	In       []Tx     `json:"in"`
	Metadata Metadata `json:"metadata"`
	Out      []Tx     `json:"out"`
	Pools    []string `json:"pools"`
	Status   string   `json:"status"`
	Type     string   `json:"type"`
}

// Coin defines model for coin.
type Coin struct {
	Amount string `json:"amount"`
	Asset  string `json:"asset"`
}

// Coins defines model for coins.
type Coins []Coin

// Metadata defines model for metadata.
type Metadata struct {
	Refund *RefundMetadata `json:"refund,omitempty"`
	Swap   *SwapMetadata   `json:"swap,omitempty"`
}

// RefundMetadata defines model for refundMetadata.
type RefundMetadata struct {
	NetworkFees []Coin `json:"networkFees"`
	Reason      string `json:"reason"`
}

// SwapMetadata defines model for swapMetadata.
type SwapMetadata struct {
	LiquidityFee string `json:"liquidityFee"`
	NetworkFees  []Coin `json:"networkFees"`
	TradeSlip    string `json:"tradeSlip"`
	TradeTarget  string `json:"tradeTarget"`
}

// Tx defines model for tx.
type Tx struct {
	Address string `json:"address"`
	Coins   Coins  `json:"coins"`
	Memo    string `json:"memo"`
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

// PoolLegacyResponse defines model for PoolLegacyResponse.
type PoolLegacyResponse PoolLegacyDetail

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
	Interval *string `json:"interval,omitempty"`

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
	Interval *string `json:"interval,omitempty"`

	// Number of intervals to return. Should be between [1..100].
	Count *int `json:"count,omitempty"`

	// End time of the query as unix timestamp. If only count is given, defaults to now.
	To *int64 `json:"to,omitempty"`

	// Start time of the query as unix timestamp
	From *int64 `json:"from,omitempty"`
}

// GetLiquidityHistoryParams defines parameters for GetLiquidityHistory.
type GetLiquidityHistoryParams struct {

	// Return stats for given pool. Returns sum of all pools if missing
	Pool *string `json:"pool,omitempty"`

	// Interval of calculations
	Interval *string `json:"interval,omitempty"`

	// Number of intervals to return. Should be between [1..100]
	Count *int `json:"count,omitempty"`

	// End time of the query as unix timestamp. If only count is given, defaults to now
	To *int64 `json:"to,omitempty"`

	// Start time of the query as unix timestamp
	From *int64 `json:"from,omitempty"`
}

// GetSwapHistoryParams defines parameters for GetSwapHistory.
type GetSwapHistoryParams struct {

	// Return history given pool. Returns sum of all pools if missing.
	Pool *string `json:"pool,omitempty"`

	// Interval of calculations
	Interval *string `json:"interval,omitempty"`

	// Number of intervals to return. Should be between [1..100].
	Count *int `json:"count,omitempty"`

	// End time of the query as unix timestamp. If only count is given, defaults to now.
	To *int64 `json:"to,omitempty"`

	// Start time of the query as unix timestamp
	From *int64 `json:"from,omitempty"`
}

// GetPoolsParams defines parameters for GetPools.
type GetPoolsParams struct {

	// Filter for only pools with this status
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
	"N5oAAA==",
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
