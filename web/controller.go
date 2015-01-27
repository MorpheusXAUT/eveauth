package web

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/morpheusxaut/eveauth/database"
	"github.com/morpheusxaut/eveauth/misc"
	"github.com/morpheusxaut/eveauth/session"

	"github.com/gorilla/mux"
)

// Controller provides functionality for handling web requests and accessing session and backend data
type Controller struct {
	Config    *misc.Configuration
	Database  database.Connection
	Session   *session.Controller
	Templates *Templates
	Checksums *AssetChecksums

	router *mux.Router
}

// SetupController prepares the web controller and initialises the router and handled routes
func SetupController(config *misc.Configuration, db database.Connection, sessions *session.Controller, templates *Templates, checksums *AssetChecksums) *Controller {
	controller := &Controller{
		Config:    config,
		Database:  db,
		Session:   sessions,
		Templates: templates,
		Checksums: checksums,
		router:    mux.NewRouter().StrictSlash(true),
	}

	routes := SetupRoutes(controller)

	for _, route := range routes {
		controller.router.Methods(route.Methods...).Path(route.Pattern).Name(route.Name).Handler(controller.ServeHTTP(route.HandlerFunc, route.Name))
	}

	controller.router.PathPrefix("/").Handler(http.FileServer(http.Dir("app/assets")))

	return controller
}

// ServeHTTP acts as a middleware between parsed requests, logging the requests and replacing the remote address with the proxy-value if needed
func (controller *Controller) ServeHTTP(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		remoteAddr := r.Header.Get("X-Forwarded-For")

		if len(remoteAddr) > 0 {
			remoteAddrs := strings.Split(remoteAddr, ", ")
			if len(remoteAddrs) > 1 {
				r.RemoteAddr = fmt.Sprintf("%s:0", remoteAddrs[0])
			} else {
				r.RemoteAddr = fmt.Sprintf("%s:0", remoteAddr)
			}
		}

		if controller.Config.DebugTemplates {
			controller.Templates.ReloadTemplates()
		}

		inner.ServeHTTP(w, r)

		misc.Logger.Debugf("ServeHTTP: [%s] %s %q {%s} - %s ", r.Method, r.RemoteAddr, r.RequestURI, name, time.Since(start))
	})
}

// HandleRequests starts the blocking call to handle web requests
func (controller *Controller) HandleRequests() {
	misc.Logger.Infof("Listening for HTTP requests on %q...", controller.Config.HTTPHost)

	http.Handle("/", controller.router)
	err := http.ListenAndServe(controller.Config.HTTPHost, nil)

	misc.Logger.Criticalf("Received error while listening for HTTP requests: [%v]", err)
}
