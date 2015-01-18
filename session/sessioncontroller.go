package session

import (
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/morpheusxaut/eveauth/database"
	"github.com/morpheusxaut/eveauth/misc"

	"net/http"
	"os"
	"path/filepath"
)

type SessionController struct {
	config   *misc.Configuration
	database database.DatabaseConnection
	store    *sessions.FilesystemStore
}

func SetupSessionController(conf *misc.Configuration, db database.DatabaseConnection) *SessionController {
	controller := &SessionController{
		config:   conf,
		database: db,
		store:    sessions.NewFilesystemStore("web/sessions", []byte(securecookie.GenerateRandomKey(128))),
	}

	return controller
}

func (controller *SessionController) CleanSessions() error {
	sessions, err := filepath.Glob("app/sessions/session_*")
	if err != nil {
		return err
	}

	for _, s := range sessions {
		err = os.Remove(s)
		if err != nil {
			return err
		}
	}

	return nil
}

func (controller *SessionController) DestroySession(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "user",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:   "auth",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}
