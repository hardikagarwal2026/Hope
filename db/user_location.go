package db

import "time"

type UserLocation struct {
	UserID    string `gorm:"primaryKey"` //foreign key to User.ID
	Latitude  float64
	Longitude float64
	Geohash   string
	Updatedat time.Time
}



