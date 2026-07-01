package middleware

import (
	"Task-Management/database/dbhelper"
	"Task-Management/models"
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
)

type ContextKeys string

const (
	userContext ContextKeys = "userContext"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("x-api-key")
		user, err := dbhelper.GetUserBySession(apiKey)

		if err != nil || user == nil {
			logrus.WithError(err).Errorf("failed to get user with token: %s", apiKey)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), userContext, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserContext(r *http.Request) *models.User {
	user := r.Context().Value(userContext).(*models.User)
	return user

}
