package middleware

import (
	"expvar"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/johnnrails/ddd_go/boilerplate/logger"
	"github.com/johnnrails/ddd_go/boilerplate/reverseproxy"
	"github.com/vardius/gorouter/v4"
	"golang.org/x/time/rate"
)

var rateLimits *expvar.Map

type visitor struct {
	*rate.Limiter

	lastSeen time.Time
}

type rateLimiter struct {
	sync.Mutex

	burst    int
	rate     rate.Limit
	visitors map[string]*visitor
}

func (l *rateLimiter) allow(ip string) bool {
	l.Lock()
	defer l.Unlock()
	v, exists := l.visitors[ip]
	if !exists {
		v = &visitor{
			Limiter: rate.NewLimiter(l.rate, l.burst),
		}
		l.visitors[ip] = v
	}
	v.lastSeen = time.Now()

	if rateLimits != nil {
		rateLimits.Add(ip, 1)
	}

	return v.Allow()
}

func (l *rateLimiter) cleanup(frequency time.Duration) {
	for {
		time.Sleep(frequency)
		l.Lock()
		for ip, v := range l.visitors {
			if time.Since(v.lastSeen) > frequency {
				delete(l.visitors, ip)
				if rateLimits != nil {
					rateLimits.Delete(ip)
				}
			}
		}
		l.Unlock()
	}
}

func RateLimit(rateLimit rate.Limit, burst int, frequency time.Duration) gorouter.MiddlewareFunc {
	rLimiter := &rateLimiter{
		rate:     rateLimit,
		burst:    burst,
		visitors: make(map[string]*visitor),
	}

	go rLimiter.cleanup(frequency)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, err := reverseproxy.GetIpAddress(r)

			if err != nil {
				logger.Error(r.Context(), fmt.Sprintf("[HTTP] RateLimit invalid IP Address: %v", err))
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			if !rLimiter.allow(string(ip)) {
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
