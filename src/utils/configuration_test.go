package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadDBConfiguration(t *testing.T) {
	c := Configuration{}
	c.Read("test_config", "yml")

	assert.Equal(t, "127.0.0.1", c.DB.Cassandra.Url, "The Cassandra DB url is incorrect.")
	assert.Equal(t, "cv", c.DB.Cassandra.Keyspace, "The Keyspace is incorrect.")
}

func TestReadDBConfiguration_use_environment(t *testing.T) {
	expectedUrl := "http://example.com"
	t.Setenv("CASSANDRA_URL", expectedUrl)

	expectedKeySpace := "keyspace-value"
	t.Setenv("CASSANDRA_KEYSPACE", expectedKeySpace)

	c := Configuration{}
	c.Read("test_config", "yml")
	assert.Equal(t, expectedUrl, c.DB.Cassandra.Url, "The Cassandra DB url is incorrect.")
	assert.Equal(t, expectedKeySpace, c.DB.Cassandra.Keyspace, "The Keyspace is incorrect.")

	c.Print()
}
