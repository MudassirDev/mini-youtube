package web

import (
	"context"
	"net/http"

	"github.com/MudassirDev/mini-youtube/internal/auth"
)

func (c *apiConfig) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// middleware logic here
		cookie, err := r.Cookie(AUTH_KEY)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err, "user unauthorized")
			return
		}

		userID, err := auth.VerifyJWT(c.JWT_SECRET, cookie.Value)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err, "user unauthorized")
			return
		}

		context := context.WithValue(r.Context(), AUTH_KEY, userID.String())
		next.ServeHTTP(w, r.WithContext(context))
	})
}

func (c *apiConfig) postMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := checkPostHeader(r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err, err.Error())
			return
		}

		next.ServeHTTP(w, r)
	})
}
