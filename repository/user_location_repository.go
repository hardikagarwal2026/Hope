package repository

import (
	"context"
	"errors"
	"time"
	"hope/db"
	"gorm.io/gorm"
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
	if loc == nil {
		return errors.New("location is nil")
	}
	loc.UpdatedAt = time.Now()

	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"latitude":   gorm.Expr("VALUES(latitude)"),
				"longitude":  gorm.Expr("VALUES(longitude)"),
				"geohash":    gorm.Expr("VALUES(geohash)"),
				"updated_at": gorm.Expr("VALUES(updated_at)"),
			}),
		}).
		Create(loc).Error
}


func (r *userLocationRepository) GetByUserID(ctx context.Context, userID string) (*db.UserLocation, error) {
	if userID == "" {
		return nil, nil
	}
	var out db.UserLocation
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Take(&out).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &out, err
}


func (r *userLocationRepository) ListNearby(ctx context.Context, geohashPrefix string, limit int) ([]db.UserLocation, error) {
	var out []db.UserLocation
	q := r.db.WithContext(ctx).
		Where("geohash LIKE ?", geohashPrefix+"%").
		Order("updated_at DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	err := q.Find(&out).Error
	return out, err
}

func (r *userLocationRepository) Delete(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("userID required")
	}
	return r.db.WithContext(ctx).
		Delete(&db.UserLocation{}, "user_id = ?", userID).Error
}
