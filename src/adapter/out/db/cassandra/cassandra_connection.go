package cassandra

import (
	"strconv"

	"github.com/ernstvorsteveld/go-cv-cassandra/src/utils"
	"github.com/gocql/gocql"
	log "github.com/sirupsen/logrus"
)

type Session struct {
	cs *gocql.Session
}

func NewCassandraSession(c *utils.Configuration) *Session {
	cluster := gocql.NewCluster(c.DB.Cassandra.Url)
	cluster.Keyspace = c.DB.Cassandra.Keyspace
	cluster.RetryPolicy = &gocql.SimpleRetryPolicy{
		NumRetries: int(c.DB.Cassandra.Retries),
	}
	cluster.Consistency = gocql.Quorum
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: c.DB.Cassandra.Username,
		Password: c.DB.Cassandra.Secret.String(),
	}
	cluster.Port, _ = strconv.Atoi(c.DB.Cassandra.Port)

	session, err := cluster.CreateSession()
	if err != nil {
		log.Errorf("error: %c", err)
	}
	return &Session{
		cs: session,
	}
}
