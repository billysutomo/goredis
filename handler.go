package main

import (
	"context"
	"strings"
)

type handler struct {
	storage Storage
}

func NewHandler(storage Storage) handler {
	hndlr := handler{storage}
	return hndlr
}

func (h *handler) handleCommand(ctx context.Context, writer Destination, command string, args []Value) Value {
	switch strings.ToUpper(command) {
	case "PING":
		return h.ping(args)
	case "SET":
		return h.set(args)
	case "GET":
		return h.get(args)
	case "HSET":
		return h.hset(args)
	case "HGET":
		return h.hget(args)
	case "ALL":
		return h.hgetall(args)
	case "SUBSCRIBE":
		return h.subscribe(ctx, writer, args)
	case "PUBLISH":
		return h.publish(args)
	default:
		return Value{typ: "error", str: "ERR unknown command"}
	}
}

type Store struct {
	storage Storage
}

func NewStore(s Storage) *Store {
	return &Store{storage: s}
}

func (s *handler) ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}
	return Value{typ: "string", str: args[0].bulk}
}

func (s *handler) set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].bulk
	value := args[1].bulk

	s.storage.Set(key, value)

	return Value{typ: "string", str: "OK"}
}

func (s *handler) get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].bulk

	value, err := s.storage.Get(key)

	if err != nil {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

func (s *handler) hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hset' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	s.storage.HSet(hash, key, value)

	return Value{typ: "string", str: "OK"}
}

func (s *handler) hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hget' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk

	value, err := s.storage.HGet(hash, key)
	if err != nil {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

func (s *handler) hgetall(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'getall' command"}
	}

	hash := args[0].bulk
	value, err := s.storage.HGetAll(hash)
	if err != nil {
		return Value{typ: "null"}
	}
	values := []Value{}
	for k, v := range value {
		values = append(values, Value{typ: "bulk", bulk: k})
		values = append(values, Value{typ: "bulk", bulk: v})
	}
	return Value{typ: "array", array: values}
}

func (s *handler) subscribe(ctx context.Context, writer Destination, args []Value) Value {
	if len(args) < 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'subscribe' command"}
	}

	key := args[0].bulk

	writer.Write(Value{
		typ: "array",
		array: []Value{
			{
				typ:  "bulk",
				bulk: "subscribe",
			},
			{
				typ:  "bulk",
				bulk: key,
			},
			{
				typ: "num",
				num: len(args),
			},
		},
	})

	for _, a := range args {
		go s.storage.Subscribe(ctx, a.bulk, writer)
	}

	return Value{}
}

func (s *handler) publish(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'publish' command"}
	}
	key := args[0].bulk
	value := args[1].bulk

	listenerCount := s.storage.Publish(key, Value{
		typ: "array",
		array: []Value{
			{
				typ:  "bulk",
				bulk: "message",
			},
			{
				typ:  "bulk",
				bulk: key,
			},
			{
				typ:  "bulk",
				bulk: value,
			},
		},
	})
	return Value{typ: "num", num: listenerCount}
}
