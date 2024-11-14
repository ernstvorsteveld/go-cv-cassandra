package cassandra

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	db "github.com/ernstvorsteveld/go-cv-cassandra/domain/port/out"
	"github.com/ernstvorsteveld/go-cv-cassandra/pkg/utils"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/cassandra"
)

var (
	cassandraContainer *cassandra.CassandraContainer
	session            *Session
	eDBPort            db.ExperienceDbPort
	tDBPort            db.TagDbPort
)

type CassandraExperienceDaoSuite struct {
	suite.Suite
}

type CassndarTagDaoSuite struct {
	suite.Suite
}

func Test_ExperienceDaoSuite(t *testing.T) {
	suite.Run(t, &CassandraExperienceDaoSuite{})
}

func Test_TagDaoSuite(t *testing.T) {
	suite.Run(t, &CassndarTagDaoSuite{})
}

func TestMain(m *testing.M) {
	log.Infof("Creating Cassandra Session in TestMain")
	ctx := context.Background()

	var err error
	cassandraContainer, err = cassandra.Run(ctx,
		"cassandra:4.1.3",
		testcontainers.CustomizeRequest(testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Name: "cassandratest",
			},
		}),
		cassandra.WithInitScripts(filepath.Join("testdata", "init.cql")),
	)

	if err != nil {
		log.Printf("failed to start container: %s", err)
	}

	rawPort, _ := cassandraContainer.MappedPort(ctx, "9042")
	parts := strings.Split(rawPort.Port(), "/")
	config := &utils.Configuration{
		DB: utils.DBConfiguration{
			Cassandra: utils.CassandraConfiguration{
				Keyspace: "testcv",
				Url:      "127.0.0.1",
				Port:     parts[0],
				Retries:  int8(3),
				Username: "cassandra",
				Secret:   utils.SensitiveInfo("cassandra"),
			}}}
	session = NewCassandraSession(config)
	log.Infof("Details: %v", session)
	eDBPort = NewExperiencePort(config, session)
	tDBPort = NewTagPort(config, session)

	m.Run()

	cassandraContainer.Terminate(ctx)
}

func (s *CassandraExperienceDaoSuite) Test_should_create_one_experience() {
	t := s.T()
	id := uuid.NewString()
	name := "value1"
	tags := []string{"ab", "ac"}

	err := eDBPort.Create(context.Background(), db.NewExperienceDto(id, name, tags))
	if err != nil {
		log.Printf("failed to start container: %s", err)
	}

	m := map[string]interface{}{}
	q := `SELECT * from testcv.cv_experiences;`
	itr := session.cs.Query(q).Iter()
	errors := true
	for itr.MapScan(m) {
		assert.Equal(t, m["id"].(string), id)
		assert.Equal(t, m["name"].(string), name)
		assert.Equal(t, m["tags"].([]string), tags)
		errors = false
	}

	assert.False(t, errors)
}

func insertExperience() (string, string, []string) {
	q := `INSERT INTO testcv.cv_experiences(id, name, tags) VALUES (?, ?, ?)`

	id := uuid.New().String()
	name := "value1"
	tags := []string{"ab", "ac"}
	session.cs.Query(q, id, name, tags).Exec()
	return id, name, tags
}

func (s *CassandraExperienceDaoSuite) Test_should_get_one_experience() {
	t := s.T()
	id, name, tags := insertExperience()

	d, err := eDBPort.Get(context.Background(), id)
	assert.Nil(t, err)
	assert.Equal(t, id, d.GetId())
	assert.Equal(t, name, d.GetName())
	assert.Equal(t, tags, d.GetTags())

	log.Infof("Experience: %v", d)
}

func (s *CassandraExperienceDaoSuite) Test_should_update_one_experience() {
	t := s.T()
	id, _, _ := insertExperience()

	name := "updated-value"
	tags := []string{"aaa", "bbb", "ccc", "ddd"}

	err := eDBPort.Update(context.Background(), id, db.NewExperienceDto(id, name, tags))
	if err != nil {
		log.Errorf("error while updating %v", err)
	}

	const q = "SELECT id, name, tags from testcv.cv_experiences WHERE id = ?"
	m := map[string]interface{}{}
	itr := session.cs.Query(q, id).Iter()
	errors := true
	for itr.MapScan(m) {
		assert.Equal(t, m["id"].(string), id)
		assert.Equal(t, m["name"].(string), name)
		assert.Equal(t, m["tags"].([]string), tags)
		errors = false
	}
	assert.False(t, errors)
}

func (s *CassandraExperienceDaoSuite) Test_should_not_find_experience_by_id() {
	t := s.T()
	_, err := eDBPort.Get(context.Background(), "not-existing-id")
	assert.Equal(t, "not found", err.Error())
}

func (s *CassndarTagDaoSuite) Test_should_create_tag() {
	t := s.T()
	id := uuid.NewString()
	name := "tag-name-value"

	dto, err := tDBPort.Create(context.Background(), db.NewTagDto(id, name))
	if err != nil {
		log.Printf("failed to start container: %s", err)
	}

	m := map[string]interface{}{}
	q := `SELECT * from testcv.cv_tags;`
	itr := session.cs.Query(q).Iter()
	errors := true
	for itr.MapScan(m) {
		assert.Equal(t, m["id"].(string), dto.GetId())
		assert.Equal(t, m["name"].(string), dto.GetName())
		errors = false
	}

	assert.False(t, errors)
}

func insertTag() (string, string) {
	q := `INSERT INTO testcv.cv_tags(id, name) VALUES (?, ?)`

	id := uuid.New().String()
	name := "test-value-for-get-by-id"
	session.cs.Query(q, id, name).Exec()
	return id, name
}

func (s *CassndarTagDaoSuite) Test_should_get_tag_by_id() {
	t := s.T()
	id, name := insertTag()

	d, err := tDBPort.Get(context.Background(), id)
	assert.Nil(t, err)
	assert.Equal(t, id, d.GetId())
	assert.Equal(t, name, d.GetName())
}

func (s *CassndarTagDaoSuite) Test_should_delete_tag_by_id() {
	t := s.T()
	id, name := insertTag()

	d, err := tDBPort.Delete(context.Background(), id)
	assert.Nil(t, err)
	assert.Equal(t, id, d.GetId())
	assert.Equal(t, name, d.GetName())
}
