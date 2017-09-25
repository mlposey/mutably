package main

import (
	"encoding/json"
	"errors"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"strings"
	"time"
)

// AuthLayer handles API authorization and jwt token management.
type AuthLayer struct {
	PrivateKey string
	keyFunc    jwt.Keyfunc
	middleware *jwtmiddleware.JWTMiddleware
}

// Creates and returns a new AuthLayer instance.
// AuthLayer relies on a private key to sign tokens. That key should be defined
// in the environment variable API_PRIVATE_KEY before calling this function.
func NewAuthLayer() *AuthLayer {
	// TODO: What do we do if they var isn't set?
	auth := &AuthLayer{PrivateKey: os.Getenv("API_PRIVATE_KEY")}
	auth.keyFunc = func(token *jwt.Token) (interface{}, error) {
		return []byte(auth.PrivateKey), nil
	}

	auth.middleware = jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: auth.keyFunc,
		SigningMethod:       jwt.SigningMethodHS256,
	})
	return auth
}

// GenerateToken creates a one-hour jwt token signed with auth.PrivateKey.
func (auth *AuthLayer) GenerateToken(w http.ResponseWriter) {
	auth.GenerateTokenWithClaim(w, nil)
}

// GenerateTokenWithClaim creates a one-hour jwt token that is signed with
// auth.PrivateKey and contains a set of custom claims.
// If no additional claims are required, use AuthLayer.GenerateToken() instead.
func (auth *AuthLayer) GenerateTokenWithClaim(w http.ResponseWriter,
	customClaims map[string]interface{}) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["iss"] = "mutably"
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	if customClaims != nil {
		for attr, claim := range customClaims {
			claims[attr] = claim
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(auth.sign(token))
}

// GetClaims returns the claims from a token in the Authorization header.
// The token should follow the standard 'Bearer token' format.
func (auth *AuthLayer) GetClaims(r *http.Request) (jwt.MapClaims, error) {
	header := strings.Split(r.Header.Get("Authorization"), " ")
	if len(header) != 2 || header[0] != "Bearer" {
		return nil, errors.New("bad token header")
	}

	token, err := jwt.Parse(header[1], auth.keyFunc)
	if err != nil {
		return nil, err
	}

	return token.Claims.(jwt.MapClaims), nil
}

func (auth *AuthLayer) sign(token *jwt.Token) []byte {
	signedToken, _ := token.SignedString([]byte(auth.PrivateKey))
	jsonBody := map[string]string{"token": signedToken}
	resp, _ := json.Marshal(jsonBody)
	return resp
}

// Authenticate validates the handler's jwt token and proceeds if checks pass.
func (auth *AuthLayer) Authenticate(handler http.HandlerFunc) http.Handler {
	return auth.middleware.Handler(handler)
}
