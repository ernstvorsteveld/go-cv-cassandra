// Package gin provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package gin

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

	"H4sIAAAAAAAC/8RY32/bNhD+Vwhuj3IsO0u3CiiGNDMGD1kdpF5XLPMDK51ldhLJkpRn19D/PpCUrZ92",
	"5SVt8hDYEu/43d13H4/e4ZCngjNgWuFgh1W4gpTYjzcSiIbJRoCkwEK4h08ZKG1ekSSZLXHwsMPfS1ji",
	"AH83LN0MCx/DOYmvpSTb+VYAzr3Ti8t97sg24STC+cLDEahQUqEpZzjABQKkOQotOEQYgoMhzj08kZJL",
	"A7FuaR8jCUpwpgB7GDYkFQmYlSGPwCy5m/jmb4Q9nIJSJDZPrxkCa8vDMJMSIuxh6WBMIxzg0fjyh6sX",
	"P/40MB8a/8zzlz75YHAJyQVITUGVW+5KFLXdBdEapIH9cD34a7G7zB/8wcvF7kWOPaxNMgOstKQsNp4P",
	"WKvuumC3LCtxVG17heThJZcp0TjAWUY73Bf+qdk7eKjs5bnoS+CLgy3/8BFCbctYFrU33WbWehr9H6p5",
	"/ZjcwUmT6Q1VmrK4ScXWNi1aFi/QkkvHaOOlzukaVRlJLVkOr9Eb88TDmsTKpFmTeOS+jvGiRTtnXqNd",
	"29PJMloPpwum2mHa3CG+rMSlqoE97DCN+jPvzCwsPEw1pBZXP1rYtiKbqbMa+f4hYOJY4OED2VrBHt40",
	"s0+/Tp/lHp6TuCPpSJO4muV9WuqotLNtrTpse7EbeVe+kZ6UbG6BxXqFgyv/S0wxfhcOm2ud46ywZavS",
	"oV7AnuUzScjbhaodQicbsAlj17Or3KrdV8Gft1rNPKJsyTvyeTe1cdxkUtIwS7IUvaOamMZIaAjm3Csl",
	"5Pfp3O5GtS160wZd302xh9cglXM+uvAvfGPBBTAiKA7wpX1kj6uVjW4IdRGIQbdR3oPOJFOIJElDDkxO",
	"iVllD9ZbqvSk9l4QSVLQIJVVjLrbt/QzWC6tAAkSg4dSsqFpliKqkOlgFMGSZIk238dX2CQRB/hTBnJb",
	"SkpCU6qxV0xBDry1wsH4qkIN8+XQIJTpy7FrD7NhXTAo0xCDNJVsQp4XUM04I21WjqAyi7pBjSqY7K7H",
	"Qb20fx24FqZv3Vxkqzb2fTeiMA3MTXtCJDS0pRl+VAZ7leyLKrJ++qocjZt6ZeKMENnLggBtyr4CEtma",
	"7/D7wc3s/n5yez2fzt4Mpr+02TWN9iSIQVf5hYoBpCaIK61FMBwSQS/CNRHiIuTpcD2q8njYU6XLFLRl",
	"8f3gzeT9vEuhE8r+MeU3gBlstOMDX6KyImfi/dl4eDX+O/P98QtFP8MrS/bT8O7uJ++msz/efgmikLCm",
	"PFNPBXPUH6bjS0H7/uz8RtN9zwawd5MO6mfMpCbUEDk4tiwqS1Mit4UUdujl/nSq9pU9nbjq0N2bjgsT",
	"+pfqlesWugaGRDGnNrW4eRMss/SaR9vz6nH2CNsvt8cuqx3ZrmytOSLRvuhuctEyg7wliaOOK2Uq9LZ6",
	"pawI1Y1Lx+CWu0S0rW8rbeXustGRqf/bqVR/cQ25lJC40LzzR9rTMGavf5vczDsRFLP17IDDkfpI4h6N",
	"5qlEZ1wTnaITl4QmzyU4f66A7ZNncDzinK1Q4WhVOg/fJ6hOTSW7BO6oSuZebVod7miUf3FkNbEdF0/a",
	"1s1foTLCvt66O+GpKdaMhPSQWqgJlQQtKaxhPySasbucEWlbxaq5OzvtzWQ/ckr8qhf88wfQY8eCPYL3",
	"cdrTAa1JQqMKgZ+uTeoTKvqwRbWNHtEpj672UwnfZV34eJZEiHGNlpQ1jrtnEMFTcL5NlYtp77lrXZPR",
	"KVty+zMCQUpASJc07CWoxgnI9V7XMpkcH13sz6iFn46JUOF8kf8XAAD//5jybP+ZGAAA",
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
