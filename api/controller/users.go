package controller

import (
	"mutably/api/model"
	"net/http"

	"github.com/gorilla/mux"
)

// Users is a Controller for the /users resource.
type Users struct {
	db   model.Database
	auth *model.AuthLayer
}

func (u *Users) Routes() []Route {
	return []Route{
		{ // GET /api/v1/users
			Version:     "v1",
			Path:        "/users",
			Method:      "GET",
			Handler:     u.getUsers,
			IsProtected: true,
		},
		{ // POST /api/v1/users
			Version:     "v1",
			Path:        "/users",
			Method:      "POST",
			Handler:     u.createUser,
			IsProtected: false,
		},
		{ // GET /api/v1/users/{id}
			Version:     "v1",
			Path:        "/users/{id}",
			Method:      "GET",
			Handler:     u.getUser,
			IsProtected: true,
		},
	}
}

// GET /api/v1/users
func (u *Users) getUsers(w http.ResponseWriter, r *http.Request) {
	claims, err := u.auth.GetClaims(r)
	if err != nil || u.db.IsAdmin(claims["id"].(string)) {
		users, err := u.db.GetUsers()
		respondWithAggregate(w, users, len(users), err, "users")
	} else {
		makeErrorResponse(w, http.StatusForbidden,
			"resource requires admin privileges")
	}
}

// POST /api/v1/users
func (u *Users) createUser(w http.ResponseWriter, r *http.Request) {
	username, password, err := u.auth.GetCredentials(r)
	if err != nil {
		makeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Create resource.
	var userId string
	userId, err = u.db.CreateUser(username, password)
	if err != nil {
		makeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
	u.auth.GenerateTokenWithClaim(w, map[string]interface{}{"id": userId})
}

// GET /api/v1/users/{id}
func (u *Users) getUser(w http.ResponseWriter, r *http.Request) {
	claims, err := u.auth.GetClaims(r)
	if err != nil {
		makeErrorResponse(w, http.StatusForbidden, err.Error())
		return
	}

	vars := mux.Vars(r)
	if vars["id"] == claims["id"] || u.db.IsAdmin(claims["id"].(string)) {
		user, err := u.db.GetUser(vars["id"])
		if err != nil {
			makeErrorResponse(w, http.StatusNotFound, "user not found")
		} else {
			makeJsonResponse(w, http.StatusOK, user)
		}
	} else {
		makeErrorResponse(w, http.StatusUnauthorized,
			"insufficient permissions")
	}
}
