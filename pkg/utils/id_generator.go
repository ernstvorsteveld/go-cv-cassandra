package utils

import "github.com/google/uuid"

type IdGenerator interface {
	UUIDString() string
	UUID() uuid.UUID
}

type DefaultUuidGenerator struct {
}

func NewDefaultUuidGenerator() *DefaultUuidGenerator {
	return &DefaultUuidGenerator{}
}

func (g *DefaultUuidGenerator) UUIDString() string {
	return uuid.New().String()
}

func (g *DefaultUuidGenerator) UUID() uuid.UUID {
	return uuid.New()
}

type MockUidGenerator struct {
	value uuid.UUID
}

func NewMockUidGenerator(value uuid.UUID) *MockUidGenerator {
	return &MockUidGenerator{
		value: value,
	}
}

func (g *MockUidGenerator) UUIDString() string {
	return g.value.String()
}

func (g *MockUidGenerator) UUID() uuid.UUID {
	return g.value
}
