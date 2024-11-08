package utils_mock

import "github.com/google/uuid"

type MockUuidGenerator struct {
	Uid uuid.UUID
}

func NewMockUuidGenerator(uid uuid.UUID) *MockUuidGenerator {
	return &MockUuidGenerator{
		Uid: uid,
	}
}

func (g *MockUuidGenerator) UUIDString() string {
	return g.Uid.String()
}

func (g *MockUuidGenerator) New() uuid.UUID {
	return g.Uid
}
