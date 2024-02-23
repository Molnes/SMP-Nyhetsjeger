package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var (
// GoogleOauthConfig oauth2.Config
)

const (
	OAUTH_STATE_COOKIE   = "oauthstate"
	GOOGLE_OAUTH_API_URL = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
)

type GoogleSsoConfig struct {
	RedirectUrl  string
	ClientId     string
	ClientSecret string
}

// func init() {
// 	redirect_url, ok := os.LookupEnv("GOOGLE_REDIRECT_URL")
// 	if !ok {
// 		log.Fatal("No redirect url provided. Expected GOOGLE_REDIRECT_URL")
// 	}
// 	client_id, ok := os.LookupEnv("GOOGLE_CLIENT_ID")
// 	if !ok {
// 		log.Fatal("No client id provided. Expected GOOGLE_CLIENT_ID")
// 	}
// 	client_secret, ok := os.LookupEnv("GOOGLE_CLIENT_SECRET")
// 	if !ok {
// 		log.Fatal("No client secret provided. Expected GOOGLE_CLIENT_SECRET")
// 	}
// 	GoogleOauthConfig = oauth2.Config{
// 		ClientID:     client_id,
// 		ClientSecret: client_secret,
// 		RedirectURL:  redirect_url,
// 		Endpoint:     google.Endpoint,
// 		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
// 	}

// }

// Generates a new state and sets it as a cookie
// Returns the state so it can be sent to the oauth provider
func GenerateAndSetStateOauthCookie(c echo.Context) string {
	state := uuid.New().String()
	cookie := http.Cookie{
		Name:   OAUTH_STATE_COOKIE,
		Value:  state,
		MaxAge: 3600,
	}
	c.SetCookie(&cookie)
	return state
}

// Struct to hold the user data from the Google OAuth2 API
type googleUser struct {
	Email          string `json:"email"`
	ID             string `json:"id"`
	Picture        string `json:"picture"`
	Verified_email bool   `json:"verified_email"`
}

// Gets the user data from the Google OAuth2 API
func GetGoogleUserData(accessToken string) (googleUser, error) {
	resp, err := http.Get(GOOGLE_OAUTH_API_URL + accessToken)
	if err != nil {
		return googleUser{}, fmt.Errorf("failed to get user info: %s", err.Error())
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return googleUser{}, fmt.Errorf("failed to read response body: %s", err.Error())
	}

	var usr googleUser
	err = json.Unmarshal(content, &usr)
	if err != nil {
		return googleUser{}, fmt.Errorf("failed to unmarshal user info: %s", err.Error())
	}
	return usr, nil
}
