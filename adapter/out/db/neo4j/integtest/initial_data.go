package neo4j_integtest

import (
	"context"
	"fmt"

	"github.com/ernstvorsteveld/go-cv-cassandra/domain/port/out"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

const (
	TAG_NUMBER = 5
	EXP_NUMBER = 12
)

var (
	TagUids = make([]string, TAG_NUMBER)
	ExpUids = make([]string, EXP_NUMBER)
)

func loadInitialData(ctx context.Context, driver neo4j.DriverWithContext) {
	for i := 0; i < TAG_NUMBER; i++ {
		TagUids[i] = uuid.NewString()
	}
	for i := 0; i < EXP_NUMBER; i++ {
		ExpUids[i] = uuid.NewString()
	}
	loadTags(ctx, driver)
	loadExperiences(ctx, driver)

	validate(ctx, driver)
}

func loadTags(ctx context.Context, driver neo4j.DriverWithContext) {
	tags := []*out.TagDto{
		out.NewTagDto(TagUids[0], "tag-1"),
		out.NewTagDto(TagUids[1], "tag-2"),
		out.NewTagDto(TagUids[2], "tag-3"),
		out.NewTagDto(TagUids[3], "tag-4"),
		out.NewTagDto(TagUids[4], "tag-5"),
	}

	for _, tag := range tags {
		neo4j.ExecuteQuery(ctx, driver,
			"CREATE (p:Tag {id : $id, name: $name}) RETURN p",
			map[string]any{
				"id":   tag.GetId(),
				"name": tag.GetName(),
			},
			neo4j.EagerResultTransformer,
			neo4j.ExecuteQueryWithDatabase("neo4j"))
	}
}

func loadExperiences(ctx context.Context, driver neo4j.DriverWithContext) {
	experiences := []*out.ExperienceDto{
		out.NewExperienceDto(ExpUids[0], "name-1", []string{"tag-1"}),
		out.NewExperienceDto(ExpUids[1], "name-2", []string{"tag-2"}),
		out.NewExperienceDto(ExpUids[2], "name-3", []string{"tag-2", "tag-3"}),
		out.NewExperienceDto(ExpUids[3], "name-4", []string{"tag-2", "tag-3"}),
		out.NewExperienceDto(ExpUids[4], "name-5", []string{"tag-2", "tag-3"}),
		out.NewExperienceDto(ExpUids[5], "name-6", []string{"tag-2", "tag-3", "tag-4"}),
		out.NewExperienceDto(ExpUids[6], "name-7", []string{"tag-2", "tag-3", "tag-4"}),
		out.NewExperienceDto(ExpUids[7], "name-8", []string{"tag-2", "tag-3", "tag-4"}),
		out.NewExperienceDto(ExpUids[8], "name-9", []string{"tag-3", "tag-4", "tag-5"}),
		out.NewExperienceDto(ExpUids[9], "name-10", []string{"tag-3", "tag-4", "tag-5"}),
		out.NewExperienceDto(ExpUids[10], "name-11", []string{"tag-3", "tag-4", "tag-5"}),
		out.NewExperienceDto(ExpUids[11], "name-12", []string{"tag-5"}),
	}

	for _, experience := range experiences {
		neo4j.ExecuteQuery(ctx, driver,
			"CREATE (p:Experience {id : $id, name: $name}) RETURN p",
			map[string]any{
				"id":   experience.GetId(),
				"name": experience.GetName(),
			},
			neo4j.EagerResultTransformer,
			neo4j.ExecuteQueryWithDatabase("neo4j"))

		for _, tagName := range experience.GetTags() {
			neo4j.ExecuteQuery(ctx, driver, `
				MATCH (tag:Tag {name: $tagName})
				MATCH (exp:Experience {id: $experienceId})
    			CREATE (tag)-[:HAS_EXPERIENCE]->(exp)`,
				map[string]any{
					"tagName":      tagName,
					"experienceId": experience.GetId(),
				}, neo4j.EagerResultTransformer,
				neo4j.ExecuteQueryWithDatabase("neo4j"))
		}
	}
}

func validate(ctx context.Context, driver neo4j.DriverWithContext) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		"MATCH (p:Experience) RETURN p.name AS name",
		nil,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"))
	if err != nil {
		panic(err)
	}

	if len(result.Records) != EXP_NUMBER {
		fmt.Errorf("do not have correct number of dxperience, expected %d, got %d", EXP_NUMBER, len(result.Records))
	}
}
