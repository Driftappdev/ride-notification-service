package pkg

import (
	gometrics "github.com/driftappdev/libpackage/gometrics"
	"net/http"
)

func NewRequestCounter(name, help string) *gometrics.Counter {
	return gometrics.NewCounter(name, help, "method", "path", "status")
}
func NewRequestLatency(name, help string) *gometrics.Histogram {
	return gometrics.NewHistogram(name, help, gometrics.DefBuckets, "method", "path")
}
func MetricsHandler() http.HandlerFunc { return gometrics.Handler() }
