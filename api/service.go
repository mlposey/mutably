package main

import (
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
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

func (service *Service) getLanguages_v1(w http.ResponseWriter, r *http.Request) {

}

func (service *Service) getLanguage_v1(w http.ResponseWriter, r *http.Request) {

}
