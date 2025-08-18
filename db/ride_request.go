package db

import (
	"gorm.io/gorm"
	"strings"
	"time"
)

type RideRequest struct {
	ID      string `gorm:"primaryKey;size:191"`
	UserID  string `gorm:"size:191;index"`
	FromGeo string `gorm:"size:64;index"`
	ToGeo   string `gorm:"size:64;index"`
	Fare    float64
	Time    time.Time `gorm:"index"`
	Seats   int
	Status  string `gorm:"size:32;index"`

	Rider *User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (r *RideRequest) BeforeCreate(tx *gorm.DB) (err error) {
	if strings.TrimSpace(r.Status) == "" {
		r.Status = "active"
	}
	if r.Time.IsZero() {
		r.Time = time.Now()
	}
	return nil
}

func (r *RideRequest) BeforeSave(tx *gorm.DB) (err error) {
	r.FromGeo = strings.TrimSpace(r.FromGeo)
	r.ToGeo = strings.TrimSpace(r.ToGeo)
	return nil
}
