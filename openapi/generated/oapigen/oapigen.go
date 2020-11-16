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

// MemberDetails defines model for MemberDetails.
type MemberDetails struct {
	StakeArray  []string `json:"stakeArray"`
	TotalStaked string   `json:"totalStaked"`
}

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

// PoolDetail defines model for PoolDetail.
type PoolDetail struct {
	Asset string `json:"asset"`

	// Int64, the amount of Asset in the pool.
	AssetDepth string `json:"assetDepth"`

	// Float, Average Percentage Yield: annual return estimated using last weeks income, taking compound interest into account.
	PoolAPY string `json:"poolAPY"`

	// Float, price of asset in rune. I.e. rune amount / asset amount.
	Price string `json:"price"`

	// Int64, the amount of Rune in the pool.
	RuneDepth string `json:"runeDepth"`

	// The state of the pool, e.g. Enabled, Bootstrap.
	Status string `json:"status"`

	// Int64, Liquidity Units in the pool.
	Units string `json:"units"`

	// Int64, the total volume of swaps in the last 24h to and from Rune denoted in Rune.
	Volume24h string `json:"volume24h"`
}

// PoolSummary defines model for PoolSummary.
type PoolSummary struct {
	Asset string `json:"asset"`

	// Int64, the amount of Asset in the pool.
	AssetDepth string `json:"assetDepth"`

	// Int64, unix timestamp (second).
	DateCreated string `json:"dateCreated"`

	// Float, price of asset in rune. I.e. rune amount / asset amount.
	Price string `json:"price"`

	// Int64, the amount of Rune in the pool.
	RuneDepth string `json:"runeDepth"`
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

// Volume defines model for Volume.
type Volume struct {

	// Sum of "rune_amount" of buy swap events.
	BuyVolume string `json:"buyVolume"`

	// Sum of "rune_amount" of sell swap events.
	SellVolume string `json:"sellVolume"`

	// Int64, The beginning time of bucket in unix timestamp.
	Time string `json:"time"`

	// buyVolume + sellVolume
	TotalVolume string `json:"totalVolume"`
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

// HealthResponse defines model for HealthResponse.
type HealthResponse struct {

	// True means healthy, connected to database.
	Database bool `json:"database"`

	// True means healthy. False means Midgard is still catching up to the chain.
	InSync bool `json:"inSync"`

	// Int64, the current block count.
	ScannerHeight string `json:"scannerHeight"`
}

// MemberDetailsResponse defines model for MemberDetailsResponse.
type MemberDetailsResponse MemberDetails

// MembersResponse defines model for MembersResponse.
type MembersResponse []string

// NetworkResponse defines model for NetworkResponse.
type NetworkResponse Network

// NodeKeyResponse defines model for NodeKeyResponse.
type NodeKeyResponse []NodeKey

// PoolDetailResponse defines model for PoolDetailResponse.
type PoolDetailResponse PoolDetail

// PoolsResponse defines model for PoolsResponse.
type PoolsResponse []PoolSummary

// StatsResponse defines model for StatsResponse.
type StatsResponse StatsData

// TxResponse defines model for TxResponse.
type TxResponse struct {

	// Int64, count of txs matching the filters.
	Count string      `json:"count"`
	Txs   []TxDetails `json:"txs"`
}

// VolumeResponse defines model for VolumeResponse.
type VolumeResponse []Volume

// GetTotalVolumeParams defines parameters for GetTotalVolume.
type GetTotalVolumeParams struct {

	// Return volume for this single pool. Returns volume for all pools if missing.
	Pool *int64 `json:"pool,omitempty"`

	// Interval of calculations
	Interval string `json:"interval"`

	// Start time of the query as unix timestamp
	From int64 `json:"from"`

	// End time of the query as unix timestamp
	To int64 `json:"to"`
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