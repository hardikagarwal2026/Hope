package db

import "time"

type UserLocation struct {
	UserID    string `gorm:"primaryKey;size:191"`
	Latitude  float64
	Longitude float64
	Geohash   string    `gorm:"size:64;index"`
	UpdatedAt time.Time `gorm:"index"`

	User *User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
