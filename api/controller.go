package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/arekbor/oauth/types"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_GORILLA_MUX_KEY")))
)

func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := provideOAuth().AuthCodeURL("state")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")

	token, err := provideOAuth().Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := provideOAuth().Client(r.Context(), token)
	res, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, fmt.Errorf("error while authorizing: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	defer res.Body.Close()

	var user types.User

	err = json.NewDecoder(res.Body).Decode(&user)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not parse user info: %s", err.Error()), http.StatusBadRequest)
		return
	}

	session, err := store.Get(r, "user-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["email"] = user.Email
	session.Values["name"] = user.Name
	session.Values["picture"] = user.Picture
	if err := session.Save(r, w); err != nil {
		http.Error(w, errors.New("error while saving session").Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleUserinfo(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	email, ok := session.Values["email"].(string)
	if !ok {
		http.Error(w, errors.New("error while reading email from session").Error(), http.StatusInternalServerError)
		return
	}

	name, ok := session.Values["name"].(string)
	if !ok {
		http.Error(w, errors.New("error while reading name from session").Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(email, name)
}

func provideOAuth() *oauth2.Config {
	var (
		clientID     = os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
		clientSecret = os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
		redirectURL  = os.Getenv("GOOGLE_OAUTH_REDIRECT_URL")
	)

	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}
