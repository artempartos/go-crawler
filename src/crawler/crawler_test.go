package crawler

import (
	. "launchpad.net/gocheck"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestGrabber(c *C) {
	//ParseLink("http://insoftretail.ru")
	//c.Assert(crawler.ParseLink('https://www.google.ru'), Equals, "42")
	//c.Assert(crawler.ParseLink('https://www.google.ru'), Equals, "42")
	//c.Assert(42, Equals, "42")
	//c.Check(os.Errno(13), Matches, "perm.*accepted")
}
