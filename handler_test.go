package main

import (
	"context"
	"testing"
)

type MockWriter struct {
	Value Value
}

func (m *MockWriter) Write(v Value) error {
	m.Value = v
	return nil
}

func TestHandlerSetGet(t *testing.T) {
	t.Run("SET GET", func(t *testing.T) {
		storage := NewInMemoryStorage()
		handler := NewHandler(storage)

		writer := &MockWriter{}
		handler.handleCommand(context.TODO(), writer, "SET", []Value{
			{
				typ:  "bulk",
				bulk: "key",
			},
			{
				typ:  "bulk",
				bulk: "value",
			},
		})

		got, err := storage.Get("key")
		mustRun(t, err)
		compareString(t, "value", got)
	})
	t.Run("HSET HGET", func(t *testing.T) {
		storage := NewInMemoryStorage()
		handler := NewHandler(storage)

		writer := &MockWriter{}
		handler.handleCommand(context.TODO(), writer, "HSET", []Value{
			{
				typ:  "bulk",
				bulk: "hash",
			},
			{
				typ:  "bulk",
				bulk: "key",
			},
			{
				typ:  "bulk",
				bulk: "value",
			},
		})

		got, err := storage.HGet("hash", "key")
		mustRun(t, err)
		compareString(t, "value", got)
	})
}
