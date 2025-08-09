package service

import (
	"context"
	"errors"
	"hope/db"
	"hope/repository"
	"time"
) 

type ChatService interface {
	SendMessage(ctx context.Context, msg *db.ChatMessage) error
    ListMessagesByRide(ctx context.Context, rideID string, limit int, before time.Time) ([]db.ChatMessage, error)
    ListMessagesBySender(ctx context.Context, senderID string, limit int, before time.Time) ([]db.ChatMessage, error)
    ListChatsForUser(ctx context.Context, userID string, limit int, before time.Time) ([]db.ChatMessage, error)
    DeleteMessage(ctx context.Context, id string) error
}

type chatService struct {
	chatrepo repository.ChatMessageRepository
	matchrepo repository.MatchRepository
}

func NewChatService(chatrepo repository.ChatMessageRepository, matchrepo repository.MatchRepository) ChatService {
	return &chatService{chatrepo: chatrepo, matchrepo: matchrepo}
}

func (s chatService) SendMessage(ctx context.Context, msg *db.ChatMessage) error {
	if msg.RideID == "" || msg.SenderID == "" || msg.Content == "" {
		return errors.New("rideID, senderID and content required")
	}

	//checking only matched users for this ride can chat
	matches, err := s.matchrepo.FindByRideID(ctx, msg.RideID) //getting the match by rideID
	if err!= nil{
		return err
	}

	//default allowed false
	//then checking 
	allowed := false
	for _, m := range matches{
		// both ridder and rider can be a sender in the chat
		/// AND match status must be accepted or completed
		if (m.RiderID == msg.SenderID || m.DriverID == msg.SenderID) && (m.Status == "accepted" || m.Status == "completed") {
			allowed = true
			break
		}
	}
	if !allowed{
		return errors.New("user notallowed to chat for this ride")
	}

	//if alllowed then sending message
	msg.Timestamp = time.Now()
	return s.chatrepo.Create(ctx, msg)

}

// list messages asssociated with a particular ride using rideID between rider and driver both
func (s chatService) ListMessagesByRide(ctx context.Context, rideID string, limit int, before time.Time) ([]db.ChatMessage, error){
	return s.chatrepo.ListByRide(ctx, rideID, limit, before)
}

//list all the messages by sender by senderID, whether its anyone rider or driver
func (s chatService)ListMessagesBySender(ctx context.Context, senderID string, limit int, before time.Time) ([]db.ChatMessage, error){
	return s.chatrepo.ListBySender(ctx, senderID, limit, before)
}

//list all the chat for user
//have to go through again through this function
func (s chatService) ListChatsForUser(ctx context.Context, userID string, limit int, before time.Time) ([]db.ChatMessage, error){
	return s.chatrepo.ListChatsForUser(ctx, userID, limit, before)

}

// to delete the particular chat by its id
func (s chatService) DeleteMessage(ctx context.Context, id string) error {
	return s.chatrepo.Delete(ctx, id)
}



