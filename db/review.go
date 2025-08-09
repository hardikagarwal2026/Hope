package db

import "time"

type Review struct {
	ID  		string      `gorm:"primaryKey"` // unique review id
	RideID      string                          // RideOffer.ID of the ride being reviewed
	FromUserID  string                          // User.ID who is writing
	ToUserID    string                          // User.ID who is being reviewed
	Score       int                             // 1-5 based on 
	Comment     string                          // feedback of the ride
	CreatedAt   time.Time                       // when the review was created
}


