package utils

import "github.com/google/uuid"

type IdGenerator interface {
	UUIDString() string
	New() uuid.UUID
}

type DefaultUuidGenerator struct {
}

func NewDefaultUuidGenerator() *DefaultUuidGenerator {
	return &DefaultUuidGenerator{}
}

func (g *DefaultUuidGenerator) UUIDString() string {
	return uuid.New().String()
}

func (g *DefaultUuidGenerator) New() uuid.UUID {
	return uuid.New()
}
