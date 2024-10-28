package utils

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_should_correlation_id(t *testing.T) {
	correlationId := uuid.New().String()
	parentCorrelationId := uuid.New().String()

	w := NewDefaultContext()
	w.AddCorrelationId(correlationId)
	w.AddParentCorrelationId(parentCorrelationId)
	ctx := w.Build()

	assert.Equal(t, correlationId, GetCorrelationId(ctx), "The CorrelationId is incorrect.")
	assert.Equal(t, parentCorrelationId, GetParentCorrelationId(ctx), "The ParentCorrelationId is incorrect.")
}
