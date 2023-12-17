package main

import (
	"context"
	"fmt"
	"sync"
)

type Storage interface {
	Set(key string, value string)
	Get(key string) (string, error)
	HSet(hash string, key string, value string)
	HGet(hash string, key string) (string, error)
	HGetAll(hash string) (map[string]string, error)
	Publish(channel string, value Value) int
	Subscribe(ctx context.Context, channelName string, writer Destination)
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		setS:           make(map[string]string),
		hSetS:          make(map[string]map[string]string),
		subscribeChans: make(map[string]chan Value),
	}
}

type PubSub struct {
	Middleware chan Value
	Targets    map[string]chan Value
}

type InMemoryStorage struct {
	setS             map[string]string
	setSMu           sync.RWMutex
	hSetS            map[string]map[string]string
	hSetSMu          sync.RWMutex
	subscribeChans   map[string](chan Value)
	subscribeChansMu sync.RWMutex
}

func (i *InMemoryStorage) Set(key string, value string) {
	i.setSMu.Lock()
	defer i.setSMu.Unlock()
	i.setS[key] = value
}

func (i *InMemoryStorage) Get(key string) (string, error) {
	i.setSMu.RLock()
	defer i.setSMu.RUnlock()
	value, ok := i.setS[key]
	if !ok {
		return "", fmt.Errorf("")
	}
	return value, nil
}

func (i *InMemoryStorage) HSet(hash string, key string, value string) {
	i.hSetSMu.Lock()
	defer i.hSetSMu.Unlock()

	if _, ok := i.hSetS[hash]; !ok {
		i.hSetS[hash] = make(map[string]string)
	}
	i.hSetS[hash][key] = value
}

func (i *InMemoryStorage) HGet(hash string, key string) (string, error) {
	i.hSetSMu.RLock()
	defer i.hSetSMu.RUnlock()
	value, ok := i.hSetS[hash][key]
	if !ok {
		return "", fmt.Errorf("")
	}
	return value, nil
}

func (i *InMemoryStorage) HGetAll(hash string) (map[string]string, error) {
	i.hSetSMu.RLock()
	defer i.hSetSMu.RUnlock()
	value, ok := i.hSetS[hash]
	if !ok {
		return nil, fmt.Errorf("")
	}
	return value, nil
}

func (i *InMemoryStorage) getChan(key string) (bool, chan Value) {
	i.subscribeChansMu.Lock()
	defer i.subscribeChansMu.Unlock()
	isExist := true
	channel, ok := i.subscribeChans[key]
	if !ok {
		isExist = false
		i.subscribeChans[key] = make(chan Value)
		channel = i.subscribeChans[key]
	}
	return isExist, channel
}

func (i *InMemoryStorage) Publish(channel string, value Value) int {
	isExist, targetChannel := i.getChan(channel)
	if !isExist {
		return 0
	}
	targetChannel <- value
	return 1
}

func (i *InMemoryStorage) Subscribe(ctx context.Context, channelName string, writer Destination) {
	_, channel := i.getChan(channelName)
	for {
		select {
		case msg := <-channel:
			writer.Write(msg)
		case <-ctx.Done():
			delete(i.subscribeChans, channelName)
		}
	}
}
