package services

import (
	"context"
	"errors"

	"github.com/ernstvorsteveld/go-cv-cassandra/src/domain/model"
	"github.com/ernstvorsteveld/go-cv-cassandra/src/port/in"
	"github.com/ernstvorsteveld/go-cv-cassandra/src/port/out"
)

type TagServices struct {
	db out.ExperienceDbPort
}

func (c *InServices) ListTags(ctx context.Context, command *in.ListTagsCommand) (*model.Tags, error) {
	return nil, errors.New("not implemeted yet")
}

func (c *InServices) GetTagById(ctx context.Context, command *in.GetTagByIdCommand) (*model.Tag, error) {
	return nil, errors.New("not implemeted yet")
}
