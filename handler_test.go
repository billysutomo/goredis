package main

import (
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
	// storage := &InMemoryStorage{setS: map[string]string{}}
	// t.Run("SET", func(t *testing.T) {
	// 	handler := NewHandler(storage)
	// 	handlerSET, ok := handler["SET"]
	// 	assert.Equal(t, true, ok)
	// 	writer := &MockWriter{}
	// 	handlerSET(writer, []Value{
	// 		{
	// 			typ:  "bulk",
	// 			bulk: "key",
	// 		},
	// 		{
	// 			typ:  "bulk",
	// 			bulk: "value",
	// 		},
	// 	})
	// 	assert.Equal(t, Value{typ: "string", str: "OK"}.Marshal(), writer.Value.Marshal())
	// })
	// t.Run("GET", func(t *testing.T) {
	// 	handler := NewHandler(storage)
	// 	handlerGET, ok := handler["GET"]
	// 	assert.Equal(t, true, ok)
	// 	writer := &MockWriter{}
	// 	handlerGET(writer, []Value{
	// 		{
	// 			typ:  "bulk",
	// 			bulk: "key",
	// 		},
	// 	})

	// 	want := Value{
	// 		typ:  "bulk",
	// 		bulk: "value",
	// 	}
	// 	assert.Equal(t, want, writer.Value)
	// })
}

// func TestHandlerHSetHGet(t *testing.T) {
// 	handler, ok := Handlers["HSET"]
// 	assert.Equal(t, true, ok)
// 	result := handler([]Value{
// 		{
// 			typ:  "bulk",
// 			bulk: "hash",
// 		},
// 		{
// 			typ:  "bulk",
// 			bulk: "key",
// 		},
// 		{
// 			typ:  "bulk",
// 			bulk: "value",
// 		},
// 	})
// 	assert.Equal(t, "OK", result.str)

// 	handler, ok = Handlers["HGET"]
// 	assert.Equal(t, true, ok)
// 	result = handler([]Value{
// 		{
// 			typ:  "bulk",
// 			bulk: "hash",
// 		},
// 		{
// 			typ:  "bulk",
// 			bulk: "key",
// 		},
// 	})
// 	assert.Equal(t, "value", result.bulk)
// }
