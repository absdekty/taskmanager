package service

import (
	"github.com/absdekty/taskmanager/internal/repository"
)

type Service struct {
	repo repository.RepositoryI
}

func NewService(repo repository.RepositoryI) *Service {
	return &Service{repo: repo}
}
