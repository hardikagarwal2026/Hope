package repository

import (
	"context"
	"errors"
	"hope/db"
	"gorm.io/gorm"
)

type MatchRepository interface {
	Create(ctx context.Context, match *db.Match) error
	FindByID(ctx context.Context, id string) (*db.Match, error)
	UpdateStatus(ctx context.Context, matchID string, status string) error
	FindByRideID(ctx context.Context, rideID string) ([]db.Match, error)
	FindByRiderID(ctx context.Context, riderID string) ([]db.Match, error)
	FindActiveByRide(ctx context.Context, rideID string) (*db.Match, error)
	ListByDriverID(ctx context.Context, driverID string, limit int) ([]db.Match, error)
}

type matchRepository struct {
	db *gorm.DB
}

func NewMatchRepository(db *gorm.DB) MatchRepository {
	return &matchRepository{db: db}
}

func (r *matchRepository) Create(ctx context.Context, match *db.Match) error {
	if match == nil {
		return errors.New("match is nil")
	}
	return r.db.WithContext(ctx).Create(match).Error
}

func (r *matchRepository) FindByID(ctx context.Context, id string) (*db.Match, error) {
	if id == "" {
		return nil, nil
	}
	var out db.Match
	err := r.db.WithContext(ctx).Where("id = ?", id).Take(&out).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &out, err
}

func (r *matchRepository) FindByRideID(ctx context.Context, rideID string) ([]db.Match, error) {
	if rideID == "" {
		return []db.Match{}, nil
	}
	var out []db.Match
	err := r.db.WithContext(ctx).
		Where("ride_id = ?", rideID).
		Order("created_at DESC").
		Find(&out).Error
	return out, err
}

func (r *matchRepository) FindByRiderID(ctx context.Context, riderID string) ([]db.Match, error) {
	if riderID == "" {
		return []db.Match{}, nil
	}
	var out []db.Match
	err := r.db.WithContext(ctx).
		Where("rider_id = ?", riderID).
		Order("created_at DESC").
		Find(&out).Error
	return out, err
}

func (r *matchRepository) UpdateStatus(ctx context.Context, matchID string, status string) error {
	if matchID == "" || status == "" {
		return errors.New("matchID and status required")
	}
	return r.db.WithContext(ctx).
		Model(&db.Match{}).
		Where("id = ?", matchID).
		Update("status", status).Error
}


func (r *matchRepository) FindActiveByRide(ctx context.Context, rideID string) (*db.Match, error) {
	if rideID == "" {
		return nil, nil
	}
	var out db.Match
	err := r.db.WithContext(ctx).
		Where("ride_id = ? AND status = ?", rideID, "accepted").
		Order("created_at DESC").
		Take(&out).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &out, err
}

func (r *matchRepository) ListByDriverID(ctx context.Context, driverID string, limit int) ([]db.Match, error) {
	if driverID == "" {
		return []db.Match{}, nil
	}
	var out []db.Match
	q := r.db.WithContext(ctx).
		Where("driver_id = ?", driverID).
		Order("created_at DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	err := q.Find(&out).Error
	return out, err
}
