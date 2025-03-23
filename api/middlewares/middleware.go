package middlewares

import (
	"api/models"
	"api/utils"
	"context"
	"net/http"
)

func Chain(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}
	return h
}

type contextKey string

var UserContextKey contextKey = "user"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// We need to remove the Bearer prefix from the token
		token = token[7:]
		retrievedUserID, err := utils.VerifyJWT(token)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		user := models.User{
			DiscordID: retrievedUserID,
		}

		models.Database.Where("discord_id = ?", retrievedUserID).First(&user)
		if user.ID == 0 {
			utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		UserContext := context.WithValue(r.Context(), UserContextKey, &user)

		r = r.WithContext(UserContext)
		next(w, r)
	}
}
