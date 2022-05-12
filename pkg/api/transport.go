package api

import (
	"net/http"

	"github.com/jinzhu/copier"
)

// UserAgentTransport implements a http.RoundTripper which adds a User-Agent header to requests
type UserAgentTransport struct {
	ua string
	rt http.RoundTripper
}

func NewUserAgentTransport(ua string, rt http.RoundTripper) *UserAgentTransport {
	t := &UserAgentTransport{ua: ua, rt: rt}

	return t
}

func (t *UserAgentTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	req := &http.Request{}
	if err := copier.Copy(req, r); err != nil {
		return t.rt.RoundTrip(r)
	}

	req.Header.Set("User-Agent", t.ua)

	return t.rt.RoundTrip(r)
}
