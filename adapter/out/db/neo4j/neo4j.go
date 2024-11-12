package neo4j

import (
	"github.com/ernstvorsteveld/go-cv-cassandra/pkg/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func NewNeo4jConnection(c *utils.Configuration) neo4j.DriverWithContext {
	dbUri := c.DB.Neo4j.Url + ":" + c.DB.Neo4j.Port
	driver, err := neo4j.NewDriverWithContext(dbUri,
		neo4j.BasicAuth(c.DB.Neo4j.Username, c.DB.Neo4j.Secret.Value(), ""))
	if err != nil {
		panic(err)
	}
	return driver
}
