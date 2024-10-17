package cassandra

import (
	"context"
	"strconv"

	"github.com/ernstvorsteveld/go-cv-cassandra/src/clients/db"
	"github.com/ernstvorsteveld/go-cv-cassandra/src/utils"
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type CassandraSession struct {
	config *utils.Configuration
	cs     *gocql.Session
}

func NewCassandraConnection(c *utils.Configuration) db.ExperienceDbAdapter {
	return ConnectDatabase(c)
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

const stmt_insert string = "INSERT INTO experiences(id,name,tags) VALUES(?,?,?)"
const stmt_select_by_id string = "SELECT id, name, tags FROM experiences WHERE id = ?"

var QryErrorNotFound = errors.Errorf("Not Found")

func (cc *CassandraSession) Create(ctx context.Context, dto *db.ExperienceDto) (*db.ExperienceDto, error) {
	log.Debugf("About to Create Experience %v", dto)
	dto = db.NewExperienceDto(dto.GetId(), dto.GetName(), dto.GetTags())
	if err := cc.cs.Query(stmt_insert, dto.GetId(), dto.GetName(), dto.GetTags()).Exec(); err != nil {
		return nil, err
	}

	return dto, nil
}

func (cc *CassandraSession) Get(ctx context.Context, id string) (*db.ExperienceDto, error) {
	log.Debugf("About to Get Experience by id %s", id)
	var _id string
	var name string
	var tags []string
	if err := cc.cs.Query(stmt_select_by_id, id).Consistency(gocql.One).Scan(&_id, &name, &tags); err != nil {
		return nil, err
	}
	e := db.NewExperienceDto(_id, name, tags)
	return e, nil
}

const stmt_update string = "UPDATE experiences SET name = ?, tags = ? WHERE id = ?"

func (cc *CassandraSession) Update(ctx context.Context, id string, dto *db.ExperienceDto) error {
	log.Debugf("About to Update Experience with id %s with value %v", id, dto)
	if err := cc.cs.Query(stmt_update, dto.GetName(), dto.GetTags(), id).Exec(); err != nil {
		return err
	}
	return nil
}

func (cc *CassandraSession) GetPage(ctx context.Context, page int32, size int16) ([]db.ExperienceDto, error) {
	return nil, nil
}

func (cc *CassandraSession) Delete(ctx context.Context, id string) (*db.ExperienceDto, error) {
	return nil, nil
}
