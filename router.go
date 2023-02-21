package router

import "net/http"

type Router struct {
	routes []RouteEntry
}

type RouteEntry struct {
	Path    string
	Method  string
	Handler http.HandlerFunc
}

func (r *Router) Route(m, p string, h http.HandlerFunc) {
	e := RouteEntry{
		Method:  m,
		Path:    p,
		Handler: h,
	}
	r.routes = append(r.routes, e)
}

func (rtr *RouteEntry) MatchRoute(w http.ResponseWriter, r *http.Request) bool {
	allowed := true

	if r.URL.Path != rtr.Path {
		allowed = false
		http.NotFound(w, r)
	}

	if rtr.Method != r.Method {
		allowed = false
		http.Error(w, r.Method+" method not allowed for route "+rtr.Path, http.StatusMethodNotAllowed)
	}

	return allowed
}

func (rtr *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, e := range rtr.routes {
		match := e.MatchRoute(w, r)

		if match {
			e.Handler.ServeHTTP(w, r)
		}
	}

}
