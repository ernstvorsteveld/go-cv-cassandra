package out

import "context"

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

type ExperienceDbPort interface {
	Create(ctx context.Context, dto *ExperienceDto) (*ExperienceDto, error)
	Get(ctx context.Context, id string) (*ExperienceDto, error)
	GetPage(ctx context.Context, page int32, size int16) ([]ExperienceDto, error)
	Update(ctx context.Context, id string, dto *ExperienceDto) error
	Delete(ctx context.Context, id string) (*ExperienceDto, error)
}
