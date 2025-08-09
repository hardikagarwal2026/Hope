package db

import "time"

type RideRequest struct {
	ID 			string     `gorm:"primaryKey"` // unique id for any ride requested
	UserID 		string      // the guy who requested ride, F.K points to User.ID
	FromGeo 	string		// departure geohash
	ToGeo 		string		//arrival geohash
	Time 		time.Time   //time for ride start
	Seats 		int			//no. of seats available
	Status 		string      // status of ride whether - active, completed, matched 
}




