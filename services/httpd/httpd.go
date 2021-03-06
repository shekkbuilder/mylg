// Package httpd provides web service
package httpd

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"net/http"

	"github.com/mehrdadrad/mylg/icmp"
	// statik is single binary including all web stuff
	_ "github.com/mehrdadrad/mylg/services/dashboard/statik"
)

// Route represents a HTTP route
type Route struct {
	Name        string
	Method      string
	Path        string
	HandlerFunc http.HandlerFunc
}

var routes = []Route{
	Route{
		"API",
		"POST",
		"/api/{name}",
		API,
	},
}

// API handles API routes
func API(w http.ResponseWriter, r *http.Request) {
	var f = mux.Vars(r)["name"]
	switch f {
	case "ping":
		r.ParseForm()
		host := r.FormValue("host")
		p, err := icmp.NewPing(host + " -c 1")
		if err != nil {
			return
		}
		resp := p.Run()
		r := <-resp
		w.Write([]byte(fmt.Sprintf(`{"rtt": %.2f,"pl":0}`, r.RTT)))
	}
}

// Run starts web service
func Run() {
	statikFS, _ := fs.New()
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Path(route.Path).
			Methods(route.Method).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	router.PathPrefix("/").Handler(http.FileServer(statikFS))
	http.ListenAndServe("127.0.0.1:8080", router)
}
