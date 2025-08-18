package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"hope/db"
	"hope/repository"
)

var (
	errUserNotFound      = errors.New("user not found")
	errEmailAndNameReq   = errors.New("name and email required")
	errEmailAlreadyInUse = errors.New("user already exists with this email")
)

type UserService interface {
	CreateUser(ctx context.Context, user *db.User) error
	GetUserByID(ctx context.Context, id string) (*db.User, error)
	GetUserByEmail(ctx context.Context, email string) (*db.User, error)
	UpdateUser(ctx context.Context, user *db.User) error
	DeleteUser(ctx context.Context, id string) error
	UpdateLastSeen(ctx context.Context, id string) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s userService) CreateUser(ctx context.Context, user *db.User) error {
	if user == nil {
		return errors.New("user payload is nil")
	}
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	if user.Email == "" || user.Name == "" {
		return errEmailAndNameReq
	}

	existing, err := s.userRepo.FindByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	if existing != nil && existing.ID != "" {
		return errEmailAlreadyInUse
	}

	user.LastSeen = time.Now()
	return s.userRepo.Create(ctx, user)
}

func (s userService) GetUserByID(ctx context.Context, id string) (*db.User, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, errUserNotFound
	}
	u, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if u == nil || u.ID == "" {
		return nil, errUserNotFound
	}
	return u, nil
}


func (s userService) GetUserByEmail(ctx context.Context, email string) (*db.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return nil, errUserNotFound
	}
	u, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if u == nil || u.Email == "" {
		return nil, errUserNotFound
	}
	return u, nil
}


func (s userService) UpdateUser(ctx context.Context, user *db.User) error {
	if user == nil || strings.TrimSpace(user.ID) == "" {
		return errors.New("user id required")
	}
	curr, err := s.userRepo.FindByID(ctx, strings.TrimSpace(user.ID))
	if err != nil {
		return err
	}
	if curr == nil || curr.ID == "" {
		return errUserNotFound
	}
	if strings.TrimSpace(user.Name) != "" {
		curr.Name = strings.TrimSpace(user.Name)
	}
	if strings.TrimSpace(user.Email) != "" {
		curr.Email = strings.TrimSpace(strings.ToLower(user.Email))
	}
	if strings.TrimSpace(user.Geohash) != "" {
		curr.Geohash = strings.TrimSpace(user.Geohash)
	}

	return s.userRepo.Update(ctx, curr)
}

func (s userService) DeleteUser(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return errors.New("user id required")
	}
	return s.userRepo.Delete(ctx, id)
}

func (s userService) UpdateLastSeen(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return errors.New("user id required")
	}
	u, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if u == nil || u.ID == "" {
		return errUserNotFound
	}
	u.LastSeen = time.Now()
	return s.userRepo.Update(ctx, u)
}
