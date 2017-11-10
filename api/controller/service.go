package controller

import (
	"encoding/json"
	"errors"
	"log"
	"mutably/api/model"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Service controls communication between outside applications (that send HTTP
// requests) and main application logic such as authentication or database
// manipulation.
type Service struct {
	db model.Database

	Router *mux.Router
	// routers specific to API versions
	versionRouters map[string]*mux.Router
	port           string
	auth           *model.AuthLayer
}

// NewService creates and returns a Service instance.
//
// database should be initialized and pointed at the application data store.
// port indicates where this service will listen for HTTP requests.
func NewService(database model.Database, port string) (*Service, error) {
	if database == nil {
		return nil, errors.New("Service given nil database")
	}
	if port == "" {
		return nil, errors.New("Service port can't be blank")
	}

	service := &Service{
		db:             database,
		Router:         mux.NewRouter(),
		versionRouters: make(map[string]*mux.Router),
		port:           port,
		auth:           model.NewAuthLayer(),
	}

	service.AddController(&Users{db: database, auth: service.auth})
	service.AddController(&Languages{db: database})
	service.AddController(&Words{db: database})
	// Todo: Probably create a Controller if another endpoint is created for
	//       this resource.
	service.versionRouters["v1"].HandleFunc("/tokens",
		service.getToken).Methods("GET")

	return service, nil
}

// AddController registers a controller's routes under the /api path.
func (service *Service) AddController(controllers Controller) {
	for _, route := range controllers.Routes() {
		router, exists := service.versionRouters[route.Version]
		if !exists {
			router = service.Router.PathPrefix("/api/" + route.Version).Subrouter()
			service.versionRouters[route.Version] = router
		}

		if route.IsProtected {
			router.HandleFunc(route.Path,
				service.auth.Authenticate(route.Handler)).Methods(route.Method)
		} else {
			router.HandleFunc(route.Path, route.Handler).Methods(route.Method)
		}
	}
}

// Start makes service begin listening for connections on the specified port.
func (service *Service) Start() error {
	log.Println("Starting service...")
	go func() {
		net.Dial("tcp", "localhost:"+service.port)
		log.Println("And we're live.")
	}()
	corsHandler := cors.Default().Handler(service.Router)
	//return http.ListenAndServe(":"+service.port, service.Router)
	return http.ListenAndServe(":"+service.port, corsHandler)
}

// makeJsonResponse creates and sends a json response to a writer.
func makeJsonResponse(w http.ResponseWriter, code int,
	respBody interface{}) {
	marshaledBody, _ := json.Marshal(respBody)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(marshaledBody)
}

// makeErrorResponse places in the response body a JSON error message.
func makeErrorResponse(w http.ResponseWriter, code int,
	errMsg string) {
	makeJsonResponse(w, code, map[string]string{"error": errMsg})
}

// getAggregate faciliates 'get all' functionality on a resource.
// Its arguments are:
//   w - the writer to output results to
//   objects - a slice of the objects requested
//   length - the number of objects
//   err - any error returned when acquiring the objects; can be nil
//   resource - the name of the resource (e.g., 'words')
func respondWithAggregate(w http.ResponseWriter, objects interface{},
	length int, err error, resource string) {
	if length == 0 {
		if err != nil {
			log.Println(err)
		}
		makeErrorResponse(w, http.StatusNotFound, "no "+resource+" exist")
	} else {
		makeJsonResponse(w, http.StatusOK, objects)
	}
}

// GET /api/v1/tokens
func (service *Service) getToken(w http.ResponseWriter, r *http.Request) {
	username, password, err := service.auth.GetCredentials(r)
	if err != nil {
		makeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	userId := service.db.GetUserId(username, password)
	if userId == "" {
		makeErrorResponse(w, http.StatusUnauthorized,
			"invalid user credentials")
	} else {
		w.WriteHeader(http.StatusOK)
		service.auth.GenerateTokenWithClaim(w, map[string]interface{}{"id": userId})
	}
}
