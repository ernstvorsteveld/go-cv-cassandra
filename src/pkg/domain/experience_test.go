package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_should_have_non_null_experience_name(t *testing.T) {
	_, err := NewExperience("", []string{"a", "b"}, 3, "cassandra", "cassandra")
	assert.NotEmpty(t, err, "Should have a name value")
}

func Test_should_have_valid_experience_name(t *testing.T) {
	e, err := NewExperience("experience name", []string{"a", "b"}, 3, "cassandra", "cassandra")
	assert.Equal(t, nil, err)
	assert.Equal(t, "experience name", e.Name, "Name is not valid")
}

func Test_should_have_invalid_short_experience_name(t *testing.T) {
	_, err := NewExperience("sho", []string{"a", "b"}, 3, "cassandra", "cassandra")
	assert.NotEmpty(t, err.Error(), "Should have error: name is too short")
}

func Test_should_have_invalid_too_long_experience_name(t *testing.T) {
	_, err := NewExperience("This experience has a name that is way too long, it is longer than 100 characters. Therefor this goes on and on, so that the error will appear.", []string{"a", "b"}, 3, "cassandra", "cassandra")
	assert.NotEmpty(t, err, "Should have error: name is too long")
}
