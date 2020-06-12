// Package http provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package http

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

// AssetDetail defines model for AssetDetail.
type AssetDetail struct {
	Asset       *Asset  `json:"asset,omitempty"`
	DateCreated *int64  `json:"dateCreated,omitempty"`
	PriceRune   *string `json:"priceRune,omitempty"`
}

// BlockRewards defines model for BlockRewards.
type BlockRewards struct {
	BlockReward *string `json:"blockReward,omitempty"`
	BondReward  *string `json:"bondReward,omitempty"`
	StakeReward *string `json:"stakeReward,omitempty"`
}

// BondMetrics defines model for BondMetrics.
type BondMetrics struct {

	// Average bond of active nodes
	AverageActiveBond *string `json:"averageActiveBond,omitempty"`

	// Average bond of standby nodes
	AverageStandbyBond *string `json:"averageStandbyBond,omitempty"`

	// Maxinum bond of active nodes
	MaximumActiveBond *string `json:"maximumActiveBond,omitempty"`

	// Maximum bond of standby nodes
	MaximumStandbyBond *string `json:"maximumStandbyBond,omitempty"`

	// Median bond of active nodes
	MedianActiveBond *string `json:"medianActiveBond,omitempty"`

	// Median bond of standby nodes
	MedianStandbyBond *string `json:"medianStandbyBond,omitempty"`

	// Minumum bond of active nodes
	MinimumActiveBond *string `json:"minimumActiveBond,omitempty"`

	// Minumum bond of standby nodes
	MinimumStandbyBond *string `json:"minimumStandbyBond,omitempty"`

	// Total bond of active nodes
	TotalActiveBond *string `json:"totalActiveBond,omitempty"`

	// Total bond of standby nodes
	TotalStandbyBond *string `json:"totalStandbyBond,omitempty"`
}

// Error defines model for Error.
type Error struct {
	Error string `json:"error"`
}

// NetworkInfo defines model for NetworkInfo.
type NetworkInfo struct {

	// Array of Active Bonds
	ActiveBonds *[]string `json:"activeBonds,omitempty"`

	// Number of Active Nodes
	ActiveNodeCount *int          `json:"activeNodeCount,omitempty"`
	BlockRewards    *BlockRewards `json:"blockRewards,omitempty"`
	BondMetrics     *BondMetrics  `json:"bondMetrics,omitempty"`
	BondingROI      *string       `json:"bondingROI,omitempty"`
	NextChurnHeight *string       `json:"nextChurnHeight,omitempty"`

	// The remaining time of pool activation (in blocks)
	PoolActivationCountdown *int64  `json:"poolActivationCountdown,omitempty"`
	PoolShareFactor         *string `json:"poolShareFactor,omitempty"`
	StakingROI              *string `json:"stakingROI,omitempty"`

	// Array of Standby Bonds
	StandbyBonds *[]string `json:"standbyBonds,omitempty"`

	// Number of Standby Nodes
	StandbyNodeCount *int `json:"standbyNodeCount,omitempty"`

	// Total left in Reserve
	TotalReserve *string `json:"totalReserve,omitempty"`

	// Total Rune Staked in Pools
	TotalStaked *string `json:"totalStaked,omitempty"`
}

// NodeKey defines model for NodeKey.
type NodeKey struct {

	// ed25519 public key
	Ed25519 *string `json:"ed25519,omitempty"`

	// secp256k1 public key
	Secp256k1 *string `json:"secp256k1,omitempty"`
}

// PoolDetail defines model for PoolDetail.
type PoolDetail struct {
	Asset *Asset `json:"asset,omitempty"`

	// Total current Asset balance
	AssetDepth *string `json:"assetDepth,omitempty"`

	// Asset return on investment
	AssetROI *string `json:"assetROI,omitempty"`

	// Total Asset staked
	AssetStakedTotal *string `json:"assetStakedTotal,omitempty"`

	// Number of RUNE->ASSET transactions
	BuyAssetCount *string `json:"buyAssetCount,omitempty"`

	// Average sell Asset fee size for RUNE->ASSET (in ASSET)
	BuyFeeAverage *string `json:"buyFeeAverage,omitempty"`

	// Total fees (in Asset)
	BuyFeesTotal *string `json:"buyFeesTotal,omitempty"`

	// Average trade slip for RUNE->ASSET in %
	BuySlipAverage *string `json:"buySlipAverage,omitempty"`

	// Average Asset buy transaction size for (RUNE->ASSET) (in ASSET)
	BuyTxAverage *string `json:"buyTxAverage,omitempty"`

	// Total Asset buy volume (RUNE->ASSET) (in Asset)
	BuyVolume *string `json:"buyVolume,omitempty"`

	// Total depth of both sides (in RUNE)
	PoolDepth *string `json:"poolDepth,omitempty"`

	// Average pool fee
	PoolFeeAverage *string `json:"poolFeeAverage,omitempty"`

	// Total fees
	PoolFeesTotal *string `json:"poolFeesTotal,omitempty"`

	// Pool ROI (average of RUNE and Asset ROI)
	PoolROI *string `json:"poolROI,omitempty"`

	// Pool ROI over 12 months
	PoolROI12 *string `json:"poolROI12,omitempty"`

	// Average pool slip
	PoolSlipAverage *string `json:"poolSlipAverage,omitempty"`

	// Rune value staked Total
	PoolStakedTotal *string `json:"poolStakedTotal,omitempty"`

	// Average pool transaction
	PoolTxAverage *string `json:"poolTxAverage,omitempty"`

	// Total pool units outstanding
	PoolUnits *string `json:"poolUnits,omitempty"`

	// Two-way volume of all-time (in RUNE)
	PoolVolume *string `json:"poolVolume,omitempty"`

	// Two-way volume in 24hrs (in RUNE)
	PoolVolume24hr *string `json:"poolVolume24hr,omitempty"`

	// Price of Asset (in RUNE).
	Price *string `json:"price,omitempty"`

	// Total current Rune balance
	RuneDepth *string `json:"runeDepth,omitempty"`

	// RUNE return on investment
	RuneROI *string `json:"runeROI,omitempty"`

	// Total RUNE staked
	RuneStakedTotal *string `json:"runeStakedTotal,omitempty"`

	// Number of ASSET->RUNE transactions
	SellAssetCount *string `json:"sellAssetCount,omitempty"`

	// Average buy Asset fee size for ASSET->RUNE (in RUNE)
	SellFeeAverage *string `json:"sellFeeAverage,omitempty"`

	// Total fees (in RUNE)
	SellFeesTotal *string `json:"sellFeesTotal,omitempty"`

	// Average trade slip for ASSET->RUNE in %
	SellSlipAverage *string `json:"sellSlipAverage,omitempty"`

	// Average Asset sell transaction size (ASSET>RUNE) (in RUNE)
	SellTxAverage *string `json:"sellTxAverage,omitempty"`

	// Total Asset sell volume (ASSET>RUNE) (in RUNE).
	SellVolume *string `json:"sellVolume,omitempty"`

	// Number of stake transactions
	StakeTxCount *string `json:"stakeTxCount,omitempty"`

	// Number of unique stakers
	StakersCount *string `json:"stakersCount,omitempty"`

	// Number of stake & withdraw transactions
	StakingTxCount *string `json:"stakingTxCount,omitempty"`
	Status         *string `json:"status,omitempty"`

	// Number of unique swappers interacting with pool
	SwappersCount *string `json:"swappersCount,omitempty"`

	// Number of swapping transactions in the pool (buys and sells)
	SwappingTxCount *string `json:"swappingTxCount,omitempty"`

	// Number of withdraw transactions
	WithdrawTxCount *string `json:"withdrawTxCount,omitempty"`
}

// Stakers defines model for Stakers.
type Stakers string

// StakersAddressData defines model for StakersAddressData.
type StakersAddressData struct {
	PoolsArray *[]Asset `json:"poolsArray,omitempty"`

	// Total value of earnings (in RUNE) across all pools.
	TotalEarned *string `json:"totalEarned,omitempty"`

	// Average of all pool ROIs.
	TotalROI *string `json:"totalROI,omitempty"`

	// Total staked (in RUNE) across all pools.
	TotalStaked *string `json:"totalStaked,omitempty"`
}

// StakersAssetData defines model for StakersAssetData.
type StakersAssetData struct {
	Asset *Asset `json:"asset,omitempty"`

	// Value of Assets earned from the pool.
	AssetEarned *string `json:"assetEarned,omitempty"`

	// ROI of the Asset side
	AssetROI *string `json:"assetROI,omitempty"`

	// Amount of Assets staked.
	AssetStaked      *string `json:"assetStaked,omitempty"`
	DateFirstStaked  *int64  `json:"dateFirstStaked,omitempty"`
	HeightLastStaked *int64  `json:"heightLastStaked,omitempty"`

	// Total value of earnings (in RUNE).
	PoolEarned *string `json:"poolEarned,omitempty"`

	// Average ROI (in RUNE) of both sides
	PoolROI *string `json:"poolROI,omitempty"`

	// RUNE value staked.
	PoolStaked *string `json:"poolStaked,omitempty"`

	// Value of RUNE earned from the pool.
	RuneEarned *string `json:"runeEarned,omitempty"`

	// ROI of the Rune side.
	RuneROI *string `json:"runeROI,omitempty"`

	// Amount of RUNE staked.
	RuneStaked *string `json:"runeStaked,omitempty"`

	// Represents ownership of a pool.
	StakeUnits *string `json:"stakeUnits,omitempty"`
}

// StatsData defines model for StatsData.
type StatsData struct {

	// Daily active users (unique addresses interacting)
	DailyActiveUsers *string `json:"dailyActiveUsers,omitempty"`

	// Daily transactions
	DailyTx *string `json:"dailyTx,omitempty"`

	// Monthly active users
	MonthlyActiveUsers *string `json:"monthlyActiveUsers,omitempty"`

	// Monthly transactions
	MonthlyTx *string `json:"monthlyTx,omitempty"`

	// Number of active pools
	PoolCount *string `json:"poolCount,omitempty"`

	// Total buying transactions
	TotalAssetBuys *string `json:"totalAssetBuys,omitempty"`

	// Total selling transactions
	TotalAssetSells *string `json:"totalAssetSells,omitempty"`

	// Total RUNE balances
	TotalDepth *string `json:"totalDepth,omitempty"`

	// Total earned (in RUNE Value).
	TotalEarned *string `json:"totalEarned,omitempty"`

	// Total staking transactions
	TotalStakeTx *string `json:"totalStakeTx,omitempty"`

	// Total staked (in RUNE Value).
	TotalStaked *string `json:"totalStaked,omitempty"`

	// Total transactions
	TotalTx *string `json:"totalTx,omitempty"`

	// Total unique swappers & stakers
	TotalUsers *string `json:"totalUsers,omitempty"`

	// Total (in RUNE Value) of all assets swapped since start.
	TotalVolume *string `json:"totalVolume,omitempty"`

	// Total (in RUNE Value) of all assets swapped in 24hrs
	TotalVolume24hr *string `json:"totalVolume24hr,omitempty"`

	// Total withdrawing transactions
	TotalWithdrawTx *string `json:"totalWithdrawTx,omitempty"`
}

// ThorchainEndpoint defines model for ThorchainEndpoint.
type ThorchainEndpoint struct {
	Address *string `json:"address,omitempty"`
	Chain   *string `json:"chain,omitempty"`
	PubKey  *string `json:"pub_key,omitempty"`
}

// ThorchainEndpoints defines model for ThorchainEndpoints.
type ThorchainEndpoints struct {
	Current *[]ThorchainEndpoint `json:"current,omitempty"`
}

// TxDetails defines model for TxDetails.
type TxDetails struct {
	Date    *int64  `json:"date,omitempty"`
	Events  *Event  `json:"events,omitempty"`
	Gas     *Gas    `json:"gas,omitempty"`
	Height  *string `json:"height,omitempty"`
	In      *Tx     `json:"in,omitempty"`
	Options *Option `json:"options,omitempty"`
	Out     *[]Tx   `json:"out,omitempty"`
	Pool    *Asset  `json:"pool,omitempty"`
	Status  *string `json:"status,omitempty"`
	Type    *string `json:"type,omitempty"`
}

// Asset defines model for asset.
type Asset string

// Coin defines model for coin.
type Coin struct {
	Amount *string `json:"amount,omitempty"`
	Asset  *Asset  `json:"asset,omitempty"`
}

// Coins defines model for coins.
type Coins []Coin

// Event defines model for event.
type Event struct {
	Fee        *string `json:"fee,omitempty"`
	Slip       *string `json:"slip,omitempty"`
	StakeUnits *string `json:"stakeUnits,omitempty"`
}

// Gas defines model for gas.
type Gas struct {
	Amount *string `json:"amount,omitempty"`
	Asset  *Asset  `json:"asset,omitempty"`
}

// Option defines model for option.
type Option struct {
	Asymmetry           *string `json:"asymmetry,omitempty"`
	PriceTarget         *string `json:"priceTarget,omitempty"`
	WithdrawBasisPoints *string `json:"withdrawBasisPoints,omitempty"`
}

// Tx defines model for tx.
type Tx struct {
	Address *string `json:"address,omitempty"`
	Coins   *Coins  `json:"coins,omitempty"`
	Memo    *string `json:"memo,omitempty"`
	TxID    *string `json:"txID,omitempty"`
}

// AssetsDetailedResponse defines model for AssetsDetailedResponse.
type AssetsDetailedResponse []AssetDetail

// GeneralErrorResponse defines model for GeneralErrorResponse.
type GeneralErrorResponse Error

// HealthResponse defines model for HealthResponse.
type HealthResponse struct {
	CatchingUp    *bool  `json:"catching_up,omitempty"`
	Database      *bool  `json:"database,omitempty"`
	ScannerHeight *int64 `json:"scannerHeight,omitempty"`
}

// NetworkResponse defines model for NetworkResponse.
type NetworkResponse NetworkInfo

// NodeKeyResponse defines model for NodeKeyResponse.
type NodeKeyResponse []NodeKey

// PoolsDetailedResponse defines model for PoolsDetailedResponse.
type PoolsDetailedResponse []PoolDetail

// PoolsResponse defines model for PoolsResponse.
type PoolsResponse []Asset

// StakersAddressDataResponse defines model for StakersAddressDataResponse.
type StakersAddressDataResponse StakersAddressData

// StakersAssetDataResponse defines model for StakersAssetDataResponse.
type StakersAssetDataResponse []StakersAssetData

// StakersResponse defines model for StakersResponse.
type StakersResponse []Stakers

// StatsResponse defines model for StatsResponse.
type StatsResponse StatsData

// ThorchainEndpointsResponse defines model for ThorchainEndpointsResponse.
type ThorchainEndpointsResponse ThorchainEndpoints

// TxsResponse defines model for TxsResponse.
type TxsResponse struct {
	Count *int64       `json:"count,omitempty"`
	Txs   *[]TxDetails `json:"txs,omitempty"`
}

// GetAssetInfoParams defines parameters for GetAssetInfo.
type GetAssetInfoParams struct {

	// One or more comma separated unique asset (CHAIN.SYMBOL)
	Asset string `json:"asset"`
}

// GetPoolsDataParams defines parameters for GetPoolsData.
type GetPoolsDataParams struct {

	// One or more comma separated unique asset (CHAIN.SYMBOL)
	Asset string `json:"asset"`
}

// GetStakersAddressAndAssetDataParams defines parameters for GetStakersAddressAndAssetData.
type GetStakersAddressAndAssetDataParams struct {

	// One or more comma separated unique asset (CHAIN.SYMBOL)
	Asset string `json:"asset"`
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

	// pagination offset
	Offset int64 `json:"offset"`

	// pagination limit
	Limit int64 `json:"limit"`
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get Asset Information
	// (GET /v1/assets)
	GetAssetInfo(ctx echo.Context, params GetAssetInfoParams) error
	// Get Documents
	// (GET /v1/doc)
	GetDocs(ctx echo.Context) error
	// Get Health
	// (GET /v1/health)
	GetHealth(ctx echo.Context) error
	// Get Network Data
	// (GET /v1/network)
	GetNetworkData(ctx echo.Context) error
	// Get Node public keys
	// (GET /v1/nodes)
	GetNodes(ctx echo.Context) error
	// Get Asset Pools
	// (GET /v1/pools)
	GetPools(ctx echo.Context) error
	// Get Pools Data
	// (GET /v1/pools/detail)
	GetPoolsData(ctx echo.Context, params GetPoolsDataParams) error
	// Get Stakers
	// (GET /v1/stakers)
	GetStakersData(ctx echo.Context) error
	// Get Staker Data
	// (GET /v1/stakers/{address})
	GetStakersAddressData(ctx echo.Context, address string) error
	// Get Staker Pool Data
	// (GET /v1/stakers/{address}/pools)
	GetStakersAddressAndAssetData(ctx echo.Context, address string, params GetStakersAddressAndAssetDataParams) error
	// Get Global Stats
	// (GET /v1/stats)
	GetStats(ctx echo.Context) error
	// Get Swagger
	// (GET /v1/swagger.json)
	GetSwagger(ctx echo.Context) error
	// Get the Proxied Pool Addresses
	// (GET /v1/thorchain/pool_addresses)
	GetThorchainProxiedEndpoints(ctx echo.Context) error
	// Get details of a tx by address, asset or tx-id
	// (GET /v1/txs)
	GetTxDetails(ctx echo.Context, params GetTxDetailsParams) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetAssetInfo converts echo context to params.
func (w *ServerInterfaceWrapper) GetAssetInfo(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetAssetInfoParams
	// ------------- Required query parameter "asset" -------------

	err = runtime.BindQueryParameter("form", true, true, "asset", ctx.QueryParams(), &params.Asset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter asset: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetAssetInfo(ctx, params)
	return err
}

// GetDocs converts echo context to params.
func (w *ServerInterfaceWrapper) GetDocs(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetDocs(ctx)
	return err
}

// GetHealth converts echo context to params.
func (w *ServerInterfaceWrapper) GetHealth(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetHealth(ctx)
	return err
}

// GetNetworkData converts echo context to params.
func (w *ServerInterfaceWrapper) GetNetworkData(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetNetworkData(ctx)
	return err
}

// GetNodes converts echo context to params.
func (w *ServerInterfaceWrapper) GetNodes(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetNodes(ctx)
	return err
}

// GetPools converts echo context to params.
func (w *ServerInterfaceWrapper) GetPools(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetPools(ctx)
	return err
}

// GetPoolsData converts echo context to params.
func (w *ServerInterfaceWrapper) GetPoolsData(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetPoolsDataParams
	// ------------- Required query parameter "asset" -------------

	err = runtime.BindQueryParameter("form", true, true, "asset", ctx.QueryParams(), &params.Asset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter asset: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetPoolsData(ctx, params)
	return err
}

// GetStakersData converts echo context to params.
func (w *ServerInterfaceWrapper) GetStakersData(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetStakersData(ctx)
	return err
}

// GetStakersAddressData converts echo context to params.
func (w *ServerInterfaceWrapper) GetStakersAddressData(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "address" -------------
	var address string

	err = runtime.BindStyledParameter("simple", false, "address", ctx.Param("address"), &address)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter address: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetStakersAddressData(ctx, address)
	return err
}

// GetStakersAddressAndAssetData converts echo context to params.
func (w *ServerInterfaceWrapper) GetStakersAddressAndAssetData(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "address" -------------
	var address string

	err = runtime.BindStyledParameter("simple", false, "address", ctx.Param("address"), &address)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter address: %s", err))
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params GetStakersAddressAndAssetDataParams
	// ------------- Required query parameter "asset" -------------

	err = runtime.BindQueryParameter("form", true, true, "asset", ctx.QueryParams(), &params.Asset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter asset: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetStakersAddressAndAssetData(ctx, address, params)
	return err
}

// GetStats converts echo context to params.
func (w *ServerInterfaceWrapper) GetStats(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetStats(ctx)
	return err
}

// GetSwagger converts echo context to params.
func (w *ServerInterfaceWrapper) GetSwagger(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetSwagger(ctx)
	return err
}

// GetThorchainProxiedEndpoints converts echo context to params.
func (w *ServerInterfaceWrapper) GetThorchainProxiedEndpoints(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetThorchainProxiedEndpoints(ctx)
	return err
}

// GetTxDetails converts echo context to params.
func (w *ServerInterfaceWrapper) GetTxDetails(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetTxDetailsParams
	// ------------- Optional query parameter "address" -------------

	err = runtime.BindQueryParameter("form", true, false, "address", ctx.QueryParams(), &params.Address)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter address: %s", err))
	}

	// ------------- Optional query parameter "txid" -------------

	err = runtime.BindQueryParameter("form", true, false, "txid", ctx.QueryParams(), &params.Txid)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter txid: %s", err))
	}

	// ------------- Optional query parameter "asset" -------------

	err = runtime.BindQueryParameter("form", true, false, "asset", ctx.QueryParams(), &params.Asset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter asset: %s", err))
	}

	// ------------- Optional query parameter "type" -------------

	err = runtime.BindQueryParameter("form", true, false, "type", ctx.QueryParams(), &params.Type)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter type: %s", err))
	}

	// ------------- Required query parameter "offset" -------------

	err = runtime.BindQueryParameter("form", true, true, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	// ------------- Required query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, true, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetTxDetails(ctx, params)
	return err
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}, si ServerInterface) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET("/v1/assets", wrapper.GetAssetInfo)
	router.GET("/v1/doc", wrapper.GetDocs)
	router.GET("/v1/health", wrapper.GetHealth)
	router.GET("/v1/network", wrapper.GetNetworkData)
	router.GET("/v1/nodes", wrapper.GetNodes)
	router.GET("/v1/pools", wrapper.GetPools)
	router.GET("/v1/pools/detail", wrapper.GetPoolsData)
	router.GET("/v1/stakers", wrapper.GetStakersData)
	router.GET("/v1/stakers/:address", wrapper.GetStakersAddressData)
	router.GET("/v1/stakers/:address/pools", wrapper.GetStakersAddressAndAssetData)
	router.GET("/v1/stats", wrapper.GetStats)
	router.GET("/v1/swagger.json", wrapper.GetSwagger)
	router.GET("/v1/thorchain/pool_addresses", wrapper.GetThorchainProxiedEndpoints)
	router.GET("/v1/txs", wrapper.GetTxDetails)

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

