package utils

import (
	"context"
	"testing"

	utils_mock "github.com/ernstvorsteveld/go-cv-cassandra/pkg/utils/mock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_should_correlation_id(t *testing.T) {
	uuid := uuid.New()

	w := NewContextWrapper(context.Background(), utils_mock.NewMockUuidGenerator(uuid))
	w.AddParentCorrelationId()
	ctx := w.Build()

	assert.Equal(t, uuid.String(), GetCorrelationId(ctx), "The CorrelationId is incorrect.")
	assert.Equal(t, uuid, GetCorrelationUuid(ctx), "The CorrelationId is incorrect.")
	assert.Equal(t, uuid.String(), GetParentCorrelationId(ctx), "The ParentCorrelationId is incorrect.")
}
