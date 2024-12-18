package in

import (
	"context"

	"github.com/ernstvorsteveld/go-cv-cassandra/domain/model"
)

type UseCasesPort interface {
	ListTags(ctx context.Context, c *ListTagsCommand) (*model.Tags, error)
	GetTagById(ctx context.Context, c *GetTagByIdCommand) (*model.Tag, error)
	CreateTag(ctx context.Context, c *CreateTagCommand) (*model.Tag, error)

	ListExperiences(ctx context.Context, c *ListExperienceCommand) (*[]model.Experience, error)
	CreateExperience(ctx context.Context, c *CreateExperienceCommand) (*model.Experience, error)
	GetExperienceById(ctx context.Context, c *GetExperienceCommand) (*model.Experience, error)
}
