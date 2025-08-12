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
    // Transport-level validation: fast fail for missing payload.
    if req == nil || req.GetIdToken() == "" {
        return nil, status.Error(codes.InvalidArgument, "id_token is required")
    }

    // Delegate verification/provisioning to service; it returns a backend JWT.
    jwtToken, err := h.authService.Login(ctx, req.GetIdToken())
    if err != nil {
        // Map all auth problems to Unauthenticated (don’t leak internals).
        return nil, status.Errorf(codes.Unauthenticated, "login failed: %v", err)
    }

    return &pb.LoginResponse{Jwt: jwtToken}, nil
}
