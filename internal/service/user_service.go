package service

import (
	"context"
	"crud/internal/models"
	"crud/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, name, email string) (*models.User, error) {
	user := &models.User{Name: name, Email: email}
	err := s.repo.Create(ctx, user)
	return user, err
}

func (s *UserService) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.GetByEmail(ctx, email)
}
