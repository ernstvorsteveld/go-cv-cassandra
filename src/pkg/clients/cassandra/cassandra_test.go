package cassandra

import (
	"context"
	"log"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ernstvorsteveld/go-cv-cassandra/src/utils"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/cassandra"
)

func Test_should_create_experience(t *testing.T) {
	ctx := context.Background()

	cr := testcontainers.ContainerRequest{
		Name: "cassandratest",
	}
	req := testcontainers.GenericContainerRequest{
		ContainerRequest: cr,
	}

	cassandraContainer, err := cassandra.Run(ctx,
		"cassandra:4.1.3",
		testcontainers.CustomizeRequest(req),
		cassandra.WithInitScripts(filepath.Join("testdata", "init.cql")),
	)

	if err != nil {
		log.Printf("failed to start container: %s", err)
		return
	}

	rawPort, _ := cassandraContainer.MappedPort(ctx, "9042")
	parts := strings.Split(rawPort.Port(), "/")
	session := ConnectDatabase(&utils.CassandraConfiguration{
		Keyspace: "testcv",
		Url:      "127.0.0.1",
		Port:     parts[0],
		Retries:  int8(3),
		Username: "cassandra",
		Secret:   utils.SensitiveInfo("cassandra"),
	})

	log.Println(session.details)
	log.Println(session.details)

	d, err := session.Create(ExperienceDto{
		name: "example1",
		tags: []string{"a", "b"},
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
		errors = false
	}
	assert.False(t, errors)

	cassandraContainer.Terminate(ctx)
}
