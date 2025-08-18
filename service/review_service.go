package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"hope/db"
	"hope/repository"
	"strings"
	"time"
)

type ReviewService interface {
	SubmitReview(ctx context.Context, review *db.Review) error
	ListReviewsByUser(ctx context.Context, userID string, limit int) ([]db.Review, error)
	ListReviewsByRide(ctx context.Context, rideID string) ([]db.Review, error)
	DeleteReview(ctx context.Context, reviewID string) error
}

type reviewService struct {
	reviewrepo repository.ReviewRepository
}

func NewReviewService(reviewrepo repository.ReviewRepository) ReviewService {
	return &reviewService{reviewrepo: reviewrepo}
}

func (s reviewService) SubmitReview(ctx context.Context, review *db.Review) error {
	if review == nil {
		return errors.New("invalid review")
	}
	review.RideID = strings.TrimSpace(review.RideID)
	review.FromUserID = strings.TrimSpace(review.FromUserID)
	review.ToUserID = strings.TrimSpace(review.ToUserID)
	review.Comment = strings.TrimSpace(review.Comment)

	if review.RideID == "" || review.FromUserID == "" || review.ToUserID == "" || review.Score < 1 || review.Score > 5 {
		return errors.New("invalid review fields")
	}

	if review.FromUserID == review.ToUserID {
		return errors.New("cannot review yourself")
	}

	review.ID = uuid.New().String()
	review.CreatedAt = time.Now().UTC()

	return s.reviewrepo.Create(ctx, review)
}

func (s reviewService) ListReviewsByUser(ctx context.Context, userID string, limit int) ([]db.Review, error) {
	return s.reviewrepo.ListByUser(ctx, strings.TrimSpace(userID), limit)
}

func (s reviewService) ListReviewsByRide(ctx context.Context, rideID string) ([]db.Review, error) {
	return s.reviewrepo.ListByRide(ctx, strings.TrimSpace(rideID))
}

func (s reviewService) DeleteReview(ctx context.Context, reviewID string) error {
	return s.reviewrepo.Delete(ctx, strings.TrimSpace(reviewID))
}
