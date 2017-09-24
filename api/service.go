package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
	"strconv"
)

// Service controls communication between outside applications (that send HTTP
// requests) and main application logic such as authentication or database
// manipulation.
type Service struct {
	Router *mux.Router
	db     Database
	port   string
}

// NewService creates and returns a Service instance.
//
// database should be initialized and pointed at the application data store.
// port indicates where this service will listen for HTTP requests.
func NewService(database Database, port string) (*Service, error) {
	if database == nil {
		return nil, errors.New("Service given nil database")
	}
	if port == "" {
		return nil, errors.New("Service port can't be blank")
	}

	service := &Service{
		Router: mux.NewRouter(),
		db:     database,
		port:   port,
	}

	service.registerV1Routes()
	return service, nil
}

// registerV1Routes assigns routes to version 1 of the API.
func (s *Service) registerV1Routes() {
	v1 := s.Router.PathPrefix("/api/v1").Subrouter()

	v1.HandleFunc("/languages", s.getLanguages_v1).Methods("GET")
	v1.HandleFunc("/languages/{id:[0-9]+}", s.getLanguage_v1).Methods("GET")

	v1.HandleFunc("/words", s.getWords_v1).Methods("GET")
	v1.HandleFunc("/words/{id:[0-9]+}", s.getWord_v1).Methods("GET")
	// TODO: GET /words/{id:[0-9]+}/inflections

	// TODO: GET /users
	// TODO: POST /users
	// TODO: GET /users/{id:[0-9]+}
}

// Start makes service begin listening for connections on the specified port.
func (service *Service) Start() error {
	log.Println("Starting service")
	go func() {
		net.Dial("tcp", "localhost:"+service.port)
		log.Println("And we're live.")
	}()
	return http.ListenAndServe(":"+service.port, service.Router)
}

// makeJsonResponse creates and sends a json response to a writer.
func (service *Service) makeJsonResponse(w http.ResponseWriter, code int,
	respBody interface{}) {
	marshaledBody, _ := json.Marshal(respBody)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(marshaledBody)
}

// GET /api/v1/languages
func (service *Service) getLanguages_v1(w http.ResponseWriter, r *http.Request) {
	languages, err := service.db.GetLanguages()
	if len(languages) == 0 {
		if err != nil {
			log.Println(err)
		}
		service.makeJsonResponse(w, http.StatusNotFound,
			NewErrorResponse("no languages exist"))
	} else {
		service.makeJsonResponse(w, http.StatusOK, languages)
	}
}

// GET /api/v1/languages/{id:[0-9]+}
func (service *Service) getLanguage_v1(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	languageId, err := strconv.Atoi(vars["id"])
	if err != nil {
		service.makeJsonResponse(w, http.StatusBadRequest,
			NewErrorResponse("invalid language id"))
	}

	language, err := service.db.GetLanguage(languageId)
	if err != nil {
		service.makeJsonResponse(w, http.StatusNotFound,
			NewErrorResponse("language not found"))
	} else {
		service.makeJsonResponse(w, http.StatusOK, language)
	}
}

// GET /api/v1/words
func (service *Service) getWords_v1(w http.ResponseWriter, r *http.Request) {
	words, err := service.db.GetWords()
	if len(words) == 0 {
		if err != nil {
			log.Println(err)
		}
		service.makeJsonResponse(w, http.StatusNotFound,
			NewErrorResponse("no words exist"))
	} else {
		service.makeJsonResponse(w, http.StatusOK, words)
	}
}

// GET /api/v1/words/{id:[0-9]+}
func (service *Service) getWord_v1(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	wordId, err := strconv.Atoi(vars["id"])
	if err != nil {
		service.makeJsonResponse(w, http.StatusBadRequest,
			NewErrorResponse("invalid word id"))
	}

	word, err := service.db.GetWord(wordId)
	if err != nil {
		service.makeJsonResponse(w, http.StatusNotFound,
			NewErrorResponse("word not found"))
	} else {
		service.makeJsonResponse(w, http.StatusOK, word)
	}
}
