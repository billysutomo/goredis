package main

import "os"

var Config *config

type config struct {
	Port string
	AOF  string
}

func init() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "6379"
	}
	aof := os.Getenv("AOF")
	if aof == "" {
		aof = "aof"
	}
	Config = &config{
		Port: port,
		AOF:  aof,
	}
}
