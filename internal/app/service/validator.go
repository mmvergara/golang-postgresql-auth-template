package service

import "github.com/go-playground/validator/v10"

// https://github.com/go-playground/validator/blob/master/_examples/simple/main.go
var Validate *validator.Validate

func InititalizeValidator() {
	// Initialize the globally accessible validator
	Validate = validator.New(validator.WithRequiredStructEnabled())
}
