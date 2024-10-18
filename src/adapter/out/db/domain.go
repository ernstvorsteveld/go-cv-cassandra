package db

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
