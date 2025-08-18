package repository

import (
	"context"
	"hope/db"
	"time"

	"gorm.io/gorm"
)

type ChatMessageRepository interface {
	Create(ctx context.Context, msg *db.ChatMessage) error
	ListByRide(ctx context.Context, rideID string, limit int, before time.Time) ([]db.ChatMessage, error)
	ListBySender(ctx context.Context, senderID string, limit int, before time.Time) ([]db.ChatMessage, error)
	ListChatsForUser(ctx context.Context, userID string, limit int, before time.Time) ([]db.ChatMessage, error)
	Delete(ctx context.Context, id string) error
}

type chatMessageRepository struct {
	db *gorm.DB
}

func NewChatMessageRepository(db *gorm.DB) ChatMessageRepository {
	return &chatMessageRepository{db: db}
}

func (r *chatMessageRepository) Create(ctx context.Context, msg *db.ChatMessage) error {
	return r.db.WithContext(ctx).Create(msg).Error
}

func (r *chatMessageRepository) ListByRide(ctx context.Context, rideID string, limit int, before time.Time) ([]db.ChatMessage, error) {
	var messages []db.ChatMessage
	q := r.db.WithContext(ctx).
		Where("ride_id = ? AND timestamp < ?", rideID, before).
		Order("timestamp DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	err := q.Find(&messages).Error
	return messages, err
}

func (r *chatMessageRepository) ListBySender(ctx context.Context, senderID string, limit int, before time.Time) ([]db.ChatMessage, error) {
	var messages []db.ChatMessage
	q := r.db.WithContext(ctx).
		Where("sender_id = ? AND timestamp < ?", senderID, before).
		Order("timestamp DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	err := q.Find(&messages).Error
	return messages, err
}

func (r *chatMessageRepository) ListChatsForUser(ctx context.Context, userID string, limit int, before time.Time) ([]db.ChatMessage, error) {
	var messages []db.ChatMessage
	q := r.db.WithContext(ctx).
		Where("sender_id = ? AND timestamp < ?", userID, before).
		Order("timestamp DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	err := q.Find(&messages).Error
	return messages, err
}

func (r *chatMessageRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&db.ChatMessage{}, "id = ?", id).Error
}
