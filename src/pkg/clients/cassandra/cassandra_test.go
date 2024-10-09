package cassandra

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ernstvorsteveld/go-cv-cassandra/src/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/cassandra"

	log "github.com/sirupsen/logrus"
)

var cassandraContainer *cassandra.CassandraContainer
var ctx context.Context
var session *CassandraSession

func TestMain(m *testing.M) {
	log.Infof("Creating Cassandra Session in TestMain")
	ctx = context.Background()

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
	session = ConnectDatabase(&utils.CassandraConfiguration{
		Keyspace: "testcv",
		Url:      "127.0.0.1",
		Port:     parts[0],
		Retries:  int8(3),
		Username: "cassandra",
		Secret:   utils.SensitiveInfo("cassandra"),
	})

	log.Infof("Details: %v", session.details)

	m.Run()

	cassandraContainer.Terminate(ctx)
}

func Test_should_create_one_experience(t *testing.T) {
	name := "value1"
	tags := []string{"ab", "ac"}

	d, err := session.Create(&ExperienceDto{
		name: name,
		tags: tags,
	})
	if err != nil {
		log.Printf("failed to start container: %s", err)
	}

	m := map[string]interface{}{}
	q := `SELECT * from testcv.experiences;`
	itr := session.session.Query(q).Iter()
	errors := true
	for itr.MapScan(m) {
		assert.Equal(t, m["id"].(string), d.id)
		assert.Equal(t, m["name"].(string), name)
		assert.Equal(t, m["tags"].([]string), tags)
		errors = false
	}
	assert.False(t, errors)
}

func insertOne() (string, string, []string) {
	q := `INSERT INTO testcv.experiences(id, name, tags) VALUES (?, ?, ?)`

	id := uuid.New().String()
	name := "value1"
	tags := []string{"ab", "ac"}
	session.session.Query(q, id, name, tags).Exec()
	return id, name, tags
}

func Test_should_get_one_experience(t *testing.T) {
	id, name, tags := insertOne()

	d, err := session.Get(id)
	assert.Nil(t, err)
	assert.Equal(t, id, d.id)
	assert.Equal(t, name, d.name)
	assert.Equal(t, tags, d.tags)

	log.Infof("Experience: %v", d)
}

func Test_should_update_one_experience(t *testing.T) {
	id, _, _ := insertOne()

	name := "updated-value"
	tags := []string{"aaa", "bbb", "ccc", "ddd"}

	err := session.Update(id, &ExperienceDto{id: id, name: name, tags: tags})
	if err != nil {
		log.Errorf("error while updating %v", err)
	}

	const q = "SELECT id, name, tags from testcv.experiences WHERE id = ?"
	m := map[string]interface{}{}
	itr := session.session.Query(q, id).Iter()
	errors := true
	for itr.MapScan(m) {
		assert.Equal(t, m["id"].(string), id)
		assert.Equal(t, m["name"].(string), name)
		assert.Equal(t, m["tags"].([]string), tags)
		errors = false
	}
	assert.False(t, errors)
}
