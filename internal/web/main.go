package web

import (
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
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
	tmpls := template.New("")
	fs := http.FileServer(http.Dir("assets"))

	apiCfg := apiConfig{
		DB:         queries,
		JWT_SECRET: jwtSecret,
		Templates:  tmpls,
	}
	apiCfg.setupTemplates()

	mux.Handle("GET /videos", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCfg.Templates.ExecuteTemplate(w, "videos.html", nil)
	}))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	mux.Handle("POST /api/users/create", apiCfg.postMiddleware(apiCfg.handleUserCreate()))
	mux.Handle("POST /api/users/login", apiCfg.postMiddleware(apiCfg.handleUserLogin()))
	mux.HandleFunc("DELETE /api/users/login", apiCfg.handleUserLogout)

	mux.Handle("POST /api/users/videos/upload", apiCfg.authMiddleware(
		apiCfg.videoUploadHandler(),
	))

	return mux
}

func (c *apiConfig) setupTemplates() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working directory %v", err)
	}
	path := path.Join(cwd, "static/templates")
	filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && strings.HasSuffix(path, ".html") {
			_, err := c.Templates.ParseFiles(path)
			if err != nil {
				log.Fatalf("failed to parse templates: %v", err)
			}
		}
		return nil
	})
}
