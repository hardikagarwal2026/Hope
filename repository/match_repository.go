package repository

import(
	"context"
	"hope/db"
	"gorm.io/gorm"
)

type MatchRepository interface {
	Create(ctx context.Context, match *db.Match) error
	FindByID(ctx context.Context, id string) (*db.Match, error)
	FindByRideID(ctx context.Context, rideID string)([]db.Match, error)
	FindByRiderID(ctx context.Context, rideID string)([]db.Match, error)
	UpdateStatus(ctx context.Context, matchID string, status string) error
}


type matchRepository struct{
	db *gorm.DB
}


func (r *matchRepository) Create(ctx context.Context, match *db.Match)error{
	return r.db.WithContext(ctx).Create(match).Error
}

func (r *matchRepository) FindByID(ctx context.Context, id string)(*db.Match, error){
	var match db.Match
	err:= r.db.WithContext(ctx).Where("id = ?", id).First(&match).Error
	return &match, err
}

func (r *matchRepository) FindByRideID(ctx context.Context, rideID string)([]db.Match, error){
	var matches []db.Match
	err := r.db.WithContext(ctx).Where("ride_id = ?", rideID).Find(&matches).Error
	return matches, err
}

func (r *matchRepository) FindByRiderID(ctx context.Context, riderID string)([]db.Match, error){
	var matches []db.Match
	err := r.db.WithContext(ctx).Where("rider_id = ?", riderID).Find(&matches).Error
	return matches, err
}

func (r *matchRepository) UpdateStatus(ctx context.Context, matchID string, status string) error {
	return r.db.WithContext(ctx).Model(&db.Match{}).Where("id = ?", matchID).Update("status", status).Error
}

