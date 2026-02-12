package web

import (
	"net/http"
	"time"

	"github.com/MudassirDev/mini-youtube/db/database"
	"github.com/go-playground/validator/v10"
)

var (
	validate = validator.New()
)

const (
	EXPIRES_IN = time.Hour * 1
	AUTH_KEY   = "AUTH_KEY"
)

func CreateMux(queries *database.Queries, jwtSecret string) *http.ServeMux {
	mux := http.NewServeMux()

	apiCfg := apiConfig{
		DB:         queries,
		JWT_SECRET: jwtSecret,
	}

	mux.HandleFunc("POST /api/users/create", apiCfg.handleUserCreate)
	mux.HandleFunc("POST /api/users/login", apiCfg.handleUserLogin)
	mux.HandleFunc("DELETE /api/users/login", apiCfg.handleUserLogout)

	return mux
}
