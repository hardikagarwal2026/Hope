package repository

import (
	"context"
	"errors"
	"hope/db"
	"gorm.io/gorm"
)

type RideOfferRepository interface {
	Create(ctx context.Context, offer *db.RideOffer) error
	FindByID(ctx context.Context, id string) (*db.RideOffer, error)
	Update(ctx context.Context, offer *db.RideOffer) error
	Delete(ctx context.Context, id string) error
	ListNearbyOffers(ctx context.Context, geohashPrefix string, limit int) ([]db.RideOffer, error)
	FindByIDWithDriver(ctx context.Context, id string) (*db.RideOffer, error)
	ListDriverActiveOffers(ctx context.Context, driverID string, limit int) ([]db.RideOffer, error)
	ListByDriver(ctx context.Context, driverID string, limit int) ([]db.RideOffer, error)
}

type rideOfferRepository struct {
	db *gorm.DB
}

func NewrideOfferRepository(db *gorm.DB) RideOfferRepository {
	return &rideOfferRepository{db: db}
}


func (r *rideOfferRepository) Create(ctx context.Context, offer *db.RideOffer) error {
	if offer == nil {
		return errors.New("offer is nil")
	}
	return r.db.WithContext(ctx).Create(offer).Error
}

func (r *rideOfferRepository) FindByID(ctx context.Context, id string) (*db.RideOffer, error) {
	if id == "" {
		return nil, nil
	}
	var out db.RideOffer
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Take(&out).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &out, err
}

func (r *rideOfferRepository) Update(ctx context.Context, offer *db.RideOffer) error {
	if offer == nil || offer.ID == "" {
		return errors.New("offer or ID missing")
	}
	return r.db.WithContext(ctx).Save(offer).Error
}

func (r *rideOfferRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id required")
	}
	return r.db.WithContext(ctx).
		Delete(&db.RideOffer{}, "id = ?", id).Error
}

func (r *rideOfferRepository) ListNearbyOffers(ctx context.Context, geohashPrefix string, limit int) ([]db.RideOffer, error) {
	var offers []db.RideOffer
	q := r.db.WithContext(ctx).
		Where("from_geo LIKE ?", geohashPrefix+"%").
		Order("time ASC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	err := q.Find(&offers).Error
	return offers, err
}

func (r *rideOfferRepository) FindByIDWithDriver(ctx context.Context, id string) (*db.RideOffer, error) {
	if id == "" {
		return nil, nil
	}
	var out db.RideOffer
	err := r.db.WithContext(ctx).
		Preload("Driver").
		Where("id = ?", id).
		Take(&out).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &out, err
}


func (r *rideOfferRepository) ListDriverActiveOffers(ctx context.Context, driverID string, limit int) ([]db.RideOffer, error) {
	var offers []db.RideOffer
	q := r.db.WithContext(ctx).
		Where("driver_id = ? AND status = ?", driverID, "active").
		Order("time DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	err := q.Find(&offers).Error
	return offers, err
}

func (r *rideOfferRepository) ListByDriver(ctx context.Context, driverID string, limit int) ([]db.RideOffer, error) {
	var offers []db.RideOffer
	q := r.db.WithContext(ctx).
		Where("driver_id = ?", driverID).
		Order("time DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	err := q.Find(&offers).Error
	return offers, err

}
