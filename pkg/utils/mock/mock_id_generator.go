package mock

import "github.com/google/uuid"

type MockUuidGenerator struct {
	Uuid uuid.UUID
}

func NewMockUuidGenerator(uuid uuid.UUID) *MockUuidGenerator {
	return &MockUuidGenerator{
		Uuid: uuid,
	}
}

func (g *MockUuidGenerator) UUIDString() string {
	return g.Uuid.String()
}

func (g *MockUuidGenerator) UUID() uuid.UUID {
	return g.Uuid
}
