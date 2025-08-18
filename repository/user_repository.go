package repository

import (
	"context"
	"errors"
	"fmt"

	"hope/db"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository interface {
	Create(ctx context.Context, user *db.User) error
	FindByEmail(ctx context.Context, email string) (*db.User, error)
	FindByID(ctx context.Context, id string) (*db.User, error)
	Update(ctx context.Context, user *db.User) error
	Delete(ctx context.Context, id string) error
	FindByIDWithLocation(ctx context.Context, id string) (*db.User, error)
	OptimisticUpdateLastSeen(ctx context.Context, id string, oldLastSeen int64, newLastSeen int64) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *db.User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*db.User, error) {
	var user db.User
	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		Take(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*db.User, error) {
	var user db.User
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Take(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByIDWithLocation(ctx context.Context, id string) (*db.User, error) {
	var user db.User
	err := r.db.WithContext(ctx).
		Preload("Location").
		Where("id = ?", id).
		Take(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *db.User) error {
	if user == nil || user.ID == "" {
		return errors.New("user or ID missing")
	}
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id required")
	}
	return r.db.WithContext(ctx).
		Delete(&db.User{}, "id = ?", id).Error
}

func (r *userRepository) OptimisticUpdateLastSeen(ctx context.Context, id string, oldLastSeen int64, newLastSeen int64) error {
	if id == "" {
		return errors.New("id required")
	}
	tx := r.db.WithContext(ctx).Model(&db.User{}).
		Where("id = ? AND UNIX_TIMESTAMP(last_seen) = ?", id, oldLastSeen).
		Update("last_seen", clause.Expr{SQL: "FROM_UNIXTIME(?)", Vars: []any{newLastSeen}})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return fmt.Errorf("conflict: user last_seen changed concurrently")
	}
	return nil
}
