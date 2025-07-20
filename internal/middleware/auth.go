package middleware

import (
	"context"
	"net/http"

	"github.com/serhappy/rssagg/internal/app"
	"github.com/serhappy/rssagg/internal/auth"
	"github.com/serhappy/rssagg/internal/db"
	"github.com/serhappy/rssagg/internal/response"
)

type ctxKey string

const userCtxKey ctxKey = "user"

func Auth(app *app.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey, err := auth.GetAPIKey(r.Header)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "missing API key")
				return
			}

			user, err := app.DB.GetUserByAPIKey(r.Context(), apiKey)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "invalid user")
				return
			}

			ctx := context.WithValue(r.Context(), userCtxKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUser(ctx context.Context) (db.User, bool) {
	user, ok := ctx.Value(userCtxKey).(db.User)
	if !ok {
		return db.User{}, false
	}
	return user, ok
}
