package main

import (
	"log"

	"grpc-todo-list/cli"
	"grpc-todo-list/config"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		panic("failed to get config")
	}

	if err = cli.Execute(config); err != nil {
		log.Fatalf("Error executing CLI: %v", err)
	}
}
