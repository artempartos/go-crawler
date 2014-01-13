package crawler

import (
	. "launchpad.net/gocheck"
	"net/url"
)

func (s *MySuite) TestCorrectNewLink(c *C) {
	testUrl := "http://correct.url"
	link, err := NewLink(testUrl)

	c.Assert(err, Equals, nil)
	c.Assert(link.url, NotNil)
	c.Assert(testUrl, Equals, link.url.String())
}

func (s *MySuite) TestFailNewLink(c *C) {
	testUrl := "/user/%USER_ID%/votes"
	link, err := NewLink(testUrl)

	c.Assert(err, NotNil)
	c.Assert(link, IsNil)
}

func (s *MySuite) TestLinkFunc(c *C) {
	testUrl := "/test"
	link, _ := NewLink(testUrl)
	c.Assert(link.isRelative(), Equals, true)

	host, _ := url.Parse("http://example.ru")
	absUrl := link.withHost(host)
	absLink, _ := NewLink(absUrl)
	c.Assert(absLink.isSameHost(host), Equals, true)
}
