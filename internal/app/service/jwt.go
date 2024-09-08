package service

import (
	"context"
	"fmt"
	"golang-postgresql-auth-template/internal/app/repository"
	"golang-postgresql-auth-template/internal/models"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AppJwt struct {
	userRepo         *repository.UserRepo
	secret           string
	tokenDurationHrs time.Duration
}

type CustomClaimsKey string

const JwtClaimsContextKey CustomClaimsKey = "claims"

// NewAppJwt creates and returns a new AppJwt instance.
//
// - userRepo: used to verify user existence during token validation.
//
// - secret: the key used to sign JWT tokens.
//
// - tokenDuration: the duration for which the tokens are valid.
//
//	appJwt := NewAppJwt(userRepo, "mysecretjwtkey", 24*time.Hour)
func NewAppJwt(userRepo *repository.UserRepo, secret string, tokenDurationHrs time.Duration) *AppJwt {
	return &AppJwt{
		userRepo:         userRepo,
		secret:           secret,
		tokenDurationHrs: tokenDurationHrs,
	}
}

func (a *AppJwt) CreateJwtWithClaims(userID models.UserID) (string, time.Time, error) {
	expirationTime := time.Now().Add(a.tokenDurationHrs)

	// https://datatracker.ietf.org/doc/html/rfc7519#section-4.1
	customClaims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		Issuer:    "golang-postgresql-auth-template",
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)
	tokenString, err := token.SignedString([]byte(a.secret))
	return tokenString, expirationTime, err
}

// Middleware that checks if the user is authenticated using a JWT cookie.
//
// This middleware extracts the JWT from the cookie, validates the token, and checks if the user exists in the system.
// If the token is valid and the user exists, it adds the JWT claims to the request context and calls the next handler.
// Otherwise, it responds with an appropriate HTTP error.
func (a *AppJwt) IsAuthenticatedMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtCookie, err := r.Cookie("jwt")
		if err != nil {
			http.Error(w, "Session expired or invalid. Please log in again.", http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(jwtCookie.Value, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid token signing method")
			}
			return []byte(a.secret), nil
		})

		if err != nil {
			log.Printf("Token validation failed: %v", err)
			http.Error(w, "Session expired or invalid. Please log in again.", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok || !token.Valid {
			log.Println("Invalid Token Claims")
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}

		if claims.ExpiresAt.Before(time.Now()) {
			log.Println("Token expired")
			http.Error(w, "Session expired or invalid.", http.StatusUnauthorized)
			return
		}

		exist, err := a.userRepo.CheckIfUserExistsByID(r.Context(), claims.Subject)
		if err != nil {
			log.Println("Server Error", err)
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}
		if !exist {
			log.Println("User does not exist, deleting JWT")
			a.DeleteJwt(w)
			http.Error(w, "User does not exist", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), JwtClaimsContextKey, claims)

		// User is authenticated, set the user ID in the request context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// SetJwtCookie sets the JWT token as an HttpOnly cookie in the response.
func (a *AppJwt) SetJwtCookie(w http.ResponseWriter, tokenString string, exp time.Time) {
	httpOnlyCookie := http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  exp,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &httpOnlyCookie)
}

// DeleteJwt deletes the JWT cookie from the response.
func (a *AppJwt) DeleteJwt(w http.ResponseWriter) {
	expirationTime := time.Now().Add(-1 * time.Hour)
	httpOnlyCookie := http.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  expirationTime,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &httpOnlyCookie)
}

// GetUserIDFromContext retrieves the user ID from the request context
func GetUserIDFromContext(r *http.Request) (models.UserID, bool) {
	claims, ok := r.Context().Value(JwtClaimsContextKey).(*jwt.RegisteredClaims)
	if !ok {
		return uuid.Nil, false
	}
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, false
	}
	return userID, true
}
