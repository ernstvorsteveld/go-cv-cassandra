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

	m "github.com/ernstvorsteveld/go-cv-cassandra/pkg/middleware"
	"github.com/ernstvorsteveld/go-cv-cassandra/pkg/utils"
	utils_mock "github.com/ernstvorsteveld/go-cv-cassandra/pkg/utils/mock"
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
	tp             *MockTagDbPort
	ig             utils.IdGenerator
	handler        *CvApiHandler
	uid            uuid.UUID = uuid.New()
	r              *gin.Engine
	experienceJson []byte
)

func readConfig() {
	c = &utils.Configuration{}
	c.Read("test_config", "yml")
}

func expectMocks() {
	ig = utils_mock.NewMockUuidGenerator(uid)
	ep = new(MockExperienceDbPort)
	tp = new(MockTagDbPort)
}

func expectHandler() {
	h := services.NewCvServices(ep, tp, ig)
	handler = NewCvApiService(h, c)
}

func expectEngine() {
	r = gin.Default()
	r.Use(m.CorrelationId(ig))
	r.POST("/experiences", func(c *gin.Context) {
		handler.CreateExperience(c)
	})
}

func expectData() {
	experience := CreateExperienceRequest{
		Name: "test-name",
		Tags: []string{"test-tag"},
	}
	experienceJson, _ = json.Marshal(experience)
}

func TestMain(m *testing.M) {
	readConfig()
	expectMocks()
	expectData()

	m.Run()
}

func Test_should_create_experince(t *testing.T) {
	expectEngine()
	expectMocks()
	expectHandler()

	gin.SetMode(gin.TestMode)
	dto := out.NewExperienceDto(uid.String(), "test-name", []string{"test-tag"})
	ep.On("Create", mock.Anything, dto).Return(nil)

	reader := strings.NewReader(string(experienceJson))
	req, err := http.NewRequest("POST", "/experiences", reader)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, uid.String(), rec.Body.String())
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func Test_should_fail_create_experince(t *testing.T) {
	expectEngine()
	expectMocks()
	expectHandler()

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
	assert.Equal(t, "internal server error: error while creating experience", e.Message)
	assert.Equal(t, uid, e.RequestId)
}
