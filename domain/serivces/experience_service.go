package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/ernstvorsteveld/go-cv-cassandra/domain/model"
	"github.com/ernstvorsteveld/go-cv-cassandra/port/in"
	"github.com/ernstvorsteveld/go-cv-cassandra/port/out"
	"github.com/ernstvorsteveld/go-cv-cassandra/utils"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func (c *InServices) ListExperiences(ctx context.Context, command *in.ListExperienceCommand) (*[]model.Experience, error) {
	log.Debugf("About to list Experiences, page %d size %d", command.Page, command.Size)
	if command.Page < 0 {
	}

	return nil, errors.New("not implemeted yet")
}

func (c *InServices) CreateExperience(ctx context.Context, command *in.CreateExperienceCommand) (*model.Experience, error) {
	dto := out.NewExperienceDto(uuid.NewString(), command.Name, command.Tags)
	dto, err := c.ep.Create(context.Background(), dto)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return model.NewExperience(dto.GetId(), dto.GetName(), dto.GetTags())
}

func (c *InServices) GetExperienceById(ctx context.Context, command *in.GetExperienceCommand) (*model.Experience, error) {
	slog.Debug("serivces.GetExperienceById", "content", "About to Get Experience by Id", "correctId", utils.Get("correlationId", ctx))
	dto, err := c.ep.Get(context.Background(), command.Id)
	if err != nil {
		return nil, err
	}

	e, _ := model.NewExperience(dto.GetId(), dto.GetName(), dto.GetTags())
	return e, nil
}
