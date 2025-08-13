package repository

import(
	"context"
	"hope/db"
	"gorm.io/gorm"
)


type RideRequestRepository interface {
	Create(ctx context.Context, req *db.RideRequest) error
    FindByID(ctx context.Context, id string) (*db.RideRequest, error)
    ListNearby(ctx context.Context, geohashPrefix string, limit int) ([]db.RideRequest, error)
    ListByUser(ctx context.Context, userID string, limit int) ([]db.RideRequest, error)
    UpdateStatus(ctx context.Context, id string, status string) error
    Delete(ctx context.Context, id string) error
}


type rideRequestRepository struct {
	db *gorm.DB
}

func NewRideRequestRepository(db *gorm.DB) RideRequestRepository {
	return &rideRequestRepository{db: db}
}

func (r *rideRequestRepository) Create(ctx context.Context, req *db.RideRequest) error{
	return r.db.WithContext(ctx).Create(req).Error
}


func (r *rideRequestRepository) FindByID(ctx context.Context, id string) (*db.RideRequest, error){
	var riderequest db.RideRequest
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&riderequest).Error
	return &riderequest, err
}

func (r *rideRequestRepository) ListNearby(ctx context.Context, geohashPrefix string, limit int)([]db.RideRequest, error){
	var riderequests []db.RideRequest
	err := r.db.WithContext(ctx).Where("from_geo LIKE ?", geohashPrefix + "%").Order("time ASC").Limit(limit).Find(&riderequests).Error
	return riderequests, err
}

func (r *rideRequestRepository) ListByUser(ctx context.Context, userID string, limit int) ([]db.RideRequest, error) {
	var riderequests []db.RideRequest
	err := r.db.WithContext(ctx).Where("user_id = ?", userID). Order("time DESC").Limit(limit).Find(&riderequests).Error
	return riderequests, err
}

func (r *rideRequestRepository) UpdateStatus(ctx context.Context, id string, status string) error{
	return r.db.WithContext(ctx).Model(&db.RideRequest{}).Where("id = ?", id).Update("status", status).Error		
}

func (r *rideRequestRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&db.RideRequest{}, "id = ?", id).Error
}