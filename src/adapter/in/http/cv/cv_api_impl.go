package cv

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"

	"log/slog"

	"github.com/ernstvorsteveld/go-cv-cassandra/src/domain/model"
	"github.com/ernstvorsteveld/go-cv-cassandra/src/port/in"
	"github.com/ernstvorsteveld/go-cv-cassandra/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	middleware "github.com/oapi-codegen/gin-middleware"
)

var idGenerator = utils.NewDefaultUuidGenerator()

const expectedHost = "localhost:8091"
const CORRELATION_ID_HEADER = "X-CORRELATION-ID"
const LOCATION_HEADER = "Location"

type CvApiServices interface {
}

type CvApiHandler struct {
	u in.UseCasesPort
}

func NewCvApiService(u in.UseCasesPort) *CvApiHandler {
	return &CvApiHandler{
		u: u,
	}
}

func (cs *CvApiHandler) ListExperiences(c *gin.Context, params ListExperiencesParams) {
	slog.Debug("cv.ListExperiences", "content", "About to List Experiences", "correlationId", Get(CORRELATION_ID_HEADER, c))
	ctx := utils.NewDefaultContextWrapper(c, Get(CORRELATION_ID_HEADER, c).(string)).Build()
	es, err := cs.u.ListExperiences(ctx, in.NewListExperienceCommand(int(*params.Page), int(*params.Limit)))
	if err == nil {
		e := Error{Code: "EXP0000002", Message: "Error while retrieving experiences.", RequestId: uuid.New()}
		c.JSON(http.StatusInternalServerError, e)
	}
	_, err = json.Marshal(es)
	if err != nil {
		e := Error{Code: "EXP0000003", Message: "Error while marshalling experiences.", RequestId: uuid.New()}
		c.JSON(http.StatusInternalServerError, e)
		return
	}

	c.JSON(http.StatusOK, es)
}

// Create an experience
// (POST /experiences)
func (cs *CvApiHandler) CreateExperience(c *gin.Context) {
	slog.Debug("cv.CreateExperience", "content", "About to Create an Experience", "correlationId", Get(CORRELATION_ID_HEADER, c))
	var e model.Experience
	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := utils.NewDefaultContextWrapper(c, Get(CORRELATION_ID_HEADER, c).(string)).Build()
	m, err := cs.u.CreateExperience(ctx, in.NewCreateExperienceCommand(e.GetName(), e.GetTags()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	location := fmt.Sprintf("http://localhost:8091/experiences/%s", m.Id)
	c.Writer.Header().Set(LOCATION_HEADER, location)
	c.JSON(http.StatusCreated, m.Id)
}

// Info for a specific experience
// (GET /experiences/{id})
func (cs *CvApiHandler) GetExperienceById(c *gin.Context, id string) {
	slog.Debug("cv.GetExperienceById", "content", "About to Get Experience by Id", "correlationId", Get(CORRELATION_ID_HEADER, c))
	ctx := utils.NewDefaultContextWrapper(c, Get(CORRELATION_ID_HEADER, c).(string)).Build()
	e, err := cs.u.GetExperienceById(ctx, in.NewGetExperienceCommand(id))
	if err == nil {
		_, err := json.Marshal(e)
		if err != nil {
			c.JSON(400, nil)
			return
		}

		c.JSON(http.StatusOK, e)

	} else {
		c.Request.Response.StatusCode = 400
	}
}

func (cs *CvApiHandler) ListTags(c *gin.Context) {
	slog.Debug("cv.ListTags", "content", "About to List Tags, using defaults page=0 and size=100")
	tags, _ := cs.u.ListTags(context.Background(), in.NewListTagsCommand(int(0), int(100)))
	c.JSON(http.StatusOK, tags)
}

func (cs *CvApiHandler) Metrics(c *gin.Context) {
	slog.Debug("Metrics", "content", "About to GET Metrics.")
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}

func NewGinCvServer(h *CvApiHandler, port string) *http.Server {
	slog.Debug("NewGinCvServer", "content", "About to create GinCVServer")
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

	r.Use(CorrelationId, Authenticate, SecurityHeaders)

	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	r.Use(middleware.OapiRequestValidator(swagger))

	RegisterHandlers(r, h)

	s := &http.Server{
		Handler: r,
		Addr:    net.JoinHostPort("0.0.0.0", port),
	}
	return s
}

func CorrelationId(c *gin.Context) {
	correlationId := idGenerator.UUIDString()
	slog.Debug("cv.CorrelationId", "content", "About to add correlationId", "correlationId", correlationId)
	c.Header(CORRELATION_ID_HEADER, correlationId)
	c.Next()
}

func Authenticate(c *gin.Context) {
	slog.Debug("cv.Authenticate", "content", "About to authenticate", "correlationId", Get(CORRELATION_ID_HEADER, c))
	c.Next()
}

func SecurityHeaders(c *gin.Context) {
	slog.Debug("cv.SecurityHeaders", "content", "About to add security headers", "correlationId", Get(CORRELATION_ID_HEADER, c))
	if c.Request.Host != expectedHost {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid host header"})
		return
	}
	c.Header("X-Frame-Options", "DENY")
	c.Header("X-XSS-Protection", "1; mode=block")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Referrer-Policy", "strict-origin")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Content-Security-Policy", "default-src 'self'; connect-src *; font-src *; script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';")
	c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	c.Next()
}

func Get(k string, c *gin.Context) any {
	value := c.Writer.Header().Get(CORRELATION_ID_HEADER)
	if value == "" {
		return "UNKNOWN"
	}
	return value
}
