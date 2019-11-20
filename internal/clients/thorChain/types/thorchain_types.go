package types

import (
	"encoding/json"

	"gitlab.com/thorchain/bepswap/chain-service/internal/common"
)

type Event struct {
	ID     int64           `json:"id,string"`
	Status string          `json:"status"`
	Height int64           `json:"height,string"`
	Type   string          `json:"type"`
	InTx   common.Tx       `json:"in_tx"`
	OutTxs common.Txs      `json:"out_txs"`
	Gas    common.Coins    `json:"gas"`
	Event  json.RawMessage `json:"event"`
}

type EventStake struct {
	Pool       common.Asset `json:"pool"`
	StakeUnits int64        `json:"stake_units,string"`
}

type EventSwap struct {
	Pool         common.Asset `json:"pool"`
	PriceTarget  int64        `json:"price_target,string"`
	TradeSlip    float64      `json:"trade_slip,string"`
	LiquidityFee int64        `json:"liquidity_fee,string"`
}

type EventUnstake struct {
	Pool       common.Asset `json:"pool"`
	StakeUnits int64        `json:"stake_units,string"`
}