package repository

import (
	"context"
	"errors"
	"hope/db"
	"gorm.io/gorm"
)


type RideRequestRepository interface {
	Create(ctx context.Context, req *db.RideRequest) error
	FindByID(ctx context.Context, id string) (*db.RideRequest, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	Delete(ctx context.Context, id string) error
	ListNearby(ctx context.Context, geohashPrefix string, limit int) ([]db.RideRequest, error)
	ListByUser(ctx context.Context, userID string, limit int) ([]db.RideRequest, error)
	FindByIDWithUser(ctx context.Context, id string) (*db.RideRequest, error)
	ListActiveByUser(ctx context.Context, userID string, limit int) ([]db.RideRequest, error)
}

type rideRequestRepository struct {
	db *gorm.DB
}

func NewRideRequestRepository(db *gorm.DB) RideRequestRepository {
	return &rideRequestRepository{db: db}
}

func (r *rideRequestRepository) Create(ctx context.Context, req *db.RideRequest) error {
	if req == nil {
		return errors.New("request is nil")
	}
	return r.db.WithContext(ctx).Create(req).Error
}

func (r *rideRequestRepository) FindByID(ctx context.Context, id string) (*db.RideRequest, error) {
	if id == "" {
		return nil, nil
	}
	var out db.RideRequest
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Take(&out).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &out, err
}

func (r *rideRequestRepository) ListNearby(ctx context.Context, geohashPrefix string, limit int) ([]db.RideRequest, error) {
	var reqs []db.RideRequest
	q := r.db.WithContext(ctx).
		Where("from_geo LIKE ?", geohashPrefix+"%").
		Order("time ASC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	err := q.Find(&reqs).Error
	return reqs, err
}

func (r *rideRequestRepository) ListByUser(ctx context.Context, userID string, limit int) ([]db.RideRequest, error) {
	var reqs []db.RideRequest
	q := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("time DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	err := q.Find(&reqs).Error
	return reqs, err
}

func (r *rideRequestRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	if id == "" || status == "" {
		return errors.New("id and status required")
	}
	return r.db.WithContext(ctx).
		Model(&db.RideRequest{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *rideRequestRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id required")
	}
	return r.db.WithContext(ctx).
		Delete(&db.RideRequest{}, "id = ?", id).Error
}

func (r *rideRequestRepository) FindByIDWithUser(ctx context.Context, id string) (*db.RideRequest, error) {
	if id == "" {
		return nil, nil
	}
	var out db.RideRequest
	err := r.db.WithContext(ctx).
		Preload("Rider").
		Where("id = ?", id).
		Take(&out).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &out, err
}

func (r *rideRequestRepository) ListActiveByUser(ctx context.Context, userID string, limit int) ([]db.RideRequest, error) {
	var reqs []db.RideRequest
	q := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ?", userID, "active").
		Order("time ASC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	err := q.Find(&reqs).Error
	return reqs, err
}
