package db

import "time"

// user is the student here, who is using the platform
// primary key is the main identifier
type User struct {
	ID       string    `gorm:"primaryKey"`           //will get it unique from uuid
	Name     string    `gorm:"size:191"`             // name of the student
	Email    string    `gorm:"uniqueIndex;size:191"` //unique for all students, it is not the main identifier still no two rows can have same email
	PhotoURL string    `gorm:"size:191"`             // url of the photo, will store in mongo later
	Geohash  string    `gorm:"size:64"`              // lat, long into geohash for proximity search
	LastSeen time.Time //last time online
}
