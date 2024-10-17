package db

import "context"

type ExperienceDao interface {
	Create(ctx context.Context, dto *ExperienceDto) (*ExperienceDto, error)
	Get(ctx context.Context, id string) (*ExperienceDto, error)
	GetPage(ctx context.Context, page int32, size int16) ([]ExperienceDto, error)
	Update(ctx context.Context, id string, dto *ExperienceDto) error
	Delete(ctx context.Context, id string) (*ExperienceDto, error)
}
