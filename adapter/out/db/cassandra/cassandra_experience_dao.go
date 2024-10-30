package cassandra

import (
	"context"

	"github.com/ernstvorsteveld/go-cv-cassandra/port/out"
	"github.com/ernstvorsteveld/go-cv-cassandra/utils"
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type CassandraExperienceSession struct {
	config *utils.Configuration
	cs     *gocql.Session
}

func NewExperiencePort(c *utils.Configuration, s *Session) out.ExperienceDbPort {
	return &CassandraExperienceSession{
		config: c,
		cs:     s.cs,
	}
}

const stmt_insert string = "INSERT INTO experiences(id,name,tags) VALUES(?,?,?)"
const stmt_select_by_id string = "SELECT id, name, tags FROM experiences WHERE id = ?"

var QryErrorNotFound = errors.Errorf("Not Found")

func (cc *CassandraExperienceSession) Create(ctx context.Context, dto *out.ExperienceDto) (*out.ExperienceDto, error) {
	log.Debugf("About to Create Experience %v", dto)
	dto = out.NewExperienceDto(dto.GetId(), dto.GetName(), dto.GetTags())
	if err := cc.cs.Query(stmt_insert, dto.GetId(), dto.GetName(), dto.GetTags()).Exec(); err != nil {
		return nil, err
	}

	return dto, nil
}

func (cc *CassandraExperienceSession) Get(ctx context.Context, id string) (*out.ExperienceDto, error) {
	log.Debugf("About to Get Experience by id %s", id)
	var _id string
	var name string
	var tags []string
	if err := cc.cs.Query(stmt_select_by_id, id).Consistency(gocql.One).Scan(&_id, &name, &tags); err != nil {
		return nil, err
	}
	e := out.NewExperienceDto(_id, name, tags)
	return e, nil
}

const stmt_update string = "UPDATE experiences SET name = ?, tags = ? WHERE id = ?"

func (cc *CassandraExperienceSession) Update(ctx context.Context, id string, dto *out.ExperienceDto) error {
	log.Debugf("About to Update Experience with id %s with value %v", id, dto)
	if err := cc.cs.Query(stmt_update, dto.GetName(), dto.GetTags(), id).Exec(); err != nil {
		return err
	}
	return nil
}

func (cc *CassandraExperienceSession) GetPage(ctx context.Context, page int32, size int16) ([]out.ExperienceDto, error) {
	return nil, nil
}

func (cc *CassandraExperienceSession) Delete(ctx context.Context, id string) (*out.ExperienceDto, error) {
	return nil, nil
}
