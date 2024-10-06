package domain

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Experience struct {
	Name     string `validate:"required,min=5,max=100"`
	Tags     []string
	Retries  int8
	username string
	secret   string
}

func NewExperience(name string, tags []string, retries int8, username string, secret string) (Experience, error) {
	e := Experience{name, tags, retries, username, secret}
	validate := validator.New()
	err := validate.Struct(e)
	if err != nil {
		fmt.Println(err.Error())
		return Experience{"error", []string{}, 0, "error", "error"}, err
	}
	return e, err
}
