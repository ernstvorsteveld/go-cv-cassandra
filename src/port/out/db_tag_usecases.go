package out

import "context"

type TagDto struct {
	id   string
	name string
}

func (e *TagDto) GetId() string {
	return e.id
}

func (e *TagDto) GetName() string {
	return e.name
}

func NewTagDto(id string, name string) *TagDto {
	return &TagDto{
		id:   id,
		name: name,
	}
}

type TagDbPort interface {
	Create(ctx context.Context, dto *TagDto) (*TagDto, error)
	Get(ctx context.Context, id string) (*TagDto, error)
	GetPage(ctx context.Context, page int32, size int16) ([]TagDto, error)
	Update(ctx context.Context, id string, dto *TagDto) error
	Delete(ctx context.Context, id string) (*TagDto, error)
}
