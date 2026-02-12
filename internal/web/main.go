package web

import (
	"net/http"

	"github.com/MudassirDev/mini-youtube/db/database"
	"github.com/go-playground/validator/v10"
)

var (
	validate = validator.New()
)

func CreateMux(queries *database.Queries) *http.ServeMux {
	mux := http.NewServeMux()

	apiCfg := apiConfig{
		DB: queries,
	}

	mux.HandleFunc("POST /users", apiCfg.HandleUserCreate)

	return mux
}
