package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var apiCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "api_request_count",
		Help: "API request count by path",
	},
	[]string{"path"},
)

func init() {
	prometheus.MustRegister(apiCounter)
}

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		apiCounter.WithLabelValues(c.FullPath()).Inc()
	}
}

func PrometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
