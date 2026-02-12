package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/MudassirDev/mini-youtube/db/database"
	"github.com/MudassirDev/mini-youtube/internal/web"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

var (
	//go:embed db/schema/*.sql
	embeddedFS embed.FS
	CONN       *sql.DB
	PORT       string
	HANDLER    *http.ServeMux
)

const (
	PRODUCTION_ENVIRONMENT = "production"
)

func init() {
	godotenv.Load()

	envs := map[string]string{
		"PORT":   "",
		"DB_URL": "",
		"ENV":    "",
	}

	for key := range envs {
		env := os.Getenv(key)
		validateEnv(env, key)
		envs[key] = env
	}
	log.Println("environment variables loaded")

	isEnvProduction := envs["ENV"] == PRODUCTION_ENVIRONMENT
	isSSLModeDisabled := strings.Contains(envs["DB_URL"], "?sslmode=disable")

	if !isEnvProduction && !isSSLModeDisabled {
		envs["DB_URL"] = fmt.Sprintf("%v?sslmode=disable", envs["DB_URL"])
	}

	conn, err := sql.Open("postgres", envs["DB_URL"])
	if err != nil {
		log.Fatal("error: ", err)
		return
	}

	fsys, err := fs.Sub(embeddedFS, "db/schema")
	if err != nil {
		log.Fatal("error: ", err)
		return
	}

	provider, err := goose.NewProvider(goose.DialectPostgres, conn, fsys)
	if err != nil {
		log.Fatal("error: ", err)
		return
	}

	result, err := provider.Up(context.Background())
	if err != nil {
		log.Fatal("error: ", err)
		return
	}

	log.Println("ran migrations successfully: ", result)

	queries := database.New(conn)

	handler := web.CreateMux(queries)

	CONN = conn
	PORT = envs["PORT"]
	HANDLER = handler
}

func main() {
	defer CONN.Close()

	srv := http.Server{
		Addr:    ":" + PORT,
		Handler: HANDLER,
	}

	log.Printf("server is listening at http://localhost:%v\n", PORT)
	log.Fatal(srv.ListenAndServe())
}

func validateEnv(env, name string) {
	if env == "" {
		log.Fatalf("error: env with key '%v' is empty\n", name)
	}
}
