package handler

import (
	"github.com/absdekty/taskmanager/internal/service"
)

type Handler struct {
	service service.ServiceI
}

func NewHandler(service service.ServiceI) *Handler {
	return &Handler{service: service}
}
