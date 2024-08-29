package server

import "net/http"

type route struct {
	path    string
	handler http.HandlerFunc
}

var routes = [...]route{
	{
		path:    "POST /auth",
		handler: authHandler,
	},
	{
		path:    "POST /refresh",
		handler: refreshHandler,
	},
}

func Start() error {
	for _, r := range routes {
		http.HandleFunc(r.path, r.handler)
	}
	return http.ListenAndServe(":88", nil)
}
