package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
)

func main() {

	reqDuration := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "A histogram of the HTTP request duration in seconds.",
		Buckets: prometheus.ExponentialBuckets(0.1, 1.5, 5),
	})

	// Create non-global registry
	registry := prometheus.NewRegistry()

	registry.MustRegister(
		collectors.NewGoCollector(
			collectors.WithGoCollectorRuntimeMetrics(collectors.GoRuntimeMetricsRule{
				Matcher: regexp.MustCompile("/."),
			}),
		),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		reqDuration,
		opsProcessed,
	)
	registry.MustRegister(collectors.NewBuildInfoCollector())

	go func() {
		for {
			// fictional latency
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
			now := time.Now()
			reqDuration.(prometheus.ExemplarObserver).ObserveWithExemplar(
				time.Since(now).Seconds(),
				prometheus.Labels{
					"dummyID": fmt.Sprint(rand.Intn(10000)),
				},
			)
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// exponse /Metrics HTTP endpoint using the created custom registry.
	http.Handle("/metrics", promhttp.HandlerFor(
		registry,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	))

	// To test: curl -H 'Accept: application/openmetrics-text' localhost:2112/metrics
	http.ListenAndServe(":2112", nil)
}
