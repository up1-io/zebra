package zebra

import (
	"net/http"
	"strings"
)

type Router struct {
	Middlewares map[string]func(ctx Request, callback Callback)
}

// Request is a wrapper around http.Request with additional path variables map.
type Request struct {
	http.Request
	PathVariables PathVariables
}

type Result struct {
	Redirect string
	Data     map[string]interface{}
}

type Callback func(err error, result Result)

func NewRouter() Router {
	return Router{
		Middlewares: make(map[string]func(r Request, callback Callback)),
	}
}

func (r *Router) On(path string, middleware func(r Request, callback Callback)) {
	r.Middlewares[path] = middleware
}

func (r *Router) getMiddlewareByURL(url string) func(r Request, callback Callback) {
	return r.Middlewares[url]
}

// PathVariables is a map of path variables
type PathVariables map[string]string

func (v PathVariables) Get(key string) string {
	return v[key]
}

func getPathVars(url string, requestURL string) PathVariables {
	pathVars := make(PathVariables)

	urlParts := strings.Split(url, "/")
	requestURLParts := strings.Split(requestURL, "/")

	for i, part := range urlParts {
		if strings.HasPrefix(part, "{") {
			key := part[1 : len(part)-1]
			pathVars[key] = requestURLParts[i]
		}
	}

	return pathVars
}
