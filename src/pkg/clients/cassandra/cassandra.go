package cassandra

import (
	"github.com/gocql/gocql"
	"github.com/sirupsen/logrus"
)

type CassandraSession struct {
	details ConnectionDetails
	session *gocql.Session
}

type ConnectionDetails struct {
	url      string
	keyspace string
}

func ConnectDatabase(c ConnectionDetails) *CassandraSession {

	cluster := gocql.NewCluster(c.url)
	cluster.Keyspace = c.keyspace
	cluster.RetryPolicy = &gocql.SimpleRetryPolicy{
		NumRetries: c.retries,
	}
	cluster.Consistency = gocql.Quorum
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: cassandraUser,
		Password: cassandraPassword,
	}	

	session, err := cluster.CreateSession()
	if err != nil {
		logrus.Fatal(err)
	}

	cs := &CassandraSession{c, session}
	return cs
}

type ExperienceDto struct {
	name string
	tags []string
}

type ExperienceDao interface {
	Create(dto ExperienceDto) (string, error)
	Get(id string) (ExperienceDto, error)
	GetPage(page int32, size int16) ([]ExperienceDto, error)
	Update(id string, dto ExperienceDto) error
	Delete(id string) (ExperienceDto, error)
}

func (cc *CassandraSession) Create(dto ExperienceDto) (string, error) {
	cc.session.
}
