package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/mosesbenjamin/sre-from-local-to-prod/internal/database"
	"github.com/mosesbenjamin/sre-from-local-to-prod/sql/schema"
)

type config struct {
	PSQL   schema.PostgresConfig
	DB     *database.Queries
	Server struct {
		Address string
	}
	JWTSecret string
}

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}

	cfg.PSQL = schema.PostgresConfig{
		Host:     os.Getenv("PSQL_HOST"),
		Port:     os.Getenv("PSQL_PORT"),
		User:     os.Getenv("PSQL_USER"),
		Password: os.Getenv("PSQL_PASSWORD"),
		Database: os.Getenv("PSQL_DATABASE"),
		SSLMode:  os.Getenv("PSQL_SSLMODE"),
	}
	if cfg.PSQL.Host == "" && cfg.PSQL.Port == "" {
		return cfg, fmt.Errorf("no PSQL config provided")
	}
	cfg.Server.Address = os.Getenv("SERVER_ADDRESS")

	cfg.JWTSecret = os.Getenv("JWT_SECRET")
	if cfg.JWTSecret == "" {
		return cfg, fmt.Errorf("no JWT_SECRET config provided")
	}

	return cfg, nil
}

func run(cfg config) error {
	db, err := schema.Open(cfg.PSQL)
	if err != nil {
		return err
	}
	defer db.Close()

	err = schema.MigrateFS(db, schema.FS, ".")
	if err != nil {
		return err
	}

	dbQueries := database.New(db)
	cfg.DB = dbQueries

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	apiRouter := chi.NewRouter()

	apiRouter.Get("/healthcheck", handleReadiness)

	apiRouter.Post("/login", cfg.handlerLogin)

	apiRouter.Post("/students", cfg.handlerStudentsCreate)
	apiRouter.Get("/students", cfg.handlerGetStudents)
	apiRouter.Get("/students/{studentID}", cfg.middlewareAuth(cfg.handlerStudentGet))
	apiRouter.Delete("/students/{studentID}", cfg.middlewareAuth(cfg.handlerStudentDelete))
	apiRouter.Put("/students/{studentID}", cfg.middlewareAuth(cfg.handlerUpdateStudentPassword))

	// Logical nesting on endpoints and API versioning
	router.Mount("/api/v1", apiRouter)

	adminRouter := chi.NewRouter()
	router.Mount("/admin", adminRouter)

	srv := &http.Server{
		Addr:    cfg.Server.Address,
		Handler: router,
	}

	log.Printf("Starting server on port %s...\n", cfg.Server.Address)
	log.Fatal(srv.ListenAndServe())
	return nil
}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}
	err = run(cfg)
	if err != nil {
		panic(err)
	}
}
