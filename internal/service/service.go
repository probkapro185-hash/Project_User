package service

import (
	"context"
	"errors"
	"project/internal/models"
	"project/internal/repository"
	"project/internal/utils"
	"strings"
)

type UserService struct {
	repo *repository.UserRepo
}

func NewUserService(repo *repository.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (h *UserService) Registration(name string, email string, password string) (models.User, error) {
	if name == "" {
		return models.User{}, errors.New("Имя не может быть пустым!")
	}
	if !strings.Contains(email, "@") {
		return models.User{}, errors.New("Email должен содержать @")
	}
	if len(password) < 6 {
		return models.User{}, errors.New("Пароль должен быть больше 6 символов")
	}
	hashpas, err := utils.HashPassword(password)
	if err != nil {
		return models.User{}, err
	}

	return h.repo.Registration(name, email, hashpas)
}

func (h *UserService) Login(email, password string) (models.User, error) {
	if !strings.Contains(email, "@") {
		return models.User{}, errors.New("Email должен содержать @")
	}
	if len(password) < 6 {
		return models.User{}, errors.New("Пароль должен быть больше 6 символов")
	}
	hashpass, err := h.repo.GetByEmail(email)
	if err != nil {
		return models.User{}, errors.New("Пользователя с таким email не найдено")
	}
	if !utils.CheckPassword(hashpass.PasswordHash, password) {
		return models.User{}, errors.New("Пароль не совпадает")
	}
	return models.User{Email: hashpass.Email, Name: hashpass.Name, Id: hashpass.Id}, err
}

func (s *UserService) Create(ctx context.Context, name, email string) (models.User, error) {
	if name == "" {
		return models.User{}, errors.New("Имя не может быть пустым!")
	}
	if !strings.Contains(email, "@") {
		return models.User{}, errors.New("Email должен содержать @")
	}
	return s.repo.Create(ctx, name, email)
}

func (s *UserService) GetByID(ctx context.Context, id int) (models.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) Update(ctx context.Context, id int, name, email string) (models.User, error) {
	if name == "" {
		return models.User{}, errors.New("Имя не может быть пустым")
	}
	if !strings.Contains(email, "@") {
		return models.User{}, errors.New("Email должен содержать @")
	}
	return s.repo.Update(ctx, id, name, email)
}

func (s *UserService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
