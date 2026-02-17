package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/MudassirDev/mini-youtube/db/database"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func getUserIDFromContext(w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {
	rawID := r.Context().Value(AUTH_KEY)
	stringID, ok := rawID.(string)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, errors.New("invalid user id"), "user unauthorized")
		return uuid.Nil, errors.New("invalid user id")
	}
	userID, err := uuid.Parse(stringID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err, "user unauthorized")
		return uuid.Nil, err
	}
	return userID, nil
}

func checkPostHeader(r *http.Request) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return errors.New("invalid content type header")
	}

	return nil
}

func getValidatorErrMsg(err error) error {
	if errs, ok := err.(validator.ValidationErrors); ok {
		e := errs[0]
		message := fmt.Sprintf("invalid field: %v", e.Field())
		return errors.New(message)
	}
	return err
}

func makeResponseUser(userResult database.User) user {
	return user{
		ID:          userResult.ID,
		Email:       userResult.Email,
		Username:    userResult.Username,
		DisplayName: userResult.DisplayName,
		CreatedAt:   userResult.CreatedAt,
		UpdatedAt:   userResult.UpdatedAt,
	}
}
