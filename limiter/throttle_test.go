package limiter

import (
	"fmt"
	. "github.com/stretchr/testify/assert"
	"net/url"
	"sync"
	"testing"
	"time"
)

func TestThrottle(t *testing.T) {
	throttle, err := NewThrottler(time.Millisecond*500, 100)

	if !NoError(t, err) {
		return
	}

	uri, _ := url.Parse("http://domain.com/foo")

	start := time.Now()
	for i := 0; i < 5; i++ {
		if throttle.WaitUrl(uri) {
			//fmt.Println("locked")
		} else {
			False(t, true)
		}
	}

	Greater(t, time.Now().Sub(start).Nanoseconds(), (time.Millisecond * 490 * 4).Nanoseconds())
}

func TestThrottleParallel(t *testing.T) {
	delay := time.Millisecond * 200

	throttle, err := NewThrottler(delay, 100)

	if !NoError(t, err) {
		return
	}

	mainStart := time.Now()
	threads := 50
	requests := 10

	wg := &sync.WaitGroup{}
	for thread := 0; thread < threads; thread++ {
		wg.Add(1)
		go func(thread int) {
			defer wg.Done()

			uri, _ := url.Parse(fmt.Sprintf("http://domain_%d.com/foo", thread))

			wgg := &sync.WaitGroup{}

			start := time.Now()
			for i := 0; i < requests; i++ {
				wgg.Add(1)

				go func(index int) {
					defer wgg.Done()

					if throttle.WaitUrl(uri) {
						//fmt.Println(time.Now(), "locked", uri.Host)
					} else {
						False(t, true)
					}
				}(i)
			}

			wgg.Wait()

			diff := time.Now().Sub(start)
			Greater(t, diff.Nanoseconds(), (delay * time.Duration(requests-1)).Nanoseconds())
		}(thread)
	}

	wg.Wait()

	Greater(t, time.Now().Sub(mainStart).Nanoseconds(), (delay * time.Duration(requests-1)).Nanoseconds())
}
