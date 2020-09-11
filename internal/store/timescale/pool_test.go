package timescale

import (
	"gitlab.com/thorchain/midgard/internal/common"
	"gitlab.com/thorchain/midgard/internal/models"
	. "gopkg.in/check.v1"
)

func (s *TimeScaleSuite) TestPool(c *C) {
	assetBolt, _ := common.NewAsset("BOLT-014")
	assetTcan, _ := common.NewAsset("TCAN-014")

	// No pool status (default value)
	poolStatus, err := s.Store.GetPoolStatus(assetBolt)
	c.Assert(err, IsNil)
	c.Assert(poolStatus, Equals, models.Unknown)

	// First pool status
	err = s.Store.CreatePoolRecord(&poolStatusEvent0)
	c.Assert(err, IsNil)
	poolStatus, err = s.Store.GetPoolStatus(assetBolt)
	c.Assert(err, IsNil)
	c.Assert(poolStatus, Equals, models.Bootstrap)

	// Unchanged pool status
	poolStatus, err = s.Store.GetPoolStatus(assetTcan)
	c.Assert(err, IsNil)
	c.Assert(poolStatus, Equals, models.Unknown)

	// Second pool status
	err = s.Store.CreatePoolRecord(&poolStatusEvent1)
	c.Assert(err, IsNil)
	poolStatus, err = s.Store.GetPoolStatus(assetBolt)
	c.Assert(err, IsNil)
	c.Assert(poolStatus, Equals, models.Enabled)
}
