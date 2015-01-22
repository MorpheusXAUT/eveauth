package web

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/morpheusxaut/eveauth/database"
	"github.com/morpheusxaut/eveauth/misc"
	"github.com/morpheusxaut/eveauth/session"

	"github.com/gorilla/mux"
)

type Controller struct {
	Config    *misc.Configuration
	Database  database.Connection
	Session   *session.Controller
	Templates *Templates
	Checksums *AssetChecksums

	router *mux.Router
}

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

	controller.router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/assets")))

	return controller
}

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

func (controller *Controller) HandleRequests() {
	misc.Logger.Infof("Listening for HTTP requests on %q...", net.JoinHostPort(controller.Config.HTTPHost, strconv.Itoa(controller.Config.HTTPPort)))

	http.Handle("/", controller.router)
	err := http.ListenAndServe(net.JoinHostPort(controller.Config.HTTPHost, strconv.Itoa(controller.Config.HTTPPort)), nil)

	misc.Logger.Criticalf("Received error while listening for HTTP requests: [%v]", err)
}
