package middleware

import (
	"context"
	"github.com/golang-jwt/jwt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// context key to be injected along with token,
// user_id and mail will be injected to context so that handler can gets its
// value for authenticated requests
type ctxKey string

const (
	ctxUserIDKey ctxKey = "user_id"
	ctxEmailKey  ctxKey = "email"
)

// config containd=s JWTSecret and public metods
// public methids are omiitted by middleware, they dont neeed
// tobe passed through the middleware
type Config struct {
	JWTSecret     []byte
	PublicMethods map[string]bool //map["proto/v1/auth.AuthService/Login"]=true
}

// Identity extracted after validating backend JWT
type Identity struct {
	UserID string
	Email  string
}

// Validate token does HS256 verification and extracts the identity
// from the JWT
func ValidateToken(_ context.Context, tokenStr string, secret []byte) (Identity, error) {
	//if token is empty, then error
	if tokenStr == "" {
		return Identity{}, status.Error(codes.Unauthenticated, "Unexpected signing method")
	}

	// Reads a JWT string, Decodes it into a structured token, Verifies its signature using a provided key, Tells you whether it’s valid.
	tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, status.Error(codes.Unauthenticated, "unexpected signing method")
		}
		return secret, nil
	})

	//check if token is vaid
	if err != nil || !tok.Valid {
		return Identity{}, status.Error(codes.Unauthenticated, "invalid or expired token")
	}

	//to get the claims from the token
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return Identity{}, status.Error(codes.Unauthenticated, "invalid token claims")
	}

	//
	sub, _ := claims["sub"].(string)
	email, _ := claims["email"].(string)
	if sub == "" {
		return Identity{}, status.Error(codes.Unauthenticated, "missing subject in token")
	}
	return Identity{UserID: sub, Email: email}, nil

}

// UserIDFromContext is used by handlers to read identity set byb the intercptor
func UserIDFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(ctxUserIDKey)
	s, ok := v.(string)
	return s, ok && s != ""
}
func EmailFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(ctxEmailKey)
	s, ok := v.(string)
	return s, ok && s != ""
}

// AuthInterceptor is like a central gatekeeper for all non-public RPC
func AuthInterceptor(cfg Config) grpc.UnaryServerInterceptor {
	secret := cfg.JWTSecret

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		//bypass for public method
		if cfg.PublicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		//extracting authorization
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}
		vals := md.Get("authorization")
		if len(vals) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization header required")
		}

		parts := strings.Fields(vals[0])
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization header")
		}
		tokenStr := parts[1]

		id, err := ValidateToken(ctx, tokenStr, secret)
		if err != nil {
			return nil, err
		}

		ctx = context.WithValue(ctx, ctxUserIDKey, id.UserID)
		if id.Email != "" {
			ctx = context.WithValue(ctx, ctxEmailKey, id.Email)
		}
		return handler(ctx, req)

	}
}
