package api

import (
	"context"
	"hope/db"
	pb "hope/proto/v1/chat"
	"hope/service"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)


type ChatHandler struct {
	chatService service.ChatService
	pb.UnimplementedChatServiceServer
}

func NewChatHandler(chatService service.ChatService)*ChatHandler {
	return &ChatHandler{chatService: chatService}
}

func toChatPB(c *db.ChatMessage)*pb.ChatMessage{
	if c == nil {
		return nil
	}

	var ts *timestamppb.Timestamp
	if !c.Timestamp.IsZero(){
		ts = timestamppb.New(c.Timestamp)
	}

	return &pb.ChatMessage{
		Id: c.ID,
		RideId: c.RideID,
		SenderId: c.SenderID,
		Content: c.Content,
		Timestamp: ts,
	}
}


func(h *ChatHandler)SendMessage(ctx context.Context, req *pb.SendMessageRequest)(*pb.SendMessageResponse, error){
	if req == nil || req.GetRideId() == "" || req.GetSenderId() == "" || req.Content == "" {
		return nil, status.Error(codes.InvalidArgument, "rideid, senderid and contnt required")
	}

	msg := &db.ChatMessage{
		RideID: req.GetRideId(),
		SenderID: req.GetSenderId(),
		Content: req.GetContent(),
	}

	err := h.chatService.SendMessage(ctx, msg)
	if err != nil{
		return nil, status.Errorf(codes.PermissionDenied, "srnd failed:%v", err)
	}

	return &pb.SendMessageResponse{
		Message: toChatPB(msg),
	}, nil
}


func(h *ChatHandler)ListMessagesByRide(ctx context.Context, req *pb.ListMessagesByRideRequest)(*pb.ListMessagesByRideResponse, error){
	if req == nil || req.GetRideId() == "" {
		return nil, status.Error(codes.InvalidArgument, "ride id required")
	}
	before := time.Now()
	if req.GetBefore() != nil {
		before = req.GetBefore().AsTime()   //converting protobuf timestamp into time.Time
	}
	msgs, err := h.chatService.ListMessagesByRide(ctx, req.GetRideId(), int(req.GetLimit()), before)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list failed:%v", err)
	}
	out := make([]*pb.ChatMessage, 0, len(msgs))
	for i := range msgs{
		out = append(out, toChatPB(&msgs[i]))
	}

	return &pb.ListMessagesByRideResponse{
		Messages: out,
	}, nil
}


func (h *ChatHandler)ListMessagesBySender(ctx context.Context, req *pb.ListMessagesBySenderRequest)(*pb.ListMessagesBySenderResponse, error){
	if req == nil || req.GetSenderId() == "" {
		return nil, status.Error(codes.InvalidArgument, "sender id required")
	}
	before := time.Now()
	if req.GetBefore() != nil {
		before = req.GetBefore().AsTime()  //converting protobuf time into time.Time
	}
	msgs, err := h.chatService.ListMessagesBySender(ctx, req.GetSenderId(), int(req.GetLimit()), before)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list failed:%v", err)
	}

	out := make([]*pb.ChatMessage, 0, len(msgs))
	for i := range msgs {
		out = append(out, toChatPB(&msgs[i]))
	}

	return &pb.ListMessagesBySenderResponse{
		Messages: out,
	}, nil
}


func (h *ChatHandler) ListChatsForUser(ctx context.Context, req *pb.ListChatsForUserRequest) (*pb.ListChatsForUserResponse, error) {
	if req == nil || req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
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


func (h *ChatHandler)DeleteMessage(ctx context.Context, req *pb.DeleteMessageRequest)(*pb.DeleteMessageResponse, error){
	if req == nil || req.GetMessageId() == "" {
		return nil, status.Error(codes.InvalidArgument, "message id required")
	}

	err := h.chatService.DeleteMessage(ctx, req.GetMessageId())
	if err != nil{
		return nil, status.Errorf(codes.Internal, "delete failed:%v", err)
	}

	return &pb.DeleteMessageResponse{
		Success: true,
	}, nil
}