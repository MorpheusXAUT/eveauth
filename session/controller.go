package session

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/morpheusxaut/eveauth/database"
	"github.com/morpheusxaut/eveauth/misc"
	"github.com/morpheusxaut/eveauth/models"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

// Controller provides functionality to handle sessions and cached values as well as retrieval of data
type Controller struct {
	config   *misc.Configuration
	database database.Connection
	store    *sessions.FilesystemStore
}

// SetupSessionController prepares the controller's session store and sets a default session lifespan
func SetupSessionController(conf *misc.Configuration, db database.Connection) *Controller {
	controller := &Controller{
		config:   conf,
		database: db,
		store:    sessions.NewFilesystemStore("app/sessions", []byte(securecookie.GenerateRandomKey(128))),
	}

	controller.store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
	}

	return controller
}

// CleanSessions removes all old session files from disk
func (controller *Controller) CleanSessions() error {
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

// DestroySession destroys a user's session by invalidating the cookies used for storage
func (controller *Controller) DestroySession(w http.ResponseWriter, r *http.Request) {
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

// IsLoggedIn checks whether the user is currently logged in and has an appropriate timestamp set
func (controller *Controller) IsLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	session, _ := controller.store.Get(r, "eveauth_login")

	if session.IsNew {
		return false
	}

	timeStamp, ok := session.Values["timestamp"].(int64)
	if !ok {
		return false
	}

	if time.Now().Sub(time.Unix(timeStamp, 0)).Minutes() >= 168 {
		controller.DestroySession(w, r)
		return false
	}

	return true
}

// SetLoginRedirect saves the given path as a redirect after successful login
func (controller *Controller) SetLoginRedirect(w http.ResponseWriter, r *http.Request, redirect string) error {
	session, _ := controller.store.Get(r, "eveauth_login")

	session.Values["loginRedirect"] = redirect

	return session.Save(r, w)
}

// GetLoginRedirect retrieves the previously set path for redirection after login
func (controller *Controller) GetLoginRedirect(r *http.Request) string {
	session, _ := controller.store.Get(r, "eveauth_login")

	if session.IsNew {
		return "/"
	}

	redirect, ok := session.Values["loginRedirect"].(string)
	if !ok {
		return "/"
	}

	return redirect
}

// Authenticate validates the given username and password against the database and creates a new session with timestamp if successful
func (controller *Controller) Authenticate(w http.ResponseWriter, r *http.Request, username string, password string) error {
	storedPassword, err := controller.database.LoadPasswordForUser(username)

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		return err
	}

	session, _ := controller.store.Get(r, "eveauth_login")

	session.Values["username"] = username
	session.Values["timestamp"] = time.Now().Unix()

	return session.Save(r, w)
}

// CreateNewUser creates a new user in the database and saves the user's data in the current session
func (controller *Controller) CreateNewUser(w http.ResponseWriter, r *http.Request, username string, email string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.NewUser(username, string(hashedPassword), email, true)

	user, err = controller.database.SaveUser(user)
	if err != nil {
		return err
	}

	session, _ := controller.store.Get(r, "eveauth_login")

	session.Values["username"] = user.Username
	session.Values["timestamp"] = time.Now().Unix()

	return session.Save(r, w)
}
