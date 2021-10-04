// End to end tests here are checking lookup functionality from Database to HTTP Api.
package timeseries_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/thorchain/midgard/internal/db/testdb"
	"gitlab.com/thorchain/midgard/internal/util"
	"gitlab.com/thorchain/midgard/openapi/generated/oapigen"
)

func TestActionsE2E(t *testing.T) {
	blocks := testdb.InitTestBlocks(t)

	blocks.NewBlock(t, "2020-09-01 00:00:00",
		testdb.AddLiquidity{
			Pool:        "BNB.TWT-123",
			RuneAddress: "thoraddr1",
			AssetAmount: 1000,
			RuneAmount:  2000,
		},
		testdb.PoolActivate{Pool: "BNB.TWT-123"},
	)

	blocks.NewBlock(t, "2020-09-02 00:00:00",
		testdb.Withdraw{
			Pool:                   "BNB.TWT-123",
			EmitAsset:              10,
			EmitRune:               20,
			Coin:                   "10 BNB.TWT-123",
			ToAddress:              "thoraddr4",
			LiquidityProviderUnits: 1,
			ImpLossProtection:      7,
		},
	)

	blocks.NewBlock(t, "2020-09-03 00:00:00",
		testdb.Swap{
			Coin:      "100000 BNB.BNB",
			EmitAsset: "10 THOR.RUNE",
			Pool:      "BNB.BNB",
		},
		testdb.PoolActivate{Pool: "BNB.BNB"},
	)

	// Basic request with no filters (should get all events ordered by height)
	body := testdb.CallJSON(t, "http://localhost:8080/v2/actions?limit=50&offset=0")

	var v oapigen.ActionsResponse
	testdb.MustUnmarshal(t, body, &v)

	require.Equal(t, "3", v.Count)

	basicTx0 := v.Actions[0]
	basicTx1 := v.Actions[1]
	basicTx2 := v.Actions[2]

	if basicTx0.Type != "swap" || basicTx0.Height != "3" {
		t.Fatal("Results of results changed.")
	}
	if basicTx1.Type != "withdraw" || basicTx1.Height != "2" {
		t.Fatal("Results of results changed.")
	}
	assert.Equal(t, "7", basicTx1.Metadata.Withdraw.ImpermanentLossProtection)
	if basicTx2.Type != "addLiquidity" || basicTx2.Height != "1" {
		t.Fatal("Results of results changed.")
	}

	// Filter by type request
	body = testdb.CallJSON(t, "http://localhost:8080/v2/actions?limit=50&offset=0&type=swap")

	testdb.MustUnmarshal(t, body, &v)

	require.Equal(t, "1", v.Count)
	typeTx0 := v.Actions[0]

	if typeTx0.Type != "swap" {
		t.Fatal("Results of results changed.")
	}

	// Filter by asset request
	body = testdb.CallJSON(t, "http://localhost:8080/v2/actions?limit=50&offset=0&asset=BNB.TWT-123")

	testdb.MustUnmarshal(t, body, &v)

	require.Equal(t, "2", v.Count)
	assetTx0 := v.Actions[0]
	assetTx1 := v.Actions[1]

	if assetTx0.Type != "withdraw" {
		t.Fatal("Results of results changed.")
	}
	if assetTx1.Type != "addLiquidity" {
		t.Fatal("Results of results changed.")
	}
}

func txResponseCount(t *testing.T, url string) string {
	body := testdb.CallJSON(t, url)

	var v oapigen.ActionsResponse
	testdb.MustUnmarshal(t, body, &v)
	return v.Count
}

func TestDepositStakeByTxIds(t *testing.T) {
	blocks := testdb.InitTestBlocks(t)

	blocks.NewBlock(t, "2020-09-01 00:00:00",
		testdb.AddLiquidity{
			Pool:        "BNB.TWT-123",
			RuneAddress: "thoraddr1",
			RuneTxID:    "RUNETX1",
			AssetTxID:   "ASSETTX1",
			AssetAmount: 1000,
			RuneAmount:  2000,
		},
		testdb.PoolActivate{Pool: "BNB.TWT-123"})

	require.Equal(t, "1", txResponseCount(t,
		"http://localhost:8080/v2/actions?limit=50&offset=0"))
	require.Equal(t, "0", txResponseCount(t,
		"http://localhost:8080/v2/actions?txid=NOSUCHID&limit=50&offset=0"))
	require.Equal(t, "1", txResponseCount(t,
		"http://localhost:8080/v2/actions?txid=ASSETTX1&limit=50&offset=0"))
	require.Equal(t, "1", txResponseCount(t,
		"http://localhost:8080/v2/actions?txid=RUNETX1&limit=50&offset=0"))
}

func TestPendingAlone(t *testing.T) {
	blocks := testdb.InitTestBlocks(t)

	blocks.NewBlock(t, "2020-01-01 00:00:00",
		testdb.PoolActivate{Pool: "BTC.BTC"},
		testdb.PoolActivate{Pool: "LTC.LTC"})

	blocks.NewBlock(t, "2020-09-01 00:00:00",
		testdb.PendingLiquidity{
			Pool:         "BTC.BTC",
			RuneAddress:  "thoraddr1",
			AssetAddress: "btcaddr1",
			RuneTxID:     "RUNETX1",
			AssetAmount:  0,
			RuneAmount:   20,
		})

	body := testdb.CallJSON(t, "http://localhost:8080/v2/actions?limit=50&offset=0")

	var v oapigen.ActionsResponse
	testdb.MustUnmarshal(t, body, &v)

	require.Equal(t, 1, len(v.Actions))
	add := v.Actions[0]

	require.Equal(t, "addLiquidity", string(add.Type))
	require.Equal(t, "pending", string(add.Status))
}

func TestPendingWithAdd(t *testing.T) {
	blocks := testdb.InitTestBlocks(t)

	blocks.NewBlock(t, "2020-01-01 00:00:00",
		testdb.PoolActivate{Pool: "BTC.BTC"},
		testdb.PoolActivate{Pool: "LTC.LTC"})

	blocks.NewBlock(t, "2020-09-01 00:00:00",
		testdb.PendingLiquidity{
			Pool:         "BTC.BTC",
			RuneAddress:  "thoraddr1",
			AssetAddress: "btcaddr1",
			RuneTxID:     "RUNETX1",
			AssetAmount:  0,
			RuneAmount:   20,
		},
		testdb.AddLiquidity{
			Pool:        "BTC.BTC",
			RuneAddress: "thoraddr1",
			RuneTxID:    "RUNETX1",
			AssetTxID:   "ASSETTX1",
			AssetAmount: 10,
			RuneAmount:  20,
		})

	body := testdb.CallJSON(t, "http://localhost:8080/v2/actions?limit=50&offset=0")

	var v oapigen.ActionsResponse
	testdb.MustUnmarshal(t, body, &v)

	require.Equal(t, 1, len(v.Actions))
	add := v.Actions[0]

	require.Equal(t, "addLiquidity", string(add.Type))
	require.Equal(t, "success", string(add.Status))
}

func TestPendingWithdrawn(t *testing.T) {
	// TODO(muninn): report withdraws of pending too. Currently we simply don't show anything.
	blocks := testdb.InitTestBlocks(t)

	blocks.NewBlock(t, "2020-01-01 00:00:00",
		testdb.PoolActivate{Pool: "BTC.BTC"},
		testdb.PoolActivate{Pool: "LTC.LTC"})

	blocks.NewBlock(t, "2020-09-01 00:00:00",
		testdb.PendingLiquidity{
			Pool:         "BTC.BTC",
			RuneAddress:  "thoraddr1",
			AssetAddress: "btcaddr1",
			RuneTxID:     "RUNETX1",
			AssetAmount:  0,
			RuneAmount:   20,
		})

	blocks.NewBlock(t, "2020-09-01 00:00:01",
		testdb.PendingLiquidity{
			Pool:         "BTC.BTC",
			RuneAddress:  "thoraddr1",
			AssetAddress: "btcaddr1",
			RuneTxID:     "RUNETX1",
			AssetAmount:  0,
			RuneAmount:   20,
			PendingType:  testdb.PendingWithdraw,
		})

	body := testdb.CallJSON(t, "http://localhost:8080/v2/actions?limit=50&offset=0")

	var v oapigen.ActionsResponse
	testdb.MustUnmarshal(t, body, &v)

	require.Equal(t, 0, len(v.Actions))
}

// TestSingleSwap swaps BNB.BNB -> THOR.RUNE and THOR.RUNE -> BTC.BTC
func TestSingleSwap(t *testing.T) {
	blocks := testdb.InitTestBlocks(t)

	blocks.NewBlock(t, "2020-08-02 00:00:00",
		testdb.PoolActivate{Pool: "BNB.BNB"},
		testdb.PoolActivate{Pool: "BTC.BTC"},
	)
	blocks.NewBlock(t, "2020-08-02 00:01:00",
		testdb.Outbound{
			TxID:      "00000",
			InTxID:    "12345",
			Coin:      "100 THOR.RUNE",
			ToAddress: "THOR1",
		},
		testdb.Swap{
			TxID:               "12345",
			Coin:               "100000 BNB.BNB",
			EmitAsset:          "100 THOR.RUNE",
			Pool:               "BNB.BNB",
			Slip:               100,
			LiquidityFeeInRune: 10000,
			FromAddress:        "BNB1",
			ToAddress:          "VAULT",
		},
	)
	blocks.NewBlock(t, "2020-08-02 00:02:00",
		testdb.Outbound{
			TxID:      "12121",
			InTxID:    "67890",
			Coin:      "55000 BTC.BTC",
			ToAddress: "BTC1",
		},
		testdb.Swap{
			TxID:               "67890",
			Coin:               "100 THOR.RUNE",
			EmitAsset:          "55000 BTC.BTC",
			Pool:               "BTC.BTC",
			Slip:               200,
			LiquidityFeeInRune: 20000,
			PriceTarget:        50000,
			FromAddress:        "THOR1",
			ToAddress:          "VAULT",
		},
	)

	body := testdb.CallJSON(t, "http://localhost:8080/v2/actions?limit=50&offset=0&type=swap")

	var v oapigen.ActionsResponse
	testdb.MustUnmarshal(t, body, &v)

	require.Equal(t, 2, len(v.Actions))

	swap1 := v.Actions[0]
	metadata1 := swap1.Metadata.Swap
	in1 := swap1.In[0]
	out1 := swap1.Out[0]
	pools1 := swap1.Pools
	swap2 := v.Actions[1]
	metadata2 := swap2.Metadata.Swap
	in2 := swap2.In[0]
	out2 := swap2.Out[0]
	pools2 := swap2.Pools

	require.Equal(t, "200", metadata1.SwapSlip)
	require.Equal(t, "20000", metadata1.LiquidityFee)
	require.Equal(t, "50000", metadata1.SwapTarget)
	require.Equal(t, "thor1", in1.Address)
	require.Equal(t, "THOR.RUNE", in1.Coins[0].Asset)
	require.Equal(t, "100", in1.Coins[0].Amount)
	require.Equal(t, "btc1", out1.Address)
	require.Equal(t, "BTC.BTC", out1.Coins[0].Asset)
	require.Equal(t, "55000", out1.Coins[0].Amount)
	require.Equal(t, 1, len(pools1))
	require.Equal(t, "BTC.BTC", pools1[0])

	require.Equal(t, "100", metadata2.SwapSlip)
	require.Equal(t, "10000", metadata2.LiquidityFee)
	require.Equal(t, "0", metadata2.SwapTarget)
	require.Equal(t, "bnb1", in2.Address)
	require.Equal(t, "BNB.BNB", in2.Coins[0].Asset)
	require.Equal(t, "100000", in2.Coins[0].Amount)
	require.Equal(t, "thor1", out2.Address)
	require.Equal(t, "THOR.RUNE", out2.Coins[0].Asset)
	require.Equal(t, "100", out2.Coins[0].Amount)
	require.Equal(t, 1, len(pools2))
	require.Equal(t, "BNB.BNB", pools2[0])
}

// TestDoubleSwap swaps BNB.BNB -> BTC.BTC
func TestDoubleSwap(t *testing.T) {
	blocks := testdb.InitTestBlocks(t)

	blocks.NewBlock(t, "2020-09-02 00:00:00",
		testdb.Outbound{
			TxID:      "2345",
			InTxID:    "1234",
			Coin:      "55000 BTC.BTC",
			ToAddress: "BTC1",
		},
		testdb.Outbound{
			TxID:      "00000",
			InTxID:    "1234",
			Coin:      "10 THOR.RUNE",
			ToAddress: "BNB1",
		},
		testdb.Swap{
			TxID:               "1234",
			Coin:               "100000 BNB.BNB",
			EmitAsset:          "10 THOR.RUNE",
			Pool:               "BNB.BNB",
			Slip:               100,
			LiquidityFeeInRune: 10000,
			FromAddress:        "BNB1",
			ToAddress:          "VAULT",
		},
		testdb.Swap{
			TxID:               "1234",
			Coin:               "10 THOR.RUNE",
			EmitAsset:          "55000 BTC.BTC",
			Pool:               "BTC.BTC",
			Slip:               200,
			LiquidityFeeInRune: 20000,
			PriceTarget:        50000,
			FromAddress:        "BNB1",
			ToAddress:          "VAULT",
		},
		testdb.PoolActivate{Pool: "BNB.BNB"},
		testdb.PoolActivate{Pool: "BTC.BTC"},
	)

	body := testdb.CallJSON(t, "http://localhost:8080/v2/actions?limit=50&offset=0&type=swap")

	var v oapigen.ActionsResponse
	testdb.MustUnmarshal(t, body, &v)

	doubleSwap := v.Actions[0]
	metadata := doubleSwap.Metadata.Swap
	in := doubleSwap.In[0]
	out := doubleSwap.Out[0]
	pools := doubleSwap.Pools
	require.Equal(t, "298", metadata.SwapSlip) // 100+200-(100*200)/10000
	require.Equal(t, "30000", metadata.LiquidityFee)
	require.Equal(t, "50000", metadata.SwapTarget)
	require.Equal(t, "bnb1", in.Address)
	require.Equal(t, "BNB.BNB", in.Coins[0].Asset)
	require.Equal(t, "100000", in.Coins[0].Amount)
	require.Equal(t, "btc1", out.Address)
	require.Equal(t, "BTC.BTC", out.Coins[0].Asset)
	require.Equal(t, "55000", out.Coins[0].Amount)
	require.Equal(t, 2, len(pools))
	require.Equal(t, "BNB.BNB", pools[0])
	require.Equal(t, "BTC.BTC", pools[1])
}

// TestDoubleSwapSynthToNativeSamePool swaps BTC/BTC -> BTC.BTC
func TestDoubleSwapSynthToNativeSamePool(t *testing.T) {
	blocks := testdb.InitTestBlocks(t)

	blocks.NewBlock(t, "2020-09-02 00:00:00",
		testdb.Outbound{
			TxID:      "22222",
			InTxID:    "11111",
			Coin:      "490 BTC.BTC",
			ToAddress: "BTC1",
		},
		testdb.Outbound{
			TxID:      "00000",
			InTxID:    "11111",
			Coin:      "10 THOR.RUNE",
			ToAddress: "THOR1",
		},
		testdb.Swap{
			TxID:               "11111",
			Coin:               "500 BTC/BTC",
			EmitAsset:          "10 THOR.RUNE",
			Pool:               "BTC.BTC",
			Slip:               100,
			LiquidityFeeInRune: 10000,
			FromAddress:        "THOR1",
			ToAddress:          "VAULT",
		},
		testdb.Swap{
			TxID:               "11111",
			Coin:               "10 THOR.RUNE",
			EmitAsset:          "490 BTC.BTC",
			Pool:               "BTC.BTC",
			Slip:               200,
			LiquidityFeeInRune: 20000,
			PriceTarget:        50000,
			FromAddress:        "THOR1",
			ToAddress:          "VAULT",
		},
		testdb.PoolActivate{Pool: "BTC.BTC"},
	)

	body := testdb.CallJSON(t, "http://localhost:8080/v2/actions?limit=50&offset=0&type=swap")

	var v oapigen.ActionsResponse
	testdb.MustUnmarshal(t, body, &v)

	doubleSwap := v.Actions[0]
	metadata := doubleSwap.Metadata.Swap
	in := doubleSwap.In[0]
	out := doubleSwap.Out[0]
	pools := doubleSwap.Pools
	require.Equal(t, "298", metadata.SwapSlip) // 100+200-(100*200)/10000
	require.Equal(t, "30000", metadata.LiquidityFee)
	require.Equal(t, "50000", metadata.SwapTarget)
	require.Equal(t, "thor1", in.Address)
	require.Equal(t, "BTC/BTC", in.Coins[0].Asset)
	require.Equal(t, "500", in.Coins[0].Amount)
	require.Equal(t, "btc1", out.Address)
	require.Equal(t, "BTC.BTC", out.Coins[0].Asset)
	require.Equal(t, "490", out.Coins[0].Amount)
	require.Equal(t, 1, len(pools))
	require.Equal(t, "BTC.BTC", pools[0])
}

// TestDoubleSwapNativeToSynthSamePool swaps BNB.BNB -> BNB/BNB
func TestDoubleSwapNativeToSynthSamePool(t *testing.T) {
	blocks := testdb.InitTestBlocks(t)

	blocks.NewBlock(t, "2020-09-02 00:00:00",
		testdb.Outbound{
			TxID:      "00000",
			InTxID:    "33333",
			Coin:      "190 BNB/BNB",
			ToAddress: "THOR1",
		},
		testdb.Outbound{
			TxID:      "00000",
			InTxID:    "33333",
			Coin:      "10 THOR.RUNE",
			ToAddress: "BNB1",
		},
		testdb.Swap{
			TxID:               "33333",
			Coin:               "200 BNB.BNB",
			EmitAsset:          "10 THOR.RUNE",
			Pool:               "BNB.BNB",
			Slip:               100,
			LiquidityFeeInRune: 10000,
			FromAddress:        "BNB1",
			ToAddress:          "VAULT",
		},
		testdb.Swap{
			TxID:               "33333",
			Coin:               "10 THOR.RUNE",
			EmitAsset:          "190 BNB/BNB",
			Pool:               "BNB.BNB",
			Slip:               200,
			LiquidityFeeInRune: 20000,
			PriceTarget:        50000,
			FromAddress:        "BNB1",
			ToAddress:          "VAULT",
		},
		testdb.PoolActivate{Pool: "BNB.BNB"},
	)

	body := testdb.CallJSON(t, "http://localhost:8080/v2/actions?limit=50&offset=0&type=swap")

	var v oapigen.ActionsResponse
	testdb.MustUnmarshal(t, body, &v)

	doubleSwap := v.Actions[0]
	metadata := doubleSwap.Metadata.Swap
	in := doubleSwap.In[0]
	out := doubleSwap.Out[0]
	pools := doubleSwap.Pools
	require.Equal(t, "298", metadata.SwapSlip) // 100+200-(100*200)/10000
	require.Equal(t, "30000", metadata.LiquidityFee)
	require.Equal(t, "50000", metadata.SwapTarget)
	require.Equal(t, "bnb1", in.Address)
	require.Equal(t, "BNB.BNB", in.Coins[0].Asset)
	require.Equal(t, "200", in.Coins[0].Amount)
	require.Equal(t, "thor1", out.Address)
	require.Equal(t, "BNB/BNB", out.Coins[0].Asset)
	require.Equal(t, "190", out.Coins[0].Amount)
	require.Equal(t, 1, len(pools))
	require.Equal(t, "BNB.BNB", pools[0])
}

// TestDoubleSwapSynths swaps BNB/BNB -> BTC/BTC
func TestDoubleSwapSynths(t *testing.T) {
	blocks := testdb.InitTestBlocks(t)

	blocks.NewBlock(t, "2020-09-02 00:00:00",
		testdb.Outbound{
			TxID:      "00000",
			InTxID:    "44444",
			Coin:      "190 BTC/BTC",
			ToAddress: "THOR2",
		},
		testdb.Outbound{
			TxID:      "00000",
			InTxID:    "44444",
			Coin:      "10 THOR.RUNE",
			ToAddress: "THOR1",
		},
		testdb.Swap{
			TxID:               "44444",
			Coin:               "200 BNB/BNB",
			EmitAsset:          "10 THOR.RUNE",
			Pool:               "BNB.BNB",
			Slip:               100,
			LiquidityFeeInRune: 10000,
			FromAddress:        "THOR1",
			ToAddress:          "VAULT",
		},
		testdb.Swap{
			TxID:               "44444",
			Coin:               "10 THOR.RUNE",
			EmitAsset:          "190 BTC/BTC",
			Pool:               "BTC.BTC",
			Slip:               200,
			LiquidityFeeInRune: 20000,
			PriceTarget:        50000,
			FromAddress:        "THOR1",
			ToAddress:          "VAULT",
		},
		testdb.PoolActivate{Pool: "BNB.BNB"},
		testdb.PoolActivate{Pool: "BTC.BTC"},
	)

	body := testdb.CallJSON(t, "http://localhost:8080/v2/actions?limit=50&offset=0&type=swap")

	var v oapigen.ActionsResponse
	testdb.MustUnmarshal(t, body, &v)

	doubleSwap := v.Actions[0]
	metadata := doubleSwap.Metadata.Swap
	in := doubleSwap.In[0]
	out := doubleSwap.Out[0]
	pools := doubleSwap.Pools
	require.Equal(t, "298", metadata.SwapSlip) // 100+200-(100*200)/10000
	require.Equal(t, "30000", metadata.LiquidityFee)
	require.Equal(t, "50000", metadata.SwapTarget)
	require.Equal(t, "thor1", in.Address)
	require.Equal(t, "BNB/BNB", in.Coins[0].Asset)
	require.Equal(t, "200", in.Coins[0].Amount)
	require.Equal(t, "thor2", out.Address)
	require.Equal(t, "BTC/BTC", out.Coins[0].Asset)
	require.Equal(t, "190", out.Coins[0].Amount)
	require.Equal(t, 2, len(pools))
	require.Equal(t, "BNB.BNB", pools[0])
	require.Equal(t, "BTC.BTC", pools[1])
}

func TestSwitch(t *testing.T) {
	blocks := testdb.InitTestBlocks(t)

	blocks.NewBlock(t, "2020-09-02 00:00:00",
		testdb.Switch{
			FromAddress: "B2",
			ToAddress:   "THOR2",
			Burn:        "200 BNB.RUNE-B1A",
		})

	blocks.NewBlock(t, "2020-09-03 00:00:00",
		testdb.Switch{
			FromAddress: "A1",
			ToAddress:   "THOR1",
			Burn:        "100 BNB.RUNE-B1A",
			TxID:        "txa1",
		})

	body := testdb.CallJSON(t, "http://localhost:8080/v2/actions?limit=50&offset=0&type=switch")

	var v oapigen.ActionsResponse
	testdb.MustUnmarshal(t, body, &v)

	require.Len(t, v.Actions, 2)

	switch0 := v.Actions[0]
	require.Equal(t, "switch", string(switch0.Type))
	require.Equal(t, "a1", switch0.In[0].Address)
	require.Equal(t, "txa1", switch0.In[0].TxID)
	require.Equal(t, "100", switch0.In[0].Coins[0].Amount)
	require.Equal(t, "thor1", switch0.Out[0].Address)
	require.Equal(t, "THOR.RUNE", switch0.Out[0].Coins[0].Asset)
	require.Equal(t, "100", switch0.Out[0].Coins[0].Amount)

	switch2 := v.Actions[1]
	require.Equal(t, "b2", switch2.In[0].Address)
	require.Equal(t, "200", switch2.In[0].Coins[0].Amount)

	// address filter
	body = testdb.CallJSON(t,
		"http://localhost:8080/v2/actions?limit=50&offset=0&type=switch&address=B2")
	testdb.MustUnmarshal(t, body, &v)
	require.Len(t, v.Actions, 1)
	require.Equal(t, "b2", v.Actions[0].In[0].Address)

	// address filter 2
	body = testdb.CallJSON(t,
		"http://localhost:8080/v2/actions?limit=50&offset=0&type=switch&address=THOR2")
	testdb.MustUnmarshal(t, body, &v)
	require.Len(t, v.Actions, 1)
	require.Equal(t, "b2", v.Actions[0].In[0].Address)
}

func checkFilter(t *testing.T, urlPostfix string, expectedResultsPool []string) {
	body := testdb.CallJSON(t,
		"http://localhost:8080/v2/actions?limit=50&offset=0"+urlPostfix)
	var v oapigen.ActionsResponse
	testdb.MustUnmarshal(t, body, &v)

	require.Equal(t, strconv.Itoa(len(expectedResultsPool)), v.Count)
	for i, pool := range expectedResultsPool {
		require.Equal(t, []string{pool}, v.Actions[i].Pools)
	}
}

func TestAddressFilter(t *testing.T) {
	blocks := testdb.InitTestBlocks(t)

	blocks.NewBlock(t, "2020-09-01 00:00:00",
		testdb.AddLiquidity{
			Pool: "POOL1.A", AssetAmount: 1000, RuneAmount: 2000, AssetAddress: "thoraddr1",
		},
		testdb.PoolActivate{Pool: "POOL1.A"})

	blocks.NewBlock(t, "2020-09-02 00:00:00",
		testdb.Swap{
			Pool:        "POOL2.A",
			Coin:        "20 POOL2.A",
			EmitAsset:   "10 THOR.RUNE",
			FromAddress: "thoraddr2",
			ToAddress:   "thoraddr3",
		},
		testdb.PoolActivate{Pool: "POOL2.A"})

	blocks.NewBlock(t, "2020-09-03 00:00:00",
		testdb.Withdraw{
			Pool:                   "POOL3.A",
			EmitAsset:              10,
			EmitRune:               20,
			ToAddress:              "thoraddr4",
			LiquidityProviderUnits: 1,
		},
		testdb.PoolActivate{Pool: "POOL3.A"})

	checkFilter(t, "", []string{"POOL3.A", "POOL2.A", "POOL1.A"})
	checkFilter(t, "&address=thoraddr1", []string{"POOL1.A"})
	checkFilter(t, "&address=thoraddr2", []string{"POOL2.A"})
	checkFilter(t, "&address=thoraddr4", []string{"POOL3.A"})

	checkFilter(t, "&address=thoraddr1,thoraddr4", []string{"POOL3.A", "POOL1.A"})
}

func TestActionsAddressCaseInsensitive(t *testing.T) {
	blocks := testdb.InitTestBlocks(t)

	blocks.NewBlock(t, "2020-09-01 00:00:00",
		testdb.AddLiquidity{
			Pool: "POOL1.A", AssetAmount: 1000, RuneAmount: 2000,
			AssetAddress: "AddR1",
		},
		testdb.PoolActivate{Pool: "POOL1.A"})

	checkFilter(t, "&address=aDDr1", []string{"POOL1.A"})
}

func TestAddLiquidityAddress(t *testing.T) {
	blocks := testdb.InitTestBlocks(t)

	blocks.NewBlock(t, "2020-09-01 00:00:00",
		testdb.AddLiquidity{
			Pool: "POOL1.A", AssetAmount: 1000, RuneAmount: 2000, AssetAddress: "thoraddr1",
		},
		testdb.PoolActivate{Pool: "POOL1.A"})

	checkFilter(t, "&address=thoraddr1", []string{"POOL1.A"})
}

func TestAddLiquidityUnits(t *testing.T) {
	blocks := testdb.InitTestBlocks(t)

	blocks.NewBlock(t, "2020-09-01 00:00:00",
		testdb.AddLiquidity{
			Pool: "POOL1.A", LiquidityProviderUnits: 42, RuneAmount: 2000, RuneTxID: "tx1",
		},
		testdb.PoolActivate{Pool: "POOL1.A"})

	body := testdb.CallJSON(t, "http://localhost:8080/v2/actions?limit=50&offset=0")

	var v oapigen.ActionsResponse
	testdb.MustUnmarshal(t, body, &v)

	addLiquidity := v.Actions[0]
	metadata := addLiquidity.Metadata.AddLiquidity
	require.Equal(t, "42", metadata.LiquidityUnits)
}

func TestAddLiquidityFields(t *testing.T) {
	blocks := testdb.InitTestBlocks(t)

	blocks.NewBlock(t, "2020-09-01 00:00:00",
		testdb.AddLiquidity{
			Pool:                   "BTC.BTC",
			LiquidityProviderUnits: 42,
			RuneAmount:             2000,
			AssetAmount:            1000,
			RuneTxID:               "rune_tx1",
			RuneAddress:            "runeaddr",
			AssetTxID:              "asset_tx1",
			AssetAddress:           "assetaddr",
		},
		testdb.PoolActivate{Pool: "BTC.BTC"})

	body := testdb.CallJSON(t, "http://localhost:8080/v2/actions?limit=50&offset=0")

	var v oapigen.ActionsResponse
	testdb.MustUnmarshal(t, body, &v)

	require.Equal(t, 1, len(v.Actions))
	action := v.Actions[0]
	require.Equal(t, "addLiquidity", string(action.Type))
	require.Equal(t, "success", string(action.Status))

	require.Equal(t, "1", action.Height)

	require.Equal(t, []oapigen.Transaction{
		{
			Address: "runeaddr",
			Coins:   []oapigen.Coin{{Amount: "2000", Asset: "THOR.RUNE"}},
			TxID:    "rune_tx1",
		},
		{
			Address: "assetaddr",
			Coins:   []oapigen.Coin{{Amount: "1000", Asset: "BTC.BTC"}},
			TxID:    "asset_tx1",
		},
	}, action.In)

	metadata := action.Metadata.AddLiquidity
	require.Equal(t, "42", metadata.LiquidityUnits)

	require.Equal(t, 0, len(action.Out))
	require.Equal(t, []string{"BTC.BTC"}, action.Pools)
}

func TestWithdrawFields(t *testing.T) {
	blocks := testdb.InitTestBlocks(t)

	blocks.NewBlock(t, "2020-09-01 00:00:00",
		testdb.AddLiquidity{
			Pool:                   "BTC.BTC",
			LiquidityProviderUnits: 42,
			RuneAmount:             2000,
			AssetAmount:            1000,
			RuneTxID:               "tx1",
			RuneAddress:            "runeaddr",
			AssetAddress:           "assetaddr",
		},
		testdb.PoolActivate{Pool: "BTC.BTC"})
	blocks.NewBlock(t, "2020-09-01 00:00:05",
		testdb.Withdraw{
			ID:                     "12345",
			Pool:                   "BTC.BTC",
			Coin:                   "1 THOR.RUNE",
			LiquidityProviderUnits: 4,
			EmitRune:               200,
			EmitAsset:              100,
			ImpLossProtection:      3,
			FromAddress:            "runeaddr",
			ToAddress:              "oficialaddr",
			Assymetry:              "0.042",
			BasisPoints:            1000,
		},
		testdb.Fee{
			TxID:       "12345",
			Coins:      "10 BTC.BTC",
			PoolDeduct: 20,
		},
		testdb.Fee{
			TxID:       "12345",
			Coins:      "2 THOR.RUNE",
			PoolDeduct: 2,
		},
	)
	blocks.NewBlock(t, "2020-09-01 00:00:10",
		testdb.Outbound{
			TxID:      "99999",
			InTxID:    "12345",
			Coin:      "90 BTC.BTC",
			ToAddress: "assetaddr",
		},
		testdb.Outbound{
			TxID:      "",
			InTxID:    "12345",
			Coin:      "198 THOR.RUNE",
			ToAddress: "runeaddr",
		},
	)
	body := testdb.CallJSON(t, "http://localhost:8080/v2/actions?type=withdraw")

	var v oapigen.ActionsResponse
	testdb.MustUnmarshal(t, body, &v)

	require.Equal(t, 1, len(v.Actions))
	action := v.Actions[0]
	require.Equal(t, "withdraw", string(action.Type))
	require.Equal(t, "success", string(action.Status))

	require.Equal(t, "2", action.Height)

	require.Equal(t, 1, len(action.In))
	in := action.In[0]
	require.Equal(t, "runeaddr", in.Address)
	require.Equal(t, "12345", in.TxID)
	require.Equal(t, 1, len(in.Coins))
	coin := in.Coins[0]
	require.Equal(t, "1", coin.Amount)
	require.Equal(t, "THOR.RUNE", coin.Asset)

	metadata := action.Metadata.Withdraw
	require.Equal(t, "0.042", metadata.Asymmetry)
	require.Equal(t, "1000", metadata.BasisPoints)
	require.Equal(t, "3", metadata.ImpermanentLossProtection)
	require.Equal(t, "-4", metadata.LiquidityUnits)

	require.Equal(t, oapigen.NetworkFees{
		oapigen.Coin{
			Amount: "10",
			Asset:  "BTC.BTC",
		},
		oapigen.Coin{
			Amount: "2",
			Asset:  "THOR.RUNE",
		},
	}, metadata.NetworkFees)

	require.Equal(t, []oapigen.Transaction{
		{
			Address: "assetaddr",
			TxID:    "99999",
			Coins:   oapigen.Coins{{Amount: "90", Asset: "BTC.BTC"}},
		},
		{
			Address: "runeaddr",
			TxID:    "",
			Coins:   oapigen.Coins{{Amount: "198", Asset: "THOR.RUNE"}},
		},
	}, action.Out)

	require.Equal(t, []string{"BTC.BTC"}, action.Pools)
}

func TestDonateFields(t *testing.T) {
	blocks := testdb.InitTestBlocks(t)

	blocks.NewBlock(t, "2020-09-01 00:00:00",
		testdb.AddLiquidity{
			Pool:        "BTC.BTC",
			RuneAmount:  2000,
			AssetAmount: 1000,
		},
		testdb.PoolActivate{Pool: "BTC.BTC"})
	blocks.NewBlock(t, "2020-09-01 00:00:05",
		testdb.Donate{
			TxID:        "999",
			Pool:        "BTC.BTC",
			Coin:        "12345 THOR.RUNE",
			FromAddress: "runeaddr",
			ToAddress:   "oficialaddr",
			Memo:        "TODOREMOVE",
		},
	)
	body := testdb.CallJSON(t, "http://localhost:8080/v2/actions?type=donate")

	var v oapigen.ActionsResponse
	testdb.MustUnmarshal(t, body, &v)

	require.Equal(t, []oapigen.Action{{
		Date:   util.IntStr(testdb.StrToSec("2020-09-01 00:00:05").ToNano().ToI()),
		Height: "2",
		In: []oapigen.Transaction{{
			Address: "runeaddr",
			Coins:   []oapigen.Coin{{Amount: "12345", Asset: "THOR.RUNE"}},
			TxID:    "999",
		}},
		Out:    []oapigen.Transaction{},
		Pools:  []string{"BTC.BTC"},
		Status: "success",
		Type:   "donate",
	}}, v.Actions)
}

func TestRefundFields(t *testing.T) {
	blocks := testdb.InitTestBlocks(t)

	blocks.NewBlock(t, "2020-09-01 00:00:00",
		testdb.AddLiquidity{
			Pool:                   "BTC.BTC",
			LiquidityProviderUnits: 42,
			RuneAmount:             2000,
			AssetAmount:            1000,
			RuneTxID:               "tx1",
			RuneAddress:            "runeaddr",
			AssetAddress:           "assetaddr",
		},
		testdb.PoolActivate{Pool: "BTC.BTC"})
	blocks.NewBlock(t, "2020-09-01 00:00:05",
		testdb.Refund{
			TxID:        "12345",
			Coin:        "1000 BTC.BTC",
			FromAddress: "userassteaddr",
			ToAddress:   "officialaddr",
			Reason:      "emit asset 100 less than price limit 200",
		},
		testdb.Fee{
			TxID:       "12345",
			Coins:      "10 BTC.BTC",
			PoolDeduct: 20,
		},
	)
	blocks.NewBlock(t, "2020-09-01 00:00:10",
		testdb.Outbound{
			TxID:        "99999",
			InTxID:      "12345",
			Coin:        "990 BTC.BTC",
			FromAddress: "officialaddr",
			ToAddress:   "userassteaddr",
		},
	)
	body := testdb.CallJSON(t, "http://localhost:8080/v2/actions?type=refund")

	var v oapigen.ActionsResponse
	testdb.MustUnmarshal(t, body, &v)
	require.Equal(t, []oapigen.Action{{
		Date:   util.IntStr(testdb.StrToSec("2020-09-01 00:00:05").ToNano().ToI()),
		Height: "2",
		In: []oapigen.Transaction{{
			Address: "userassteaddr",
			Coins:   []oapigen.Coin{{Amount: "1000", Asset: "BTC.BTC"}},
			TxID:    "12345",
		}},
		Metadata: oapigen.Metadata{
			Refund: &oapigen.RefundMetadata{
				NetworkFees: []oapigen.Coin{{Amount: "10", Asset: "BTC.BTC"}},
				Reason:      "emit asset 100 less than price limit 200",
			},
		},
		Out: []oapigen.Transaction{{
			Address: "userassteaddr",
			Coins:   []oapigen.Coin{{Amount: "990", Asset: "BTC.BTC"}},
			TxID:    "99999",
		}},
		Pools:  []string{},
		Status: "success",
		Type:   "refund",
	}}, v.Actions)
}

func TestSwitchFields(t *testing.T) {
	blocks := testdb.InitTestBlocks(t)

	blocks.NewBlock(t, "2020-09-01 00:00:00",
		testdb.Switch{
			TxID:        "12345",
			FromAddress: "bnbaddr",
			ToAddress:   "thoraddr",
			Burn:        "42 THOR.RUNE",
		})

	body := testdb.CallJSON(t, "http://localhost:8080/v2/actions?type=switch")

	var v oapigen.ActionsResponse
	testdb.MustUnmarshal(t, body, &v)
	require.Equal(t, []oapigen.Action{{
		Date:   util.IntStr(testdb.StrToSec("2020-09-01 00:00:00").ToNano().ToI()),
		Height: "1",
		In: []oapigen.Transaction{{
			Address: "bnbaddr",
			Coins:   []oapigen.Coin{{Amount: "42", Asset: "THOR.RUNE"}},
			TxID:    "12345",
		}},
		Out: []oapigen.Transaction{{
			Address: "thoraddr",
			Coins:   []oapigen.Coin{{Amount: "42", Asset: "THOR.RUNE"}},
			TxID:    "",
		}},
		Pools:  []string{},
		Status: "success",
		Type:   "switch",
	}}, v.Actions)
}
