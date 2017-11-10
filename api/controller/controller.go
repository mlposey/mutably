package controller

import "net/http"

// Route defines how a REST resource is accessed.
type Route struct {
	Version     string
	Path        string
	Method      string
	Handler     http.HandlerFunc
	IsProtected bool
}

// A Controller connects a view (HTTP responses) with the
// application model (data and manipulation methods).
type Controller interface {
	Routes() []Route
}
