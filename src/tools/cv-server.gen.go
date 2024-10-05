// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"github.com/oapi-codegen/runtime"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// List all experiences
	// (GET /experiences)
	ListExperiences(c *gin.Context, params ListExperiencesParams)
	// Create an experience
	// (POST /experiences)
	CreateExperience(c *gin.Context)
	// Info for a specific experience
	// (GET /experiences/{id})
	GetExperienceById(c *gin.Context, id string)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// ListExperiences operation middleware
func (siw *ServerInterfaceWrapper) ListExperiences(c *gin.Context) {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params ListExperiencesParams

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", c.Request.URL.Query(), &params.Limit)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter limit: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "page" -------------

	err = runtime.BindQueryParameter("form", true, false, "page", c.Request.URL.Query(), &params.Page)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter page: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.ListExperiences(c, params)
}

// CreateExperience operation middleware
func (siw *ServerInterfaceWrapper) CreateExperience(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.CreateExperience(c)
}

// GetExperienceById operation middleware
func (siw *ServerInterfaceWrapper) GetExperienceById(c *gin.Context) {

	var err error

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithOptions("simple", "id", c.Param("id"), &id, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetExperienceById(c, id)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router gin.IRouter, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router gin.IRouter, si ServerInterface, options GinServerOptions) {
	errorHandler := options.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.GET(options.BaseURL+"/experiences", wrapper.ListExperiences)
	router.POST(options.BaseURL+"/experiences", wrapper.CreateExperience)
	router.GET(options.BaseURL+"/experiences/:id", wrapper.GetExperienceById)
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/8xWTY/bNhD9K8S0R64l73YDRLdNsCgMpE3QGj3U8IGVRjZTiWRIyrAr6L8XQ8qWtfK6",
	"RhsU2dOa4ny8eW+e1EKua6MVKu8ga8HlW6xF+Pe9ReHxeW/QSlQ5/oJfGnSeHomq+lhCtmrhe4slZPBd",
	"MqRJ+hzJUmyerBWH5cEgdPz65aHOJ3GotCigW3ccnq3Vlkoaqw1aLzH0lusCw6nwHq2CDFZPd7+v24du",
	"ld69XbdvOuDgqW4GzlupNtBxqNE5sQmBk2c2olsU9LTUthYeMmgaWUwz9delxQKy1Vkoj40NldanWP3H",
	"Z8w9VRqQ3j7KjyF6UfybMfLbWOrnPYmfzF6J+tIIXwwl3LoOP2STHuvwz22oAo1iv4hR8zQ9VRARBIfT",
	"rCadyzG5Uvk3PwzsSuVxgxY6grIUm2m8j4enBLN2zh9Tklot9h9QbfwWssf0n/RCadaxRpz8rWOgprop",
	"4NGeXWra/bcK3YRFOpKq1JSukjkqh4My4KfFMiSRvqKf7xtrZd5UTc1+k14ge/q0AA47tE5qWt35LJ2l",
	"FKENKmEkZPAQjjjt9zY0neBYNhsMRkRAhZdaEd3wQTp/Li+Kt6JGj9aFFSvQ5VYaH+v+Kv9Cpkvmt8iM",
	"2CBntdjLuqmZdIykxQosRVN5+n3/CIQZMvjSoD0AP6KtZC098N45qak+CrL7Rz5S28N9lArVGIt30B5/",
	"2eWy7455zSz6xqpXGqFLl/uYX2vjbfi70MmaVOuMVi6O/D5No/Uqjyq+BoypZB7mn3x21G17Vv62hXZR",
	"TWPMTwFxwYIAiSKDnujcoigCly3s7xTufYQ5Dq2k+pOGRbTSnTg9XbIBzfmUXu5q7Kaf3NeCG95iF4A2",
	"inSdeywY9nc4uKauhT30gmaiqhiOVB1XejWyUjIUo92FrXj5HofTy+6dLg5fDeNrnwvd2P68bbCbKGs+",
	"JfLnpqpOlME3REoEyoQ6Y+VVUjo+sq6klUX3qn/9iGf29e4QPiiuOhh5gyyOHjbU6b3CStzh0S3ISgez",
	"CJ81Y1auLcX/4wWXyHg+UnGsT9gE24lK0lmvsW9HHQtValZqywRzBnNZyvwmnVAStLsjz42tIIOt9yZL",
	"EmHkLN8JY2a5rpPdnD7V/g4AAP//V9D8o7kLAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
