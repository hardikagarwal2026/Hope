package db

import "time"

// chatmessage is a message/chat between rider and driver after matching
type ChatMessage struct {
	ID 			string     `gorm:"primaryKey"`   // evey message has unique messageID
	RideID 		string      // RideOffer.ID - which ride this message is about
	SenderID 	string      // who sent the message,User.ID -  who sent the message
	Content 	string      //content of the chat
	Timestamp 	time.Time   //when the mesage was sennt
}