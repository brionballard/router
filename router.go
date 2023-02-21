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

func (rtr *Router) Route(m, p string, h http.HandlerFunc) {
	e := RouteEntry{
		Method:  m,
		Path:    p,
		Handler: h,
	}
	rtr.routes = append(rtr.routes, e)
}

func (re *RouteEntry) MatchRoute(w http.ResponseWriter, r *http.Request) bool {
	allowed := true

	if r.URL.Path != re.Path {
		allowed = false
		http.NotFound(w, r)
	}

	if re.Method != r.Method {
		allowed = false
		http.Error(w, r.Method+" method not allowed for route "+re.Path, http.StatusMethodNotAllowed)
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
