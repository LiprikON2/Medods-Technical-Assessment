package chi

import (
	"net/http"

	"github.com/go-chi/chi"
)

type Router interface {
	http.Handler
	// Routes

	Route(pattern string, fn func(r Router))
	Group(fn func(r Router))
	Use(middlewares ...func(http.Handler) http.Handler)
	With(middlewares ...func(http.Handler) http.Handler) Router
	Handle(pattern string, h http.Handler)
	HandleFunc(pattern string, h http.HandlerFunc)
	Method(method, pattern string, h http.Handler)
	MethodFunc(method, pattern string, h http.HandlerFunc)
	Get(pattern string, h http.HandlerFunc)
	Post(pattern string, h http.HandlerFunc)
	Put(pattern string, h http.HandlerFunc)
	Delete(pattern string, h http.HandlerFunc)
}

// ChiRouter wraps chi.Router to implement Router interface
type ChiRouter struct {
	chi.Router
}

func NewChiRouter() Router {
	return &ChiRouter{chi.NewRouter()}
}

// Methods that take Router as a parameter must be wrapped in order to convert between Router and chi.Router
func (r *ChiRouter) Route(pattern string, fn func(r Router)) {
	r.Router.Route(pattern, func(chiRouter chi.Router) {
		fn(&ChiRouter{chiRouter})
	})
}

func (r *ChiRouter) Group(fn func(r Router)) {
	r.Router.Group(func(chiRouter chi.Router) {
		fn(&ChiRouter{chiRouter})
	})
}

func (r *ChiRouter) With(middlewares ...func(http.Handler) http.Handler) Router {
	return &ChiRouter{r.Router.With(middlewares...)}
}
