package out

import (
	"context"
	"net/url"
)

type ExperienceDto struct {
	id   string
	name string
	tags []string
}

type IExperienceDto interface {
	GetId() string
	GetName() string
	GetTags() []string
}

func (e *ExperienceDto) GetId() string {
	return e.id
}

func (e *ExperienceDto) GetName() string {
	return e.name
}

func (e *ExperienceDto) GetTags() []string {
	return e.tags
}

func NewExperienceDto(id string, name string, tags []string) *ExperienceDto {
	return &ExperienceDto{
		id:   id,
		name: name,
		tags: tags,
	}
}

type GetParams struct {
	Limit *int32
	Page  *string
	Tag   *string
	Name  *string
}

type ExperiencePageReslt struct {
	Next *url.URL
	Prev *url.URL
	Data []ExperienceDto
}

func NewExperiencePageReslt(next *url.URL, prev *url.URL, data []ExperienceDto) *ExperiencePageReslt {
	return &ExperiencePageReslt{
		Next: next,
		Prev: prev,
		Data: data,
	}
}

type ExperienceDbPort interface {
	Create(ctx context.Context, dto *ExperienceDto) error
	Get(ctx context.Context, id string) (*ExperienceDto, error)
	GetPage(ctx context.Context, params *GetParams) (*ExperiencePageReslt, error)
	Update(ctx context.Context, id string, dto *ExperienceDto) error
	Delete(ctx context.Context, id string) (*ExperienceDto, error)
}
