package service

import (
	"context"
	"errors"
	"hope/db"
	"hope/repository"
	"time"
)

type ReviewService interface{
	SubmitReview(ctx context.Context, review *db.Review) error
    ListReviewsByUser(ctx context.Context, userID string, limit int) ([]db.Review, error)
    ListReviewsByRide(ctx context.Context, rideID string) ([]db.Review, error)
    DeleteReview(ctx context.Context, reviewID string) error

}


type reviewService struct {
    reviewrepo repository.ReviewRepository
}

//you know why we write Constructor,  so we can wire dependencies later and can use these services, later by creating its instance
func NewReviewService(reviewrepo repository.ReviewRepository) ReviewService {
	return &reviewService{reviewrepo: reviewrepo}
}

// to review the ride , to submit
func (s reviewService) SubmitReview(ctx context.Context, review *db.Review) error {
	if review.RideID == "" || review.FromUserID == "" || review.ToUserID == "" || review.Score < 1 || review.Score > 5 {
		return errors.New("invalid review fields")
	}

	//cannot review yourself baby
	if review.FromUserID == review.ToUserID {
		return errors.New("cannot review yourself")
	}

	//setting review time
	review.CreatedAt = time.Now()
	return s.reviewrepo.Create(ctx, review)

}

// toget all the reviews by a particular user
func (s reviewService) ListReviewsByUser(ctx context.Context, userID string, limit int) ([]db.Review, error) {
	return s.reviewrepo.ListByUser(ctx, userID, limit)
}

//to list all the review of a particular ride by its rideID
func (s reviewService) ListReviewsByRide(ctx context.Context, rideID string) ([]db.Review, error){
	return s.reviewrepo.ListByRide(ctx, rideID)
}


func (s reviewService) DeleteReview(ctx context.Context, reviewID string) error {
	return s.reviewrepo.Delete(ctx, reviewID)
}

