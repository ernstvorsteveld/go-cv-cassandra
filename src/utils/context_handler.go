package utils

import (
	"context"
	"log/slog"
)

const correlationId string = "correlationId"
const parentCorrelationId string = "parentCorrelationId"

type ContextWrapper struct {
	g IdGenerator
	c context.Context
	m map[string]any
}

func NewDefaultContextWrapper() *ContextWrapper {
	return &ContextWrapper{
		g: NewDefaultUuidGenerator(),
		c: context.Background(),
		m: make(map[string]any),
	}
}

func NewContextWrapper(ig IdGenerator) *ContextWrapper {
	return &ContextWrapper{
		g: ig,
		c: context.Background(),
		m: make(map[string]any),
	}
}

func (w *ContextWrapper) AddParentCorrelationId() *ContextWrapper {
	w.m = add(w.m, parentCorrelationId, w.g.UUIDString())
	return w
}

func GetCorrelationId(c context.Context) string {
	return c.Value(correlationId).(string)
}

func (w *ContextWrapper) AddCorrelationId() *ContextWrapper {
	w.m = add(w.m, correlationId, w.g.UUIDString())
	return w
}

func GetParentCorrelationId(c context.Context) string {
	return c.Value(parentCorrelationId).(string)
}

func add(m map[string]any, k string, v any) map[string]any {
	m[k] = v
	return m
}

func (w *ContextWrapper) Build() context.Context {
	slog.Debug("Build", "content", "About to build context", "attributes", w.m)
	for k, v := range w.m {
		w.c = context.WithValue(w.c, k, v)
	}
	return w.c
}

func Get(k string, c context.Context) any {
	v := c.Value(k)
	if v == nil {
		return "UNKNOWN"
	}
	return v.(string)
}
