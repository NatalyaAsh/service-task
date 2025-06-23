package main

import (
	"service-task/internal/config"
	"service-task/internal/server"
	"service-task/internal/tasks"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		return
	}

	tasks.Init()

	server.Start(cfg)
}
