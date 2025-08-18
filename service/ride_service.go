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

var (
	errOfferNotFound   = errors.New("offer not found")
	errRequestNotFound = errors.New("request not found")
	errMissingFields   = errors.New("missing required fields")
	errPastTime        = errors.New("time cannot be in the past")
	errSeatsPositive   = errors.New("seats must be positive")
	errInvalidDriver   = errors.New("invalid driver")
	errInvalidUser     = errors.New("invalid user")
)

type RideService interface {
	CreateOffer(ctx context.Context, offer *db.RideOffer) error
	ListNearbyOffers(ctx context.Context, geohashPrefix string, limit int) ([]db.RideOffer, error)
	GetOfferByID(ctx context.Context, id string) (*db.RideOffer, error)
	UpdateOffer(ctx context.Context, offer *db.RideOffer) error
	DeleteOffer(ctx context.Context, id string) error
	ListMyOffers(ctx context.Context, driverID string, limit int) ([]db.RideOffer, error)

	CreateRequest(ctx context.Context, req *db.RideRequest) error
	ListNearbyRequests(ctx context.Context, geohashPrefix string, limit int) ([]db.RideRequest, error)
	GetRequestByID(ctx context.Context, id string) (*db.RideRequest, error)
	UpdateRequestStatus(ctx context.Context, id string, status string) error
	DeleteRequest(ctx context.Context, id string) error
	ListMyRequests(ctx context.Context, userID string, limit int) ([]db.RideRequest, error)
}

type rideService struct {
	rideofferepo    repository.RideOfferRepository
	riderequestrepo repository.RideRequestRepository
	userrepo        repository.UserRepository
}

func NewRideService(rideofferepo repository.RideOfferRepository, riderequestrepo repository.RideRequestRepository, userrepo repository.UserRepository) RideService {
	return &rideService{rideofferepo: rideofferepo, riderequestrepo: riderequestrepo, userrepo: userrepo}
}

func (s rideService) CreateOffer(ctx context.Context, offer *db.RideOffer) error {
	if offer == nil {
		return errMissingFields
	}
	offer.DriverID = strings.TrimSpace(offer.DriverID)
	offer.FromGeo = strings.TrimSpace(offer.FromGeo)
	offer.ToGeo = strings.TrimSpace(offer.ToGeo)

	if offer.DriverID == "" || offer.FromGeo == "" || offer.ToGeo == "" {
		return errMissingFields
	}
	now := time.Now().UTC()
	if offer.Time.Before(now) {
		return errPastTime
	}
	if offer.Seats <= 0 {
		return errSeatsPositive
	}
	if strings.TrimSpace(offer.Status) == "" {
		offer.Status = "active"
	}

	offer.ID = uuid.New().String()

	driver, _ := s.userrepo.FindByID(ctx, offer.DriverID)
	if driver == nil || driver.ID == "" {
		return errInvalidDriver
	}

	return s.rideofferepo.Create(ctx, offer)
}

func (s rideService) ListNearbyOffers(ctx context.Context, geohashPrefix string, limit int) ([]db.RideOffer, error) {
	return s.rideofferepo.ListNearbyOffers(ctx, strings.TrimSpace(geohashPrefix), limit)
}

func (s rideService) GetOfferByID(ctx context.Context, id string) (*db.RideOffer, error) {
	id = strings.TrimSpace(id)
	o, err := s.rideofferepo.FindByID(ctx, id)
	if err != nil || o == nil || o.ID == "" {
		return nil, errOfferNotFound
	}
	return o, nil
}

func (s rideService) UpdateOffer(ctx context.Context, offer *db.RideOffer) error {
	if offer == nil || strings.TrimSpace(offer.ID) == "" {
		return errOfferNotFound
	}
	current, err := s.rideofferepo.FindByID(ctx, strings.TrimSpace(offer.ID))
	if err != nil || current == nil || current.ID == "" {
		return errOfferNotFound
	}

	if offer.Seats > 0 {
		current.Seats = offer.Seats
	}

	if strings.TrimSpace(offer.Status) != "" {
		current.Status = strings.TrimSpace(offer.Status)
	}

	return s.rideofferepo.Update(ctx, current)
}

func (s rideService) DeleteOffer(ctx context.Context, id string) error {
	return s.rideofferepo.Delete(ctx, strings.TrimSpace(id))
}

func (s rideService) ListMyOffers(ctx context.Context, driverID string, limit int) ([]db.RideOffer, error) {
	if driverID == "" {
		return nil, errors.New("driverID required")
	}
	return s.rideofferepo.ListByDriver(ctx, driverID, limit)
}

func (s rideService) CreateRequest(ctx context.Context, req *db.RideRequest) error {
	if req == nil {
		return errMissingFields
	}
	req.UserID = strings.TrimSpace(req.UserID)
	req.FromGeo = strings.TrimSpace(req.FromGeo)
	req.ToGeo = strings.TrimSpace(req.ToGeo)

	if req.UserID == "" || req.FromGeo == "" || req.ToGeo == "" {
		return errMissingFields
	}
	now := time.Now().UTC()
	if req.Time.Before(now) {
		return errPastTime
	}
	if req.Seats <= 0 {
		return errSeatsPositive
	}
	if strings.TrimSpace(req.Status) == "" {
		req.Status = "active"
	}

	req.ID = uuid.New().String()

	u, _ := s.userrepo.FindByID(ctx, req.UserID)
	if u == nil || u.ID == "" {
		return errInvalidUser
	}

	return s.riderequestrepo.Create(ctx, req)
}

func (s rideService) ListNearbyRequests(ctx context.Context, geohashPrefix string, limit int) ([]db.RideRequest, error) {
	return s.riderequestrepo.ListNearby(ctx, strings.TrimSpace(geohashPrefix), limit)
}

func (s rideService) GetRequestByID(ctx context.Context, id string) (*db.RideRequest, error) {
	id = strings.TrimSpace(id)
	r, err := s.riderequestrepo.FindByID(ctx, id)
	if err != nil || r == nil || r.ID == "" {
		return nil, errRequestNotFound
	}
	return r, nil
}

func (s rideService) UpdateRequestStatus(ctx context.Context, id string, status string) error {
	return s.riderequestrepo.UpdateStatus(ctx, strings.TrimSpace(id), strings.TrimSpace(status))
}

func (s rideService) DeleteRequest(ctx context.Context, id string) error {
	return s.riderequestrepo.Delete(ctx, strings.TrimSpace(id))
}

func (s rideService) ListMyRequests(ctx context.Context, userID string, limit int) ([]db.RideRequest, error) {
	if userID == "" {
		return nil, errors.New("userID required")
	}
	return s.riderequestrepo.ListByUser(ctx, userID, limit)
}
