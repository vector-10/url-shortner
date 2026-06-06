package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"
	"github.com/google/uuid"
	"github.com/vector-10/url-shortner/internal/store"
	"github.com/vector-10/url-shortner/internal/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)


type OAuthHandler struct {
	userStore store.UserStore
	oauthConfig *oauth2.Config
}

func NewOAuthHandler(us store.UserStore) *OAuthHandler {
	config := &oauth2.Config{
        ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
        ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
        RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
        Scopes:       []string{"email", "profile"},
        Endpoint:     google.Endpoint,
    }
    return &OAuthHandler{userStore: us, oauthConfig: config}
}

func (h *OAuthHandler) GoogleLogin(W http.ResponseWriter, r *http.Request) {
	url := h.oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOnline)
	http.Redirect(W, r, url, http.StatusTemporaryRedirect)
}

func (h *OAuthHandler) GoogleCallback(W http.ResponseWriter, r *http.Request) {
     code := r.URL.Query().Get("code")
	 if code == "" {
		http.Error(W, "Code not found", http.StatusBadRequest)
		return
	}

	token, err := h.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(W, "failed to exchange token", http.StatusInternalServerError)
		return
	}

	client := h.oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(W, "failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var googleUser struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		http.Error(W, "failed to decode user info", http.StatusInternalServerError)
		return
	}

	user, err := h.userStore.GetUserByEmail(googleUser.Email)
	if err != nil {
		user = &models.User{
			ID: uuid.New().String(),
			Email: googleUser.Email,
			CreatedAt: time.Now(),
		}
		if err := h.userStore.CreateUser(user); err != nil {
			http.Error(W, "failed to create user", http.StatusInternalServerError)
			return
		}
	}

	jwtToken, err := generateToken(user.ID)
	if err != nil {
		http.Error(W, "failed to generate token", http.StatusInternalServerError)
		return
	}

	http.Redirect(W, r, os.Getenv("FRONTEND_URL")+"?token="+jwtToken, http.StatusTemporaryRedirect)



}

