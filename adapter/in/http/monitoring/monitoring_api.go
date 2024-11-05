package monitoring

import (
	"log/slog"
	"net"
	"net/http"

	"github.com/ernstvorsteveld/go-cv-cassandra/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MonitoringApiHandler struct {
	c *utils.Configuration
}

func NewMonitoringApiService(c *utils.Configuration) *MonitoringApiHandler {
	return &MonitoringApiHandler{
		c: c,
	}
}

func (cs *MonitoringApiHandler) Metrics(c *gin.Context) {
	slog.Debug("Metrics", "content", "About to GET Metrics.")
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}

func NewGinMonitoringServer(h *MonitoringApiHandler, c *utils.Configuration) *http.Server {
	slog.Debug("NewGinMonitoringServer", "content", "About to create GinCVServer", "port", c.Api.Monitoring.Port)
	// swagger, err := GetSwagger()
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
	// 	os.Exit(1)
	// }

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	// swagger.Servers = nil

	// This is how you set up a basic gin router
	r := gin.Default()

	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	// r.Use(middleware.OapiRequestValidator(swagger))

	RegisterHandlers(r, h)

	s := &http.Server{
		Handler: r,
		Addr:    net.JoinHostPort("0.0.0.0", c.Api.Monitoring.Port),
	}
	return s
}
