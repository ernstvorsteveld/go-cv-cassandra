package model

type Tag struct {
	Id   string `form:"id" json:"id" xml:"id"`
	Name string `form:"name" json:"name" xml:"name" validate:"required,min=5,max=100"`
}

type Tags []Tag
