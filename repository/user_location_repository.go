package repository

import (
    "context"
    "hope/db"
	"gorm.io/gorm"
    "time"
	"gorm.io/gorm/clause"
)

type UserLocationRepository interface {
    Upsert(ctx context.Context, loc *db.UserLocation) error
    GetByUserID(ctx context.Context, userID string) (*db.UserLocation, error)
    ListNearby(ctx context.Context, geohashPrefix string, limit int) ([]db.UserLocation, error)
    Delete(ctx context.Context, userID string) error
}

type userLocationRepository struct {
    db *gorm.DB
}

func NewUserLocationRepository(db *gorm.DB) UserLocationRepository {
    return &userLocationRepository{db: db}
}

func (r *userLocationRepository) Upsert(ctx context.Context, loc *db.UserLocation) error {
    loc.Updatedat = time.Now()
    return r.db.WithContext(ctx).
        Clauses(clause.OnConflict{
            Columns:   []clause.Column{{Name: "user_id"}},
            DoUpdates: clause.AssignmentColumns([]string{"latitude", "longitude", "geohash", "updated_at"}),
        }).
        Create(loc).Error
}

func (r *userLocationRepository) GetByUserID(ctx context.Context, userID string) (*db.UserLocation, error) {
    var loc db.UserLocation
    err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&loc).Error
    if err != nil {
        return nil, err
    }
    return &loc, nil
}

func (r *userLocationRepository) ListNearby(ctx context.Context, geohashPrefix string, limit int) ([]db.UserLocation, error) {
    var locations []db.UserLocation
    err := r.db.WithContext(ctx).
        Where("geohash LIKE ?", geohashPrefix+"%").
        Order("updated_at DESC").
        Limit(limit).
        Find(&locations).Error
    return locations, err
}

func (r *userLocationRepository) Delete(ctx context.Context, userID string) error {
    return r.db.WithContext(ctx).Delete(&db.UserLocation{}, "user_id = ?", userID).Error
}
