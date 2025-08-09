package db
import "time"

type RideOffer struct {
	ID 			 string 	`gorm:"primaryKey"` // unique Id for any ride offer
	DriverID 	 string	 	// the guy driving , Foreign key points to User.ID	
	FromGeo 	 string		// departure geo hash
	ToGeo 	     string		// arrival geohash
	Fare		 string		// amount for ride offer
	Time		 time.Time  //time for departure
	Seats		 int        // no. of seats available 
	Status 		 string     // status of ride- whether active , matched, completed
}
