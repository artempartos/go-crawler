package crawler

import "net/url"

type Link struct {
	url *url.URL
}

func NewLink(str string) (*Link, error) {
	u, err := url.Parse(str)
	if err != nil {
		return nil, err
	}
	return &Link{u}, nil
}

func (l *Link) isSameHost(host *url.URL) bool {
	return l.url.IsAbs() && l.url.Host == host.Host
}

func (l *Link) isRelative() bool {
	return !l.url.IsAbs()
}

func (l *Link) withHost(host *url.URL) string {
	l.url.Host = host.Host
	l.url.Scheme = "http"
	return l.Unify()
}

func (l *Link) Unify() string {
	l.url.Fragment = ""
	return l.url.String()
}
