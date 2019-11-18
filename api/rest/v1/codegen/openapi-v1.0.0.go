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
	"time"
)

// AssetDetail defines model for AssetDetail.
type AssetDetail struct {
	Asset       *Asset     `json:"asset,omitempty"`
	DateCreated *time.Time `json:"dateCreated,omitempty"`
	Logo        *string    `json:"logo,omitempty"`
	Name        *string    `json:"name,omitempty"`
	PriceRune   *float64   `json:"priceRune,omitempty"`
	PriceUSD    *float64   `json:"priceUSD,omitempty"`
}

// BEPSwapData defines model for BEPSwapData.
type BEPSwapData struct {
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
	Statue           *string  `json:"statue,omitempty"`
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
	Asset           *Asset     `json:"asset,omitempty"`
	AssetEarned     *int64     `json:"assetEarned,omitempty"`
	AssetROI        *float64   `json:"assetROI,omitempty"`
	AssetStaked     *int64     `json:"assetStaked,omitempty"`
	DateFirstStaked *time.Time `json:"dateFirstStaked,omitempty"`
	PoolEarned      *int64     `json:"poolEarned,omitempty"`
	PoolROI         *float64   `json:"poolROI,omitempty"`
	PoolStaked      *int64     `json:"poolStaked,omitempty"`
	RuneEarned      *int64     `json:"runeEarned,omitempty"`
	RuneROI         *float64   `json:"runeROI,omitempty"`
	RuneStaked      *int64     `json:"runeStaked,omitempty"`
	StakeUnits      *int64     `json:"stakeUnits,omitempty"`
}

// TxDetails defines model for TxDetails.
type TxDetails struct {
	Date   *time.Time `json:"date,omitempty"`
	Events *struct {
		Fee        *int64   `json:"fee,omitempty"`
		Slip       *float64 `json:"slip,omitempty"`
		StakeUnits *int64   `json:"stakeUnits,omitempty"`
	} `json:"events,omitempty"`
	Gas *struct {
		Amount *string `json:"amount,omitempty"`
		Asset  *Asset  `json:"asset,omitempty"`
	} `json:"gas,omitempty"`
	Height *int64 `json:"height,omitempty"`
	In     *struct {
		Address *string `json:"address,omitempty"`
		Coin    *struct {
			Amount *string `json:"amount,omitempty"`
			Asset  *Asset  `json:"asset,omitempty"`
		} `json:"coin,omitempty"`
		Memo *string `json:"memo,omitempty"`
		TxID *string `json:"txID,omitempty"`
	} `json:"in,omitempty"`
	Options *struct {
		Stake    *string `json:"Stake,omitempty"`
		Swap     *string `json:"Swap,omitempty"`
		Withdraw *string `json:"Withdraw,omitempty"`
	} `json:"options,omitempty"`
	Out *struct {
		Address *string `json:"address,omitempty"`
		Coin    *struct {
			Amount *string `json:"amount,omitempty"`
			Asset  *Asset  `json:"asset,omitempty"`
		} `json:"coin,omitempty"`
		Memo *string `json:"memo,omitempty"`
		TxID *string `json:"txID,omitempty"`
	} `json:"out,omitempty"`
	Pool   *Asset  `json:"pool,omitempty"`
	Status *string `json:"status,omitempty"`
	Type   *string `json:"type,omitempty"`
}

// Asset defines model for asset.
type Asset struct {
	Chain  *string `json:"chain,omitempty"`
	Symbol *string `json:"symbol,omitempty"`
	Ticker *string `json:"ticker,omitempty"`
}

// AssetsDetailedResponse defines model for AssetsDetailedResponse.
type AssetsDetailedResponse AssetDetail

// AssetsResponse defines model for AssetsResponse.
type AssetsResponse []Asset

// BEPSwapResponse defines model for BEPSwapResponse.
type BEPSwapResponse BEPSwapData

// GeneralErrorResponse defines model for GeneralErrorResponse.
type GeneralErrorResponse Error

// PoolsDetailedResponse defines model for PoolsDetailedResponse.
type PoolsDetailedResponse PoolDetail

// StakersAddressDataResponse defines model for StakersAddressDataResponse.
type StakersAddressDataResponse StakersAddressData

// StakersAssetDataResponse defines model for StakersAssetDataResponse.
type StakersAssetDataResponse StakersAssetData

// StakersResponse defines model for StakersResponse.
type StakersResponse []Stakers

// TxDetailedResponse defines model for TxDetailedResponse.
type TxDetailedResponse []TxDetails

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get Assets// (GET /v1/assets)
	GetAssets(ctx echo.Context) error
	// Get Asset Information// (GET /v1/assets/{asset})
	GetAssetInfo(ctx echo.Context, asset string) error
	// Get BEPSwap Data// (GET /v1/bepswap)
	GetBEPSwapData(ctx echo.Context) error
	// Get Documents// (GET /v1/doc)
	GetDocs(ctx echo.Context) error
	// Get Health// (GET /v1/health)
	GetHealth(ctx echo.Context) error
	// Get Pools Data// (GET /v1/pools/{asset})
	GetPoolsData(ctx echo.Context, asset string) error
	// Get Stakers// (GET /v1/stakers)
	GetStakersData(ctx echo.Context) error
	// Get Staker Data// (GET /v1/stakers/{address})
	GetStakersAddressData(ctx echo.Context, address string) error
	// Get Staker Pool Data// (GET /v1/stakers/{address}/{asset})
	GetStakersAddressAndAssetData(ctx echo.Context, address string, asset string) error
	// Get Swagger// (GET /v1/swagger.json)
	GetSwagger(ctx echo.Context) error
	// Get transaction// (GET /v1/tx/{address})
	GetTxDetails(ctx echo.Context, address string) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetAssets converts echo context to params.
func (w *ServerInterfaceWrapper) GetAssets(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetAssets(ctx)
	return err
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

// GetBEPSwapData converts echo context to params.
func (w *ServerInterfaceWrapper) GetBEPSwapData(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetBEPSwapData(ctx)
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

// GetSwagger converts echo context to params.
func (w *ServerInterfaceWrapper) GetSwagger(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetSwagger(ctx)
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

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router runtime.EchoRouter, si ServerInterface) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET("/v1/assets", wrapper.GetAssets)
	router.GET("/v1/assets/:asset", wrapper.GetAssetInfo)
	router.GET("/v1/bepswap", wrapper.GetBEPSwapData)
	router.GET("/v1/doc", wrapper.GetDocs)
	router.GET("/v1/health", wrapper.GetHealth)
	router.GET("/v1/pools/:asset", wrapper.GetPoolsData)
	router.GET("/v1/stakers", wrapper.GetStakersData)
	router.GET("/v1/stakers/:address", wrapper.GetStakersAddressData)
	router.GET("/v1/stakers/:address/:asset", wrapper.GetStakersAddressAndAssetData)
	router.GET("/v1/swagger.json", wrapper.GetSwagger)
	router.GET("/v1/tx/:address", wrapper.GetTxDetails)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+Qb227bOPZXCO0ukACuc+llijxt0nRmAmybIklnsZjtAy0dW+xIpEJSvkyR39of2B9b",
	"nENJliPKppwGe3uKY5LnfiWPv0WxygslQVoTnX2LNJhCSQP0z7kxYM0lWC4ySG6qJVyJlbQgLX7kRZGJ",
	"mFuh5NFXoyR+Z+IUco6f/qhhGp1Ffzhaozlyq+aIwDvo0cPDwyhKwMRaFAgqOovU5CvEliEqLqSQM5ZU",
	"lDCOJ5mQU6Vzwhw9jCpq96JSWMjNLnIJKSKyqwKis4hrzVc+wmmBqakj0+CRi/efbhe8+O4irOBecsvD",
	"RDjL1IRnrDrHEjo4in4CCZpn77VW+rsTSVB95AEusByM4TNAMj4plT2fvSH0IeZWKJWRhNhUaWZTbllj",
	"A7eW/wbanCeJBmNQ/t+d3C6K7caWZcymQFQb+mQIABOGPiFHQrZpJ/d7TsprBGHirolsJM6ZKSAWUxHX",
	"rHCZdFTwfA5fIRjm8pUWzPrs3fJJFh1Eao0jiNgbsKWWhnHJGrqt5tLwGHcwmFM2wIMV/CYZVO5z9i0q",
	"tCpAW+EyhVNKaABNuIV3GriFBM+4KB6d0fdRQ7+xWsgZ7s/UTOHGzoLkOXgXCi1iuCklbMJX5SRrYZBl",
	"PgHd7P98exm0fS1iZ8StAE/W3pFOwkW2Oo+tmMNng1aB322o5BJ3ME5bWIl72EEpxX0JjDv/B8OEtKBx",
	"j5wdRqM1nULaN6/WZOK2mWOLEN8t+/C1lG7CAOZK2nQHLx/cng1uBkH3EVwDHU4yxsN3qnSutgn0IymU",
	"3NbRSrEzDKxVlmfkExflyiOFO1xnk3KFIW041Wvwt5BlvfANZNkTEFxCYdM+2DefP75nE55xGcMQoO+5",
	"ls6xfVCBVtmBkA7BLzwr4XA8AD6FZZ+NVCKp0sieIiHoyTbgT6O+n/A9Ce7xQgeziiFmwYsCg8rfy+Pj",
	"0zdNdgrH8ovKShdrfWgeyaNOhK7+rbAnzAgZU2bUdjwY9emrVD8NvZAMgQzh+q/Cponmi36lLaod+1ic",
	"L4+4WrmTQaD++lGmexhFGu5LodFkf622ffHAbdW+T0ze3NUBWwJHXGoN0jIKX3UECZM6Ab+5vuqCdsA0",
	"lS5MSSbkHIzNsYQahWR4guycm6jsI97hcX4eRvOkXNGhnUkGTfQFeuBLOL+9fX+3h8dPytWPAOdz0Ngx",
	"dYXkFigtVJxMAZgRvwNV1B0K0HPo0+EQ/GarBKcAxgFGAsIB32ai2MmZ1TwBZjJR+BkSkv0pzCAm5epu",
	"uRNfZcPlRuWxFujBYwIO9xHp9uC6JmFO+3qRDhB3QeFgiw8nuIhGO1E2ZUYklU4R9QAcIcZKbfYUYBDU",
	"3SYYDs4bbzBgspvrK3bAKzorD6Ym1Onk5vrqMMzYKjQnp1sQqTlodnLKqAo24XCD/IaEjG4zAO62YIm9",
	"FZtjuq1rIrcxWOgBvkc0t9wuHPhnKWxvTURgS9zBVGmN5TLBXBoMvNdbF+rFgjde6oqQF1agxw73nK1V",
	"zyamurAZjAabXo9B4tdIvjPyBug4zHR0KSGoPCATGlQdIGivs5Jf7l8bINyA0oCwDKkMMA2HlQYUxquQ",
	"TmiGlwaILCTcYiLxlAYdCgYaU4U+sDQYCHif0qDDUHhpgCiDawOqtTrFwQGhX2M/3IftkLqA0NeFQT/S",
	"wIbLuN56p73Svn2M1DWeO+HXfeuQPrVq/IOpr3rhuoPbjxtbkn5Aljk2YBOlrLGaGs5oFIHkk4w+JcK4",
	"j18895t1ex4slrqdb90IEh+U2wJpRxiB4qq2bkgI3al+bmAHk3JlqDBCczSBJr5ouuudJOytJF+PXd/r",
	"d/DdVq8M7sYV9bfkeZHhaTuRk5Pp19Ps/uvbZK5fF2U+jdP4B2mz6X1yOn/ze7K8X3yFxfS17wbb85bT",
	"acPp8vGcLu2f+CQZdBHnijc1pRs5IWetwMx4rJUxdIlCVA25rvF37+saugaKNa8JrCmG38ztycMWc1m/",
	"Zn2PC5Q+3fxSa8W9ZtfXpVOt8sbdxk+9S6FuY0rwqiQiEhh8jeLRco5+3KLeKWUc+lhh4UehTQv8xssQ",
	"ldI+50KZ7G3r4yf2ibVlU6vYGN9G4zy05eqpb9vd1ji8YN5paQT8CXbWX5WvzYxqfRTGeGhFvs3KWgX5",
	"kPqmpzu8gUKDQZ9laiFBm1QUFLDCheELH+unWc+7oIVwM6+eZTtQprAJZAv32P+HPYhuCmovxmfcQyzP",
	"61Tf4W9IFPXhS0HMUhsoCiE9tFV530dcrLwnno+bHHL/m7ddXl36XwI6MFTh6qQO3eRaXui3C154F+qn",
	"kFDUpf0/EDCV26F5n7oF0+4WTBnHrtDUMC0l9gl0J/OOyxiyvm7BfdGCAhnW/JNyFY0iKno+uS6AfDha",
	"19keaA9tQW6KPk6500i3XVnlE8d3lzYR/wY6SHoP5IVTVQ/D8JiIgJzeqKIE5ubPNlWa6BgrTfAfJfcU",
	"2AeRzLhO2KdykomYnX+6YvclaAGG3f18ffMOT7u5IbliBMuwTEisFeeCU3K6EFP9z38YS9sKDQXXNG3R",
	"TBcyPlGlpb0S7ELp35hVbAJMA09ocGPORYb9Hd08FI4UGnwYMyQSqSq4NmC6gzbVvBMm3k2CjVVIh00h",
	"x2aLM0wML4zjDQ9NuAEkJKdeHBcTKEAmCLSWAXCzGjdCShQYJpVlqcoSFmthRcyzNqtjdqeattL1lPWs",
	"oLv3RjiwHFVDKiZVZZYQtlWL/ERoiG22wqRphaUOqquoaBTNQRuny5Px8fj4heLmpYtcIHkhorPoJX6P",
	"jsFtSnZ5ND85qgYrz75FM7C+TP5owKk1ZVYPaNWPw2VRKG0hYUrWnCLV6AkkkKskOot+AusqWnLV1oTs",
	"6fFxn/83+44eDaY+jKJXIce8M5k0lFXmOdcrR1dVatPCWjRH3+jvQ6+I6pE0j5m3xu4IyJjV8gSpylm6",
	"ccQqlghTZHzFeH09UU/nsjnXQpWGbMUZ1ZTHYEZMyDgrE1RHxi0YyyjsOV/J1EzRzKIqdVyXpFy685Jn",
	"jen3q+kKwwqajOY5WGr1f30sgM9tWg/e/Xx+9XF8+7cPF9d/Odxo/C8+Xozvrj9cX7w4eX8SucKBjDGq",
	"x8+qANp+ibe6hFFrhu9xMPyyvxl1Zgmfw5zYVXu0urKsCRSmqg52eV3Q3LGb8azbcxejN69SjVfF7VG7",
	"feT4eBa7K4Oaymp21bGfqLiX9dsFn81AH1Vhi70cHzdO5PxkRtLHOJOouMyRJC9zlyruCzKbKE0Pyk1M",
	"xsPdZU0ABmc+Q+eI6u+cxr/UPKfAM/eQM1zjGGTdeVZz01w6fLryMv+zQxfC/jbIXZYrwDVbZHA7Q+RW",
	"7trT1r4Z8TE7X08grVMLS/mcgo6KBemouZrqCMPNwjsj/x+LZP4x/67aaN+mF5r13e3g1E9pv5mmraNP",
	"BdGrg+rqb+9Q83hEvctiM2O+yd/Rt4rQhycVOdt/C7CN5fZNdZj9mf7L84mcnHxdTtPT2dvX9y/nxza5",
	"f/1mKmG+fLOMlzaWqTV5XL55lfeYZQPzmQ1zy686+lTnNc+1+p4WZTo/iSBVusQCSftXEfUV1Q59nstk",
	"fY/9X6nX0bbw9x8a73p/cNNrVDQUtGlZLt+P6x+IbDWmtMy5ayJzHqdCuk6VGtTHdcNGmeK3IHciLC3v",
	"idgniQZtXaTcbpxoihS7HBYtv0upsr7V3eFF5/I7vCf+O2Oi5+dLXWW1p7ToZseAntcCKXWGpmFtcXZ0",
	"dHL6A3b045Ozt8dvjyP05/W68Wz48vCvAAAA//8DM7+0pDoAAA==",
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

