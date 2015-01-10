package web

import (
	"github.com/gorilla/mux"
	"github.com/morpheusxaut/eveauth/misc"
	"net"
	"net/http"
	"strconv"
)

var (
	Router *HTTPRouter
)

type HTTPRouter struct {
	router *mux.Router
}

func SetupRouter() {
	Router = &HTTPRouter{
		router: mux.NewRouter().StrictSlash(true),
	}

	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = misc.WebLogging(handler, route.Name)

		Router.router.Methods(route.Methods...).Path(route.Pattern).Name(route.Name).Handler(handler)
	}

	Router.router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/assets")))
}

func (router *HTTPRouter) HandleRequests() {
	misc.Logger.Infof("Listening for HTTP requests on %q...", net.JoinHostPort(misc.Config.HTTPHost, strconv.Itoa(misc.Config.HTTPPort)))

	http.Handle("/", router.router)
	err := http.ListenAndServe(net.JoinHostPort(misc.Config.HTTPHost, strconv.Itoa(misc.Config.HTTPPort)), nil)

	misc.Logger.Criticalf("Received error while listening for HTTP requests: [%v]", err)
}
