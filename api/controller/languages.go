package controller

import (
	"net/http"
	"mutably/api/model"
	"strconv"
	"github.com/gorilla/mux"
)

// Languages is a Controller that handles the /languages resource.
type Languages struct {
	db model.Database
}

func (lang *Languages) Routes() []Route {
	return []Route{
		{ // GET /v1/languages
			Version:     "v1",
			Path:        "/languages",
			Method:      "GET",
			Handler:     lang.getLanguages,
			IsProtected: false,
		},
		{ // GET /v1/languages/{id:[0-9]+}
			Version:     "v1",
			Path:        "/languages/{id:[0-9]+}",
			Method:      "GET",
			Handler:     lang.getLanguage,
			IsProtected: false,
		},
	}
}

// GET /api/v1/languages
func (lang *Languages) getLanguages(w http.ResponseWriter, r *http.Request) {
	languages, err := lang.db.GetLanguages()
	respondWithAggregate(w, languages, len(languages), err, "languages")
}

// GET /api/v1/languages/{id:[0-9]+}
func (lang *Languages) getLanguage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	languageId, err := strconv.Atoi(vars["id"])
	if err != nil {
		makeErrorResponse(w, http.StatusBadRequest, "invalid language id")
	}

	language, err := lang.db.GetLanguage(languageId)
	if err != nil {
		makeErrorResponse(w, http.StatusNotFound, "language not found")
	} else {
		makeJsonResponse(w, http.StatusOK, language)
	}
}