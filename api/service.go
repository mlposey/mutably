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
	auth   *AuthLayer
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
		auth:   NewAuthLayer(),
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
	v1.HandleFunc("/words/{word}/inflections", s.getInflections_v1).Methods("GET")

	v1.Handle("/users", s.auth.Authenticate(s.getUsers_v1)).Methods("GET")
	v1.HandleFunc("/users", s.createUser_v1).Methods("POST")
	v1.Handle("/users/{id}", s.auth.Authenticate(s.getUser_v1)).Methods("GET")

	v1.HandleFunc("/sessions", s.createSessionTok_v1).Methods("GET")
}

// Start makes service begin listening for connections on the specified port.
func (service *Service) Start() error {
	log.Println("Starting service...")
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

// makeErrorResponse places in the response body a JSON error message.
func (service *Service) makeErrorResponse(w http.ResponseWriter, code int,
	errMsg string) {
	service.makeJsonResponse(w, code, map[string]string{"error": errMsg})
}

// getAggregate faciliates 'get all' functionality on a resource.
// Its arguments are:
//   w - the writer to output results to
//   objects - a slice of the objects requested
//   length - the number of objects
//   err - any error returned when acquiring the objects; can be nil
//   resource - the name of the resource (e.g., 'words')
func (service *Service) respondWithAggregate(w http.ResponseWriter,
	objects interface{}, length int, err error, resource string) {
	if length == 0 {
		if err != nil {
			log.Println(err)
		}
		service.makeErrorResponse(w, http.StatusNotFound, "no "+resource+" exist")
	} else {
		service.makeJsonResponse(w, http.StatusOK, objects)
	}
}

// GET /api/v1/languages
func (service *Service) getLanguages_v1(w http.ResponseWriter, r *http.Request) {
	languages, err := service.db.GetLanguages()
	service.respondWithAggregate(w, languages, len(languages), err, "languages")
}

// GET /api/v1/languages/{id:[0-9]+}
func (service *Service) getLanguage_v1(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	languageId, err := strconv.Atoi(vars["id"])
	if err != nil {
		service.makeErrorResponse(w, http.StatusBadRequest, "invalid language id")
	}

	language, err := service.db.GetLanguage(languageId)
	if err != nil {
		service.makeErrorResponse(w, http.StatusNotFound, "language not found")
	} else {
		service.makeJsonResponse(w, http.StatusOK, language)
	}
}

// GET /api/v1/words
func (service *Service) getWords_v1(w http.ResponseWriter, r *http.Request) {
	words, err := service.db.GetWords()
	service.respondWithAggregate(w, words, len(words), err, "words")
}

// GET /api/v1/words/{id:[0-9]+}
func (service *Service) getWord_v1(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	wordId, err := strconv.Atoi(vars["id"])
	if err != nil {
		service.makeErrorResponse(w, http.StatusBadRequest, "invalid word id")
	}

	word, err := service.db.GetWord(wordId)
	if err != nil {
		service.makeErrorResponse(w, http.StatusNotFound, "word not found")
	} else {
		service.makeJsonResponse(w, http.StatusOK, word)
	}
}

// GET /api/v1/words/{word}/inflections
func (service *Service) getInflections_v1(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	table, err := service.db.GetConjugationTable(vars["word"])
	if err != nil {
		service.makeErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}
	service.makeJsonResponse(w, http.StatusOK, table)
}

// GET /api/v1/users
func (service *Service) getUsers_v1(w http.ResponseWriter, r *http.Request) {
	claims, err := service.auth.GetClaims(r)
	if err != nil || service.db.IsAdmin(claims["id"].(string)) {
		users, err := service.db.GetUsers()
		service.respondWithAggregate(w, users, len(users), err, "users")
	} else {
		service.makeErrorResponse(w, http.StatusForbidden,
			"resource requires admin privileges")
	}
}

// POST /api/v1/users
func (service *Service) createUser_v1(w http.ResponseWriter, r *http.Request) {
	username, password, err := service.auth.GetCredentials(r)
	if err != nil {
		service.makeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Create resource.
	var userId string
	userId, err = service.db.CreateUser(username, password)
	if err != nil {
		service.makeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
	service.auth.GenerateTokenWithClaim(w, map[string]interface{}{"id": userId})
}

// GET /api/v1/users/{id}
func (service *Service) getUser_v1(w http.ResponseWriter, r *http.Request) {
	claims, err := service.auth.GetClaims(r)
	if err != nil {
		service.makeErrorResponse(w, http.StatusForbidden, err.Error())
		return
	}

	vars := mux.Vars(r)
	if vars["id"] == claims["id"] || service.db.IsAdmin(claims["id"].(string)) {
		user, err := service.db.GetUser(vars["id"])
		if err != nil {
			service.makeErrorResponse(w, http.StatusNotFound, "user not found")
		} else {
			service.makeJsonResponse(w, http.StatusOK, user)
		}
	} else {
		service.makeErrorResponse(w, http.StatusUnauthorized,
			"insufficient permissions")
	}
}

// GET /api/v1/sessions
func (service *Service) createSessionTok_v1(w http.ResponseWriter, r *http.Request) {
	username, password, err := service.auth.GetCredentials(r)
	if err != nil {
		service.makeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	userId := service.db.GetUserId(username, password)

	if userId == "" {
		service.makeErrorResponse(w, http.StatusUnauthorized,
			"invalid user credentials")
	} else {
		w.WriteHeader(http.StatusOK)
		service.auth.GenerateTokenWithClaim(w, map[string]interface{}{"id": userId})
	}
}
