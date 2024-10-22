package model

import "github.com/google/uuid"

type Error struct {
	Code      string
	Message   string
	RequestId uuid.UUID
}

var ErrorCodes = map[string]string{
	"EXP0000001": "Error while retrieving experience by id.",
	"EXP0000002": "Error while retrieving a page of experiences.",
	"EXP0000100": "Error while marshalling experiences.",
	"EXP0000500": "Error while creating a single experience.",
	"EXP0000501": "Error while updating a single experience.",

	"TAG0000001": "Error while retrieving tag by id.",
	"TAG0000002": "Error while retrieving a page of tags.",
	"TAG0000100": "Error while marshalling tags.",
	"TAG0000500": "Error while creating a single tag.",
	"TAG0000501": "Error while updating a single tag.",
}
