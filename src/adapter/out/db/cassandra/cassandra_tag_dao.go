package cassandra

import (
	"context"

	"github.com/ernstvorsteveld/go-cv-cassandra/src/port/out"
	"github.com/ernstvorsteveld/go-cv-cassandra/src/utils"
	"github.com/gocql/gocql"
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
	return nil, nil
}

func (c *CassandraTagSession) Get(ctx context.Context, id string) (*out.TagDto, error) {
	return nil, nil
}

func (c *CassandraTagSession) GetPage(ctx context.Context, page int32, size int16) ([]out.TagDto, error) {
	return nil, nil
}

func (c *CassandraTagSession) Update(ctx context.Context, id string, dto *out.TagDto) error {
	return nil
}

func (c *CassandraTagSession) Delete(ctx context.Context, id string) (*out.TagDto, error) {
	return nil, nil
}
