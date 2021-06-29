// This tool checks if pool units reported reported for each member sum by ThorNode sums up to the
// total units of the pool.
// At the time of writing it does.
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"gitlab.com/thorchain/midgard/config"
	"gitlab.com/thorchain/midgard/internal/api"
	"gitlab.com/thorchain/midgard/internal/db"
	"gitlab.com/thorchain/midgard/internal/fetch/record"
	"gitlab.com/thorchain/midgard/internal/timeseries"
)

const usageStr = `Check pool units share of each member
Usage:
$ go run ./cmd/membercheck config pool heightOrBlockTimestamp
or
$ go run ./cmd/membercheck --allpools  config heightOrBlockTimestamp
`

func init() {
	flag.Usage = func() {
		fmt.Println(usageStr)
		flag.PrintDefaults()
	}
}

var AllPoolsStructured = flag.Bool("allpools", false,
	"No binary search, only the latest depth differences in structured form.")

type ThorNodeSummary struct {
	TotalUnits int64 `json:"pool_units,string"`
}

// MemberChange may represent the state of a Member, an Add or a Withdraw.
// In case of withdraw the units are negative.
type MemberChange struct {
	Units        int64  `json:"units,string"`
	RuneAddress  string `json:"rune_address"`
	AssetAddress string `json:"asset_address"`
}

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true})
	logrus.SetLevel(logrus.InfoLevel)

	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Println("Not enough arguments!")
		flag.Usage()
		return
	}

	var c config.Config = config.ReadConfigFrom(flag.Arg(0))

	db.Setup(&c.TimeScale)

	db.LoadFirstBlockFromDB(context.Background())

	if *AllPoolsStructured {
		CheckAllPoolsStructured(c)
	} else {
		CheckOnePool(c)
	}
}

func findHeight(param string) (height int64, timestamp db.Nano) {
	idStr := param

	heightOrTimestamp, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logrus.Fatal("Couldn't parse height or timestamp: ", idStr)
	}
	height, timestamp, err = api.TimestampAndHeight(context.Background(), heightOrTimestamp)
	if err != nil {
		logrus.Fatal("Couldn't find height or timestamp. ", err)
	}
	return
}

func CheckOnePool(c config.Config) {
	if flag.NArg() != 3 {
		fmt.Println("provide 3 args!")
		flag.Usage()
		return
	}

	ctx := context.Background()
	pool := flag.Arg(1)
	idStr := flag.Arg(2)

	height, timestamp := findHeight(idStr)
	thorNodeMembers := getThorNodeMembers(c, pool, height)
	logrus.Debug("Thornode rune addresses: ", len(thorNodeMembers.RuneMemberUnits),
		" assetOnly addresses: ", len(thorNodeMembers.AssetMemberUnits),
		" assetToRuneMap: ", len(thorNodeMembers.AssetToRuneMap))

	midgardMembers := getMidgardMembers(ctx, pool, timestamp)
	logrus.Debug("Midgard rune addresses: ", len(midgardMembers.RuneMemberUnits),
		" assetOnly addresses: ", len(midgardMembers.AssetMemberUnits),
		" assetToRuneMap: ", len(midgardMembers.AssetToRuneMap))

	memberDiff(thorNodeMembers, midgardMembers)
}

func CheckAllPoolsStructured(c config.Config) {
	if flag.NArg() != 2 {
		fmt.Println("provide 2 args!")
		flag.Usage()
		return
	}

	ctx := context.Background()
	idStr := flag.Arg(1)

	height, timestamp := findHeight(idStr)

	poolsWithStatus, err := timeseries.GetPoolsStatuses(ctx, timestamp)
	if err != nil {
		logrus.Fatal(err)
	}
	sortedPools := []string{}
	for k := range poolsWithStatus {
		sortedPools = append(sortedPools, k)
	}
	sort.Strings(sortedPools)
	for _, pool := range sortedPools {
		status := poolsWithStatus[pool]
		if status == "suspended" {
			continue
		}
		thorNodeMembers := getThorNodeMembers(c, pool, height)
		logrus.Debug("Thornode rune addresses: ", len(thorNodeMembers.RuneMemberUnits),
			" assetOnly addresses: ", len(thorNodeMembers.AssetMemberUnits),
			" assetToRuneMap: ", len(thorNodeMembers.AssetToRuneMap))

		midgardMembers := getMidgardMembers(ctx, pool, timestamp)
		logrus.Debug("Midgard rune addresses: ", len(midgardMembers.RuneMemberUnits),
			" assetOnly addresses: ", len(midgardMembers.AssetMemberUnits),
			" assetToRuneMap: ", len(midgardMembers.AssetToRuneMap))

		saveStructuredDiffs(pool, thorNodeMembers, midgardMembers)
	}
	printStructuredDiffs()
}

type MemberMap struct {
	// RuneMemberUnits is keyed by rune address if the member has one, otherwise asset address
	RuneMemberUnits map[string]int64

	// Only if it doesn't have a rune address
	AssetMemberUnits map[string]int64

	// Disjoint with AssetMemberUnits
	AssetToRuneMap map[string]string
}

func NewMemberMap() MemberMap {
	return MemberMap{
		RuneMemberUnits:  map[string]int64{},
		AssetMemberUnits: map[string]int64{},
		AssetToRuneMap:   map[string]string{},
	}
}

func check(ok bool, v ...interface{}) {
	if !ok {
		log.Fatal(v...)
	}
}

func (x *MemberMap) AddMemberSimple(m MemberChange) {
	rAddr := m.RuneAddress
	aAddr := m.AssetAddress

	if rAddr != "" {
		x.RuneMemberUnits[rAddr] += m.Units
	} else {
		check(aAddr != "", "Empty rune and asset address")
		x.AssetMemberUnits[aAddr] += m.Units
	}
}

// If there is an asset address without rune address
// then it looks for previous usage of the asset address to find an adequate rune adresses.
func (x *MemberMap) AddMemberClustered(m MemberChange) {
	rAddr := m.RuneAddress
	aAddr := m.AssetAddress

	if rAddr != "" {
		x.RuneMemberUnits[rAddr] += m.Units

		if aAddr != "" {
			assetUnits, previouslyAssetOnly := x.AssetMemberUnits[aAddr]
			if previouslyAssetOnly {
				x.RuneMemberUnits[rAddr] += assetUnits
				delete(x.AssetMemberUnits, aAddr)
			}

			previousRuneAddr, assetAddrAlreadyRegistered := x.AssetToRuneMap[aAddr]
			if assetAddrAlreadyRegistered {
				if previousRuneAddr != rAddr {
					logrus.Fatal("AssetAddress registered with multiple rune addresses",
						rAddr, previousRuneAddr)
				}
			} else {
				x.AssetToRuneMap[aAddr] = m.RuneAddress
			}
		}
	} else {
		check(aAddr != "", "Empty rune and asset address")
		previousRuneAddr, hasRunePair := x.AssetToRuneMap[aAddr]
		if hasRunePair {
			x.RuneMemberUnits[previousRuneAddr] += m.Units
		} else {
			x.AssetMemberUnits[aAddr] += m.Units
		}
	}
}

func (x *MemberMap) RemoveZero() {
	for k, v := range x.RuneMemberUnits {
		if v == 0 {
			delete(x.RuneMemberUnits, k)
		}
	}
	for k, v := range x.AssetMemberUnits {
		if v == 0 {
			delete(x.AssetMemberUnits, k)
		}
	}
}

func mapSum(m map[string]int64) int64 {
	var ret int64
	for _, v := range m {
		ret += v
	}
	return ret
}

func (x *MemberMap) TotalUnits() int64 {
	return mapSum(x.RuneMemberUnits) + mapSum(x.AssetMemberUnits)
}

func getThorNodeMembers(c config.Config, pool string, height int64) MemberMap {
	logrus.Info("Checking pool units sum. Pool: ", pool, " Height: ", height)

	var summary ThorNodeSummary
	queryThorNode(c.ThorChain.ThorNodeURL, "/pool/"+pool, height, &summary)
	logrus.Info("ThorNode global units: ", summary.TotalUnits)

	var thornodeBreakdown []MemberChange
	queryThorNode(c.ThorChain.ThorNodeURL, "/pool/"+pool+"/liquidity_providers", height, &thornodeBreakdown)

	ret := NewMemberMap()

	var sum2 int64
	for _, member := range thornodeBreakdown {
		sum2 += member.Units
		ret.AddMemberSimple(member)
	}
	logrus.Info("ThorNode units per member summed up: ", sum2)
	if sum2 == summary.TotalUnits {
		logrus.Info("thornode is consistent")
	} else {
		logrus.Fatal("thornode INCONSISTENT")
	}

	ret.RemoveZero()
	return ret
}

func queryThorNode(thorNodeUrl string, urlPath string, height int64, dest interface{}) {
	url := thorNodeUrl + urlPath
	if 0 < height {
		url += "?height=" + strconv.FormatInt(height, 10)
	}
	logrus.Debug("Querying thornode: ", url)
	resp, err := http.Get(url)
	if err != nil {
		logrus.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, dest)
	if err != nil {
		logrus.Fatal(err)
	}
}

func getMidgardMembers(ctx context.Context, pool string, timestamp db.Nano) MemberMap {
	ret := NewMemberMap()

	addQ := `
		SELECT rune_addr, asset_addr, stake_units
		FROM stake_events
		WHERE pool = $1 and block_timestamp <= $2
		ORDER BY block_timestamp
	`
	addRows, err := db.Query(ctx, addQ, pool, timestamp)
	if err != nil {
		logrus.Fatal(err)
	}
	defer addRows.Close()

	for addRows.Next() {
		var runeAddress, assetAddress sql.NullString
		var add MemberChange
		err := addRows.Scan(
			&runeAddress,
			&assetAddress,
			&add.Units)
		if err != nil {
			logrus.Fatal(err)
		}
		if runeAddress.Valid {
			add.RuneAddress = runeAddress.String
		}
		if assetAddress.Valid {
			add.AssetAddress = assetAddress.String
		}
		ret.AddMemberSimple(add)
	}

	withdrawQ := `
		SELECT from_addr, stake_units
		FROM unstake_events
		WHERE pool = $1 and block_timestamp <= $2
		ORDER BY block_timestamp
	`

	withdrawRows, err := db.Query(ctx, withdrawQ, pool, timestamp)
	if err != nil {
		logrus.Fatal(err)
	}
	defer withdrawRows.Close()

	for withdrawRows.Next() {
		var fromAddr string
		var units int64
		err := withdrawRows.Scan(
			&fromAddr,
			&units)
		if err != nil {
			logrus.Fatal(err)
		}
		withdraw := MemberChange{Units: -units}
		if record.AddressIsRune(fromAddr) {
			withdraw.RuneAddress = fromAddr
		} else {
			withdraw.AssetAddress = fromAddr
		}
		ret.AddMemberClustered(withdraw)
	}

	ret.RemoveZero()
	return ret
}

func mapDiff(thorNodeMap map[string]int64, midgardMap map[string]int64) {
	diffCount := 0
	for k, tValue := range thorNodeMap {
		mValue, mOk := midgardMap[k]
		if !mOk {
			logrus.Warn("Missing address in Midgard: ", k, " ThorNode units: ", tValue)
			diffCount++
		} else if mValue != tValue {
			logrus.Warn(
				"Mismatch units for address: ", k, " ThorNode: ", tValue, " Midgard: ", mValue,
				" diff: ", tValue-mValue)
			diffCount++
		}
	}
	for k, mValue := range midgardMap {
		_, tOk := thorNodeMap[k]
		if !tOk {
			logrus.Warn("Extra address in Midgard: ", k, " Midgard units: ", mValue)
			diffCount++
		}
	}
	if diffCount == 0 {
		logrus.Info("No difference")
	}
}

func memberDiff(thorNodeMembers MemberMap, midgardMembers MemberMap) {
	logrus.Info("Checking Rune adresses")
	mapDiff(thorNodeMembers.RuneMemberUnits, midgardMembers.RuneMemberUnits)
	logrus.Info("Checking Asset adresses")
	mapDiff(thorNodeMembers.AssetMemberUnits, midgardMembers.AssetMemberUnits)

	thorNodeUnits := thorNodeMembers.TotalUnits()
	midgardUnits := midgardMembers.TotalUnits()
	if thorNodeUnits != midgardUnits {
		logrus.Warn("Total units mismatch. ThorNode: ", thorNodeUnits, " Midgard: ", midgardUnits)
	} else {
		logrus.Info("Total units are equal")
	}
}

var structuredBuff strings.Builder

func saveStructuredDiffs(pool string, thorNodeMembers MemberMap, midgardMembers MemberMap) {
	diffValue := map[string]int64{}
	accumulate := func(thorNodeMap map[string]int64, midgardMap map[string]int64) {
		for k, tValue := range thorNodeMap {
			mValue, mOk := midgardMap[k]
			if !mOk {
				diffValue[k] = tValue
			} else if mValue != tValue {
				diffValue[k] = tValue - mValue
			}
		}
		for k, mValue := range midgardMap {
			_, tOk := thorNodeMap[k]
			if !tOk {
				diffValue[k] = -mValue
			}
		}
	}
	accumulate(thorNodeMembers.RuneMemberUnits, midgardMembers.RuneMemberUnits)
	accumulate(thorNodeMembers.AssetMemberUnits, midgardMembers.AssetMemberUnits)

	keys := []string{}
	for k := range diffValue {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := diffValue[k]
		fmt.Fprintf(&structuredBuff, `{"%s", "%s", %d},`+"\n", pool, k, v)
	}
}

func printStructuredDiffs() {
	logrus.Info("Needed changes to Midgard:\n", structuredBuff.String())
}
