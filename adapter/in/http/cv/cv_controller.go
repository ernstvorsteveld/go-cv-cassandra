package cv

import (
	"encoding/json"
	"fmt"
	"net/http"

	"log/slog"

	"github.com/ernstvorsteveld/go-cv-cassandra/domain/port/in"
	m "github.com/ernstvorsteveld/go-cv-cassandra/pkg/middleware"
	"github.com/ernstvorsteveld/go-cv-cassandra/pkg/utils"
	"github.com/gin-gonic/gin"
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

func (cs *CvApiHandler) ListExperiences(c *gin.Context, params ListExperiencesParams) {
	cId := m.GetCorrelationIdHeader(c)
	slog.Debug("cv.ListExperiences", "content", "About to List Experiences", "correlationId", cId)
	ctx := utils.NewDefaultContextWrapper(c, cId).Build()
	es, err := cs.u.ListExperiences(ctx, in.NewListExperienceCommand(int(*params.Page), int(*params.Limit)))
	if err != nil {
		NewListExperienceError(c, err)
		return
	}
	_, err = json.Marshal(es)
	if err != nil {
		NewListExperienceMarshalError(c, err)
		return
	}

	c.JSON(http.StatusOK, es)
}

// Create an experience
// (POST /experiences)
func (cs *CvApiHandler) CreateExperience(c *gin.Context) {
	cId := m.GetCorrelationIdHeader(c)
	slog.Debug("cv.CreateExperience", "content", "About to Create an Experience", "correlationId", cId)
	ctx := utils.NewDefaultContextWrapper(c, cId).AddUrl(cs.c.Api.Url).Build()

	var e CreateExperienceRequest
	if err := c.ShouldBindJSON(&e); err != nil {
		NewExperienceBindError(c, err)
		return
	}
	m, err := cs.u.CreateExperience(ctx, in.NewCreateExperienceCommand(e.Name, e.Tags))
	if err != nil {
		NewCreateExperienceError(c, err)
		return
	}
	c.Writer.Header().Set(LOCATION_HEADER, location(c, m.GetId()))
	c.Writer.Header().Set(OBJECT_ID_HEADER, m.Id)

	c.Data(http.StatusCreated, "application/json", []byte(m.Id))
}

func location(c *gin.Context, id string) string {
	host := c.Request.Host
	path := c.Request.URL.String()

	base := host + path
	return fmt.Sprintf("%s/%s", base, id)
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
			NewGetExperienceByIdMarshalError(c, err)
			return
		}
		c.JSON(http.StatusOK, e)
	} else {
		NewGetExperienceByIdNotFoundError(c, err)
	}
}

func (cs *CvApiHandler) ListTags(c *gin.Context) {
	cId := m.GetCorrelationIdHeader(c)
	slog.Debug("cv.ListTags", "content", "About to List Tags, using defaults page=0 and size=100")
	ctx := utils.NewDefaultContextWrapper(c, cId).Build()
	tags, err := cs.u.ListTags(ctx, in.NewListTagsCommand(int(0), int(100)))
	if err != nil {
		NewListTagsError(c, err)
		return
	}
	c.JSON(http.StatusOK, tags)
}

// Create a tag
// (POST /tags)
func (cs *CvApiHandler) CreateTag(c *gin.Context) {
	cId := m.GetCorrelationIdHeader(c)
	slog.Debug("cv.CreateTag", "content", "About to create Tag")
	ctx := utils.NewDefaultContextWrapper(c, cId).AddUrl(cs.c.Api.Url).Build()

	var t Tag
	if err := c.ShouldBindJSON(&t); err != nil {
		NewCreateTagMarshalError(c, err)
		return
	}
	m, err := cs.u.CreateTag(ctx, in.NewCreateTagCommand(t.Tag))
	if err != nil {
		NewCreateTagError(c, err)
		return
	}
	c.Writer.Header().Set(LOCATION_HEADER, fmt.Sprintf("%s/tags/%s", utils.GetHostUrl(ctx), m.GetId()))
	c.Writer.Header().Set(OBJECT_ID_HEADER, m.Id)
	c.JSON(http.StatusCreated, m.Id)
}