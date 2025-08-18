package api

import (
	"context"
	pb "hope/proto/v1/auth"
	"hope/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if req == nil || req.GetIdToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "id_token is required")
	}

	jwtToken, user, err := h.authService.Login(ctx, req.GetIdToken())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "login failed: %v", err)
	}

	return &pb.LoginResponse{
		Jwt:      jwtToken,
		Userid:   user.ID,
		Email:    user.Email,
		PhotoUrl: user.PhotoURL,
	}, nil
}
