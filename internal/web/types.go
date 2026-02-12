package web

import (
	"time"

	"github.com/MudassirDev/mini-youtube/db/database"
	"github.com/google/uuid"
)

type apiConfig struct {
	DB *database.Queries
}

type createUserRequest struct {
	Email       string `json:"email" validate:"required"`
	Username    string `json:"username" validate:"required"`
	Password    string `json:"password" validate:"required"`
	DisplayName string `json:"display_name" validate:"required"`
}

type user struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	DisplayName string    `json:"display_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
