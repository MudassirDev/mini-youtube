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

	mux.Handle("POST /api/users/create", apiCfg.postMiddleware(apiCfg.handleUserCreate()))
	mux.Handle("POST /api/users/login", apiCfg.postMiddleware(apiCfg.handleUserLogin()))
	mux.HandleFunc("DELETE /api/users/login", apiCfg.handleUserLogout)

	mux.Handle("POST /api/users/videos/create", apiCfg.authMiddleware(
		apiCfg.postMiddleware(apiCfg.videoUploadHandler()),
	))

	return mux
}
