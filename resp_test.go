package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResp(t *testing.T) {
	// bulkString := "$5\r\nAhmed\r\n"
	arrayString := "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"
	resp := NewResp(strings.NewReader(arrayString))
	val, err := resp.Read()
	assert.NoError(t, err)
	fmt.Println(val)
}
