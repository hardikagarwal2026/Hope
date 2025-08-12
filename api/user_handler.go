package api

import (
	"context"
	"hope/db"
	"hope/middleware"
	pb "hope/proto/v1/user"
	"hope/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	userService service.UserService
	pb.UnimplementedUserServiceServer
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}


func toUserPB(u *db.User) *pb.User{
	if u== nil {
		return nil
	}

	return &pb.User{
		Id: u.ID,
		Name: u.Name,
		Email: u.Email,
		PhotoUrl: u.PhotoURL,
		Geohash: u.Geohash,
		LastSeen: u.LastSeen.Unix(),
	}
}


func(h *UserHandler)GetMe(ctx context.Context, _ *pb.GetMeRequest)(*pb.GetMeResponse, error){
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok || userID == "" {
		return nil, status.Error(codes.Unauthenticated, "missing auth")
	}

	u, err := h.userService.GetUserByID(ctx, userID)
	if err != nil || u == nil || u.ID == "" {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &pb.GetMeResponse{
		User: toUserPB(u),
	}, nil
}


func(h *UserHandler)GetUser(ctx context.Context, req *pb.GetUserRequest)(*pb.GetUserResponse, error){
	if req == nil || req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	u, err := h.userService.GetUserByID(ctx, req.GetUserId())
	if err != nil || u == nil ||u.ID == "" {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &pb.GetUserResponse{User: toUserPB(u)}, nil
}


func(h *UserHandler)UpdateMe(ctx context.Context, req *pb.UpdateMeRequest)(*pb.UpdateMeResponse, error){
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok || userID == "" {
		return nil, status.Error(codes.Unauthenticated, "missing auth")
	}
	if req == nil || req.GetUser() == nil || req.GetFieldMask() == nil {
		return nil, status.Error(codes.InvalidArgument, "user and field maask required")
	}

	//loading the current user
	curr, err := h.userService.GetUserByID(ctx, userID)
	if err != nil || curr == nil || curr.ID == "" {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	for _, path := range req.FieldMask.Paths {
		switch path {
		case "name":
			curr.Name = req.User.GetName()
		case "photo_url":
			curr.PhotoURL = req.User.GetPhotoUrl()
		case "geohash":
			curr.Geohash = req.User.GetGeohash()
		default:
		}
	}

	if err := h.userService.UpdateUser(ctx, curr); err != nil {
		return nil, status.Errorf(codes.Internal, "update failed: %v", err)
	}

	return &pb.UpdateMeResponse{User: toUserPB(curr)}, nil
}


func(h *UserHandler)ListUsers(ctx context.Context, req *pb.ListUsersRequest)(*pb.ListUsersResponse, error){
	if req == nil || len(req.GetUserIds()) == 0 {
		return &pb.ListUsersResponse{Users: nil}, nil
	}

	out := make([]*pb.User,0, len(req.UserIds))
	for _, id := range req.UserIds {
		if id == ""{
			continue
		}
		u, err := h.userService.GetUserByID(ctx, id)
		if err == nil && u != nil && u.ID != "" {
			out = append(out, toUserPB(u))
		}
	}

	return &pb.ListUsersResponse{Users: out}, nil
}
