package services

import (
	"testing"

	"github.com/ernstvorsteveld/go-cv-cassandra/src/domain/model"
	"github.com/ernstvorsteveld/go-cv-cassandra/src/utils"
	"github.com/stretchr/testify/assert"
)

func Test_should_create_experience(t *testing.T) {
	_, err := model.NewExperience("", "", []string{"a", "b"})

	c := utils.Configuration{}
	c.Read("test_config", "yml")

	assert.NotEmpty(t, err, "Should have a name value")
}
