package domain

import (
	"errors"

	"github.com/ernstvorsteveld/go-cv-cassandra/src/clients/cassandra"
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

func (c *CvServices) GetExperienceById(id string) (*Experience, error) {
	dto, err := c.cs.Get(id)
	if err != nil {
		return nil, err
	}

	e, _ := NewExperience(dto.GetId(), dto.GetName(), dto.GetTags())
	return e, nil
}
