package cv

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"

	"log/slog"

	"github.com/ernstvorsteveld/go-cv-cassandra/domain/port/in"
	"github.com/ernstvorsteveld/go-cv-cassandra/pkg/utils"
	"github.com/gin-gonic/gin"

	middleware "github.com/oapi-codegen/gin-middleware"
)

var idGenerator = utils.NewDefaultUuidGenerator()
var expecectedHosts StringArray

const expectedHost = "localhost:8091"
const CORRELATION_ID_HEADER = "X-CORRELATION-ID"
const OBJECT_ID_HEADER = "X-OBJECT-ID"
const LOCATION_HEADER = "Location"

type CvApiHandler struct {
	u in.UseCasesPort
	c *utils.Configuration
}

func NewCvApiService(u in.UseCasesPort, c *utils.Configuration) *CvApiHandler {
	return &CvApiHandler{
		u: u,
		c: c,
	}
}

func newError(ctx context.Context, code string) *Error {
	return &Error{
		Code:      code,
		Message:   experienceErrors[code],
		RequestId: utils.GetCorrelationUuid(ctx),
	}
}

func (cs *CvApiHandler) ListExperiences(c *gin.Context, params ListExperiencesParams) {
	slog.Debug("cv.ListExperiences", "content", "About to List Experiences", "correlationId", Get(CORRELATION_ID_HEADER, c))
	ctx := utils.NewDefaultContextWrapper(c, Get(CORRELATION_ID_HEADER, c).(string)).Build()
	es, err := cs.u.ListExperiences(ctx, in.NewListExperienceCommand(int(*params.Page), int(*params.Limit)))
	if err != nil {
		slog.Debug("cv.ListExperiences", "content", "Error while retrieving Experiences", "correlationId", Get(CORRELATION_ID_HEADER, c), "error-code", "EXP0000001", "error", err.Error())
		c.JSON(http.StatusInternalServerError, newError(ctx, "EXP0000001"))
		return
	}
	_, err = json.Marshal(es)
	if err != nil {
		slog.Debug("cv.ListExperiences", "content", "Error while retrieving Experiences", "correlationId", Get(CORRELATION_ID_HEADER, c), "error-code", "EXP0000002", "error", err.Error())
		c.JSON(http.StatusInternalServerError, newError(ctx, "EXP0000002"))
		return
	}

	c.JSON(http.StatusOK, es)
}

// Create an experience
// (POST /experiences)
func (cs *CvApiHandler) CreateExperience(c *gin.Context) {
	slog.Debug("cv.CreateExperience", "content", "About to Create an Experience", "correlationId", Get(CORRELATION_ID_HEADER, c))
	ctx := utils.NewDefaultContextWrapper(c, Get(CORRELATION_ID_HEADER, c).(string)).AddUrl(cs.c.Api.CV.Url).Build()

	var e CreateExperienceRequest
	if err := c.ShouldBindJSON(&e); err != nil {
		slog.Debug("cv.CreateExperience", "content", "Error while creating Experience", "correlationId", Get(CORRELATION_ID_HEADER, c), "error-code", "EXP0000003", "error", err.Error())
		c.JSON(http.StatusBadRequest, newError(ctx, "EXP0000003"))
		return
	}
	m, err := cs.u.CreateExperience(ctx, in.NewCreateExperienceCommand(e.Name, e.Tags))
	if err != nil {
		slog.Debug("cv.CreateExperience", "content", "Error while creating Experience", "correlationId", Get(CORRELATION_ID_HEADER, c), "error-code", "EXP0000004", "error", err.Error())
		c.JSON(http.StatusInternalServerError, newError(ctx, "EXP0000004"))
		return
	}
	c.Writer.Header().Set(LOCATION_HEADER, fmt.Sprintf("%s/experiences/%s", utils.GetHostUrl(ctx), m.GetId()))
	c.Writer.Header().Set(OBJECT_ID_HEADER, m.Id)

	c.Data(http.StatusCreated, "application/json", []byte(m.Id))
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
			slog.Debug("cv.GetExperienceById", "content", "Error while retrieving Experience by Id", "correlationId", Get(CORRELATION_ID_HEADER, c), "error-code", "EXP0000005", "error", err.Error())
			c.JSON(http.StatusInternalServerError, newError(ctx, "EXP0000005"))
			return
		}
		c.JSON(http.StatusOK, e)
	} else {
		if isNotFound(err) {
			slog.Debug("cv.GetExperienceById", "content", "Error while retrieving Experience by Id", "correlationId", Get(CORRELATION_ID_HEADER, c), "error-code", "EXP0000006", "error", err.Error())
			c.JSON(http.StatusNotFound, newError(ctx, "EXP0000006"))
			return
		} else {
			c.JSON(http.StatusInternalServerError, newError(ctx, "EXP0000007"))
			return
		}
	}
}

func isNotFound(err error) bool {
	return err != nil && err.Error() == "not found"
}

func (cs *CvApiHandler) ListTags(c *gin.Context) {
	slog.Debug("cv.ListTags", "content", "About to List Tags, using defaults page=0 and size=100")
	ctx := utils.NewDefaultContextWrapper(c, Get(CORRELATION_ID_HEADER, c).(string)).Build()
	tags, err := cs.u.ListTags(context.Background(), in.NewListTagsCommand(int(0), int(100)))
	if err != nil {
		slog.Debug("cv.ListTags", "content", "Error while retrieving Tags", "correlationId", Get(CORRELATION_ID_HEADER, c), "error-code", "TAG0000001", "error", err.Error())
		c.JSON(http.StatusInternalServerError, newError(ctx, "TAG0000001"))
		return
	}
	c.JSON(http.StatusOK, tags)
}

// Create a tag
// (POST /tags)
func (cs *CvApiHandler) CreateTag(c *gin.Context) {
	slog.Debug("cv.CreateTag", "content", "About to create Tag")
	ctx := utils.NewDefaultContextWrapper(c, Get(CORRELATION_ID_HEADER, c).(string)).AddUrl(cs.c.Api.CV.Url).Build()

	var t Tag
	if err := c.ShouldBindJSON(&t); err != nil {
		slog.Debug("cv.CreateTag", "content", "Error while creating Tag", "correlationId", Get(CORRELATION_ID_HEADER, c), "error-code", "TAG0000002", "error", err.Error())
		c.JSON(http.StatusBadRequest, newError(ctx, "EXP0000002"))
		return
	}
	m, err := cs.u.CreateTag(ctx, in.NewCreateTagCommand(t.Tag))
	if err != nil {
		slog.Debug("cv.CreateTag", "content", "Error while creating Tag", "correlationId", Get(CORRELATION_ID_HEADER, c), "error-code", "TAG0000003", "error", err.Error())
		c.JSON(http.StatusInternalServerError, newError(ctx, "EXP0000004"))
		return
	}
	c.Writer.Header().Set(LOCATION_HEADER, fmt.Sprintf("%s/tags/%s", utils.GetHostUrl(ctx), m.GetId()))
	c.Writer.Header().Set(OBJECT_ID_HEADER, m.Id)
	c.JSON(http.StatusCreated, m.Id)
}

func NewGinCvServer(h *CvApiHandler, c *utils.Configuration) *http.Server {
	slog.Debug("NewGinCvServer", "content", "About to create GinCVServer", "port", c.Api.CV.Port)
	expecectedHosts = c.Api.CV.Expectedhosts
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
		Addr:    net.JoinHostPort("0.0.0.0", c.Api.CV.Port),
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
	if !expecectedHosts.contains(c.Request.Host) {
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

type StringArray []string

func (v StringArray) contains(s string) bool {
	// iterate using the for loop
	for i := 0; i < len(v); i++ {
		if v[i] == s {
			return true
		}
	}
	return false
}

var experienceErrors = map[string]string{
	"EXP0000001": "error while retrieving experiences",
	"EXP0000002": "error while marshalling experiences",
	"EXP0000003": "error while marshalling input payload Experience",
	"EXP0000004": "Internal Server error: error while creating experience",
	"EXP0000005": "Internal Server error: error while marshalling experience",
	"EXP0000006": "Experience with the provided id does not exist",
	"EXP0000007": "Other error wjile retrieving experience",
	"TAG0000001": "Error while retrieving tags",
	"TAG0000002": "error while marshalling input payload Tag",
	"TAG0000003": "Internal Server error: error while creating Tag",
}
