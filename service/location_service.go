package service

import (
	"context"
	"errors"
	"hope/db"
	"hope/repository"
	"time"
)

// all methods defined for location service
type LocationService interface {
	UpsertLocation(ctx context.Context, loc *db.UserLocation) error
    GetLocationByUser(ctx context.Context, userID string) (*db.UserLocation, error)
    ListNearby(ctx context.Context, geohashPrefix string, limit int) ([]db.UserLocation, error)
    DeleteLocation(ctx context.Context, userID string) error
}

// struct implementing location service
type locationService struct {
	locationrepo repository.UserLocationRepository
}


//constructor to get the new location service
func NewLocationService(locationrepo repository.UserLocationRepository) LocationService {
	return &locationService{locationrepo: locationrepo}
}

// to insert or update the location of the user
func (s locationService) UpsertLocation(ctx context.Context, loc *db.UserLocation) error {
    loc.Updatedat = time.Now()
	if loc.Latitude < -90 || loc.Latitude > 90 || loc.Longitude < -180 || loc.Longitude >  180 {
		return errors.New("invalid latitude and longitude")
	}

	return s.locationrepo.Upsert(ctx, loc)
}


//to get the location of the user
func (s locationService) GetLocationByUser(ctx context.Context, userID string) (*db.UserLocation, error) {
	loc, err := s.locationrepo.GetByUserID(ctx, userID)
	if err != nil || loc == nil {
		return nil, errors.New("location not found")
	}
	return loc, nil
}

// to get the ist of all the nearby users from geohash with limit
func(s locationService) ListNearby(ctx context.Context, geohashPrefix string, limit int) ([]db.UserLocation, error) {
	return s.locationrepo.ListNearby(ctx, geohashPrefix, limit)
}


// to delete the location of the user with userID
func(s locationService) DeleteLocation(ctx context.Context, userID string) error {
	return s.locationrepo.Delete(ctx, userID)
}

