package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/ernstvorsteveld/go-cv-cassandra/domain/model"
	"github.com/ernstvorsteveld/go-cv-cassandra/domain/port/in"
	"github.com/ernstvorsteveld/go-cv-cassandra/domain/port/out"
	"github.com/ernstvorsteveld/go-cv-cassandra/pkg/utils"
)

type TagServices struct {
	expDB out.ExperienceDbPort
	tagEb out.TagDbPort
}

func (c *InServices) ListTags(ctx context.Context, command *in.ListTagsCommand) (*model.Tags, error) {
	return nil, errors.New("not implemeted yet")
}

func (c *InServices) GetTagById(ctx context.Context, command *in.GetTagByIdCommand) (*model.Tag, error) {
	return nil, errors.New("not implemeted yet")
}

func (c *InServices) CreateTag(ctx context.Context, command *in.CreateTagCommand) (*model.Tag, error) {
	slog.Debug("serivces.CreateTag", "content", "About to Create Tag", "correlationId", utils.GetCorrelationId(ctx))
	dto := out.NewTagDto(c.ig.UUIDString(), command.Name)
	_, err := c.tp.Create(ctx, dto)
	if err != nil {
		slog.Info("serivces.CreateTag", "content", "Error while creating tag", "correlationId", utils.GetCorrelationId(ctx), "error", err.Error())
		return nil, err
	}
	return model.NewTag(dto.GetId(), dto.GetName())
}
