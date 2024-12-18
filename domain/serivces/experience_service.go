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

func (c *InServices) ListExperiences(ctx context.Context, command *in.ListExperienceCommand) (*[]model.Experience, error) {
	slog.Debug("serivces.ListExperiences", "content", "About to List Experience", "command", command, "correlationId", utils.GetCorrelationId(ctx))
	c.ep.GetPage(ctx, NewGetParams(command))
	return nil, errors.New("not implemeted yet")
}

func NewGetParams(command *in.ListExperienceCommand) *out.GetParams {
	return &out.GetParams{
		Limit: command.Limit,
		Page:  command.Page,
		Tag:   command.Tag,
		Name:  command.Name}
}

func (c *InServices) CreateExperience(ctx context.Context, command *in.CreateExperienceCommand) (*model.Experience, error) {
	slog.Debug("serivces.CreateExperience", "content", "About to Create Experience", "correlationId", utils.GetCorrelationId(ctx))
	dto := out.NewExperienceDto(c.ig.UUIDString(), command.Name, command.Tags)
	err := c.ep.Create(ctx, dto)
	if err != nil {
		slog.Info("serivces.CreateExperience", "content", "Error while creating experience", "correlationId", utils.GetCorrelationId(ctx), "error", err.Error())
		return nil, err
	}
	return model.NewExperience(dto.GetId(), dto.GetName(), dto.GetTags())
}

func (c *InServices) GetExperienceById(ctx context.Context, command *in.GetExperienceCommand) (*model.Experience, error) {
	slog.Debug("serivces.GetExperienceById", "content", "About to Get Experience by Id", "id", command.Id, "correlationId", utils.GetCorrelationId(ctx))
	dto, err := c.ep.Get(ctx, command.Id)
	if err != nil {
		return nil, err
	}
	e, _ := model.NewExperience(dto.GetId(), dto.GetName(), dto.GetTags())
	return e, nil
}
