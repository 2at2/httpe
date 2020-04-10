package limiter

import (
	lru "github.com/hashicorp/golang-lru"
	"net/url"
	"sync"
	"time"
)

type Throttler interface {
	Wait(u string) bool
	WaitUrl(u *url.URL) bool
}

type throttler struct {
	delay time.Duration
	cache *lru.Cache

	m sync.Mutex
}

func NewThrottler(delay time.Duration, size int) (Throttler, error) {
	cache, err := lru.New(size)
	if err != nil {
		return nil, err
	}

	return &throttler{
		delay: delay,
		cache: cache,
	}, nil
}

func (t *throttler) Wait(u string) bool {
	if uri, err := url.Parse(u); err != nil {
		return true
	} else {
		return t.throttle(uri)
	}
}

func (t *throttler) WaitUrl(u *url.URL) bool {
	return t.throttle(u)
}

func (t *throttler) throttle(uri *url.URL) bool {
	t.m.Lock()
	val, ok := t.cache.Get(uri.Host)

	if !ok {
		t.cache.Add(uri.Host, time.Now().Add(t.delay))
		t.m.Unlock()
	} else {
		x := val.(time.Time)
		t.cache.Add(uri.Host, x.Add(t.delay))
		t.m.Unlock()

		if tick := time.Tick(x.Sub(time.Now())); tick != nil {
			<-tick
		}
	}

	return true
}
