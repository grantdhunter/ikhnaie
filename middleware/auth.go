package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/grantdhunter/ikhnaie/models"
)

type AuthMiddleware struct {
	handler http.Handler
}

func (am *AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := authenticate(r)

	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	ctx := context.WithValue(r.Context(), 0, user)
	rWithUser := r.WithContext(ctx)

	am.handler.ServeHTTP(w, rWithUser)

}

func authenticate(r *http.Request) (*models.User, error) {
	header := r.Header.Get("Authorization")
	parts := strings.Split(header, " ")

	var token string
	if len(parts) == 2 {
		token = parts[1]
	}

	if len(token) == 0 {
		query := r.URL.Query()
		parts = query["token"]
		log.Print(parts)
		if len(parts) == 1 {
			token = parts[0]
		}
	}
	log.Print("Auth token is:", token)
	if len(token) == 0 {
		return nil, errors.New("No Auth")
	}

	user := models.GetUser(token)
	return user, nil
}

func NewAuth(h http.Handler) *AuthMiddleware {
	return &AuthMiddleware{h}
}
