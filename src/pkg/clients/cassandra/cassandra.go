package cassandra

import (
	"strconv"

	"github.com/ernstvorsteveld/go-cv-cassandra/src/utils"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	log "github.com/labstack/gommon/log"
	"github.com/pkg/errors"
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
		log.Errorf("error: %c", err)
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

const stmt_select_by_id = "SELECT id, name, tags FROM experiences WHERE id = ?"

var QryErrorNotFound = errors.Errorf("Not Found")

func (cc *CassandraSession) Get(id string) (*ExperienceDto, error) {
	var _id string
	var name string
	var tags []string
	if err := cc.session.Query(stmt_select_by_id, id).Consistency(gocql.One).Scan(&_id, &name, &tags); err != nil {
		return nil, err
	}
	e := &ExperienceDto{id: _id, name: name, tags: tags}
	return e, nil
}
