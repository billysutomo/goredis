package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResp(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		str := "+OK\r\n"
		resp := NewResp(strings.NewReader(str))
		val, err := resp.Read()
		assert.NoError(t, err)

		want := Value{
			typ: "str",
			str: "OK",
		}

		assert.Equal(t, want, val)
	})
	t.Run("array string", func(t *testing.T) {
		arrayString := "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"
		resp := NewResp(strings.NewReader(arrayString))
		val, err := resp.Read()
		assert.NoError(t, err)

		want := Value{
			typ: "array",
			array: []Value{
				{
					typ:  "bulk",
					bulk: "hello",
				},
				{
					typ:  "bulk",
					bulk: "world",
				}},
		}

		assert.Equal(t, want, val)
	})
}
