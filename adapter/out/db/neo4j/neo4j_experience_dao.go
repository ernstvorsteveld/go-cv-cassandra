package neo4j

import (
	"context"
	"errors"
	"strconv"

	"github.com/ernstvorsteveld/go-cv-cassandra/domain/port/out"
	"github.com/ernstvorsteveld/go-cv-cassandra/pkg/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	n "github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"golang.org/x/exp/slog"
)

type Neo4jExperienceSession struct {
	config *utils.Configuration
	driver n.DriverWithContext
}

func NewExperiencePort(c *utils.Configuration, driver neo4j.DriverWithContext) out.ExperienceDbPort {
	return &Neo4jExperienceSession{
		config: c,
		driver: driver,
	}
}

func (new *Neo4jExperienceSession) Create(ctx context.Context, dto *out.ExperienceDto) error {
	slog.Debug("neo4j.Create", "content", "About to Create Experience",
		"dto", dto, "correlationId", utils.GetCorrelationId(ctx))

	result, err := n.ExecuteQuery(ctx, new.driver,
		"CREATE (p:Experience {id : $id, name: $name}) RETURN p",
		map[string]any{
			"id":   dto.GetId(),
			"name": dto.GetName(),
		},
		n.EagerResultTransformer,
		n.ExecuteQueryWithDatabase("neo4j"))
	if err != nil {
		return err
	}
	slog.Debug("neo4j.Create", "content", "Neo4j result create experience",
		"result", result.Summary.Counters().ContainsUpdates(),
		"correlationId", utils.GetCorrelationId(ctx))

	for _, tagName := range dto.GetTags() {
		result, err := neo4j.ExecuteQuery(ctx, new.driver, `
			MATCH (tag:Tag {name: $tagName})
			MATCH (exp:Experience {id: $experienceId})
    		CREATE (tag)-[:HAS_EXPERIENCE]->(exp)`,
			map[string]any{
				"tagName":      tagName,
				"experienceId": dto.GetId(),
			}, neo4j.EagerResultTransformer,
			neo4j.ExecuteQueryWithDatabase("neo4j"))
		if err != nil {
			return err
		}
		slog.Debug("neo4j.Create", "content", "Neo4j result create tag-experience relation",
			"result", result.Summary.Counters().ContainsUpdates(),
			"correlationId", utils.GetCorrelationId(ctx))
	}
	return nil
}

func (new *Neo4jExperienceSession) Get(ctx context.Context, id string) (*out.ExperienceDto, error) {
	slog.Debug("neo4j.Get", "content", "About to Get Experience by Id",
		"id", id, "correlationId", utils.GetCorrelationId(ctx))

	expResult, err := neo4j.ExecuteQuery(ctx, new.driver, `
		MATCH (p:Experience {id: $id}) RETURN p.id AS id, p.name AS name`,
		map[string]any{
			"id": id,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"))
	if err != nil {
		return nil, err
	}
	slog.Debug("neo4j.Get", "content", "Neo4j result get Experience",
		"result", expResult.Summary.Counters().ContainsUpdates(),
		"correlationId", utils.GetCorrelationId(ctx))

	if len(expResult.Records) <= 0 {
		return nil, errors.New("experience not found")
	}
	if len(expResult.Records) > 1 {
		return nil, errors.New("experience not unique")
	}

	result, err := neo4j.ExecuteQuery(ctx, new.driver, `
		MATCH (p:Experience {id: $id})<-[:HAS_EXPERIENCE]-(tag:Tag)
		RETURN tag.name as tagName`,
		map[string]any{
			"id": id,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"))

	slog.Debug("neo4j.GetTags", "content", "Neo4j result get Experience",
		"result", result.Summary.Counters().ContainsUpdates(),
		"correlationId", utils.GetCorrelationId(ctx))

	x, _ := expResult.Records[0].Get("id")
	y, _ := expResult.Records[0].Get("name")
	tags := getTags(result.Records)
	dto := out.NewExperienceDto(
		x.(string),
		y.(string),
		tags)
	return dto, nil
}

func getTags(records []*n.Record) []string {
	var tags []string
	for _, record := range records {
		tag, _ := record.Get("tagName")
		tags = append(tags, tag.(string))
	}
	return tags
}

func (cc *Neo4jExperienceSession) Update(ctx context.Context, id string, dto *out.ExperienceDto) error {
	slog.Debug("neo4j.Update", "content", "About to Update Experience", "id", id,
		"dto", dto, "correlationId", utils.GetCorrelationId(ctx))
	return nil
}

// Limit *int32
// Page  *string
// Tag   *string
// Name  *string
func (new *Neo4jExperienceSession) GetPage(ctx context.Context, params *out.GetParams) (*out.ExperiencePageReslt, error) {
	slog.Debug("neo4j.GetPage", "content", "About to GetPage Experience",
		"params", params, "correlationId", utils.GetCorrelationId(ctx))

	skip := getSkip(params)

	result, err := neo4j.ExecuteQuery(ctx, new.driver, `
		MATCH (p:Experience {name: $name})
		ORDER BY p.name
		LIMIT $limit
		SKIP $skip
		RETURN p.id AS id, p.name AS name`,
		map[string]any{
			"name":  *params.Name,
			"limit": *params.Limit,
			"skip":  skip,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"))
	if err != nil {
		return nil, err
	}
	if len(result.Records) <= 0 {
		return nil, errors.New("experience not found")
	}

	slog.Debug("neo4j.GetPage", "content", "Neo4j result get Experience by page",
		"result", result.Summary.Counters().ContainsUpdates(),
		"correlationId", utils.GetCorrelationId(ctx))

	dtos := make([]out.ExperienceDto, len(result.Records))
	for index, _ := range result.Records {
		x, _ := result.Records[index].Get("id")
		y, _ := result.Records[0].Get("name")
		tags, err := getTagsFromDb(ctx, new.driver, x.(string))
		if err != nil {
			return nil, err
		}
		dtos[index] = *out.NewExperienceDto(
			x.(string),
			y.(string),
			tags)
	}

	return out.NewExperiencePageReslt(nil, nil, dtos), nil
}

func getTagsFromDb(ctx context.Context, driver neo4j.DriverWithContext, id string) ([]string, error) {
	tags, err := neo4j.ExecuteQuery(ctx, driver, `
		MATCH (p:Experience {id: $id})<-[:HAS_EXPERIENCE]-(tag:Tag)
		RETURN tag.name as tagName`,
		map[string]any{
			"id": id,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"))
	if err != nil {
		return nil, err
	}
	result := make([]string, len(tags.Records))
	for index, _ := range tags.Records {
		t, _ := tags.Records[index].Get("tagName")
		result[index] = t.(string)
	}
	return result, nil
}

func getSkip(p *out.GetParams) int {
	if p.Page == nil {
		return 0
	}
	page, err := strconv.Atoi(*p.Page)
	if err != nil {
		return 0
	}
	return page * int(*p.Limit)
}

func (cc *Neo4jExperienceSession) Delete(ctx context.Context, id string) (*out.ExperienceDto, error) {
	return nil, nil
}
