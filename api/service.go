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
	router *mux.Router
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
		router: mux.NewRouter(),
		db:     database,
		port:   port,
	}

	return service, nil
}

// Start makes service begin listening for connections on the specified port.
func (service *Service) Start() error {
	log.Println("Starting service...")
	go func() {
		net.Dial("tcp", "localhost:"+service.port)
		log.Println("And we're live.")
	}()
	return http.ListenAndServe(":"+service.port, service.router)
}
