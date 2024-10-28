package utils

import (
	"context"
	"log/slog"
)

const correlationId string = "correlationId"
const parentCorrelationId string = "parentCorrelationId"

type ContextWrapper struct {
	c context.Context
	m map[string]any
}

func NewDefaultContext() *ContextWrapper {
	return &ContextWrapper{
		c: context.Background(),
		m: make(map[string]any),
	}
}

func (w *ContextWrapper) AddParentCorrelationId(v string) {
	w.m = add(w.m, parentCorrelationId, v)
}

func GetCorrelationId(c context.Context) string {
	return c.Value(correlationId).(string)
}

func (w *ContextWrapper) AddCorrelationId(v string) {
	w.m = add(w.m, correlationId, v)
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
