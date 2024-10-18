package services

import (
	"context"
	"errors"

	"github.com/ernstvorsteveld/go-cv-cassandra/src/domain/model"
	"github.com/ernstvorsteveld/go-cv-cassandra/src/port/in"
	"github.com/ernstvorsteveld/go-cv-cassandra/src/port/out"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type CvServices struct {
	db out.ExperienceDbPort
}

func NewCvService(db out.ExperienceDbPort) in.ExperienceUseCases {
	return &CvServices{
		db: db,
	}
}

func (c *CvServices) ListExperiences(ctx context.Context, command *in.ListExperienceCommand) (*[]model.Experience, error) {
	log.Debugf("About to list Experiences, page %d size %d", command.Page, command.Size)
	if command.Page < 0 {
	}

	return nil, errors.New("not implemeted yet")
}

func (c *CvServices) CreateExperience(ctx context.Context, command *in.CreateExperienceCommand) (*model.Experience, error) {
	dto := out.NewExperienceDto(uuid.NewString(), command.Name, command.Tags)
	dto, err := c.db.Create(context.Background(), dto)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return model.NewExperience(dto.GetId(), dto.GetName(), dto.GetTags())
}

func (c *CvServices) GetExperienceById(ctx context.Context, command *in.GetExperienceCommand) (*model.Experience, error) {
	dto, err := c.db.Get(context.Background(), command.Id)
	if err != nil {
		return nil, err
	}

	e, _ := model.NewExperience(dto.GetId(), dto.GetName(), dto.GetTags())
	return e, nil
}
