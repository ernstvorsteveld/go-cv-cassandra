package experiences

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/ernstvorsteveld/go-cv-cassandra/src/domain/model"
	"github.com/ernstvorsteveld/go-cv-cassandra/src/port/in"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	middleware "github.com/oapi-codegen/gin-middleware"
)

type CvApiServices interface {
}

type CvApiHandler struct {
	h in.ExperienceUseCases
}

func NewCvApiService(s in.ExperienceUseCases) *CvApiHandler {
	return &CvApiHandler{
		h: s,
	}
}

func (cs *CvApiHandler) ListExperiences(c *gin.Context, params ListExperiencesParams) {
	log.Debugf("About to List Experiences")
	cs.h.ListExperiences(context.Background(), in.NewListExperienceCommand(int(*params.Page), int(*params.Limit)))
}

// Create an experience
// (POST /experiences)
func (cs *CvApiHandler) CreateExperience(c *gin.Context) {
	log.Debugf("About to Create an Experience")
	var e model.Experience
	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cs.h.CreateExperience(context.Background(), in.NewCreateExperienceCommand(e.GetName(), e.GetTags()))
}

// Info for a specific experience
// (GET /experiences/{id})
func (cs *CvApiHandler) GetExperienceById(c *gin.Context, id string) {
	log.Debugf("About to Get an Experience by Id")
	e, err := cs.h.GetExperienceById(context.Background(), in.NewGetExperienceCommand(id))
	if err != nil {
		c.Request.Response.StatusCode = 400
	} else {
		body, err := json.Marshal(e)
		if err != nil {
			c.JSON(400, nil)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": string(body),
		})
	}
}

func (cs *CvApiHandler) ListTags(c *gin.Context) {
	log.Debugf("About to List Tags")
}

func NewGinCvServer(h *CvApiHandler, port string) *http.Server {
	log.Debugf("About to create GinCVServer")
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
