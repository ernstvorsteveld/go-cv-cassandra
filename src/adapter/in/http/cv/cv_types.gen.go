// Package cv provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package cv

import (
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// CreateExperienceRequest defines model for CreateExperienceRequest.
type CreateExperienceRequest struct {
	Name string `json:"name"`
	Tags *[]Tag `json:"tags,omitempty"`
}

// Error Error response
type Error struct {
	Code      string             `json:"code"`
	Message   string             `json:"message"`
	RequestId openapi_types.UUID `json:"requestId"`
}

// Experience defines model for Experience.
type Experience struct {
	Id   *openapi_types.UUID `json:"id,omitempty"`
	Name string              `json:"name"`
}

// ExperiencePayload Payload for creating an experience
type ExperiencePayload struct {
	Name string `json:"name"`
}

// Experiences Array of experiences
type Experiences = []Experience

// MetricsResponse Metrics for Prometheus response
type MetricsResponse = string

// ObjectId ObjectId
type ObjectId struct {
	Id *openapi_types.UUID `json:"id,omitempty"`
}

// Tag A tag
type Tag struct {
	Tag string `json:"tag"`
}

// TagArray Array of tags
type TagArray = []Tag

// TagArrayType Payload for tags
type TagArrayType struct {
	Tags *[]Tag `json:"tags,omitempty"`
}

// TagResponse A single tag
type TagResponse struct {
	Id   openapi_types.UUID `json:"id"`
	Name string             `json:"name"`
}

// TagsResponse List of tags
type TagsResponse = []TagResponse

// ListExperiencesParams defines parameters for ListExperiences.
type ListExperiencesParams struct {
	// Limit Size of the page, maximum is 100, default is 25
	Limit *int32 `form:"limit,omitempty" json:"limit,omitempty"`

	// Page The page to return
	Page *int32 `form:"page,omitempty" json:"page,omitempty"`
}

// CreateExperienceJSONRequestBody defines body for CreateExperience for application/json ContentType.
type CreateExperienceJSONRequestBody = CreateExperienceRequest
