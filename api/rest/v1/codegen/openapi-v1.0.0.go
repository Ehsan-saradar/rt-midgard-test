// Package api provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package api

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
	Asset       *Asset   `json:"asset,omitempty"`
	DateCreated *int64   `json:"dateCreated,omitempty"`
	Logo        *string  `json:"logo,omitempty"`
	Name        *string  `json:"name,omitempty"`
	PriceRune   *float64 `json:"priceRune,omitempty"`
}

// Error defines model for Error.
type Error struct {
	Error string `json:"error"`
}

// PoolDetail defines model for PoolDetail.
type PoolDetail struct {
	Asset            *Asset   `json:"asset,omitempty"`
	AssetDepth       *int64   `json:"assetDepth,omitempty"`
	AssetROI         *float64 `json:"assetROI,omitempty"`
	AssetStakedTotal *int64   `json:"assetStakedTotal,omitempty"`
	BuyAssetCount    *int64   `json:"buyAssetCount,omitempty"`
	BuyFeeAverage    *int64   `json:"buyFeeAverage,omitempty"`
	BuyFeesTotal     *int64   `json:"buyFeesTotal,omitempty"`
	BuySlipAverage   *float64 `json:"buySlipAverage,omitempty"`
	BuyTxAverage     *int64   `json:"buyTxAverage,omitempty"`
	BuyVolume        *int64   `json:"buyVolume,omitempty"`
	PoolDepth        *int64   `json:"poolDepth,omitempty"`
	PoolFeeAverage   *int64   `json:"poolFeeAverage,omitempty"`
	PoolFeesTotal    *int64   `json:"poolFeesTotal,omitempty"`
	PoolROI          *float64 `json:"poolROI,omitempty"`
	PoolROI12        *float64 `json:"poolROI12,omitempty"`
	PoolSlipAverage  *float64 `json:"poolSlipAverage,omitempty"`
	PoolStakedTotal  *int64   `json:"poolStakedTotal,omitempty"`
	PoolTxAverage    *int64   `json:"poolTxAverage,omitempty"`
	PoolUnits        *int64   `json:"poolUnits,omitempty"`
	PoolVolume       *int64   `json:"poolVolume,omitempty"`
	PoolVolume24hr   *int64   `json:"poolVolume24hr,omitempty"`
	Price            *float64 `json:"price,omitempty"`
	RuneDepth        *int64   `json:"runeDepth,omitempty"`
	RuneROI          *float64 `json:"runeROI,omitempty"`
	RuneStakedTotal  *int64   `json:"runeStakedTotal,omitempty"`
	SellAssetCount   *int64   `json:"sellAssetCount,omitempty"`
	SellFeeAverage   *int64   `json:"sellFeeAverage,omitempty"`
	SellFeesTotal    *int64   `json:"sellFeesTotal,omitempty"`
	SellSlipAverage  *float64 `json:"sellSlipAverage,omitempty"`
	SellTxAverage    *int64   `json:"sellTxAverage,omitempty"`
	SellVolume       *int64   `json:"sellVolume,omitempty"`
	StakeTxCount     *int64   `json:"stakeTxCount,omitempty"`
	StakersCount     *int64   `json:"stakersCount,omitempty"`
	StakingTxCount   *int64   `json:"stakingTxCount,omitempty"`
	Status           *string  `json:"status,omitempty"`
	SwappersCount    *int64   `json:"swappersCount,omitempty"`
	SwappingTxCount  *int64   `json:"swappingTxCount,omitempty"`
	WithdrawTxCount  *int64   `json:"withdrawTxCount,omitempty"`
}

// Stakers defines model for Stakers.
type Stakers string

// StakersAddressData defines model for StakersAddressData.
type StakersAddressData struct {
	PoolsArray  *[]Asset `json:"poolsArray,omitempty"`
	TotalEarned *int64   `json:"totalEarned,omitempty"`
	TotalROI    *float64 `json:"totalROI,omitempty"`
	TotalStaked *int64   `json:"totalStaked,omitempty"`
}

// StakersAssetData defines model for StakersAssetData.
type StakersAssetData struct {
	Asset           *Asset   `json:"asset,omitempty"`
	AssetEarned     *int64   `json:"assetEarned,omitempty"`
	AssetROI        *float64 `json:"assetROI,omitempty"`
	AssetStaked     *int64   `json:"assetStaked,omitempty"`
	DateFirstStaked *int64   `json:"dateFirstStaked,omitempty"`
	PoolEarned      *int64   `json:"poolEarned,omitempty"`
	PoolROI         *float64 `json:"poolROI,omitempty"`
	PoolStaked      *int64   `json:"poolStaked,omitempty"`
	RuneEarned      *int64   `json:"runeEarned,omitempty"`
	RuneROI         *float64 `json:"runeROI,omitempty"`
	RuneStaked      *int64   `json:"runeStaked,omitempty"`
	StakeUnits      *int64   `json:"stakeUnits,omitempty"`
}

// StatsData defines model for StatsData.
type StatsData struct {
	DailyActiveUsers   *int64 `json:"dailyActiveUsers,omitempty"`
	DailyTx            *int64 `json:"dailyTx,omitempty"`
	MonthlyActiveUsers *int64 `json:"monthlyActiveUsers,omitempty"`
	MonthlyTx          *int64 `json:"monthlyTx,omitempty"`
	PoolCount          *int64 `json:"poolCount,omitempty"`
	TotalAssetBuys     *int64 `json:"totalAssetBuys,omitempty"`
	TotalAssetSells    *int64 `json:"totalAssetSells,omitempty"`
	TotalDepth         *int64 `json:"totalDepth,omitempty"`
	TotalEarned        *int64 `json:"totalEarned,omitempty"`
	TotalStakeTx       *int64 `json:"totalStakeTx,omitempty"`
	TotalStaked        *int64 `json:"totalStaked,omitempty"`
	TotalTx            *int64 `json:"totalTx,omitempty"`
	TotalUsers         *int64 `json:"totalUsers,omitempty"`
	TotalVolume        *int64 `json:"totalVolume,omitempty"`
	TotalVolume24hr    *int64 `json:"totalVolume24hr,omitempty"`
	TotalWithdrawTx    *int64 `json:"totalWithdrawTx,omitempty"`
}

// ThorchainEndpoint defines model for ThorchainEndpoint.
type ThorchainEndpoint struct {
	Address *string `json:"address,omitempty"`
	Chain   *string `json:"chain,omitempty"`
	PubKey  *string `json:"pub_key,omitempty"`
}

// TxDetails defines model for TxDetails.
type TxDetails struct {
	Date    *int64  `json:"date,omitempty"`
	Events  *Event  `json:"events,omitempty"`
	Gas     *Gas    `json:"gas,omitempty"`
	Height  *int64  `json:"height,omitempty"`
	In      *Tx     `json:"in,omitempty"`
	Options *Option `json:"options,omitempty"`
	Out     *Tx     `json:"out,omitempty"`
	Pool    *Asset  `json:"pool,omitempty"`
	Status  *string `json:"status,omitempty"`
	Type    *string `json:"type,omitempty"`
}

// Asset defines model for asset.
type Asset struct {
	Chain  *string `json:"chain,omitempty"`
	Symbol *string `json:"symbol,omitempty"`
	Ticker *string `json:"ticker,omitempty"`
}

// Coin defines model for coin.
type Coin struct {
	Amount *int64 `json:"amount,omitempty"`
	Asset  *Asset `json:"asset,omitempty"`
}

// Coins defines model for coins.
type Coins []Coin

// Event defines model for event.
type Event struct {
	Fee        *int64   `json:"fee,omitempty"`
	Slip       *float64 `json:"slip,omitempty"`
	StakeUnits *int64   `json:"stakeUnits,omitempty"`
}

// Gas defines model for gas.
type Gas struct {
	Amount *int64 `json:"amount,omitempty"`
	Asset  *Asset `json:"asset,omitempty"`
}

// Option defines model for option.
type Option struct {
	Asymmetry           *float64 `json:"asymmetry,omitempty"`
	PriceTarget         *int64   `json:"priceTarget,omitempty"`
	WithdrawBasisPoints *int64   `json:"withdrawBasisPoints,omitempty"`
}

// Tx defines model for tx.
type Tx struct {
	Address *string `json:"address,omitempty"`
	Coins   *Coins  `json:"coins,omitempty"`
	Memo    *string `json:"memo,omitempty"`
	TxID    *string `json:"txID,omitempty"`
}

// AssetsDetailedResponse defines model for AssetsDetailedResponse.
type AssetsDetailedResponse AssetDetail

// GeneralErrorResponse defines model for GeneralErrorResponse.
type GeneralErrorResponse Error

// PoolsDetailedResponse defines model for PoolsDetailedResponse.
type PoolsDetailedResponse PoolDetail

// PoolsResponse defines model for PoolsResponse.
type PoolsResponse []Asset

// StakersAddressDataResponse defines model for StakersAddressDataResponse.
type StakersAddressDataResponse StakersAddressData

// StakersAssetDataResponse defines model for StakersAssetDataResponse.
type StakersAssetDataResponse StakersAssetData

// StakersResponse defines model for StakersResponse.
type StakersResponse []Stakers

// StatsResponse defines model for StatsResponse.
type StatsResponse StatsData

// ThorchainEndpointsResponse defines model for ThorchainEndpointsResponse.
type ThorchainEndpointsResponse struct {
	Current *[]ThorchainEndpoint `json:"current,omitempty"`
}

// TxDetailedResponse defines model for TxDetailedResponse.
type TxDetailedResponse []TxDetails

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get Asset Information// (GET /v1/assets/{asset})
	GetAssetInfo(ctx echo.Context, asset string) error
	// Get Documents// (GET /v1/doc)
	GetDocs(ctx echo.Context) error
	// Get Health// (GET /v1/health)
	GetHealth(ctx echo.Context) error
	// Get Asset Pools// (GET /v1/pools)
	GetPools(ctx echo.Context) error
	// Get Pools Data// (GET /v1/pools/{asset})
	GetPoolsData(ctx echo.Context, asset string) error
	// Get Stakers// (GET /v1/stakers)
	GetStakersData(ctx echo.Context) error
	// Get Staker Data// (GET /v1/stakers/{address})
	GetStakersAddressData(ctx echo.Context, address string) error
	// Get Staker Pool Data// (GET /v1/stakers/{address}/{asset})
	GetStakersAddressAndAssetData(ctx echo.Context, address string, asset string) error
	// Get Global Stats// (GET /v1/stats)
	GetStats(ctx echo.Context) error
	// Get Swagger// (GET /v1/swagger.json)
	GetSwagger(ctx echo.Context) error
	// Get the Proxied Pool Addresses// (GET /v1/thorchain/pool_addresses)
	GetThorchainProxiedEndpoints(ctx echo.Context) error
	// Get transaction// (GET /v1/tx/asset/{asset})
	GetTxDetailsByAsset(ctx echo.Context, asset string) error
	// Get transaction// (GET /v1/tx/{address})
	GetTxDetails(ctx echo.Context, address string) error
	// Get transaction// (GET /v1/tx/{address}/asset/{asset})
	GetTxDetailsByAddressAsset(ctx echo.Context, address string, asset string) error
	// Get transaction// (GET /v1/tx/{address}/txid/{txid})
	GetTxDetailsByAddressTxId(ctx echo.Context, address string, txid string) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetAssetInfo converts echo context to params.
func (w *ServerInterfaceWrapper) GetAssetInfo(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "asset" -------------
	var asset string

	err = runtime.BindStyledParameter("simple", false, "asset", ctx.Param("asset"), &asset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter asset: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetAssetInfo(ctx, asset)
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
	// ------------- Path parameter "asset" -------------
	var asset string

	err = runtime.BindStyledParameter("simple", false, "asset", ctx.Param("asset"), &asset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter asset: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetPoolsData(ctx, asset)
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

	// ------------- Path parameter "asset" -------------
	var asset string

	err = runtime.BindStyledParameter("simple", false, "asset", ctx.Param("asset"), &asset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter asset: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetStakersAddressAndAssetData(ctx, address, asset)
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

// GetTxDetailsByAsset converts echo context to params.
func (w *ServerInterfaceWrapper) GetTxDetailsByAsset(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "asset" -------------
	var asset string

	err = runtime.BindStyledParameter("simple", false, "asset", ctx.Param("asset"), &asset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter asset: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetTxDetailsByAsset(ctx, asset)
	return err
}

// GetTxDetails converts echo context to params.
func (w *ServerInterfaceWrapper) GetTxDetails(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "address" -------------
	var address string

	err = runtime.BindStyledParameter("simple", false, "address", ctx.Param("address"), &address)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter address: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetTxDetails(ctx, address)
	return err
}

// GetTxDetailsByAddressAsset converts echo context to params.
func (w *ServerInterfaceWrapper) GetTxDetailsByAddressAsset(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "address" -------------
	var address string

	err = runtime.BindStyledParameter("simple", false, "address", ctx.Param("address"), &address)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter address: %s", err))
	}

	// ------------- Path parameter "asset" -------------
	var asset string

	err = runtime.BindStyledParameter("simple", false, "asset", ctx.Param("asset"), &asset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter asset: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetTxDetailsByAddressAsset(ctx, address, asset)
	return err
}

// GetTxDetailsByAddressTxId converts echo context to params.
func (w *ServerInterfaceWrapper) GetTxDetailsByAddressTxId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "address" -------------
	var address string

	err = runtime.BindStyledParameter("simple", false, "address", ctx.Param("address"), &address)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter address: %s", err))
	}

	// ------------- Path parameter "txid" -------------
	var txid string

	err = runtime.BindStyledParameter("simple", false, "txid", ctx.Param("txid"), &txid)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter txid: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetTxDetailsByAddressTxId(ctx, address, txid)
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

	router.GET("/v1/assets/:asset", wrapper.GetAssetInfo)
	router.GET("/v1/doc", wrapper.GetDocs)
	router.GET("/v1/health", wrapper.GetHealth)
	router.GET("/v1/pools", wrapper.GetPools)
	router.GET("/v1/pools/:asset", wrapper.GetPoolsData)
	router.GET("/v1/stakers", wrapper.GetStakersData)
	router.GET("/v1/stakers/:address", wrapper.GetStakersAddressData)
	router.GET("/v1/stakers/:address/:asset", wrapper.GetStakersAddressAndAssetData)
	router.GET("/v1/stats", wrapper.GetStats)
	router.GET("/v1/swagger.json", wrapper.GetSwagger)
	router.GET("/v1/thorchain/pool_addresses", wrapper.GetThorchainProxiedEndpoints)
	router.GET("/v1/tx/asset/:asset", wrapper.GetTxDetailsByAsset)
	router.GET("/v1/tx/:address", wrapper.GetTxDetails)
	router.GET("/v1/tx/:address/asset/:asset", wrapper.GetTxDetailsByAddressAsset)
	router.GET("/v1/tx/:address/txid/:txid", wrapper.GetTxDetailsByAddressTxId)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"jINDAAA=",
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

