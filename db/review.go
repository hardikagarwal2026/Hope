package db

import "time"

type Review struct {
	ID         string `gorm:"primaryKey;size:191"`
	RideID     string `gorm:"size:191;index"`
	FromUserID string `gorm:"size:191;index"`
	ToUserID   string `gorm:"size:191;index"`
	Score      int
	Comment    string    `gorm:"type:text"`
	CreatedAt  time.Time `gorm:"index"`

	Ride     *RideOffer `gorm:"foreignKey:RideID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	FromUser *User      `gorm:"foreignKey:FromUserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	ToUser   *User      `gorm:"foreignKey:ToUserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
