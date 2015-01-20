package session

import (
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/morpheusxaut/eveauth/database"
	"github.com/morpheusxaut/eveauth/misc"
	"time"

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
		Name:   "eveauth_user",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:   "eveauth_login",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}

func (controller *SessionController) IsLoggedIn(w http.ResponseWriter, r *http.Request) (bool, error) {
	session, err := controller.store.Get(r, "eveauth_login")
	if err != nil {
		return false, err
	}

	if session.IsNew {
		return false, nil
	}

	timeStamp, ok := session.Values["timestamp"].(time.Time)
	if !ok {
		return false, fmt.Errorf("Failed to load timestamp from login session")
	}

	if time.Now().Sub(timeStamp).Minutes() >= 168 {
		controller.DestroySession(w, r)
		return false, fmt.Errorf("Login session expired")
	}

	return true, nil
}

func (controller *SessionController) SetLoginRedirect(w http.ResponseWriter, r *http.Request, redirect string) error {
	session, err := controller.store.Get(r, "eveauth_login")
	if err != nil {
		return err
	}

	session.Values["loginRedirect"] = redirect

	err = session.Save(r, w)
	if err != nil {
		return err
	}

	return nil
}

func (controller *SessionController) GetLoginRedirect(r *http.Request) (string, error) {
	session, err := controller.store.Get(r, "eveauth_login")
	if err != nil {
		return "/", err
	}

	if session.IsNew {
		return "/", nil
	}

	redirect, ok := session.Values["loginRedirect"].(string)
	if !ok {
		return "/", err
	}

	return redirect, nil
}
