package web

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/MudassirDev/mini-youtube/db/database"
	"github.com/MudassirDev/mini-youtube/internal/auth"
)

func (c *apiConfig) handleUserCreate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestBody createUserRequest
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			respondWithError(w, http.StatusBadRequest, err, "invalid payload")
			return
		}

		err := validate.Struct(requestBody)
		if err != nil {
			message := getValidatorErrMsg(err)
			respondWithError(w, http.StatusBadRequest, message, message.Error())
			return
		}

		password, err := auth.HashPassword(requestBody.Password)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err, "failed to encrypt password")
			return
		}

		userResult, err := c.DB.CreateUser(context.Background(), database.CreateUserParams{
			Email:        requestBody.Email,
			Username:     requestBody.Username,
			PasswordHash: password,
			DisplayName:  requestBody.DisplayName,
		})

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err, "failed to create user")
			return
		}

		respondWithJSON(w, http.StatusCreated, makeResponseUser(userResult))
	})
}

func (c *apiConfig) handleUserLogin() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestBody loginUserRequest
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			respondWithError(w, http.StatusBadRequest, err, "invalid payload")
			return
		}

		err := validate.Struct(requestBody)
		if err != nil {
			message := getValidatorErrMsg(err)
			respondWithError(w, http.StatusBadRequest, message, message.Error())
			return
		}

		userResult, err := c.DB.GetUserWithUsername(context.Background(), requestBody.Username)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err, "no such user")
			return
		}

		err = auth.VerifyPassword(requestBody.Password, userResult.PasswordHash)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err, "wrong password")
			return
		}

		token, err := auth.CreateJWTToken(userResult.ID, EXPIRES_IN, c.JWT_SECRET)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err, "failed to create token")
			return
		}

		cookie := &http.Cookie{
			Name:     AUTH_KEY,
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteStrictMode,
			Expires:  time.Now().Add(EXPIRES_IN),
			MaxAge:   int(EXPIRES_IN),
		}
		http.SetCookie(w, cookie)

		respondWithJSON(w, http.StatusOK, makeResponseUser(userResult))
	})
}

func (c *apiConfig) handleUserLogout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     AUTH_KEY,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
	respondWithJSON(w, http.StatusOK, "deleted cookie")
}
