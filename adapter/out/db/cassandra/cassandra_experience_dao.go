package cassandra

import (
	"context"

	cu "github.com/ernstvorsteveld/go-cv-cassandra/adapter/out/db/cassandra/utils"
	"github.com/ernstvorsteveld/go-cv-cassandra/domain/port/out"
	"github.com/ernstvorsteveld/go-cv-cassandra/pkg/utils"
	"github.com/gocql/gocql"
	"golang.org/x/exp/slog"
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

func corrId(ctx context.Context) string {
	return utils.GetCorrelationId(ctx)
}

func (cc *CassandraExperienceSession) Create(ctx context.Context, dto *out.ExperienceDto) error {
	slog.Debug("cassandra.Create", "content", "About to Create Experience", "dto", dto, "correlationId", corrId(ctx))
	return cc.cs.Query(cu.Stmt_insert, dto.GetId(), dto.GetName(), dto.GetTags()).Exec()
}

func (cc *CassandraExperienceSession) Get(ctx context.Context, id string) (*out.ExperienceDto, error) {
	slog.Debug("cassandra.Get", "content", "About to Get Experience by Id", "id", id, "correlationId", corrId(ctx))

	var _id string
	var name string
	var tags []string
	if err := cc.cs.Query(cu.Stmt_select_by_id, id).Consistency(gocql.One).Scan(&_id, &name, &tags); err != nil {
		return nil, err
	}
	e := out.NewExperienceDto(_id, name, tags)
	return e, nil
}

func (cc *CassandraExperienceSession) Update(ctx context.Context, id string, dto *out.ExperienceDto) error {
	slog.Debug("cassandra.Update", "content", "About to Update Experience", "id", id, "dto", dto, "correlationId", corrId(ctx))

	return cc.cs.Query(cu.Stmt_update, dto.GetName(), dto.GetTags(), id).Exec()
}

// Limit *int32
// Page  *string
// Tag   *string
// Name  *string
func (cc *CassandraExperienceSession) GetPage(ctx context.Context, params *out.GetParams) (*out.ExperiencePageReslt, error) {
	slog.Debug("cassandra.GetPage", "content", "About to GetPage Experience", "params", params, "correlationId", corrId(ctx))

	stmt := cu.GetStatement(params)
	slog.Debug(stmt)

	// iter := cc.cs.Query(stmt, limit).Iter()

	return nil, nil
}

func (cc *CassandraExperienceSession) Delete(ctx context.Context, id string) (*out.ExperienceDto, error) {
	return nil, nil
}
