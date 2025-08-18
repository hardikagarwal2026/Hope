package db

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type Match struct {
	ID        string    `gorm:"primaryKey;size:191" json:"id"`
	RiderID   string    `gorm:"size:191;index"      json:"rider_id"`
	DriverID  string    `gorm:"size:191;index"      json:"driver_id"`
	RideID    string    `gorm:"size:191;index"      json:"ride_id"`
	Status    string    `gorm:"size:32;index"       json:"status"`
	CreatedAt time.Time `gorm:"index"               json:"created_at"`

	Rider  *User      `gorm:"foreignKey:RiderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	Driver *User      `gorm:"foreignKey:DriverID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	Ride   *RideOffer `gorm:"foreignKey:RideID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"   json:"-"`
}

func (m *Match) BeforeCreate(tx *gorm.DB) (err error) {

	if strings.TrimSpace(m.Status) == "" {
		m.Status = "requested"
	}

	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	return nil
}

func (m *Match) BeforeSave(tx *gorm.DB) (err error) {
	m.RiderID = strings.TrimSpace(m.RiderID)
	m.DriverID = strings.TrimSpace(m.DriverID)
	m.RideID = strings.TrimSpace(m.RideID)
	m.Status = strings.TrimSpace(m.Status)
	return nil
}
