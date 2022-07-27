package pymath

//go test -v comparison_test.go comparison.go

import (
	"testing"

	"github.com/go-check/check"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type AssertSuite struct{}

func init() {
	check.Suite(&AssertSuite{})
}

func (suite *AssertSuite) TestMaxer(c *check.C) {
	c.Assert(Maxer(1, 4), check.Not(check.Equals), 2)
	c.Assert(Maxer(1, 4), check.Equals, int64(2))
	c.Assert(Maxer(1, 1), check.Equals, int64(1))
	c.Assert(Maxer(4, 2), check.Equals, int64(3))
}

func (suite *AssertSuite) TestMiner(c *check.C) {
	c.Assert(Miner(1, 4), check.Not(check.Equals), 1)
	c.Assert(Miner(1, 4), check.Equals, int64(1))
	c.Assert(Miner(1, 1), check.Equals, int64(1))
	c.Assert(Miner(4, 2), check.Equals, int64(2))
}
