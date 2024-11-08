package middleware

import (
	"log/slog"
	"net/http"

	"github.com/ernstvorsteveld/go-cv-cassandra/pkg/utils"
	"github.com/gin-gonic/gin"
)

const CORRELATION_ID_HEADER = "X-CORRELATION-ID"

var ExpectedHosts StringArray

type StringArray []string

func CorrelationId(ig utils.IdGenerator) gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationId := ig.UUIDString()
		slog.Debug("cv.CorrelationId", "content", "About to add correlationId", "correlationId", correlationId)
		c.Header(CORRELATION_ID_HEADER, correlationId)
		c.Next()
	}
}

func Authenticate(c *gin.Context) {
	slog.Debug("cv.Authenticate", "content", "About to authenticate", "correlationId", GetCorrelationIdHeader(c))
	c.Next()
}

func ValidHostHeaders(c *gin.Context) {
	slog.Debug("cv.HostHeaders", "content", "About to add security headers", "correlationId", GetCorrelationIdHeader(c))
	if !ExpectedHosts.contains(c.Request.Host) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid host header"})
		return
	}
}

func SecurityHeaders(c *gin.Context) {
	slog.Debug("cv.SecurityHeaders", "content", "About to add security headers", "correlationId", GetCorrelationIdHeader(c))
	c.Header("X-Frame-Options", "DENY")
	c.Header("X-XSS-Protection", "1; mode=block")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Referrer-Policy", "strict-origin")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Content-Security-Policy", "default-src 'self'; connect-src *; font-src *; script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';")
	c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	c.Next()
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			c.Next()

			if err := recover(); err != nil {
				slog.Info("cv.ErrorHandler", "content", "Recovered Internal Server Error", "correlationId", GetCorrelationIdHeader(c), "error", err)
				e := Error{
					Code:      "ALG0000001",
					Message:   "Internal Server Error",
					RequestId: utils.GetCorrelationUuid(c).String(),
				}
				c.JSON(http.StatusInternalServerError, e)
			}
		}()
	}
}

func (v StringArray) contains(s string) bool {
	// iterate using the for loop
	for i := 0; i < len(v); i++ {
		if v[i] == s {
			return true
		}
	}
	return false
}

func GetCorrelationIdHeader(c *gin.Context) string {
	return c.Writer.Header().Get(CORRELATION_ID_HEADER)
}
