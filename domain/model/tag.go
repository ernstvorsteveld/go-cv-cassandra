package model

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
)

type Tag struct {
	Id   string `form:"id" json:"id" xml:"id"`
	Name string `form:"name" json:"name" xml:"name" validate:"required,min=2,max=100"`
}

type Tags []Tag

type ITag interface {
	GetId() string
	GetName() string
}

func (e *Tag) GetId() string {
	return e.Id
}

func (e *Tag) GetName() string {
	return e.Name
}

func NewTag(id string, name string) (*Tag, error) {
	slog.Debug("model.NewTag", "content", "About to construct Tag", "id", id, "name", name)
	validate := validator.New(validator.WithRequiredStructEnabled())

	t := &Tag{Id: id, Name: name}
	err := validate.Struct(t)
	if err != nil {
		slog.Info("model.NewTag", "content", "Could mot create Tag", "error", err.Error())
		return &Tag{"error", "error"}, err
	}
	return t, err
}
