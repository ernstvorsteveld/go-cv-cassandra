package services

import (
	"context"
	"testing"

	"github.com/ernstvorsteveld/go-cv-cassandra/adapter/out/db/mock"
	"github.com/ernstvorsteveld/go-cv-cassandra/domain/port/in"
	"github.com/stretchr/testify/assert"
)

func Test_should_create_experience(t *testing.T) {
	eDB := mock.NewMockExpDb()
	tDB := mock.NewMockTagDb()
	service := NewCvServices(eDB, tDB)

	e := in.NewCreateExperienceCommand("test-name", []string{"a", "b"})
	newExp, _ := service.CreateExperience(context.Background(), e)
	assert.Equal(t, newExp.Name, "test-name")
	assert.NotEmpty(t, newExp.Id)
	assert.Equal(t, len(eDB.Items), 1)
	assert.Equal(t, newExp.Id, eDB.Items[newExp.Id].GetId())
}

func Test_should_get_experience_by_id(t *testing.T) {
	eDB := mock.NewMockExpDb()
	tDB := mock.NewMockTagDb()
	service := NewCvServices(eDB, tDB)

	e := in.NewCreateExperienceCommand("test-name", []string{"a", "b"})
	newExp, _ := service.CreateExperience(context.Background(), e)
	assert.Equal(t, newExp.Name, "test-name")
	assert.NotEmpty(t, newExp.Id)
	exp, err := service.GetExperienceById(context.Background(), in.NewGetExperienceCommand(newExp.Id))
	assert.Nil(t, err, "Experience not found by id")
	assert.Equal(t, exp.Name, "test-name")
}
