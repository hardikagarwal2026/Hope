package repository

import (
	"context"
	"hope/db"
	"gorm.io/gorm"
)


type RideOfferRepository interface {
	Create(ctx context.Context, offer *db.RideOffer) error
	ListNearbyOffers(ctx context.Context, geohashprefix string, limit int) ([]db.RideOffer, error)
	FindByID(ctx context.Context, id string)(*db.RideOffer, error)
	Update(ctx context.Context, offer *db.RideOffer) error
	Delete(ctx context.Context, id string)error
}

type rideOfferRepository struct{
	db *gorm.DB // gorm.DB object, it embeds a pointer to the gorm DB object
}

func NewrideOfferRepository(db *gorm.DB) RideOfferRepository {
	return &rideOfferRepository{db: db}
}


// to save the new ride offer to the database
func (r *rideOfferRepository) Create(ctx context.Context, offer *db.RideOffer) error {
    return r.db.WithContext(ctx).Create(offer).Error
}


// to list all the nearby ride offers near starting point 
func (r *rideOfferRepository) ListNearbyOffers(ctx context.Context, geohashprefix string, limit int)([]db.RideOffer, error){
	var offers []db.RideOffer
	err := r.db.WithContext(ctx).Where("from_geo LIKE ?", geohashprefix+"%").Limit(limit).Find(&offers).Error
	// SELECT * FROM ride_offers WHERE from_geo LIKE 'prefix%'
	return offers, err
}


// to find the particular offer by its id
func (r *rideOfferRepository) FindByID(ctx context.Context, id string)(*db.RideOffer, error){
	var offer db.RideOffer
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&offer).Error
	return &offer, err
}

// to update the existing offer in the db
func (r *rideOfferRepository) Update(ctx context.Context, offer *db.RideOffer) error {
	return r.db.WithContext(ctx).Save(offer).Error
}


// to remove the ride offer from the db
func (r *rideOfferRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&db.RideOffer{}, "id = ?", id).Error
}

