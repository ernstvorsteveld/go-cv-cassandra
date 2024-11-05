package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_should_read_configuration(t *testing.T) {
	c := Configuration{}
	c.Read("test_config", "yml")

	assert.Equal(t, "DEBUG", c.DebugLevel, "Debug level incorrect, not DEBUG.")
	assert.Equal(t, "8091", c.Api.CV.Port, "The CV Api Ports incorrect.")
	assert.Equal(t, []string{"localhost:8091", "localhost:8092"}, c.Api.CV.Expectedhosts, "The CV Api Expectehosts incorrect.")
	assert.Equal(t, "8092", c.Api.Monitoring.Port, "The Monitoring Api Ports incorrect.")
	assert.Equal(t, []string{"localhost:8091"}, c.Api.Monitoring.Expectedhosts, "The Monitoring Api Expectehosts incorrect.")

	assert.Equal(t, "127.0.0.1", c.DB.Cassandra.Url, "The Cassandra DB url is incorrect.")
	assert.Equal(t, "cv", c.DB.Cassandra.Keyspace, "The Keyspace is incorrect.")
	assert.Equal(t, int8(3), c.DB.Cassandra.Retries, "The Retries is incorrect.")
	assert.Equal(t, "cassandra", c.DB.Cassandra.Username, "The Username is incorrect.")
	assert.Equal(t, "cassandra", c.DB.Cassandra.Secret.Value(), "The Secret is incorrect.")
}

func Test_should_read_configuration_use_environment(t *testing.T) {
	expectedUrl := "http://example.com"
	t.Setenv("CASSANDRA_URL", expectedUrl)

	expectedKeySpace := "keyspace-value"
	t.Setenv("CASSANDRA_KEYSPACE", expectedKeySpace)

	expectedSecret := "secret"
	t.Setenv("CASSANDRA_SECRET", expectedSecret)

	c := Configuration{}
	c.Read("test_config", "yml")
	assert.Equal(t, expectedUrl, c.DB.Cassandra.Url, "The Cassandra DB url is incorrect.")
	assert.Equal(t, expectedKeySpace, c.DB.Cassandra.Keyspace, "The Keyspace is incorrect.")
	assert.Equal(t, expectedSecret, c.DB.Cassandra.Secret.Value(), "The Secret is incorrect.")

	c.Print()
}
