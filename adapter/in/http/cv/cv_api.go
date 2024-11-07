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
	m "github.com/ernstvorsteveld/go-cv-cassandra/pkg/middleware"
	"github.com/ernstvorsteveld/go-cv-cassandra/pkg/utils"
	"github.com/gin-gonic/gin"

	middleware "github.com/oapi-codegen/gin-middleware"
)

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
		Message:   m.ExperienceErrors[code],
		RequestId: utils.GetCorrelationUuid(ctx),
	}
}

func (cs *CvApiHandler) ListExperiences(c *gin.Context, params ListExperiencesParams) {
	cId := m.GetCorrelationIdHeader(c)
	slog.Debug("cv.ListExperiences", "content", "About to List Experiences", "correlationId", cId)
	ctx := utils.NewDefaultContextWrapper(c, cId).Build()
	es, err := cs.u.ListExperiences(ctx, in.NewListExperienceCommand(int(*params.Page), int(*params.Limit)))
	if err != nil {
		slog.Debug("cv.ListExperiences", "content", "Error while retrieving Experiences", "correlationId", cId, "error-code", "EXP0000001", "error", err.Error())
		c.JSON(http.StatusInternalServerError, newError(ctx, "EXP0000001"))
		return
	}
	_, err = json.Marshal(es)
	if err != nil {
		slog.Debug("cv.ListExperiences", "content", "Error while retrieving Experiences", "correlationId", cId, "error-code", "EXP0000002", "error", err.Error())
		c.JSON(http.StatusInternalServerError, newError(ctx, "EXP0000002"))
		return
	}

	c.JSON(http.StatusOK, es)
}

// Create an experience
// (POST /experiences)
func (cs *CvApiHandler) CreateExperience(c *gin.Context) {
	cId := m.GetCorrelationIdHeader(c)
	slog.Debug("cv.CreateExperience", "content", "About to Create an Experience", "correlationId", cId)
	ctx := utils.NewDefaultContextWrapper(c, cId).AddUrl(cs.c.Api.CV.Url).Build()

	var e CreateExperienceRequest
	if err := c.ShouldBindJSON(&e); err != nil {
		slog.Debug("cv.CreateExperience", "content", "Error while creating Experience", "correlationId", cId, "error-code", "EXP0000003", "error", err.Error())
		c.JSON(http.StatusBadRequest, newError(ctx, "EXP0000003"))
		return
	}
	m, err := cs.u.CreateExperience(ctx, in.NewCreateExperienceCommand(e.Name, e.Tags))
	if err != nil {
		slog.Debug("cv.CreateExperience", "content", "Error while creating Experience", "correlationId", cId, "error-code", "EXP0000004", "error", err.Error())
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
	cId := m.GetCorrelationIdHeader(c)
	slog.Debug("cv.GetExperienceById", "content", "About to Get Experience by Id", "correlationId", cId)
	ctx := utils.NewDefaultContextWrapper(c, cId).Build()
	e, err := cs.u.GetExperienceById(ctx, in.NewGetExperienceCommand(id))
	if err == nil {
		_, err := json.Marshal(e)
		if err != nil {
			slog.Debug("cv.GetExperienceById", "content", "Error while retrieving Experience by Id", "correlationId", cId, "error-code", "EXP0000005", "error", err.Error())
			c.JSON(http.StatusInternalServerError, newError(ctx, "EXP0000005"))
			return
		}
		c.JSON(http.StatusOK, e)
	} else {
		if isNotFound(err) {
			slog.Debug("cv.GetExperienceById", "content", "Error while retrieving Experience by Id", "correlationId", cId, "error-code", "EXP0000006", "error", err.Error())
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
	cId := m.GetCorrelationIdHeader(c)
	slog.Debug("cv.ListTags", "content", "About to List Tags, using defaults page=0 and size=100")
	ctx := utils.NewDefaultContextWrapper(c, cId).Build()
	tags, err := cs.u.ListTags(context.Background(), in.NewListTagsCommand(int(0), int(100)))
	if err != nil {
		slog.Debug("cv.ListTags", "content", "Error while retrieving Tags", "correlationId", cId, "error-code", "TAG0000001", "error", err.Error())
		c.JSON(http.StatusInternalServerError, newError(ctx, "TAG0000001"))
		return
	}
	c.JSON(http.StatusOK, tags)
}

// Create a tag
// (POST /tags)
func (cs *CvApiHandler) CreateTag(c *gin.Context) {
	cId := m.GetCorrelationIdHeader(c)
	slog.Debug("cv.CreateTag", "content", "About to create Tag")
	ctx := utils.NewDefaultContextWrapper(c, cId).AddUrl(cs.c.Api.CV.Url).Build()

	var t Tag
	if err := c.ShouldBindJSON(&t); err != nil {
		slog.Debug("cv.CreateTag", "content", "Error while creating Tag", "correlationId", cId, "error-code", "TAG0000002", "error", err.Error())
		c.JSON(http.StatusBadRequest, newError(ctx, "EXP0000002"))
		return
	}
	m, err := cs.u.CreateTag(ctx, in.NewCreateTagCommand(t.Tag))
	if err != nil {
		slog.Debug("cv.CreateTag", "content", "Error while creating Tag", "correlationId", cId, "error-code", "TAG0000003", "error", err.Error())
		c.JSON(http.StatusInternalServerError, newError(ctx, "EXP0000004"))
		return
	}
	c.Writer.Header().Set(LOCATION_HEADER, fmt.Sprintf("%s/tags/%s", utils.GetHostUrl(ctx), m.GetId()))
	c.Writer.Header().Set(OBJECT_ID_HEADER, m.Id)
	c.JSON(http.StatusCreated, m.Id)
}

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

	r.Use(m.CorrelationId, m.Authenticate, m.ValidHostHeaders, m.SecurityHeaders, m.ErrorHandler())

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
