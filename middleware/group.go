package middleware

import "net/http"

// Group routes that use the same middlewares
type Group struct {
	middlewares []func(http.Handler) http.Handler
	routes      map[string]http.Handler
}

func NewGroup() *Group {
	return &Group{
		routes: make(map[string]http.Handler),
	}
}

// register middleware
func (g *Group) Use(middlewares ...func(http.Handler) http.Handler) *Group {
	g.middlewares = append(g.middlewares, middlewares...)
	return g
}

// register a route
func (g *Group) Route(pattern string, handler http.HandlerFunc) *Group {
	g.routes[pattern] = handler
	return g // for chaining
}

// handle requests
func (g *Group) Handle(mux *http.ServeMux) {
	for pattern, handler := range g.routes {
		// Apply middleware in reverse order for correct execution sequence
		finalHandler := handler
		for i := len(g.middlewares) - 1; i >= 0; i-- {
			finalHandler = g.middlewares[i](finalHandler)
		}
		mux.Handle(pattern, finalHandler)
	}
}
