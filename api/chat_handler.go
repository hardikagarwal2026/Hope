package api

import (
	"context"
	"time"
	"hope/db"
	"hope/middleware"
	pb "hope/proto/v1/chat"
	"hope/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ChatHandler struct {
	chatService service.ChatService
	pb.UnimplementedChatServiceServer
}

func NewChatHandler(chatService service.ChatService) *ChatHandler {
	return &ChatHandler{chatService: chatService}
}

func toChatPB(c *db.ChatMessage) *pb.ChatMessage {
	if c == nil {
		return nil
	}
	var ts *timestamppb.Timestamp
	if !c.Timestamp.IsZero() {
		ts = timestamppb.New(c.Timestamp)
	}
	return &pb.ChatMessage{
		Id:        c.ID,
		RideId:    c.RideID,
		SenderId:  c.SenderID,
		Content:   c.Content,
		Timestamp: ts,
	}
}

func (h *ChatHandler) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	if req == nil || req.GetRideId() == "" || req.GetContent() == "" {
		return nil, status.Error(codes.InvalidArgument, "ride_id and content are required")
	}

	senderID, ok := middleware.UserIDFromContext(ctx)
	if !ok || senderID == "" {
		return nil, status.Error(codes.Unauthenticated, "missing auth")
	}

	msg := &db.ChatMessage{
		RideID:   req.GetRideId(),
		SenderID: senderID, 
		Content:  req.GetContent(),
	}

	if err := h.chatService.SendMessage(ctx, msg); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "send failed: %v", err)
	}

	return &pb.SendMessageResponse{Message: toChatPB(msg)}, nil
}

func (h *ChatHandler) ListMessagesByRide(ctx context.Context, req *pb.ListMessagesByRideRequest) (*pb.ListMessagesByRideResponse, error) {
	if req == nil || req.GetRideId() == "" {
		return nil, status.Error(codes.InvalidArgument, "ride_id required")
	}
	before := time.Now()
	if req.GetBefore() != nil {
		before = req.GetBefore().AsTime()
	}
	msgs, err := h.chatService.ListMessagesByRide(ctx, req.GetRideId(), int(req.GetLimit()), before)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list failed: %v", err)
	}
	out := make([]*pb.ChatMessage, 0, len(msgs))
	for i := range msgs {
		out = append(out, toChatPB(&msgs[i]))
	}
	return &pb.ListMessagesByRideResponse{Messages: out}, nil
}

func (h *ChatHandler) ListMessagesBySender(ctx context.Context, req *pb.ListMessagesBySenderRequest) (*pb.ListMessagesBySenderResponse, error) {
	if req == nil || req.GetSenderId() == "" {
		return nil, status.Error(codes.InvalidArgument, "sender_id required")
	}
	before := time.Now()
	if req.GetBefore() != nil {
		before = req.GetBefore().AsTime()
	}
	msgs, err := h.chatService.ListMessagesBySender(ctx, req.GetSenderId(), int(req.GetLimit()), before)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list failed: %v", err)
	}
	out := make([]*pb.ChatMessage, 0, len(msgs))
	for i := range msgs {
		out = append(out, toChatPB(&msgs[i]))
	}
	return &pb.ListMessagesBySenderResponse{Messages: out}, nil
}

func (h *ChatHandler) ListChatsForUser(ctx context.Context, req *pb.ListChatsForUserRequest) (*pb.ListChatsForUserResponse, error) {
	if req == nil || req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id required")
	}
	before := time.Now()
	if req.GetBefore() != nil {
		before = req.GetBefore().AsTime()
	}
	msgs, err := h.chatService.ListChatsForUser(ctx, req.GetUserId(), int(req.GetLimit()), before)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list failed: %v", err)
	}
	out := make([]*pb.ChatMessage, 0, len(msgs))
	for i := range msgs {
		out = append(out, toChatPB(&msgs[i]))
	}
	return &pb.ListChatsForUserResponse{Messages: out}, nil
}