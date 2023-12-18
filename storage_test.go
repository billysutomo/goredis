package main

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"
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
	t.Run("subscribe", func(t *testing.T) {
		storage := NewInMemoryStorage()
		channel := "channel"
		writer := &MockWriter{}
		want := Value{
			typ:  "bulk",
			bulk: "message",
		}
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			storage.Subscribe(context.Background(), channel, writer)
		}()

		time.Sleep(10 * time.Millisecond)
		storage.Publish(channel, want)
		if !waitFor(&wg, 10*time.Millisecond) {
			t.Errorf("timeout waiting for Subscribe goroutine to finish")
			return
		}

		if writer.Value.bulk != want.bulk {
			t.Errorf("got %+v, want %+v", writer.Value, want)
		}
	})
	t.Run("publish even without subscriber", func(t *testing.T) {
		storage := NewInMemoryStorage()
		channel := "channel"
		want := Value{
			typ:  "bulk",
			bulk: "message",
		}
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			storage.Publish(channel, want)
		}()

		if !waitFor(&wg, 10*time.Millisecond) {
			t.Errorf("timeout waiting for Publish goroutine to finish")
		}
	})
}

func waitFor(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return true
	case <-time.After(timeout):
		return false
	}
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
