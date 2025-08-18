package db

import (
	"gorm.io/gorm"
	"strings"
	"time"
)

type User struct {
	ID       string    `gorm:"primaryKey;size:191"`
	Name     string    `gorm:"size:191"`
	Email    string    `gorm:"uniqueIndex;size:191"`
	PhotoURL string    `gorm:"size:191"`
	Geohash  string    `gorm:"size:64;index"`
	LastSeen time.Time `gorm:"index"`

	Location *UserLocation `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	RideOffers   []RideOffer   `gorm:"foreignKey:DriverID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	RideRequests []RideRequest `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {

	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
	if u.LastSeen.IsZero() {
		u.LastSeen = time.Now()
	}
	return nil
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
	return nil
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	return nil
}
