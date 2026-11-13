package main

import (
	"log"
	"tasks-api/internal/api"
	"tasks-api/internal/storage/memory"
)

func main() {
	store := memory.NewMemStorage()

	handler := api.NewHandler(store)

	server := api.NewServer(handler)

	server.RegisterHandlers()

	if err := server.Start(); err != nil {
		log.Fatalf("error:%v", err)
		return
	}
}
