package handler

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"su-server/internal/service"
)

type OAuthHandler struct {
	oauthService *service.OAuthService
	jwtService   *service.JWTService
}

func NewOAuthHandler(oauthService *service.OAuthService, jwtService *service.JWTService) *OAuthHandler {
	return &OAuthHandler{
		oauthService: oauthService,
		jwtService:   jwtService,
	}
}

// GoogleLogin handles GET /auth/google
// redirects the user to Google's OAuth2 login page
func (h *OAuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	state := generateState()

	// store state in cookie to verify later in callback
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == "production",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   300, // 5 minutes
	})

	url := h.oauthService.GetAuthURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// GoogleCallback handles GET /auth/google/callback
// Google redirects here after the user logs in
func (h *OAuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	// verify state to prevent CSRF
	cookie, err := r.Cookie("oauth_state")
	if err != nil {
		http.Error(w, "missing state cookie", http.StatusBadRequest)
		return
	}

	if r.URL.Query().Get("state") != cookie.Value {
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}

	// exchange code for user info + upsert user
	code := r.URL.Query().Get("code")
	user, err := h.oauthService.ExchangeCode(r.Context(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// generate JWT
	token, err := h.jwtService.Generate(user.ID, user.UserType)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"token": token,
		"user":  user,
	})
}

// GoogleVerify handles POST /auth/google/verify
// for Flutter mobile — Flutter does the Google login itself
// and sends the ID token here to verify
func (h *OAuthHandler) GoogleVerify(w http.ResponseWriter, r *http.Request) {
	var body struct {
		IDToken string `json:"id_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if body.IDToken == "" {
		http.Error(w, "id_token is required", http.StatusBadRequest)
		return
	}

	user, err := h.oauthService.VerifyIDToken(r.Context(), body.IDToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := h.jwtService.Generate(user.ID, user.UserType)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"token": token,
		"user":  user,
	})
}

// generateState generates a random state string for CSRF protection
func generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
