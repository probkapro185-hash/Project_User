package service

import (
	"context"
	"errors"
	"project/internal/models"
	"project/internal/repository"
	"strings"
)

type UserService struct {
	repo *repository.UserRepo
}

func NewUserService(repo *repository.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, name, email string) (models.User, error) {
	if name == "" {
		return models.User{}, errors.New("Имя не может быть пустым!")
	}
	if !strings.Contains(email, "@") {
		return models.User{}, errors.New("email должен содержать @")
	}
	return s.repo.Create(ctx, name, email)
}

func (s *UserService) GetByID(ctx context.Context, id int) (models.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) Update(ctx context.Context, id int, name, email string) (models.User, error) {
	if name == "" {
		return models.User{}, errors.New("имя не может быть пустым")
	}
	if !strings.Contains(email, "@") {
		return models.User{}, errors.New("email должен содержать @")
	}
	return s.repo.Update(ctx, id, name, email)
}

func (s *UserService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
