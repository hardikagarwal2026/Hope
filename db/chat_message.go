package db

import "time"

type ChatMessage struct {
	ID        string    `gorm:"primaryKey;size:191"`
	RideID    string    `gorm:"size:191"`
	SenderID  string    `gorm:"size:191"`
	Content   string    `gorm:"type:text"`
	Timestamp time.Time `gorm:"index"`

	Ride   *RideOffer `gorm:"foreignKey:RideID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Sender *User      `gorm:"foreignKey:SenderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
