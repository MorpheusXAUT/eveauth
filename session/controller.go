package session

import (
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/morpheusxaut/eveauth/database"
	"github.com/morpheusxaut/eveauth/misc"
	"github.com/morpheusxaut/eveauth/models"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/boj/redistore.v1"
)

// Controller provides functionality to handle sessions and cached values as well as retrieval of data
type Controller struct {
	config   *misc.Configuration
	database database.Connection
	store    *redistore.RediStore
}

// SetupSessionController prepares the controller's session store and sets a default session lifespan
func SetupSessionController(conf *misc.Configuration, db database.Connection) (*Controller, error) {
	controller := &Controller{
		config:   conf,
		database: db,
	}

	store, err := redistore.NewRediStore(10, "tcp", net.JoinHostPort(controller.config.RedisHost, strconv.Itoa(controller.config.RedisPort)), controller.config.RedisPassword, securecookie.GenerateRandomKey(128))
	if err != nil {
		return nil, err
	}

	controller.store = store

	controller.store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
	}

	return controller, nil
}

// DestroySession destroys a user's session by setting a negative maximum age
func (controller *Controller) DestroySession(w http.ResponseWriter, r *http.Request) {
	loginSession, _ := controller.store.Get(r, "eveauthLogin")
	userSession, _ := controller.store.Get(r, "eveauthUser")

	loginSession.Options.MaxAge = -1
	userSession.Options.MaxAge = -1

	err := sessions.Save(r, w)
	if err != nil {
		misc.Logger.Errorf("Failed to destroy session: [%v]", err)
	}
}

// IsLoggedIn checks whether the user is currently logged in and has an appropriate timestamp set
func (controller *Controller) IsLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	loginSession, _ := controller.store.Get(r, "eveauthLogin")

	if loginSession.IsNew {
		return false
	}

	timeStamp, ok := loginSession.Values["timestamp"].(int64)
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
	loginSession, _ := controller.store.Get(r, "eveauthLogin")

	loginSession.Values["loginRedirect"] = redirect

	return loginSession.Save(r, w)
}

// GetLoginRedirect retrieves the previously set path for redirection after login
func (controller *Controller) GetLoginRedirect(r *http.Request) string {
	loginSession, _ := controller.store.Get(r, "eveauthLogin")

	if loginSession.IsNew {
		return "/"
	}

	redirect, ok := loginSession.Values["loginRedirect"].(string)
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

	loginSession, _ := controller.store.Get(r, "eveauthLogin")

	loginSession.Values["username"] = username
	loginSession.Values["timestamp"] = time.Now().Unix()

	return loginSession.Save(r, w)
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

	loginSession, _ := controller.store.Get(r, "eveauthLogin")

	loginSession.Values["username"] = user.Username

	return loginSession.Save(r, w)
}

// SendEmailVerification sends an email with a verification link to the given address, currently not implemented
func (controller *Controller) SendEmailVerification(username string, email string) error {
	verification := misc.GenerateRandomString(32)

	// TODO actual implementation, skipped for now
	misc.Logger.Tracef("Sending email verification for user %q to email %q using verification code %q", username, email, verification)

	return nil
}

// VerifyEmail checks the given code and verifies the presented email address is correct, currently not implemented
func (controller *Controller) VerifyEmail(w http.ResponseWriter, r *http.Request, email string, verification string) error {
	// TODO actual implementation, skipped for now
	misc.Logger.Tracef("Verifying email %q using verification code %q", email, verification)

	loginSession, _ := controller.store.Get(r, "eveauthLogin")

	loginSession.Values["timestamp"] = time.Now().Unix()

	return loginSession.Save(r, w)
}
