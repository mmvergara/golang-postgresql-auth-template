package server

import (
	"fmt"
	"golang-postgresql-auth-template/config"
	"golang-postgresql-auth-template/internal/app"
	"net/http"
	"time"
)

func NewServer() *http.Server {
	cfg := config.NewConfig()
	app := app.NewApplication(cfg)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.PORT),
		Handler:      RegisterRoutes(app),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	fmt.Println("Server is running on port:", cfg.PORT)

	return server
}
