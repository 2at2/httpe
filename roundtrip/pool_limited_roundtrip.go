package roundtrip

import (
	"net/http"
)

type PoolLimitedRoundTripper struct {
	origin  http.RoundTripper
	limiter chan bool
}

func NewPoolLimitedRoundTripper(
	origin http.RoundTripper,
	max int,
) (http.RoundTripper, error) {
	if origin == nil {
		origin = http.DefaultTransport
	}

	return &PoolLimitedRoundTripper{
		origin:  origin,
		limiter: make(chan bool, max),
	}, nil
}

func (r *PoolLimitedRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	r.limiter <- true
	resp, err := r.origin.RoundTrip(req)
	<-r.limiter

	return resp, err
}
