package server

import (
	"golang-postgresql-auth-template/internal/app"
	"net/http"
)

func RegisterRoutes(app *app.Application) http.Handler {
	router := http.NewServeMux()
	middlewares := CreateMiddlewareStack(
		Cors,
		Timeout,
		Logger,
	)

	// Public Routes
	router.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello Mom"))
	})
	router.HandleFunc("POST /auth/signup", app.AuthHandler.HandleSignup())
	router.HandleFunc("POST /auth/signin", app.AuthHandler.HandleSignin())
	router.HandleFunc("POST /auth/signout", app.AuthHandler.HandleSignout)

	// V1 Protected Routes
	v1Protected := http.NewServeMux()
	v1Protected.HandleFunc("GET /auth-status", app.AuthHandler.HandleCheckIfSignedIn)

	// Add middleware to protected routes
	router.Handle("/v1/", http.StripPrefix("/v1", app.AppJwt.IsAuthenticatedMiddleware(v1Protected)))
	return middlewares(router)
}
