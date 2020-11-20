package timeseries_test

import (
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/assert"
	"gitlab.com/thorchain/midgard/chain/notinchain"
	"gitlab.com/thorchain/midgard/internal/graphql"
	"gitlab.com/thorchain/midgard/internal/graphql/generated"
	"gitlab.com/thorchain/midgard/internal/graphql/model"
	"gitlab.com/thorchain/midgard/openapi/generated/oapigen"
	"testing"

	"github.com/jarcoal/httpmock"
	"gitlab.com/thorchain/midgard/internal/timeseries"
	"gitlab.com/thorchain/midgard/internal/timeseries/testdb"
)

func TestNetwork(t *testing.T) {
	testdb.SetupTestDB(t)
	schema := generated.NewExecutableSchema(generated.Config{Resolvers: &graphql.Resolver{}})
	gqlClient := client.New(handler.NewDefaultServer(schema))

	testdb.MustExec(t, "DELETE FROM stake_events")
	testdb.MustExec(t, "DELETE FROM block_log")
	testdb.MustExec(t, "DELETE FROM swap_events")
	testdb.MustExec(t, "DELETE FROM block_pool_depths")
	testdb.MustExec(t, "DELETE FROM active_vault_events")
	testdb.MustExec(t, "DELETE FROM set_mimir_events")

	timeseries.SetLastTimeForTest(testdb.ToTime("2020-09-01 00:00:00"))
	timeseries.SetDepthsForTest("BNB.TWT-123", 30000000000000, 2240582804123679)
	testdb.InsertActiveVaultEvent(t, "addr", "2020-09-01 00:00:00")
	testdb.SetThornodeConstants(t)

	// Setting number of bonds, nodes  and totalReserve in the mocked ThorNode
	nodeAccounts := make([]notinchain.NodeAccount, 2)
	nodeAccounts[0] = notinchain.NodeAccount{
		Status: "active",
		Bond:   22772603677970,
	}
	nodeAccounts[1] = notinchain.NodeAccount{
		Status: "standby",
		Bond:   9999990,
	}
	testdb.MockThorNode(108915513107, nodeAccounts)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testdb.InsertBlockLog(t, 1, "2020-09-01 00:00:00")
	testdb.InsertSwapEvent(t, testdb.FakeSwap{Pool: "BNB.TWT-123", FromAsset: "BNB.RUNE", FromE8: 300000000, LiqFeeInRuneE8: 3908720129799, BlockTimestamp: "2020-09-01 00:00:00"})
	testdb.InsertStakeEvent(t, testdb.FakeStake{Pool: "BNB.TWT-123", BlockTimestamp: "2020-09-01 00:00:00"})

	body := testdb.CallV1(t, "http://localhost:8080/v2/network")

	var jsonApiResult oapigen.Network
	testdb.MustUnmarshal(t, body, &jsonApiResult)

	queryString := `{
		network {
			activeBonds,
			activeNodeCount
			standbyBonds
			standbyNodeCount
			bondMetrics {
				active {
					averageBond
					totalBond
					medianBond
					maximumBond
				}
				standby {
					averageBond
					totalBond
					medianBond
					maximumBond
				}
			}
			blockRewards {
				blockReward
				bondReward
				poolReward
			}
			liquidityAPY
			bondingAPY
			nextChurnHeight
			poolActivationCountdown
			poolShareFactor
			totalReserve
			totalPooledRune
		}
	}`

	type Result struct {
		Network model.Network
	}
	var graphqlResult Result
	gqlClient.MustPost(queryString, &graphqlResult)

	// specified in ThorNode
	assert.Equal(t, "1", jsonApiResult.ActiveNodeCount)
	assert.Equal(t, int64(1), graphqlResult.Network.ActiveNodeCount)
	assert.Equal(t, "1", jsonApiResult.StandbyNodeCount)
	assert.Equal(t, int64(1), graphqlResult.Network.StandbyNodeCount)
	assert.Equal(t, "22772603677970", jsonApiResult.BondMetrics.TotalActiveBond)
	assert.Equal(t, int64(22772603677970), graphqlResult.Network.BondMetrics.Active.TotalBond)
	assert.Equal(t, "9999990", jsonApiResult.BondMetrics.TotalStandbyBond)
	assert.Equal(t, int64(9999990), graphqlResult.Network.BondMetrics.Standby.TotalBond)
	assert.Equal(t, "108915513107", jsonApiResult.TotalReserve)
	assert.Equal(t, int64(108915513107), graphqlResult.Network.TotalReserve)

	assert.Equal(t, "17256", jsonApiResult.BlockRewards.BlockReward)
	assert.Equal(t, int64(17256), graphqlResult.Network.BlockRewards.BlockReward)

	assert.Equal(t, "0", jsonApiResult.LiquidityAPY)
	assert.Equal(t, 0.0, graphqlResult.Network.LiquidityApy)
	assert.Equal(t, "3879.8255319373584", jsonApiResult.BondingAPY)
	assert.Equal(t, 3879.8255319373584, graphqlResult.Network.BondingApy)
	assert.Equal(t, "2161", jsonApiResult.NextChurnHeight)
	assert.Equal(t, int64(2161), graphqlResult.Network.NextChurnHeight)
	assert.Equal(t, "49999", jsonApiResult.PoolActivationCountdown)
	assert.Equal(t, int64(49999), graphqlResult.Network.PoolActivationCountdown)
	assert.Equal(t, "0", jsonApiResult.PoolShareFactor)
	assert.Equal(t, 0.0, graphqlResult.Network.PoolShareFactor)
	assert.Equal(t, "108915513107", jsonApiResult.TotalReserve)
	assert.Equal(t, int64(108915513107), graphqlResult.Network.TotalReserve)
	assert.Equal(t, "2240582804123679", jsonApiResult.TotalPooledRune)
	assert.Equal(t, int64(2240582804123679), graphqlResult.Network.TotalPooledRune)
}