// This is an example of implementing the Pet Store from the OpenAPI documentation
// found at:
// https://github.com/OAI/OpenAPI-Specification/blob/master/examples/v3.0/petstore.yaml

package api

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/ernstvorsteveld/go-cv-cassandra/src/pkg/domain"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	middleware "github.com/oapi-codegen/gin-middleware"
)

type CvApiService struct {
	cvService *domain.CvServices
}

func NewCvApiService() *CvApiService {
	return &CvApiService{
		cvService: domain.NewCvService(),
	}
}

func (cs *CvApiService) ListExperiences(c *gin.Context, params ListExperiencesParams) {
	log.Debugf("About to List Experiences")
	cs.cvService.ListExperiences(int(*params.Page), int(*params.Limit))
}

// Create an experience
// (POST /experiences)
func (cs *CvApiService) CreateExperience(c *gin.Context) {
	log.Debugf("About to Create an Experience")
}

// Info for a specific experience
// (GET /experiences/{id})
func (cs *CvApiService) GetExperienceById(c *gin.Context, id string) {
	log.Debugf("About to Get an Experience by Id")
}

func NewGinCvServer(cvApiService *CvApiService, port string) *http.Server {
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

	// We now register our petStore above as the handler for the interface
	RegisterHandlers(r, cvApiService)

	s := &http.Server{
		Handler: r,
		Addr:    net.JoinHostPort("0.0.0.0", port),
	}
	return s
}
