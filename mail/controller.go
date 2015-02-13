package mail

import (
	"github.com/morpheusxaut/eveauth/database"
	"github.com/morpheusxaut/eveauth/misc"
)

// Controller handles sending mails via a given SMTP server as required by the app
type Controller struct {
	config   *misc.Configuration
	database *database.Connection
}

// SetupMailController initialises a new mail controller
func SetupMailController(conf *misc.Configuration, db *database.Connection) *Controller {
	controller := &Controller{
		config:   conf,
		database: db,
	}

	return controller
}
