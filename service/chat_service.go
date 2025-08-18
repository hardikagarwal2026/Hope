package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"hope/db"
	"hope/repository"
	"strings"
	"time"
)

var (
	errChatInvalidFields = errors.New("ride_id, sender_id and content are required")
	errChatNotAllowed    = errors.New("user not allowed to chat for this ride")
)

type ChatService interface {
	SendMessage(ctx context.Context, msg *db.ChatMessage) error
	ListMessagesByRide(ctx context.Context, rideID string, limit int, before time.Time) ([]db.ChatMessage, error)
	ListMessagesBySender(ctx context.Context, senderID string, limit int, before time.Time) ([]db.ChatMessage, error)
	ListChatsForUser(ctx context.Context, userID string, limit int, before time.Time) ([]db.ChatMessage, error)
	DeleteMessage(ctx context.Context, id string) error
}

type chatService struct {
	chatrepo  repository.ChatMessageRepository
	matchrepo repository.MatchRepository
}

func NewChatService(chatrepo repository.ChatMessageRepository, matchrepo repository.MatchRepository) ChatService {
	return &chatService{chatrepo: chatrepo, matchrepo: matchrepo}
}

func (s chatService) SendMessage(ctx context.Context, msg *db.ChatMessage) error {
	if msg == nil {
		return errChatInvalidFields
	}
	msg.ID = uuid.New().String()
	msg.RideID = strings.TrimSpace(msg.RideID)
	msg.SenderID = strings.TrimSpace(msg.SenderID)
	msg.Content = strings.TrimSpace(msg.Content)

	if msg.RideID == "" || msg.SenderID == "" || msg.Content == "" {
		return errChatInvalidFields
	}
	matches, err := s.matchrepo.FindByRideID(ctx, msg.RideID)
	if err != nil {
		return err
	}
	allowed := false
	for _, m := range matches {
		if (m.RiderID == msg.SenderID || m.DriverID == msg.SenderID) &&
			(m.Status == "accepted" || m.Status == "completed") {
			allowed = true
			break
		}
	}
	if !allowed {
		return errChatNotAllowed
	}
	msg.Timestamp = time.Now().UTC()
	return s.chatrepo.Create(ctx, msg)
}

func (s chatService) ListMessagesByRide(ctx context.Context, rideID string, limit int, before time.Time) ([]db.ChatMessage, error) {
	rideID = strings.TrimSpace(rideID)
	if before.IsZero() {
		before = time.Now().UTC()
	}
	if limit <= 0 {
		limit = 50
	}
	return s.chatrepo.ListByRide(ctx, rideID, limit, before)
}

func (s chatService) ListMessagesBySender(ctx context.Context, senderID string, limit int, before time.Time) ([]db.ChatMessage, error) {
	senderID = strings.TrimSpace(senderID)
	if before.IsZero() {
		before = time.Now().UTC()
	}
	if limit <= 0 {
		limit = 50
	}
	return s.chatrepo.ListBySender(ctx, senderID, limit, before)
}

func (s chatService) ListChatsForUser(ctx context.Context, userID string, limit int, before time.Time) ([]db.ChatMessage, error) {
	userID = strings.TrimSpace(userID)
	if before.IsZero() {
		before = time.Now().UTC()
	}
	if limit <= 0 {
		limit = 50
	}
	return s.chatrepo.ListChatsForUser(ctx, userID, limit, before)
}

func (s chatService) DeleteMessage(ctx context.Context, id string) error {
	return s.chatrepo.Delete(ctx, strings.TrimSpace(id))
}
