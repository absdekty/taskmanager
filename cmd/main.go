package main

import (
	"log"
	"github.com/absdekty/taskmanager/internal/repository/sqlite"
	"github.com/absdekty/taskmanager/internal/service"
	"github.com/absdekty/taskmanager/internal/handler"
)

func main() {
	repo, err := repository.NewRepository("...")
	if err != nil {
		log.Fatal(err)
	}
	
	service := service.NewService(repo)
	handler := handler.NewHandler(service)

	log.Printf("Started: %+v\n", handler)
}
