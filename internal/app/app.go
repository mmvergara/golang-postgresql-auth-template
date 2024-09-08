package app

import (
	"golang-postgresql-auth-template/config"
	"golang-postgresql-auth-template/internal/app/handler"
	"golang-postgresql-auth-template/internal/app/repository"
	"golang-postgresql-auth-template/internal/app/service"
	"golang-postgresql-auth-template/internal/database"
	"time"
)

type Application struct {
	AuthHandler *handler.AuthHandler
	// Add more handlers or services needed here
	// Example
	// Book Handler *handler.BookHandler
	AppJwt *service.AppJwt
}

func NewApplication(cfg *config.Config) *Application {
	db := database.NewPostgres(cfg).Sql

	// Repositories
	userRepo := repository.NewUserRepo(db)
	// bookRepo := repository.NewBookRepo(db)

	// Services
	appJwt := service.NewAppJwt(userRepo, cfg.JWT_SECRET, time.Duration(cfg.JWT_TOKEN_DURATION_HOURS)*time.Hour) // 24 hours
	// ex. cloudinary := service.NewCloudinary(cfg.CLOUDINARY_CLOUD_NAME, cfg.CLOUDINARY_API_KEY, cfg.CLOUDINARY_API_SECRET)

	// Initialize the single instance of the struct validator
	service.InititalizeValidator()

	return &Application{
		AuthHandler: handler.NewAuthHandler(userRepo, appJwt),
		// BookHandler: handler.NewBookHandler(bookRepo, cloudinary),
		AppJwt: appJwt,
	}
}
