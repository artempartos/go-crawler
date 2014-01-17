package crawler

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestCorrectNewLink(t *testing.T) {
	testUrl := "http://correct.url"
	link, err := NewLink(testUrl)

	assert.Nil(t, err)
	assert.NotNil(t, link.url)
	assert.Equal(t, testUrl, link.url.String())
}

func TestFailNewLink(t *testing.T) {
	testUrl := "/user/%USER_ID%/votes"
	link, err := NewLink(testUrl)

	assert.NotNil(t, err)
	assert.Nil(t, link)
}

func TestLinkFunc(t *testing.T) {
	testUrl := "/test"
	link, _ := NewLink(testUrl)
	assert.True(t, link.isRelative())

	host, _ := url.Parse("http://example.ru")
	absUrl := link.withHost(host)
	absLink, _ := NewLink(absUrl)
	assert.True(t, absLink.isSameHost(host))
}
