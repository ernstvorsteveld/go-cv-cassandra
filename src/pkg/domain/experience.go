package domain

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Experience struct {
	Id   string
	Name string `validate:"required,alphanum,min=5,max=100"`
	Tags []string
}

func NewExperience(id string, name string, tags []string) (*Experience, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	e := &Experience{Id: id, Name: name, Tags: tags}
	err := validate.Struct(e)
	if err != nil {
		fmt.Println(err.Error())
		return &Experience{"error", "error", []string{}}, err
	}
	return e, err
}

type IExperience interface {
	GetId() string
	GetName() string
	GetTags() []string
}

func (e *Experience) GetId() string {
	return e.Id
}

func (e *Experience) GetName() string {
	return e.Name
}

func (e *Experience) GetTags() []string {
	return e.Tags
}
