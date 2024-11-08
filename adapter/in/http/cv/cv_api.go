package cv

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"log/slog"

	m "github.com/ernstvorsteveld/go-cv-cassandra/pkg/middleware"
	"github.com/ernstvorsteveld/go-cv-cassandra/pkg/utils"
	"github.com/gin-gonic/gin"

	middleware "github.com/oapi-codegen/gin-middleware"
)

func NewGinCvServer(h *CvApiHandler, c *utils.Configuration) *http.Server {
	slog.Debug("NewGinCvServer", "content", "About to create GinCVServer", "port", c.Api.CV.Port)
	m.ExpectedHosts = c.Api.CV.Expectedhosts
	swagger, err := GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	// This is how you set up a basic gin router
	r := gin.Default()

	r.Use(m.CorrelationId(utils.NewDefaultUuidGenerator()))
	r.Use(m.Authenticate)
	r.Use(m.ValidHostHeaders, m.SecurityHeaders)
	r.Use(m.ErrorHandler())

	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	r.Use(middleware.OapiRequestValidator(swagger))

	RegisterHandlers(r, h)

	s := &http.Server{
		Handler: r,
		Addr:    net.JoinHostPort("0.0.0.0", c.Api.CV.Port),
	}
	return s
}
