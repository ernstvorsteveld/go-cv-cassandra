package neo4j_integtest

import (
	"context"
	"fmt"
	"testing"

	"log"
	"strings"

	"github.com/ernstvorsteveld/go-cv-cassandra/domain/port/out"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	n "github.com/ernstvorsteveld/go-cv-cassandra/adapter/out/db/neo4j"
	"github.com/ernstvorsteveld/go-cv-cassandra/pkg/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/testcontainers/testcontainers-go"
	tc "github.com/testcontainers/testcontainers-go/modules/neo4j"
)

type Neo4jExperienceDaoSuite struct {
	suite.Suite
}

var (
	Parts        []string
	TestPassword string = "letmein!"
	Driver       neo4j.DriverWithContext
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	neo4jContainer, err := tc.Run(ctx,
		"neo4j:latest",
		tc.WithAdminPassword(TestPassword),
		tc.WithLabsPlugin(tc.Apoc),
		tc.WithNeo4jSetting("dbms.tx_log.rotation.size", "42M"),
	)
	defer func() {
		if err := testcontainers.TerminateContainer(neo4jContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()

	if err != nil {
		log.Printf("failed to start container: %s", err)
	}

	rawPort, _ := neo4jContainer.MappedPort(ctx, "7687")
	Parts = strings.Split(rawPort.Port(), "/")

	Driver = n.NewNeo4jConnection(&utils.Configuration{
		DB: utils.DBConfiguration{
			Neo4j: utils.Neo4jConfiguration{
				Url:      "neo4j://localhost",
				Port:     Parts[0],
				Username: "neo4j",
				Secret:   utils.SensitiveInfo(TestPassword),
			},
		},
	})

	m.Run()
}

func Test_Neo4jExperienceDaoSuite(t *testing.T) {
	suite.Run(t, &Neo4jExperienceDaoSuite{})
}

func (s *Neo4jExperienceDaoSuite) Test_create_experience() {
	ctx := context.Background()
	session := (Driver).NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	var err error
	defer func() {
		err = session.Close(ctx)
	}()
	assert.Nil(s.T(), err, "failed to create session")

	port := n.NewExperiencePort(
		&utils.Configuration{
			DB: utils.DBConfiguration{
				Neo4j: utils.Neo4jConfiguration{
					Url: "bolt://localhost:7687",
				},
			},
		},
		Driver,
	)

	neo4j.ExecuteQuery(ctx, Driver,
		"CREATE (p:Tag {id : $id, name: $name}) RETURN p",
		map[string]any{
			"id":   uuid.NewString(),
			"name": "tag-1",
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"))

	neo4j.ExecuteQuery(ctx, Driver,
		"CREATE (p:Tag {id : $id, name: $name}) RETURN p",
		map[string]any{
			"id":   uuid.NewString(),
			"name": "tag-2",
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"))

	uid := uuid.NewString()
	err = port.Create(ctx, out.NewExperienceDto(uid, "name-1", []string{"tag-1", "tag-2"}))
	assert.Nil(s.T(), err, "failed to create experience")

	dto, err := port.Get(ctx, uid)
	assert.Nil(s.T(), err, "failed to get experience")
	assert.Equal(s.T(), dto.GetName(), "name-1", "experience name is not equal")
}

func (s *Neo4jExperienceDaoSuite) Test_get_experience() {
	ctx := context.Background()
	session := (Driver).NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	var err error
	defer func() {
		err = session.Close(ctx)
	}()
	assert.Nil(s.T(), err, "failed to create session")

	port := n.NewExperiencePort(
		&utils.Configuration{
			DB: utils.DBConfiguration{
				Neo4j: utils.Neo4jConfiguration{
					Url: "bolt://localhost:7687",
				},
			},
		},
		Driver,
	)

	_, err = port.Get(ctx, "test")
}

func Test_create_connection(t *testing.T) {
	ctx := context.Background()
	session := (Driver).NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})

	var err error
	defer func() {
		err = session.Close(ctx)
	}()
	assert.Nil(t, err, "failed to create session")

	result, err := neo4j.ExecuteQuery(ctx, Driver,
		"CREATE (p:Person {name: $name}) RETURN p",
		map[string]any{
			"name": "Doe",
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"))
	if err != nil {
		panic(err)
	}

	summary := result.Summary
	fmt.Printf("Created %v nodes in %+v.\n",
		summary.Counters().NodesCreated(),
		summary.ResultAvailableAfter())

	result, err = neo4j.ExecuteQuery(ctx, Driver,
		"MATCH (p:Person) RETURN p.name AS name",
		nil,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"))
	if err != nil {
		panic(err)
	}

	for _, record := range result.Records {
		name, _ := record.Get("name")
		fmt.Println(name)
	}

	fmt.Printf("The query `%v` returned %v records in %+v.\n",
		result.Summary.Query().Text(), len(result.Records),
		result.Summary.ResultAvailableAfter())
}
