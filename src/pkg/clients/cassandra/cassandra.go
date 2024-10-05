package cassandra

import (
	"github.com/gocql/gocql"
	"github.com/sirupsen/logrus"
)

type ConnectionDetails struct {
	url      string
	keyspace string
}

func ConnectDatabase(c ConnectionDetails) *gocql.Session {

	cluster := gocql.NewCluster(c.url)
	cluster.Keyspace = c.keyspace

	session, err := cluster.CreateSession()
	if err != nil {
		logrus.Fatal(err)
	}
	return session
}
