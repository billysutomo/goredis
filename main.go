package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {
	fmt.Printf("Listening on port :%s\n", Config.Port)
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", Config.Port))
	if err != nil {
		fmt.Println(err)
		return
	}

	aof, err := NewAof(fmt.Sprintf("%s.aof", Config.AOF))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer aof.Close()

	storage := NewInMemoryStorage()
	handler := NewHandler(storage)

	aof.Read(func(value Value) {
		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]
		handler.handleCommand(context.TODO(), nil, command, args)
	})
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		defer conn.Close()
		go NewConn(conn, aof, storage, handler)
	}
}

func NewConn(conn net.Conn, aof *Aof, storage Storage, handler handler) {
	writer := NewWriter(conn)
	ctx, done := context.WithCancel(context.Background())
	defer done()
	for {
		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			if err != io.EOF {
				fmt.Printf("error read %v \n", err)
			}
			return
		}

		if value.typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		val := handler.handleCommand(ctx, writer, command, args)

		if command == "SET" || command == "HSET" {
			aof.Write(value)
		}

		writer.Write(val)
	}
}
