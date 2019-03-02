package godash_test

import (
	"testing"

	. "gopkg.in/check.v1"

	"github.com/mvo5/godash"
)

func Test(t *testing.T) { TestingT(t) }

type DashSuite struct{}

var _ = Suite(&DashSuite{})

func (s *DashSuite) TestTrivial(c *C) {
	dash := godash.New()
	c.Assert(dash, FitsTypeOf, &godash.Dash{})
}
