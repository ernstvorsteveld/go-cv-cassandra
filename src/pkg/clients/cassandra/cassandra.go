package cassandra

import (
	"strconv"

	"github.com/ernstvorsteveld/go-cv-cassandra/src/utils"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type CassandraSession struct {
	details *utils.CassandraConfiguration
	session *gocql.Session
}

func ConnectDatabase(c *utils.CassandraConfiguration) *CassandraSession {

	cluster := gocql.NewCluster(c.Url)
	cluster.Keyspace = c.Keyspace
	cluster.RetryPolicy = &gocql.SimpleRetryPolicy{
		NumRetries: int(c.Retries),
	}
	cluster.Consistency = gocql.Quorum
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: c.Username,
		Password: c.Secret.String(),
	}
	cluster.Port, _ = strconv.Atoi(c.Port)

	session, err := cluster.CreateSession()
	if err != nil {
		logrus.Fatal(err)
	}

	cs := &CassandraSession{c, session}
	return cs
}

type ExperienceDto struct {
	id   string
	name string
	tags []string
}

type ExperienceDao interface {
	Create(dto ExperienceDto) (*ExperienceDto, error)
	Get(id string) (ExperienceDto, error)
	GetPage(page int32, size int16) ([]ExperienceDto, error)
	Update(id string, dto ExperienceDto) error
	Delete(id string) (ExperienceDto, error)
}

const stmt_insert string = "INSERT INTO experiences(id,name,tags) VALUES(?,?,?)"

func (cc *CassandraSession) Create(dto ExperienceDto) (*ExperienceDto, error) {
	uuid := uuid.New()
	dto.id = uuid.String()
	if err := cc.session.Query(stmt_insert, dto.id, dto.name, dto.tags).Exec(); err != nil {
		return nil, err
	}

	return &dto, nil
}
