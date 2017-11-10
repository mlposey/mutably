package model

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
)

// AuthLayer handles resource authorization.
//
// The API uses JSON Web Tokens to grant access to protected resources. Those
// tokens need to be generated and/or validated before changes to the system
// can be considered. AuthLayer encapsulates those operations and signs the
// tokens using a private key. See NewAuthLayer() for setup details.
type AuthLayer struct {
	PrivateKey string
	keyFunc    jwt.Keyfunc
	middleware *jwtmiddleware.JWTMiddleware
}

// Creates and returns a new AuthLayer instance.
// AuthLayer relies on a private key to sign tokens. That key should be defined
// in the environment variable API_PRIVATE_KEY before calling this function.
func NewAuthLayer() *AuthLayer {
	privateKey := os.Getenv("API_PRIVATE_KEY")
	if privateKey == "" {
		log.Fatal("Environment variable API_PRIVATE_KEY should not be empty")
	}

	auth := &AuthLayer{PrivateKey: privateKey}
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

// Authenticate returns a middleware handler that will use jwt to validate the
// Authorization header of a *http.Request.
func (auth *AuthLayer) Authenticate(handler http.HandlerFunc) http.HandlerFunc {
	return auth.middleware.Handler(handler).(http.HandlerFunc)
}

// GetCredentials reads a set of username:password credentials from a base64
// encoded basic authorization header. It returns the decoded username, password
// and nil error if the credentials existed and matched the colon format.
func (auth *AuthLayer) GetCredentials(r *http.Request) (string, string, error) {
	// Authorization: Basic base64gobblygoop
	encodedCreds := strings.Split(r.Header.Get("Authorization"), " ")
	if len(encodedCreds) != 2 {
		return "", "", errors.New("username:password required but missing")
	}

	// base64gobblygoop -> username:password
	decoded, err := base64.StdEncoding.DecodeString(encodedCreds[1])
	if err != nil {
		return "", "", err
	}

	// 'username:password' -> {'username', 'password'}
	credentials := strings.Split(string(decoded), ":")
	if len(credentials) != 2 {
		return "", "", errors.New("bad authorization string")
	}

	return credentials[0], credentials[1], nil
}

func (auth *AuthLayer) sign(token *jwt.Token) []byte {
	signedToken, _ := token.SignedString([]byte(auth.PrivateKey))
	jsonBody := map[string]string{"token": signedToken}
	resp, _ := json.Marshal(jsonBody)
	return resp
}
