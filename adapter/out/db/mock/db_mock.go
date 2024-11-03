package mock

import (
	"context"
	"errors"

	"github.com/ernstvorsteveld/go-cv-cassandra/domain/model"
	"github.com/ernstvorsteveld/go-cv-cassandra/domain/port/in"
	"github.com/ernstvorsteveld/go-cv-cassandra/domain/port/out"
	"github.com/google/uuid"
)

type MockExpDb struct {
	Items map[string]*out.ExperienceDto
}

func NewMockExpDb() *MockExpDb {
	return &MockExpDb{
		Items: make(map[string]*out.ExperienceDto),
	}
}

type MockTagDb struct {
	Items map[string]*model.Tag
}

func NewMockTagDb() *MockTagDb {
	return &MockTagDb{
		Items: make(map[string]*model.Tag),
	}
}

func (m *MockExpDb) Create(ctx context.Context, dto *out.ExperienceDto) (*out.ExperienceDto, error) {
	d := out.NewExperienceDto(uuid.NewString(), dto.GetName(), dto.GetTags())
	m.Items[d.GetId()] = d
	return d, nil
}

func (m *MockExpDb) Get(ctx context.Context, id string) (*out.ExperienceDto, error) {
	res, ok := m.Items[id]
	if !ok {
		return nil, errors.New("Not found")
	}
	return res, nil
}

func (m *MockExpDb) GetPage(ctx context.Context, page int32, size int16) ([]out.ExperienceDto, error) {
	return nil, nil
}

func (m *MockExpDb) Update(ctx context.Context, id string, dto *out.ExperienceDto) error {
	return nil
}
func (m *MockExpDb) Delete(ctx context.Context, id string) (*out.ExperienceDto, error) {
	return out.NewExperienceDto(uuid.NewString(), "mock-name", []string{"tag1", "tag2"}), nil
}

func (m *MockExpDb) ListTags(ctx context.Context, command *in.ListTagsCommand) (*model.Tags, error) {
	return nil, errors.New("not implemeted yet")
}

func (m *MockExpDb) GetTagById(ctx context.Context, command *in.GetTagByIdCommand) (*model.Tag, error) {
	return nil, errors.New("not implemeted yet")
}

func (c *MockTagDb) Create(ctx context.Context, dto *out.TagDto) (*out.TagDto, error) {
	return nil, nil
}

func (c *MockTagDb) Get(ctx context.Context, id string) (*out.TagDto, error) {
	return nil, nil
}

func (c *MockTagDb) GetPage(ctx context.Context, page int32, size int16) ([]out.TagDto, error) {
	return nil, nil
}

func (c *MockTagDb) Update(ctx context.Context, id string, dto *out.TagDto) error {
	return nil
}

func (c *MockTagDb) Delete(ctx context.Context, id string) (*out.TagDto, error) {
	return nil, nil
}
