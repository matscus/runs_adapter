package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	RequestsDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "unimock_http_requests_duration_ms",
			Help:       "A summary of the handling duration of requests.",
			Objectives: map[float64]float64{0.9: 0.01, 0.95: 0.01, 0.99: 0.01},
		},
		[]string{"method", "path"},
	)
)

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Content-Type", "application/json")
		start := time.Now()
		defer func() {
			if err := recover(); err != nil {
				log.Error("Middleware recover panic ", err)
				c.JSON(500, gin.H{"Status": "error", "Message": err})
			}
		}()
		c.Next()
		RequestsDuration.WithLabelValues(c.Request.Method, c.Request.URL.Path).Observe(float64(time.Since(start).Milliseconds()))
	}
}
