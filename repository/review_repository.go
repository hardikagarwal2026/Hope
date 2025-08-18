package repository

import (
	"context"
	"errors"
	"hope/db"
	"gorm.io/gorm"
)


type ReviewRepository interface {
	Create(ctx context.Context, review *db.Review) error
	ListByUser(ctx context.Context, userID string, limit int) ([]db.Review, error)
	ListByRide(ctx context.Context, rideID string) ([]db.Review, error)
	Delete(ctx context.Context, reviewID string) error
	ListReceivedByUser(ctx context.Context, userID string, limit int) ([]db.Review, error)
	GetByID(ctx context.Context, id string) (*db.Review, error)
}

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(ctx context.Context, review *db.Review) error {
	if review == nil {
		return errors.New("review is nil")
	}
	return r.db.WithContext(ctx).Create(review).Error
}

func (r *reviewRepository) ListByUser(ctx context.Context, userID string, limit int) ([]db.Review, error) {
	var out []db.Review
	q := r.db.WithContext(ctx).
		Where("from_user_id = ?", userID).
		Order("created_at DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	err := q.Find(&out).Error
	return out, err
}

func (r *reviewRepository) ListByRide(ctx context.Context, rideID string) ([]db.Review, error) {
	var out []db.Review
	err := r.db.WithContext(ctx).
		Where("ride_id = ?", rideID).
		Order("created_at DESC").
		Find(&out).Error
	return out, err
}

func (r *reviewRepository) Delete(ctx context.Context, reviewID string) error {
	if reviewID == "" {
		return errors.New("reviewID required")
	}
	// Use primary key column "id" (matches GORM default for your struct).
	return r.db.WithContext(ctx).
		Delete(&db.Review{}, "id = ?", reviewID).Error
}

func (r *reviewRepository) ListReceivedByUser(ctx context.Context, userID string, limit int) ([]db.Review, error) {
	var out []db.Review
	q := r.db.WithContext(ctx).
		Where("to_user_id = ?", userID).
		Order("created_at DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	err := q.Find(&out).Error
	return out, err
}

func (r *reviewRepository) GetByID(ctx context.Context, id string) (*db.Review, error) {
	if id == "" {
		return nil, nil
	}
	var out db.Review
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Take(&out).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &out, err
}
