package mock

import (
	"context"
	"errors"

	"github.com/ernstvorsteveld/go-cv-cassandra/src/port/out"
	"github.com/google/uuid"
)

type MockDb struct {
	Items map[string]*out.ExperienceDto
}

func NewMockDb() *MockDb {
	return &MockDb{
		Items: make(map[string]*out.ExperienceDto),
	}
}

func (m *MockDb) Create(ctx context.Context, dto *out.ExperienceDto) (*out.ExperienceDto, error) {
	d := out.NewExperienceDto(uuid.NewString(), dto.GetName(), dto.GetTags())
	m.Items[d.GetId()] = d
	return d, nil
}

func (m *MockDb) Get(ctx context.Context, id string) (*out.ExperienceDto, error) {
	res, ok := m.Items[id]
	if !ok {
		return nil, errors.New("Not found")
	}
	return res, nil
}

func (m *MockDb) GetPage(ctx context.Context, page int32, size int16) ([]out.ExperienceDto, error) {
	return nil, nil
}

func (m *MockDb) Update(ctx context.Context, id string, dto *out.ExperienceDto) error {
	return nil
}
func (m *MockDb) Delete(ctx context.Context, id string) (*out.ExperienceDto, error) {
	return out.NewExperienceDto(uuid.NewString(), "mock-name", []string{"tag1", "tag2"}), nil
}
