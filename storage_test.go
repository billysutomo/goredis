package main

import (
	"reflect"
	"testing"
)

func TestStorage(t *testing.T) {
	t.Run("Set Get", func(t *testing.T) {
		storage := NewInMemoryStorage()
		want := "value"
		storage.Set("key", want)
		got, err := storage.Get("key")
		mustRun(t, err)
		compareString(t, want, got)
	})
	t.Run("HSet HGet", func(t *testing.T) {
		storage := NewInMemoryStorage()
		want := "value"
		storage.HSet("hash", "key", want)
		got, err := storage.HGet("hash", "key")
		mustRun(t, err)
		compareString(t, want, got)
	})
	t.Run("HGetAll", func(t *testing.T) {
		storage := NewInMemoryStorage()
		want := map[string]string{
			"key1": "value1",
			"key2": "value2",
		}

		storage.HSet("hash", "key1", "value1")
		storage.HSet("hash", "key2", "value2")

		got, err := storage.HGetAll("hash")
		mustRun(t, err)

		if !reflect.DeepEqual(want, got) {
			t.Errorf("want %v got %v", want, got)
		}
	})
}

func mustRun(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("didn't expect error %v", err)
	}
}

func compareString(t testing.TB, want, got string) {
	t.Helper()
	if got != want {
		t.Errorf("want %s got %s", want, got)
	}
}
