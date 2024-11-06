package utils

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

const CORRELATION_ID_HEADER = "X-CORRELATION-ID"
const CORRELATION_ID = "correlationId"
const PARENT_CORRELATION_ID = "parentCorrelationId"

type ContextWrapper struct {
	g IdGenerator
	c context.Context
	m map[string]any
}

func NewDefaultContextWrapper(pc context.Context, correlationId string) *ContextWrapper {
	w := ContextWrapper{
		g: NewDefaultUuidGenerator(),
		c: pc,
		m: make(map[string]any),
	}
	w.m[CORRELATION_ID] = correlationId
	return &w
}

func NewContextWrapper(pc context.Context, ig IdGenerator) *ContextWrapper {
	w := ContextWrapper{
		g: ig,
		c: pc,
		m: make(map[string]any),
	}
	w.m[CORRELATION_ID] = ig.UUIDString()
	return &w
}

func (w *ContextWrapper) AddParentCorrelationId() *ContextWrapper {
	w.m = add(w.m, PARENT_CORRELATION_ID, w.g.UUIDString())
	return w
}

func GetCorrelationUuid(c *context.Context) uuid.UUID {
	return uuid.MustParse(get(CORRELATION_ID, c).(string))
}

func GetCorrelationId(c *context.Context) string {
	return get(CORRELATION_ID, c).(string)
}

func GetParentCorrelationId(c *context.Context) string {
	return get(PARENT_CORRELATION_ID, c).(string)
}

func add(m map[string]any, k string, v any) map[string]any {
	m[k] = v
	return m
}

func (w *ContextWrapper) Build() *context.Context {
	slog.Debug("Build", "content", "About to build context", "attributes", w.m)
	for k, v := range w.m {
		w.c = context.WithValue(w.c, k, v)
	}
	return &w.c
}

func get(k string, c *context.Context) any {
	v := (*c).Value(k)
	if v == nil {
		return "UNKNOWN"
	}
	return v.(string)
}
