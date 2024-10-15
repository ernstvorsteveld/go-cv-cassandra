package cassandra

import (
	"strconv"

	"github.com/ernstvorsteveld/go-cv-cassandra/src/utils"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type CassandraSession struct {
	config *utils.Configuration
	cs     *gocql.Session
}

func ConnectDatabase(c *utils.Configuration) *CassandraSession {
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

	cs := &CassandraSession{
		config: c,
		cs:     session,
	}
	return cs
}

type ExperienceDto struct {
	id   string
	name string
	tags []string
}

type IExperienceDto interface {
	GetId() string
	GetName() string
	GetTags() []string
}

func (e *ExperienceDto) GetId() string {
	return e.id
}

func (e *ExperienceDto) GetName() string {
	return e.name
}

func (e *ExperienceDto) GetTags() []string {
	return e.tags
}

func NewExperienceDto(id string, name string, tags []string) *ExperienceDto {
	return &ExperienceDto{
		id:   id,
		name: name,
		tags: tags,
	}
}

type ExperienceDao interface {
	Create(dto *ExperienceDto) (*ExperienceDto, error)
	Get(id string) (*ExperienceDto, error)
	GetPage(page int32, size int16) ([]ExperienceDto, error)
	Update(id string, dto *ExperienceDto) error
	Delete(id string) (*ExperienceDto, error)
}

const stmt_insert string = "INSERT INTO experiences(id,name,tags) VALUES(?,?,?)"
const stmt_select_by_id = "SELECT id, name, tags FROM experiences WHERE id = ?"

var QryErrorNotFound = errors.Errorf("Not Found")

func (cc *CassandraSession) Create(dto *ExperienceDto) (*ExperienceDto, error) {
	log.Debugf("About to Create Experience %v", dto)
	uuid := uuid.New()
	dto.id = uuid.String()
	if err := cc.cs.Query(stmt_insert, dto.id, dto.name, dto.tags).Exec(); err != nil {
		return nil, err
	}

	return dto, nil
}

func (cc *CassandraSession) Get(id string) (*ExperienceDto, error) {
	log.Debugf("About to Get Experience by id %s", id)
	var _id string
	var name string
	var tags []string
	if err := cc.cs.Query(stmt_select_by_id, id).Consistency(gocql.One).Scan(&_id, &name, &tags); err != nil {
		return nil, err
	}
	e := &ExperienceDto{id: _id, name: name, tags: tags}
	return e, nil
}

const stmt_update string = "UPDATE experiences SET name = ?, tags = ? WHERE id = ?"

func (cc *CassandraSession) Update(id string, dto *ExperienceDto) error {
	log.Debugf("About to Update Experience with id %s with value %v", id, dto)
	if err := cc.cs.Query(stmt_update, dto.name, dto.tags, id).Exec(); err != nil {
		return err
	}
	return nil
}
