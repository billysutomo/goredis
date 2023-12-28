package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRespString(t *testing.T) {
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
}

func TestRespError(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		errStr := "-Error message\r\n"
		resp := NewResp(strings.NewReader(errStr))
		val, err := resp.Read()
		assert.NoError(t, err)

		want := Value{
			typ: "error",
			str: "Error message",
		}

		assert.Equal(t, want, val)
	})
}

func TestRespInteger(t *testing.T) {
	t.Run("integer 0 ", func(t *testing.T) {
		intStr := ":0\r\n"
		resp := NewResp(strings.NewReader(intStr))
		val, err := resp.Read()
		assert.NoError(t, err)

		want := Value{
			typ: "integer",
			num: 0,
		}

		assert.Equal(t, want, val)
	})
	t.Run("integer 1000", func(t *testing.T) {
		intStr := ":1000\r\n"
		resp := NewResp(strings.NewReader(intStr))
		val, err := resp.Read()
		assert.NoError(t, err)

		want := Value{
			typ: "integer",
			num: 1000,
		}

		assert.Equal(t, want, val)
	})
}

func TestBulkString(t *testing.T) {
	buff := "$6\r\nfoobar\r\n"
	want := Value{
		typ:  "bulk",
		bulk: "foobar",
	}
	resp := NewResp(strings.NewReader(string(buff)))
	marshaled, err := resp.Read()
	assert.NoError(t, err)
	assert.Equal(t, want, marshaled)

	secondBuff := marshaled.Marshal()
	assert.Equal(t, buff, string(secondBuff))
}

func TestRespArray(t *testing.T) {
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
