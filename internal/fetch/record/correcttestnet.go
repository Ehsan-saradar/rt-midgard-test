package record

import "github.com/rs/zerolog/log"

// Testnet started on 2021-11-06
const ChainIDTestnet20211106 = "D4DF73AD98535DCD72BD0C9FE76B96CAF350C2FF517A61F77F5F89665A0593E7"

// ThorNode state and events diverged on testnet. We apply all these changes to be in sync with
// Thornode.
func loadTestnet202111Corrections(chainID string) {
	if chainID == ChainIDTestnet20211106 {
		log.Info().Msgf(
			"Loading corrections for testnet started on 2021-11-06 id: %s",
			chainID)

		loadTestnetMissingWithdraws()
	}
}

//////////////////////// Missing withdraws

func loadTestnetMissingWithdraws() {
	// On Pool suspension the withdraws had FromAddr=null and they were skipped by Midgard.
	// Later the pool was reactivated, so having correct units is important even at suspension.
	// There is a plan to fix ThorNode events:
	// https://gitlab.com/thorchain/thornode/-/issues/1164
	addWithdraw(10000, AdditionalWithdraw{
		Pool:     "BNB.BUSD-74E",
		FromAddr: "tthor1qkd5f9xh2g87wmjc620uf5w08ygdx4etu0u9fs",
		Reason:   "midgard correction suspended pool withdraws missing",
		RuneE8:   0,
		AssetE8:  0,
		Units:    10000000000,
	})

	addWithdraw(222784, AdditionalWithdraw{
		Pool:     "BNB.BNB",
		FromAddr: "tbnb1yc20slera2g4fhnkkyttqxf70qxa4jtm42qq4t",
		Reason:   "midgard correction",
		RuneE8:   294194696841,
		AssetE8:  106918851,
		Units:    170138465261,
	})
}
