//go:build tools
// +build tools

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=./types.cfg.yaml ../../../api/api.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=./server.cfg.yaml ../../../api/api.yaml

package api

import (
	"github.com/gin-gonic/gin"
	_ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
)

func ListExperiences(c *gin.Context, params ListExperiencesParams) {

}

// Create an experience
// (POST /experiences)
func CreateExperience(c *gin.Context) {

}

// Info for a specific experience
// (GET /experiences/{id})
func GetExperienceById(c *gin.Context, id string) {

}
