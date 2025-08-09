package service

import (
	"context"
	"errors"
	"time"
	"hope/db"
	"hope/repository"
)

// methods related to rides service
type RideService interface {
	CreateOffer(ctx context.Context, offer *db.RideOffer) error
    ListNearbyOffers(ctx context.Context, geohashPrefix string, limit int) ([]db.RideOffer, error)
    GetOfferByID(ctx context.Context, id string) (*db.RideOffer, error)
    UpdateOffer(ctx context.Context, offer *db.RideOffer) error
    DeleteOffer(ctx context.Context, id string) error

    CreateRequest(ctx context.Context, req *db.RideRequest) error
    ListNearbyRequests(ctx context.Context, geohashPrefix string, limit int) ([]db.RideRequest, error)
    GetRequestByID(ctx context.Context, id string) (*db.RideRequest, error)
    UpdateRequestStatus(ctx context.Context, id string, status string) error
    DeleteRequest(ctx context.Context, id string) error
}

// struct having all the repo dependencies, rideoffer, riderequest, user
type rideService struct {
	rideofferepo 		repository.RideOfferRepository
	riderequestrepo 	repository.RideRequestRepository
	userrepo 			repository.UserRepository
}

// constructor for getting new ride service
func NewRideService(rideofferepo  repository.RideOfferRepository,riderequestrepo repository.RideRequestRepository, userrepo repository.UserRepository) RideService {
	return &rideService{rideofferepo: rideofferepo, riderequestrepo: riderequestrepo, userrepo: userrepo}
}

func (s rideService) CreateOffer(ctx context.Context, offer *db.RideOffer) error {
	if offer.DriverID == "" || offer.FromGeo == "" || offer.ToGeo == "" {
		return errors.New("missing requiredfields")
	}

	//checking the ride offer time, as it cant be in the past
	 if offer.Time.Before(time.Now()) {
        return errors.New("ride offer time cannot be in the past")
    }

	if offer.Seats <= 0 {
        return errors.New("seats must be positive")
    }

	//checking order status, and setting as active, if nothing
	if offer.Status == "" {
        offer.Status = "active"
    }

	//to check if the driver exists, as he is also one of the user, so we check in user table
	driver, _ := s.userrepo.FindByID(ctx, offer.DriverID)
	if driver == nil || driver.ID == "" {
        return errors.New("invalid driver")
    }
    return s.rideofferepo.Create(ctx, offer)


}

//list all the nearby offers from the geohash
func (s rideService) ListNearbyOffers(ctx context.Context, geohashPrefix string, limit int) ([]db.RideOffer, error) {
	return s.rideofferepo.ListNearbyOffers(ctx, geohashPrefix, limit)
}


// to get offer by id
func (s rideService) GetOfferByID(ctx context.Context, id string) (*db.RideOffer, error) {
	offer, err := s.rideofferepo.FindByID(ctx, id)
    if err != nil || offer == nil || offer.ID == "" {
        return nil, errors.New("offer not found")
    }
    return offer, nil
}

// to update the rideoffer, like fare, seats and status that can be changed aftr as well
func (s rideService) UpdateOffer(ctx context.Context, offer *db.RideOffer) error {
	current, err := s.rideofferepo.FindByID(ctx, offer.ID)
    if err != nil || current == nil || current.ID == "" {
        return errors.New("offer not found")
    }
    // Only update mutable fields
    if offer.Seats != 0 {
        current.Seats = offer.Seats
    }
    if offer.Fare != "" {
        current.Fare = offer.Fare
    }
    if offer.Status != "" {
        current.Status = offer.Status
    }
    return s.rideofferepo.Update(ctx, current)
}

// to delete the ride offer
func (s rideService) DeleteOffer(ctx context.Context, id string) error {
	return s.rideofferepo.Delete(ctx, id)
}

// to create request or ride
func (s rideService) CreateRequest(ctx context.Context, req *db.RideRequest) error {
	if req.UserID == "" || req.FromGeo == "" || req.ToGeo == "" {
        return errors.New("missing required fields")
    }

	//ride request time checking
    if req.Time.Before(time.Now()) {
        return errors.New("ride request time cannot be in the past")
    }
    if req.Seats <= 0 {
        return errors.New("seats must be positive")
    }
    if req.Status == "" {
        req.Status = "active"
    }

	//to check if the user exists, who is requesting
    user, _ := s.userrepo.FindByID(ctx, req.UserID)
    if user == nil || user.ID == "" {
        return errors.New("invalid user")
    }
    return s.riderequestrepo.Create(ctx, req)
}

//list all the nearbyrequests of the ride
func (s rideService) ListNearbyRequests(ctx context.Context, geohashPrefix string, limit int) ([]db.RideRequest, error) {
	 return s.riderequestrepo.ListNearby(ctx, geohashPrefix, limit)
}

// to find the paricular request by its id
func (s rideService) GetRequestByID(ctx context.Context, id string) (*db.RideRequest, error){
	req, err := s.riderequestrepo.FindByID(ctx, id)
    if err != nil || req == nil || req.ID == "" {
        return nil, errors.New("request not found")
    }
    return req, nil
}

//to update the status of the ride requests
func (s rideService) UpdateRequestStatus(ctx context.Context, id string, status string) error{
	return s.riderequestrepo.UpdateStatus(ctx, id, status)
}

//to delete the ride request from the id
func (s rideService) DeleteRequest(ctx context.Context, id string) error {
	return s.riderequestrepo.Delete(ctx, id)
}


