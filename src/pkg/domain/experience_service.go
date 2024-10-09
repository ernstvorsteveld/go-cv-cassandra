package domain

import (
	"errors"

	"github.com/ernstvorsteveld/go-cv-cassandra/src/pkg/clients/cassandra"
	log "github.com/sirupsen/logrus"
)

type CvServices struct {
	cs *cassandra.CassandraSession
}

func NewCvService() (*CvServices, error) {
	return &CvServices{}, nil
}

type ExperienceServices interface {
	ListExperiences(page int, size int) (*[]Experience, error)
	CreateExperience(e Experience) (*Experience, error)
	GetExperienceById(id string) (*Experience, error)
}

func (c *CvServices) ListExperiences(page int, size int) (*[]Experience, error) {
	log.Debugf("About to list Experiences, page %d size %d", page, size)
	return nil, errors.New("not implemeted yet")
}

func (c *CvServices) CreateExperience(e Experience) (*Experience, error) {
	dto := cassandra.NewExperienceDto(e.GetId(), e.GetName(), e.GetTags())
	dto, err := c.cs.Create(dto)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return NewExperience(dto.GetId(), dto.GetName(), dto.GetTags())
}
