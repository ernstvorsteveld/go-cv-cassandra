package domain

import (
	"context"
	"errors"

	"github.com/ernstvorsteveld/go-cv-cassandra/src/clients/db"
	"github.com/ernstvorsteveld/go-cv-cassandra/src/clients/db/cassandra"
	"github.com/ernstvorsteveld/go-cv-cassandra/src/utils"
	log "github.com/sirupsen/logrus"
)

type CvServices struct {
	cs *cassandra.CassandraSession
}

func NewCvService(c *utils.Configuration) *CvServices {
	return &CvServices{
		cs: cassandra.ConnectDatabase(c),
	}
}

type ExperienceDbAdapter interface {
	Create(dto *db.ExperienceDto) (*db.ExperienceDto, error)
	Get(id string) (*db.ExperienceDto, error)
	List(page int, size int) (*[]Experience, error)
}

type ExperienceServices interface {
	ListExperiences(ctx context.Context, page int, size int) (*[]Experience, error)
	CreateExperience(ctx context.Context, e Experience) (*Experience, error)
	GetExperienceById(ctx context.Context, id string) (*Experience, error)
}

func (c *CvServices) ListExperiences(ctx context.Context, page int, size int) (*[]Experience, error) {
	log.Debugf("About to list Experiences, page %d size %d", page, size)
	if page < 0 {
	}

	return nil, errors.New("not implemeted yet")
}

func (c *CvServices) CreateExperience(ctx context.Context, e Experience) (*Experience, error) {
	dto := db.NewExperienceDto(e.GetId(), e.GetName(), e.GetTags())
	dto, err := c.cs.Create(context.Background(), dto)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return NewExperience(dto.GetId(), dto.GetName(), dto.GetTags())
}

func (c *CvServices) GetExperienceById(ctx context.Context, id string) (*Experience, error) {
	dto, err := c.cs.Get(context.Background(), id)
	if err != nil {
		return nil, err
	}

	e, _ := NewExperience(dto.GetId(), dto.GetName(), dto.GetTags())
	return e, nil
}
