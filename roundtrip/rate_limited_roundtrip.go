package roundtrip

import (
	"github.com/2at2/httpe/limiter"
	"net/http"
	"time"
)

type RateLimitedRoundTripper struct {
	origin    http.RoundTripper
	throttler limiter.Throttler
}

func NewRateLimitedRoundTripper(
	origin http.RoundTripper,
	throttler limiter.Throttler,
) (http.RoundTripper, error) {
	if origin == nil {
		origin = http.DefaultTransport
	}
	if throttler == nil {
		if x, err := limiter.NewThrottler(time.Second, 100); err != nil {
			return nil, err
		} else {
			throttler = x
		}
	}

	return &RateLimitedRoundTripper{
		origin:    origin,
		throttler: throttler,
	}, nil
}

func (r *RateLimitedRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	r.throttler.WaitUrl(req.URL)
	return r.origin.RoundTrip(req)
}
