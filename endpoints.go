package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

/*
Method HTTP

typeNames

GET
The GET method is used to retrieve information from the given server using a given URI. Requests using GET should only retrieve data and should have no other effect on the data.

HEAD
Same as GET, but transfers the status line and header section only.

POST
A POST request is used to send data to the server, for example, customer information, file upload, etc. using HTML forms.

PUT
Replaces all current representations of the target resource with the uploaded content.

DELETE
Removes all current representations of the target resource given by a URI.

CONNECT
Establishes a tunnel to the server identified by a given URI.

OPTIONS
Describes the communication options for the target resource.

TRACE
Performs a message loop-back test along the path to the target resource.
*/

// Endpoint Standard
type Endpoint struct {
	path      string
	subRoutes []Endpoint
	methods   map[string]http.Handler
	secure    bool
}

// AllRoutes for api
var AllRoutes = [...]Endpoint{
	{
		path: "/ping",
		methods: map[string]http.Handler{
			http.MethodGet:  YourHandler,
			http.MethodPost: YourHandler,
		},
		subRoutes: []Endpoint{
			{
				path: "/{id}",
				methods: map[string]http.Handler{
					http.MethodGet: getBook,
					//http.MethodDelete: deleteBook,
				},
			},
		},
	},
}

// CreateRoutes for the app
func CreateRoutes(r *mux.Router, routes []Endpoint, rootURL string) {
	// Loop through endpoints
	for _, e := range routes {
		// Loop through methods
		for k, m := range e.methods {
			r.Handle(rootURL+e.path, m).Methods(k)
		}
		// Recurse if there are subRoutes
		if len(e.subRoutes) > 0 {
			CreateRoutes(r, e.subRoutes, rootURL+e.path)
		}
	}
}

// RouteWalker goes through and prints all routes
func RouteWalker(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	pathTemplate, err := route.GetPathTemplate()
	if err == nil {
		fmt.Println("ROUTE:", pathTemplate)
	}
	pathRegexp, err := route.GetPathRegexp()
	if err == nil {
		fmt.Println("Path regexp:", pathRegexp)
	}
	queriesTemplates, err := route.GetQueriesTemplates()
	if err == nil {
		fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
	}
	queriesRegexps, err := route.GetQueriesRegexp()
	if err == nil {
		fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
	}
	methods, err := route.GetMethods()
	if err == nil {
		fmt.Println("Methods:", strings.Join(methods, ","))
	}
	fmt.Println()
	return nil
}
