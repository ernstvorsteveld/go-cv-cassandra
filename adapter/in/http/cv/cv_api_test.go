package cv

import (
	"context"
	"encoding/json"
	"errors"
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

var (
	c              *utils.Configuration
	ep             *MockExperienceDbPort
	handler        *CvApiHandler
	uid            uuid.UUID
	r              *gin.Engine
	experienceJson []byte
)

func readConfig() {
	c = &utils.Configuration{}
	c.Read("test_config", "yml")
}

func createHandler() {
	ep = new(MockExperienceDbPort)
	tp := new(MockTagDbPort)
	uid = uuid.New()
	h := services.NewCvServices(ep, tp, utils.NewMockUidGenerator(uid))
	handler = NewCvApiService(h, c)

	r = gin.Default()
	r.Use(MockCorrelationId)
	r.POST("/experiences", func(c *gin.Context) {
		handler.CreateExperience(c)
	})
}

func prepareData() {
	experience := CreateExperienceRequest{
		Name: "test-name",
		Tags: []string{"test-tag"},
	}
	experienceJson, _ = json.Marshal(experience)
}

func TestMain(m *testing.M) {
	readConfig()
	prepareData()

	m.Run()
}

func Test_should_create_experince(t *testing.T) {
	createHandler()

	gin.SetMode(gin.TestMode)
	ep.On("Create", mock.Anything, out.NewExperienceDto(uid.String(), "test-name", []string{"test-tag"})).Return(nil)

	req, err := http.NewRequest("POST", "/experiences", strings.NewReader(string(experienceJson)))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, uid.String(), rec.Body.String())
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func Test_should_fail_create_experince(t *testing.T) {
	createHandler()

	gin.SetMode(gin.TestMode)
	ep.On("Create", mock.Anything, out.NewExperienceDto(uid.String(), "test-name", []string{"test-tag"})).Return(errors.New("test-error"))

	req, err := http.NewRequest("POST", "/experiences", strings.NewReader(string(experienceJson)))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	e := &Error{}
	json.Unmarshal([]byte(rec.Body.String()), e)
	assert.Equal(t, "EXP0000004", e.Code)
	assert.Equal(t, "Internal Server error: error while creating experience", e.Message)
	assert.Equal(t, uuid.MustParse("05f4ae90-b8c9-4673-ab46-1f726e57932f"), e.RequestId)
}
