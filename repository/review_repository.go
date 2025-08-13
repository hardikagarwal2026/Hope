package repository

import (
	"context"
	"hope/db"
	"gorm.io/gorm"
)

type ReviewRepository interface {
	Create(ctx context.Context, review *db.Review) error
    ListByUser(ctx context.Context, userID string, limit int) ([]db.Review, error)
    ListByRide(ctx context.Context, rideID string) ([]db.Review, error)
    Delete(ctx context.Context, reviewID string) error
}

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(ctx context.Context, review *db.Review) error {
	return r.db.WithContext(ctx).Create(review).Error
}

func (r *reviewRepository) ListByUser(ctx context.Context, userID string, limit int) ([]db.Review, error) {
	var reviews []db.Review
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Find(&reviews).Error
	return reviews, err
}

func (r *reviewRepository) ListByRide(ctx context.Context, rideID string) ([]db.Review, error) {
	var reviews []db.Review
	err := r.db.WithContext(ctx).Where("ride_id = ?", rideID).Order("created_at DESC").Find(&reviews).Error
	return reviews, err
}

func (r *reviewRepository) Delete(ctx context.Context, reviewID string) error {
	return r.db.WithContext(ctx).Delete(&db.Review{}, "review_id = ?", reviewID).Error
}
