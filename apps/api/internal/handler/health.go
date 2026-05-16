package handler

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var startTime = time.Now()

type HealthResponse struct {
	Status  string      `json:"status"`
	Uptime  string      `json:"uptime"`
	System  *SystemInfo `json:"system,omitempty"`
}

type SystemInfo struct {
	GoVersion    string  `json:"go_version"`
	NumCPU       int     `json:"num_cpu"`
	NumGoroutine int     `json:"num_goroutine"`
	MemAllocMB   float64 `json:"mem_alloc_mb"`
}

func Healthz(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		c.JSON(http.StatusOK, HealthResponse{
			Status: "pass",
			Uptime: time.Since(startTime).Round(time.Second).String(),
			System: &SystemInfo{
				GoVersion:    runtime.Version(),
				NumCPU:       runtime.NumCPU(),
				NumGoroutine: runtime.NumGoroutine(),
				MemAllocMB:   float64(m.Alloc) / 1024 / 1024,
			},
		})
	}
}

func Readyz(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, HealthResponse{Status: "pass"})
	}
}

func Metrics() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
