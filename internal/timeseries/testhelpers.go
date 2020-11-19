// Exports functionality for tests only.
// Because the filename has _test, it won't be included in production builds.

package timeseries

import (
	"time"
)

func copyMap(m map[string]int64) map[string]int64 {
	m2 := map[string]int64{}
	for i, v := range m {
		m2[i] = v
	}
	return m2
}

func copyOfLastTrack() (ret *blockTrack) {
	height := int64(1)
	t := time.Unix(0, 0)
	hash := []byte("hash0")
	assetDepth := map[string]int64{}
	runeDepth := map[string]int64{}
	interfacePtr := lastBlockTrack.Load()
	if interfacePtr != nil {
		oldTrack := interfacePtr.(*blockTrack)
		height = oldTrack.Height
		t = oldTrack.Timestamp
		hash = oldTrack.Hash
		assetDepth = copyMap(oldTrack.aggTrack.AssetE8DepthPerPool)
		runeDepth = copyMap(oldTrack.aggTrack.RuneE8DepthPerPool)
	}
	return &blockTrack{
		Height:    height,
		Timestamp: t,
		Hash:      hash,
		aggTrack: aggTrack{
			AssetE8DepthPerPool: assetDepth,
			RuneE8DepthPerPool:  runeDepth,
		},
	}
}

// Often current height or timestamp is read from the last track, this function helps
// to set them for tests.
func SetLastTimeForTest(timestamp time.Time) {
	trackPtr := copyOfLastTrack()
	trackPtr.Timestamp = timestamp
	lastBlockTrack.Store(trackPtr)
}

func SetDepthsForTest(pool string, assetDepth, runeDepth int64) {
	resetAggTrack()
	trackPtr := copyOfLastTrack()
	trackPtr.aggTrack.AssetE8DepthPerPool[pool] = assetDepth
	trackPtr.aggTrack.RuneE8DepthPerPool[pool] = runeDepth
	lastBlockTrack.Store(trackPtr)
}

func resetAggTrack() {
	trackPtr := copyOfLastTrack()
	trackPtr.aggTrack = aggTrack{
		AssetE8DepthPerPool: make(map[string]int64),
		RuneE8DepthPerPool:  make(map[string]int64),
	}
	lastBlockTrack.Store(trackPtr)
}
