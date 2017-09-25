package main

import (
	//"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"time"
)

// AuthLayer handles API authorization and jwt token management.
type AuthLayer struct {
	PrivateKey string
}

// Creates and returns a new AuthLayer instance.
// The method expects a private key environment variable to be set.
func NewAuthLayer() *AuthLayer {
	// TODO: What do we do if they var isn't set?
	return &AuthLayer{
		PrivateKey: os.Getenv("API_PRIVATE_KEY"),
	}
}

// GenerateToken creates a one-hour jwt token signed with auth.PrivateKey.
func (auth *AuthLayer) GenerateToken(w http.ResponseWriter) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["iss"] = "mutably"
	claims["exp"] = time.Now().Add(time.Hour).Unix()

	signedToken, _ := token.SignedString([]byte(auth.PrivateKey))

	w.Write([]byte(signedToken))
}
