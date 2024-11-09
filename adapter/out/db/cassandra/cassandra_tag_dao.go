package cassandra

import (
	"context"
	"fmt"

	"github.com/ernstvorsteveld/go-cv-cassandra/domain/port/out"
	"github.com/ernstvorsteveld/go-cv-cassandra/pkg/utils"
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	"golang.org/x/exp/slog"
)

const (
	TAG_TABLE_NAME = "cv_tags"
)

var (
	tag_stmt_insert       string = fmt.Sprintf("INSERT INTO %s (id, name) VALUES (?, ?)", TAG_TABLE_NAME)
	tag_stmt_select_by_id string = fmt.Sprintf("SELECT id, name FROM %s WHERE id = ?", TAG_TABLE_NAME)
	tag_stmt_update       string = fmt.Sprintf("UPDATE %s SET name = ? WHERE id = ?", TAG_TABLE_NAME)
	tag_stmt_delete       string = fmt.Sprintf("DELETE FROM %s WHERE id = ?", TAG_TABLE_NAME)
	tag_QryErrorNotFound  error  = errors.Errorf("Not Found")
)

type CassandraTagSession struct {
	config *utils.Configuration
	cs     *gocql.Session
}

func NewTagPort(c *utils.Configuration, s *Session) out.TagDbPort {
	return &CassandraTagSession{
		config: c,
		cs:     s.cs,
	}
}

func (c *CassandraTagSession) Create(ctx context.Context, dto *out.TagDto) (*out.TagDto, error) {
	slog.Debug("cassandra.Create", "content", "About to Create Tag %v", dto)
	err := c.cs.Query(tag_stmt_insert, dto.GetId(), dto.GetName()).Exec()
	return dto, err
}

func (c *CassandraTagSession) Get(ctx context.Context, id string) (*out.TagDto, error) {
	slog.Debug("cassandra.Get", "content", "About to Get Tag by id %v", id)
	return getTagById(c.cs, id)
}

func getTagById(cs *gocql.Session, id string) (*out.TagDto, error) {
	var _id string
	var name string
	if err := cs.Query(tag_stmt_select_by_id, id).Consistency(gocql.One).Scan(&_id, &name); err != nil {
		return nil, err
	}
	return out.NewTagDto(_id, name), nil
}

func (c *CassandraTagSession) GetPage(ctx context.Context, page int32, size int16) ([]out.TagDto, error) {
	return nil, nil
}

func (c *CassandraTagSession) Update(ctx context.Context, id string, dto *out.TagDto) error {
	return nil
}

func (c *CassandraTagSession) Delete(ctx context.Context, id string) (*out.TagDto, error) {
	slog.Debug("cassandra.Delete", "content", "About to Delete Tag by id %v", id)
	dto, err := getTagById(c.cs, id)
	if err != nil {
		return nil, err
	}
	if err = c.cs.Query(tag_stmt_delete, id).Exec(); err != nil {
		return nil, err
	}
	return dto, nil
}
