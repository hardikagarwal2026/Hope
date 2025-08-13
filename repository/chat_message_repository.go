package repository

import (
	"context"
	"hope/db"
	"time"

	"gorm.io/gorm"
)


type ChatMessageRepository interface {
	Create(ctx context.Context, msg *db.ChatMessage) error
	ListByRide(ctx context.Context, rideID string, limit int, before time.Time)([]db.ChatMessage, error)
	ListBySender(ctx context.Context, senderID string, limit int, before time.Time)([]db.ChatMessage, error)
	Delete(ctx context.Context, id string)error
	ListChatsForUser(ctx context.Context, userID string, limit int, before time.Time) ([]db.ChatMessage, error)
}

// it holds the reference to db coonnection, it implements interface
type chatMessageRepository struct{
	db *gorm.DB
}

func NewChatMessageRepository(db *gorm.DB) ChatMessageRepository {
	return &chatMessageRepository{db: db}
}

func (r *chatMessageRepository) Create(ctx context.Context, msg *db.ChatMessage) error {
	return r.db.WithContext(ctx).Create(msg).Error
}

func(r *chatMessageRepository) ListByRide(ctx context.Context, rideID string, limit int, before time.Time)([]db.ChatMessage, error){
	var messages []db.ChatMessage
	err := r.db.WithContext(ctx).Where("ride_id = ? AND timestamp < ?",rideID, before).Limit(limit).Find(&messages).Error
	return messages, err
}

func (r *chatMessageRepository) ListBySender(ctx context.Context, senderID string, limit int, before time.Time)([]db.ChatMessage, error){
	var messages []db.ChatMessage
	err := r.db.WithContext(ctx).Where("sender_id = ? AND timestamp < ?", senderID, before).Order("timestamp DESC").Limit(limit).Find(&messages).Error
	return messages, err
}

func (r *chatMessageRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&db.ChatMessage{}, "id = ?", id).Error
}

// //have to go through again through this function, write now using AI for this function
func (r *chatMessageRepository) ListChatsForUser(ctx context.Context, userID string, limit int, before time.Time) ([]db.ChatMessage, error) {
    var messages []db.ChatMessage
    err := r.db.WithContext(ctx).
        Where("sender_id = ? AND timestamp < ?", userID, before).
        Or("ride_id IN (?)", r.db.Model(&db.Match{}).Select("ride_id").Where("rider_id = ? OR driver_id = ?", userID, userID)).
        Order("timestamp DESC").
        Limit(limit).
        Find(&messages).Error
    return messages, err
}