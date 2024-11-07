package cv

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ernstvorsteveld/go-cv-cassandra/domain/port/out"
	services "github.com/ernstvorsteveld/go-cv-cassandra/domain/serivces"
	"github.com/google/uuid"

	"github.com/ernstvorsteveld/go-cv-cassandra/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockExperienceDbPort struct {
	mock.Mock
}

func (m *MockExperienceDbPort) Get(ctx context.Context, id string) (*out.ExperienceDto, error) {
	args := m.Called(ctx, id)
	return args.Get(1).(*out.ExperienceDto), args.Error(1)
}

func (m *MockExperienceDbPort) GetPage(ctx context.Context, page int32, size int16) ([]out.ExperienceDto, error) {
	return nil, nil
}

func (m *MockExperienceDbPort) Update(ctx context.Context, id string, dto *out.ExperienceDto) error {
	return nil
}

func (m *MockExperienceDbPort) Create(ctx context.Context, dto *out.ExperienceDto) error {
	slog.Info("Mock Create called with ctx %d and dto %w", ctx, dto)
	args := m.Called(ctx, dto)
	return args.Error(0)
}

func (m *MockExperienceDbPort) Delete(ctx context.Context, id string) (*out.ExperienceDto, error) {
	args := m.Called(ctx, id)
	return args.Get(1).(*out.ExperienceDto), args.Error(1)
}

type MockTagDbPort struct {
	mock.Mock
}

func (m *MockTagDbPort) Create(ctx context.Context, dto *out.TagDto) (*out.TagDto, error) {
	return nil, nil
}

func (m *MockTagDbPort) Get(ctx context.Context, id string) (*out.TagDto, error) {
	return nil, nil
}

func (m *MockTagDbPort) GetPage(ctx context.Context, page int32, size int16) ([]out.TagDto, error) {
	return nil, nil
}

func (m *MockTagDbPort) Update(ctx context.Context, id string, dto *out.TagDto) error {
	return nil
}

func (m *MockTagDbPort) Delete(ctx context.Context, id string) (*out.TagDto, error) {
	return nil, nil
}

func Test_should_create_experince(t *testing.T) {
	c := &utils.Configuration{}
	c.Read("test_config", "yml")

	ep := new(MockExperienceDbPort)
	tp := new(MockTagDbPort)
	uid := uuid.New()
	h := services.NewCvServices(ep, tp, utils.NewMockUidGenerator(uid))
	handler := NewCvApiService(h, c)

	ep.On("Create", mock.Anything, out.NewExperienceDto(uid.String(), "test-name", []string{"test-tag"})).Return(nil)

	r := gin.Default()
	r.POST("/experiences", func(c *gin.Context) {
		handler.CreateExperience(c)
	})

	experience := CreateExperienceRequest{
		Name: "test-name",
		Tags: []string{"test-tag"},
	}
	experienceJson, _ := json.Marshal(experience)
	req, err := http.NewRequest("POST", "/experiences", strings.NewReader(string(experienceJson)))
	assert.NoError(t, err)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, uid.String(), rec.Body.String())
	assert.Equal(t, http.StatusCreated, rec.Code)
}
