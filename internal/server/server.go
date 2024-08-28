package server

import "net/http"

type route struct {
	path    string
	handler http.HandlerFunc
}

var routes = [...]route{
	{
		path:    "/auth",
		handler: auth,
	},
	{
		path:    "/refresh",
		handler: refresh,
	},
}

func Start() error {
	for _, r := range routes {
		http.HandleFunc(r.path, r.handler)
	}
	return http.ListenAndServe(":88", nil)
}
