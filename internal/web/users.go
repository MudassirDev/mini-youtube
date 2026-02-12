package web

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/MudassirDev/mini-youtube/db/database"
	"github.com/MudassirDev/mini-youtube/internal/auth"
)

func (c *apiConfig) HandleUserCreate(w http.ResponseWriter, r *http.Request) {
	err := checkHeader(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err, err.Error())
		return
	}

	var requestBody createUserRequest
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&requestBody); err != nil {
		respondWithError(w, http.StatusBadRequest, err, "invalid payload")
		return
	}

	err = validate.Struct(requestBody)
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

	respondWithJSON(w, http.StatusCreated, user{
		ID:          userResult.ID,
		Email:       userResult.Email,
		Username:    userResult.Username,
		DisplayName: userResult.DisplayName,
		CreatedAt:   userResult.CreatedAt,
		UpdatedAt:   userResult.UpdatedAt,
	})
}
