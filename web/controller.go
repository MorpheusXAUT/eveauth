package web

import (
	"github.com/gorilla/mux"
	"github.com/morpheusxaut/eveauth/database"
	"github.com/morpheusxaut/eveauth/misc"
	"net"
	"net/http"
	"strconv"
)

type Controller struct {
	Config   *misc.Configuration
	Database database.DatabaseConnection
	router   *mux.Router
}

func SetupController(config *misc.Configuration, db database.DatabaseConnection) *Controller {
	controller := &Controller{
		Config:   config,
		Database: db,
		router:   mux.NewRouter().StrictSlash(true),
	}

	routes := SetupRoutes(controller)

	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = misc.WebLogging(handler, route.Name)

		controller.router.Methods(route.Methods...).Path(route.Pattern).Name(route.Name).Handler(handler)
	}

	controller.router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/assets")))

	return controller
}

func (controller *Controller) HandleRequests() {
	misc.Logger.Infof("Listening for HTTP requests on %q...", net.JoinHostPort(controller.Config.HTTPHost, strconv.Itoa(controller.Config.HTTPPort)))

	http.Handle("/", controller.router)
	err := http.ListenAndServe(net.JoinHostPort(controller.Config.HTTPHost, strconv.Itoa(controller.Config.HTTPPort)), nil)

	misc.Logger.Criticalf("Received error while listening for HTTP requests: [%v]", err)
}
