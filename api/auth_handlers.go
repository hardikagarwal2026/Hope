package api

import (
	"context"

	pb "hope/proto/v1/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"hope/service"
)

// Handler depends only on service layer; transports proto <-> domain.
type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	authService service.AuthService // FIX: typo
}

// Constructor wires the dependency (useful for DI later).
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

	// Example parsing of expiry from JWT if you want expires_at to be accurate;
	// otherwise, set to 0 or time.Now().Add(24*time.Hour).Unix() as per your claims.
	var expiresAt int64 = 0
	// Optional: if your issueJWT stores expiry as exp in JWT claims, parse it here for accurate expires_at field.

	return &pb.LoginResponse{
		Jwt:       jwtToken,
		ExpiresAt: expiresAt, // assign correct value if needed
		Userid:    user.ID,
		Email:     user.Email,
		Phone:     "", // If not available from user, leave empty
		PhotoUrl:  user.PhotoURL,
	}, nil
}
