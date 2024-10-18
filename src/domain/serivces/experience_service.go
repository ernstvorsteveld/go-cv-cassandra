package services

import (
	"context"
	"errors"

	"github.com/ernstvorsteveld/go-cv-cassandra/src/adapter/db"
	"github.com/ernstvorsteveld/go-cv-cassandra/src/domain/model"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type CvServices struct {
	db db.ExperienceDbAdapter
}

func NewCvService(db db.ExperienceDbAdapter) ExperienceServices {
	return &CvServices{
		db: db,
	}
}

type ExperienceServices interface {
	ListExperiences(ctx context.Context, page int, size int) (*[]model.Experience, error)
	CreateExperience(ctx context.Context, e model.Experience) (*model.Experience, error)
	GetExperienceById(ctx context.Context, id string) (*model.Experience, error)
}

func (c *CvServices) ListExperiences(ctx context.Context, page int, size int) (*[]model.Experience, error) {
	log.Debugf("About to list Experiences, page %d size %d", page, size)
	if page < 0 {
	}

	return nil, errors.New("not implemeted yet")
}

func (c *CvServices) CreateExperience(ctx context.Context, e model.Experience) (*model.Experience, error) {
	dto := db.NewExperienceDto(uuid.NewString(), e.GetName(), e.GetTags())
	dto, err := c.db.Create(context.Background(), dto)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return model.NewExperience(dto.GetId(), dto.GetName(), dto.GetTags())
}

func (c *CvServices) GetExperienceById(ctx context.Context, id string) (*model.Experience, error) {
	dto, err := c.db.Get(context.Background(), id)
	if err != nil {
		return nil, err
	}

	e, _ := model.NewExperience(dto.GetId(), dto.GetName(), dto.GetTags())
	return e, nil
}
