package api

import (
	"database/sql"
	"fmt"
	"jwt-auth/internal/auth"
	"jwt-auth/internal/middleware"
	"jwt-auth/internal/user"
	"log"
	"net/http"
	"time"
)

type APIServer struct {
	addr      string
	db        *sql.DB
	jwtManger *auth.JWTManager
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	j := auth.NewJWTManager(
		"mon-secret",
		20*time.Minute,
	)

	return &APIServer{addr: addr, db: db, jwtManger: j}
}

func (s *APIServer) Run() {
	//* Injection de dependaces

	userRepository := user.NewRepository(s.db)

	authService := auth.NewService(userRepository, s.jwtManger)
	userService := user.NewService(userRepository)

	authHandler := auth.NewHandler(authService)
	userHandler := user.NewHandler(userService)

	//* Routing
	mux := http.NewServeMux()

	mux.HandleFunc("/auth/register", authHandler.Register)
	mux.HandleFunc("/auth/login", authHandler.Login)

	mux.Handle(
		"/users/me",
		middleware.AuthMiddleware(
			s.jwtManger,
			http.HandlerFunc(userHandler.Me),
		),
	)

	//* Server
	fmt.Println("Server start on", s.addr)

	if err := http.ListenAndServe(s.addr, mux); err != nil {
		log.Fatal("error server: ", err)
	}

}
