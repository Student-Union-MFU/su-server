package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"su-server/internal/model"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OAuthService struct {
	userService  *UserService
	oauthConfig  *oauth2.Config
}

type GoogleUserInfo struct {
	Sub        string `json:"sub"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Picture    string `json:"picture"`
}

func NewOAuthService(userService *UserService) *OAuthService {
	config := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &OAuthService{
		userService: userService,
		oauthConfig: config,
	}
}

// GetAuthURL returns the Google OAuth2 login URL
func (s *OAuthService) GetAuthURL(state string) string {
	return s.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// ExchangeCode exchanges the OAuth2 code for a Google token
// then fetches user info and upserts the user in the DB
func (s *OAuthService) ExchangeCode(ctx context.Context, code string) (*model.User, error) {
	// exchange code for token
	token, err := s.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %w", err)
	}

	// fetch user info from Google
	userInfo, err := s.fetchGoogleUserInfo(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}

	// validate MFU domain
	if !strings.HasSuffix(userInfo.Email, "@lamduan.mfu.ac.th") {
		return nil, fmt.Errorf("only MFU students are allowed")
	}

	// extract student id from email prefix
	studentID, _, _ := strings.Cut(userInfo.Email, "@")

	// upsert user
	user, err := s.userService.UpsertUser(ctx, model.User{
		UserType:     model.UserTypeStudent,
		Name:         userInfo.Name,
		Email:        userInfo.Email,
		AvatarURL:    &userInfo.Picture,
		StudentID:    &studentID,
		OAuthSubject: userInfo.Sub,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upsert user: %w", err)
	}

	return user, nil
}

// VerifyIDToken verifies a Google ID token from Flutter mobile
func (s *OAuthService) VerifyIDToken(ctx context.Context, idToken string) (*model.User, error) {
	// fetch Google's public keys and verify token
	resp, err := http.Get("https://oauth2.googleapis.com/tokeninfo?id_token=" + idToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid token")
	}

	var info struct {
		Sub     string `json:"sub"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
		Aud     string `json:"aud"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("failed to decode token info: %w", err)
	}

	// verify the token was issued for your app
	if info.Aud != os.Getenv("GOOGLE_CLIENT_ID") {
		return nil, fmt.Errorf("token audience mismatch")
	}

	// validate MFU domain
	if !strings.HasSuffix(info.Email, "@lamduan.mfu.ac.th") {
		return nil, fmt.Errorf("only MFU students are allowed")
	}

	studentID, _, _ := strings.Cut(info.Email, "@")

	user, err := s.userService.UpsertUser(ctx, model.User{
		UserType:     model.UserTypeStudent,
		Name:         info.Name,
		Email:        info.Email,
		AvatarURL:    &info.Picture,
		StudentID:    &studentID,
		OAuthSubject: info.Sub,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upsert user: %w", err)
	}

	return user, nil
}

// fetchGoogleUserInfo fetches user info from Google using the OAuth2 token
func (s *OAuthService) fetchGoogleUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
	client := s.oauthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}
