package zebra

import "net/http"

type Router struct {
	Middlewares map[string]func(ctx Request, callback Callback)
}

type Request struct {
	http.Request
	PathVariables map[string]string
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
