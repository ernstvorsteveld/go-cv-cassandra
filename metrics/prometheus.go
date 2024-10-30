package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	HttpRequestCountWithPath = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total_with_path",
			Help: "Number of HTTP requests by path.",
		},
		[]string{"url"},
	)
)

func init() {
	log.Debugf("Register Handler for Prometheus")
	prometheus.MustRegister(HttpRequestCountWithPath)
}
