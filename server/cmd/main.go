package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"project/internal/handlers"
	"project/internal/middleware"
	"project/internal/repository"
	"project/internal/service"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	URL := os.Getenv("CONN_BASE")
	if URL == "" {
		log.Fatal("CONN_BASE не задана в .env файле")
	}
	JwtSecret := os.Getenv("JWT_SECRET")
	if JwtSecret == "" {
		log.Fatal("JWT_SECRET не задана в .env файле")
	}
	Conn, err := pgx.Connect(context.Background(), URL)
	if err != nil {
		log.Fatal(err)
	}
	defer Conn.Close(context.Background())

	repo := repository.NewUserRepo(Conn)
	svc := service.NewUserService(repo)
	h := handlers.NewUserHandler(svc, JwtSecret)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /users", middleware.CheckTocken(JwtSecret, h.CreateUser))
	mux.HandleFunc("GET /users/{id}", middleware.CheckTocken(JwtSecret, h.GetUser))
	mux.HandleFunc("PUT /users/{id}", middleware.CheckTocken(JwtSecret, h.UpdateUser))
	mux.HandleFunc("DELETE /users/{id}", middleware.CheckTocken(JwtSecret, h.DeleteUser))
	mux.HandleFunc("POST /auth/register", h.Registr)
	mux.HandleFunc("POST /auth/login", h.Login)

	if err := http.ListenAndServe(":8050", mux); err != nil {
		log.Fatal(err)
	}

}
