package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hope/db"
	"hope/repository"
	"net/http"
	"strings"
	"time"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// GoogleTokenInfo represents the payload returned by Google's token info endpoint
// when validating an OAuth 2.0 ID token. It contains user and token-related data.
type GoogleTokenInfo struct {
    Aud           string `json:"aud"`             // audience is the client ID of your app that the ID token is intended for
    Email         string `json:"email"`           // Email
    EmailVerified string `json:"email_verified"`  // Whether the email address has been verified ("true" or "false")
    Exp           string `json:"exp"`             // Expiration time in  Unix timestamp (seconds) when the token will expire
    Iss           string `json:"iss"`             // Issuer is the issuer of the token,"accounts.google.com"
    Sub           string `json:"sub"`             // subject is a unique identifier for the Google account (remains constant even if email changes)
    Name          string `json:"name"`            
    Picture       string `json:"picture"`         // URL to the user's Google profile picture
}


// var jwtSecret = os.Getenv("JWT_SECRET")

type AuthService interface{
	Login(ctx context.Context, idToken string) (string, error)
}

//authservice struct containg dependencies such as userrepo and allowed domains
type authService struct {
	userrepo		 repository.UserRepository
	allowedDomains   map[string]struct{}
	jwtSecret      	 []byte
    googleClientID 	 string
}

// Constructor to get create instance of AuthService
func NewAuthService(userrepo repository.UserRepository, allowedDomains map[string]struct{}, jwtSecret []byte, googleClientID string) AuthService {
	return &authService{userrepo: userrepo, allowedDomains: allowedDomains, jwtSecret: jwtSecret, googleClientID: googleClientID}
}

// to validate google id token(it comes from frontend) and getting googletoken info
// it contains 
func(s *authService) verifyGoogleIDToken(idToken string)(*GoogleTokenInfo, error) {
	resp, err := http.Get(fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", idToken))
	if err != nil{
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("invalid token")
	}

	var tokeninfo GoogleTokenInfo
	if err := json.NewDecoder(resp.Body).Decode(&tokeninfo); err!= nil{
		return nil,err
	}

	return &tokeninfo, nil
}

func(s *authService) issueJWT(user *db.User)(string, error){
	//custom claims send along with jwt token
	claims := jwt.MapClaims{
		"sub":   user.ID,       //user ID (subject)
		"email": user.Email,	//user email
		"name":  user.Name,		//use name
		"iat":   time.Now().Unix(),    //issued at, jab isse hua in UNIX format
		"exp":   time.Now().Add(24 * time.Hour).Unix(),    //expiry time, 24 hr + current time
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)  //method is HS256 along with claims
	return token.SignedString(s.jwtSecret)    //returning signed string along with jwt secret key
}

func (s authService)Login(ctx context.Context, idToken string) (string, error){
	tokeninfo, err := s.verifyGoogleIDToken(idToken)
	if err != nil {
		return "", err
	}


	//checking audience is matching your google client id
	if tokeninfo.Aud != s.googleClientID{
		return "", errors.New("invalid audience")
	}


	//validate email verified
	if tokeninfo.EmailVerified != "true" {
		return "", errors.New("email not verified")
	}

	// restricting domain to a particular domains from the allowed list in the env
	allowed := false
	email := strings.ToLower(tokeninfo.Email)
	for domain := range s.allowedDomains {
        if strings.HasSuffix(email, "@"+domain) {
            allowed = true
            break
        }
    }
	if !allowed {
		return "", errors.New("unauthorized email domain")
	}

	//find if exist or create user from tokeninfo
	user, err := s.userrepo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		user = &db.User{
			ID: 		uuid.New().String(),
			Email: 	    email,
			Name:       tokeninfo.Name,
			PhotoURL:   tokeninfo.Picture,
			LastSeen:   time.Now(),
		}
		if err := s.userrepo.Create(ctx, user);err!= nil{
			return "", err
		}
	}

	//after login and storing user in db, returning jwt Token after
	//so user can be in session for 24hrs- expiry 
	return s.issueJWT(user)

}