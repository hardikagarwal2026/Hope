package service

import (
	"context"
	"errors"
	"hope/db"
	"hope/repository"
	"strings"
	"time"
)

var (
	errInvalidLatLon    = errors.New("invalid latitude or longitude")
	errLocationNotFound = errors.New("location not found")
)

type LocationService interface {
	UpsertLocation(ctx context.Context, loc *db.UserLocation) error
	GetLocationByUser(ctx context.Context, userID string) (*db.UserLocation, error)
	ListNearby(ctx context.Context, geohashPrefix string, limit int) ([]db.UserLocation, error)
	DeleteLocation(ctx context.Context, userID string) error
}

type locationService struct {
	locationrepo repository.UserLocationRepository
}

func NewLocationService(locationrepo repository.UserLocationRepository) LocationService {
	return &locationService{locationrepo: locationrepo}
}

func (s locationService) UpsertLocation(ctx context.Context, loc *db.UserLocation) error {
	if loc == nil {
		return errInvalidLatLon
	}

	loc.UserID = strings.TrimSpace(loc.UserID)
	loc.Geohash = strings.TrimSpace(loc.Geohash)

	if loc.Latitude < -90 || loc.Latitude > 90 || loc.Longitude < -180 || loc.Longitude > 180 {
		return errInvalidLatLon
	}
	loc.UpdatedAt = time.Now().UTC()

	return s.locationrepo.Upsert(ctx, loc)
}

func (s locationService) GetLocationByUser(ctx context.Context, userID string) (*db.UserLocation, error) {
	userID = strings.TrimSpace(userID)

	loc, err := s.locationrepo.GetByUserID(ctx, userID)
	if err != nil || loc == nil {
		return nil, errLocationNotFound
	}
	return loc, nil
}

func (s locationService) ListNearby(ctx context.Context, geohashPrefix string, limit int) ([]db.UserLocation, error) {
	return s.locationrepo.ListNearby(ctx, strings.TrimSpace(geohashPrefix), limit)
}

func (s locationService) DeleteLocation(ctx context.Context, userID string) error {
	return s.locationrepo.Delete(ctx, strings.TrimSpace(userID))
}
