package db

import "time"

// match is for linking the rideoffer and a rider's request
type Match struct {
	ID			 string     `gorm:"primaryKey;autoIncrement"` // unique id for every match, UUID so string
	RiderID      string		                    // User.ID of the rider
	DriverID     string							// User.ID of the driver
	RideId       string							// foreignKey RideOffer.ID, links to which ride is matched(pointin RideOffer Table)
	Status       string							// requested, accepted, rejected or completed 
	CreatedAt    time.Time                      // time when match created
}
