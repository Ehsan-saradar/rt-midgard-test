// In 2021-04 ThorNode had two bug when withdrawing with impermanent loss.
// Bigger pool units were reported for withdraws than the actual remove:
// https://gitlab.com/thorchain/thornode/-/issues/912
//
// Sometimes when withdrawing the pool units of a member went up, not down:
// https://gitlab.com/thorchain/thornode/-/issues/896
//
// This script fetches ThorNode at each point where we had a withdraw with impermanent loss
// and prints out the corrections to those withdraws.
//
// There would be two ways of querying ThorNode for pool units.
// 1) Go to the pool page and read the total units of the pool. Sum up the withdraws and adds
//    form MidgardDb, and if there is a discrepancy with the previous block assume it's because
//    of the withdraw with impermanent loss.
//    https://thornode.thorchain.info/thorchain/pool/BNB.BNB?height=50000
// 2) Go to the liquidity_providers page and check how much the pool units of that individual
//    member changes.
//    https://thornode.thorchain.info/thorchain/pool/BNB.BNB/liquidity_providers?height=81054
//
// Here we do 1, but mannually was checked that 2 would give consistent results.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
	"gitlab.com/thorchain/midgard/config"
	"gitlab.com/thorchain/midgard/internal/db"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true})
	logrus.SetLevel(logrus.DebugLevel)
	config.ReadGlobal()

	ctx := context.Background()

	db.Setup()
	db.SetFirstBlockFromDB(ctx)

	summaries := withdrawsWithImpermanentLoss(ctx)
	logrus.Infof("Withdraws count with impermanent loss protection: %d ", len(summaries))

	type WithdrawCorrection struct {
		TX          string
		ActualUnits int64
	}
	sortedWithdrawKeys := []int64{}
	correctWithdraws := map[int64]WithdrawCorrection{}

	type AddCorrection struct {
		RuneAddr string
		Units    int64
		Pool     string
	}
	adds := map[int64]AddCorrection{}
	sortedAddKeys := []int64{}
	for _, summary := range summaries {
		checkWithdrawsIsAlone(ctx, summary)
		readAdds(ctx, &summary)
		fetchNodeUnits(&summary)
		nodeWithdraw := summary.NodeDiff + summary.Adds
		if nodeWithdraw == summary.MidgardWithrawUnits {
			// No problems with this withdraw
			continue
		}
		if summary.Adds != 0 {
			logrus.Info("Pool has add in the same block")
		}
		diff := float64(nodeWithdraw-summary.MidgardWithrawUnits) / float64(summary.MidgardWithrawUnits)
		logrus.Infof("Diff %f Pool: %s height %d -- [ %d vs %d ]",
			diff, summary.Pool, summary.Height, summary.MidgardWithrawUnits, (nodeWithdraw))
		if -0.2 <= diff && diff <= 0 {
			correctWithdraws[summary.Height] = WithdrawCorrection{
				TX:          summary.TX,
				ActualUnits: nodeWithdraw,
			}
			sortedWithdrawKeys = append(sortedWithdrawKeys, summary.Height)
		} else {
			logrus.Warn("Big impermanent loss change, creating append")
			if -1 <= diff {
				logrus.Fatal("Big impermanent loss but not addition yet")
			}
			adds[summary.Height] = AddCorrection{
				RuneAddr: summary.FromAddr,
				Units:    summary.MidgardWithrawUnits - nodeWithdraw, // nodeWithdraw is negative
				Pool:     summary.Pool,
			}
			sortedAddKeys = append(sortedAddKeys, summary.Height)
		}
	}
	logrus.Info("correct withdraws: ", correctWithdraws)
	logrus.Info("adds: ", adds)
	withdrawString := "var withdrawUnitCorrections = map[int64]withdrawUnitCorrection{\n"

	for _, k := range sortedWithdrawKeys {
		v := correctWithdraws[k]
		withdrawString += fmt.Sprintf(
			"\t%d: {\"%s\", %d},\n", k, v.TX, v.ActualUnits)
	}
	withdrawString += "}\n"
	logrus.Warn("Correct withdraws:\n", withdrawString)

	addString := "var addInsteadWithdrawMap = map[int64]addInsteadWithdraw{\n"
	for _, k := range sortedAddKeys {
		v := adds[k]
		addString += fmt.Sprintf(
			"\t%d: {\"%s\", \"%s\", %d},\n", k, v.Pool, v.RuneAddr, v.Units)
	}
	addString += "}\n"
	logrus.Warn("Additional adds:\n", addString)
}

type UnitsSummary struct {
	Pool                string
	TX                  string
	MidgardWithrawUnits int64
	FromAddr            string
	Adds                int64
	NodeDiff            int64
	Timestamp           db.Nano
	Height              int64
}

func withdrawsWithImpermanentLoss(ctx context.Context) []UnitsSummary {
	ret := []UnitsSummary{}
	q := `
		SELECT
			x.pool,
			x.tx,
			x.stake_units,
			x.from_addr,
			x.block_timestamp,
			b.height
		FROM unstake_events AS x JOIN block_log AS b ON x.block_timestamp = b.timestamp
		WHERE imp_loss_protection_e8 <> 0
		ORDER BY height
	`
	rows, err := db.Query(ctx, q)
	if err != nil {
		logrus.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var w UnitsSummary
		err := rows.Scan(
			&w.Pool, &w.TX, &w.MidgardWithrawUnits, &w.FromAddr, &w.Timestamp, &w.Height)
		if err != nil {
			logrus.Fatal(err)
		}
		ret = append(ret, w)
	}
	return ret
}

func checkWithdrawsIsAlone(ctx context.Context, summary UnitsSummary) {
	q := `select count(*) from unstake_events where block_timestamp = $1`

	rows, err := db.Query(ctx, q, summary.Timestamp)
	if err != nil {
		logrus.Fatal(err)
	}
	defer rows.Close()

	if !rows.Next() {
		logrus.Fatal("Expected one row from count")
	}
	var count int
	err = rows.Scan(&count)
	if err != nil {
		logrus.Fatal(err)
	}
	if count != 1 {
		logrus.Fatalf(
			"Multiple withdraws at timestamp %d, height %d",
			summary.Timestamp, summary.Height)
	}
}

func readAdds(ctx context.Context, summary *UnitsSummary) {
	q := `
		SELECT COALESCE(SUM(stake_units), 0)
		FROM stake_events
		WHERE pool=$1 AND block_timestamp = $2`

	rows, err := db.Query(ctx, q, summary.Pool, summary.Timestamp)
	if err != nil {
		logrus.Fatal(err)
	}
	defer rows.Close()

	if !rows.Next() {
		logrus.Fatal("Expected one row adds")
	}
	err = rows.Scan(&summary.Adds)
	logrus.Debug("Add: ", summary.Adds)
	if err != nil {
		logrus.Fatal(err)
	}
}

type ThorNodeUnits struct {
	TotalUnits int64 `json:"pool_units,string"`
}

func NodeUnits(thorNodeUrl string, urlPath string, height int64) int64 {
	url := thorNodeUrl + urlPath + "?height=" + strconv.FormatInt(height, 10)
	logrus.Debug("Querying thornode: ", url)
	resp, err := http.Get(url)
	if err != nil {
		logrus.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var result ThorNodeUnits
	err = json.Unmarshal(body, &result)
	if err != nil {
		logrus.Fatal(err)
	}
	return result.TotalUnits
}

func fetchNodeUnits(summary *UnitsSummary) {
	thorNodeUrl := config.Global.ThorChain.ThorNodeURL
	unitsBefore := NodeUnits(thorNodeUrl, "/pool/"+summary.Pool, summary.Height-1)
	unitsAfter := NodeUnits(thorNodeUrl, "/pool/"+summary.Pool, summary.Height)
	summary.NodeDiff = unitsBefore - unitsAfter
	logrus.Debug("before ", unitsBefore, " after ", unitsAfter, " diff ", summary.NodeDiff)
}
