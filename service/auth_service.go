package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"hope/db"
	"hope/repository"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type GoogleTokenInfo struct {
	Aud           string `json:"aud"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Exp           string `json:"exp"`
	Iss           string `json:"iss"`
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

type AuthService interface {
	// Returns: jwt, user, error
	Login(ctx context.Context, idToken string) (string, *db.User, error)
}

type authService struct {
	userrepo       repository.UserRepository
	allowedDomains map[string]struct{}
	jwtSecret      []byte
	googleClientID string
}

var (
	errInvalidIDToken     = errors.New("invalid token")
	errInvalidAudience    = errors.New("invalid audience")
	errEmailNotVerified   = errors.New("email not verified")
	errUnauthorizedDomain = errors.New("unauthorized email domain")
)

func NewAuthService(
	userrepo repository.UserRepository,
	allowedDomains map[string]struct{},
	jwtSecret []byte,
	googleClientID string,
) AuthService {
	return &authService{
		userrepo:       userrepo,
		allowedDomains: allowedDomains,
		jwtSecret:      jwtSecret,
		googleClientID: googleClientID,
	}
}

func (s *authService) verifyGoogleIDToken(idToken string) (*GoogleTokenInfo, error) {
	if strings.TrimSpace(idToken) == "" {
		return nil, errInvalidIDToken
	}
	resp, err := http.Get(fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", idToken))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errInvalidIDToken
	}
	var tokeninfo GoogleTokenInfo
	if err := json.NewDecoder(resp.Body).Decode(&tokeninfo); err != nil {
		return nil, err
	}
	return &tokeninfo, nil
}

func (s *authService) issueJWT(user *db.User) (string, error) {
	if user == nil {
		return "", errors.New("nil user")
	}
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"name":  user.Name,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s authService) Login(ctx context.Context, idToken string) (string, *db.User, error) {
	tokeninfo, err := s.verifyGoogleIDToken(idToken)
	if err != nil {
		return "", nil, err
	}
	if tokeninfo.Aud != s.googleClientID {
		return "", nil, errInvalidAudience
	}
	if tokeninfo.EmailVerified != "true" {
		return "", nil, errEmailNotVerified
	}

	email := strings.ToLower(tokeninfo.Email)
	allowed := false
	for domain := range s.allowedDomains {
		if strings.HasSuffix(email, "@"+domain) {
			allowed = true
			break
		}
	}
	if !allowed {
		return "", nil, errUnauthorizedDomain
	}

	user, err := s.userrepo.FindByEmail(ctx, email)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		user = &db.User{
			ID:       uuid.New().String(),
			Email:    email,
			Name:     tokeninfo.Name,
			PhotoURL: tokeninfo.Picture,
			LastSeen: time.Now(),
		}
		if err := s.userrepo.Create(ctx, user); err != nil {
			return "", nil, err
		}
	}

	jwtStr, err := s.issueJWT(user)
	if err != nil {
		return "", nil, err
	}
	return jwtStr, user, nil
}
