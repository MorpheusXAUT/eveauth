package session

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/morpheusxaut/eveauth/database"
	"github.com/morpheusxaut/eveauth/mail"
	"github.com/morpheusxaut/eveauth/misc"
	"github.com/morpheusxaut/eveauth/models"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/morpheusxaut/redistore"
	"golang.org/x/crypto/bcrypt"
)

// Controller provides functionality to handle sessions and cached values as well as retrieval of data
type Controller struct {
	config   *misc.Configuration
	database database.Connection
	mail     *mail.Controller
	store    *redistore.RediStore
}

// SetupSessionController prepares the controller's session store and sets a default session lifespan
func SetupSessionController(conf *misc.Configuration, db database.Connection, mailer *mail.Controller) (*Controller, error) {
	controller := &Controller{
		config:   conf,
		database: db,
		mail:     mailer,
	}

	store, err := redistore.NewRediStoreWithDB(10, "tcp", controller.config.RedisHost, controller.config.RedisPassword, "0", securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))
	if err != nil {
		return nil, err
	}

	controller.store = store

	controller.store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
	}

	gob.Register(&models.User{})

	return controller, nil
}

// DestroySession destroys a user's session by setting a negative maximum age
func (controller *Controller) DestroySession(w http.ResponseWriter, r *http.Request) {
	loginSession, _ := controller.store.Get(r, "eveauthLogin")
	dataSession, _ := controller.store.Get(r, "eveauthData")

	loginSession.Options.MaxAge = -1
	dataSession.Options.MaxAge = -1

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

	verifiedEmail, ok := loginSession.Values["verifiedEmail"].(bool)
	if !ok {
		return false
	}

	if !verifiedEmail {
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

	loginAttempt := models.NewLoginAttempt(username, r.RemoteAddr, r.UserAgent(), (err == nil))

	logErr := controller.database.SaveAuthenticationAttempt(loginAttempt)
	if logErr != nil {
		misc.Logger.Errorf("Failed to log authentication attempt: [%v]", logErr)
	}

	if err != nil {
		return err
	}

	user, err := controller.database.LoadUserFromUsername(username)
	if err != nil {
		return err
	}

	user, err = controller.SetUser(w, r, user)
	if err != nil {
		return err
	}

	loginSession, _ := controller.store.Get(r, "eveauthLogin")

	loginSession.Values["username"] = user.Username
	loginSession.Values["userID"] = user.ID
	loginSession.Values["timestamp"] = time.Now().Unix()
	loginSession.Values["verifiedEmail"] = user.VerifiedEmail

	return sessions.Save(r, w)
}

// CreateNewUser creates a new user in the database and saves the user's data in the current session
func (controller *Controller) CreateNewUser(w http.ResponseWriter, r *http.Request, username string, email string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.NewUser(username, string(hashedPassword), email, false, true)

	user, err = controller.SetUser(w, r, user)
	if err != nil {
		return err
	}

	loginSession, _ := controller.store.Get(r, "eveauthLogin")

	loginSession.Values["username"] = user.Username
	loginSession.Values["userID"] = user.ID
	loginSession.Values["verifiedEmail"] = user.VerifiedEmail

	return sessions.Save(r, w)
}

// SendEmailVerification sends an email with a verification link to the given address
func (controller *Controller) SendEmailVerification(w http.ResponseWriter, r *http.Request, username string, email string) error {
	verification := misc.GenerateRandomString(32)

	err := controller.mail.SendEmailVerification(username, email, verification)
	if err != nil {
		return err
	}

	loginSession, _ := controller.store.Get(r, "eveauthLogin")

	loginSession.Values["emailVerification"] = verification

	return sessions.Save(r, w)
}

// ResendEmailVerification resends an email with a verification link to the given address
func (controller *Controller) ResendEmailVerification(w http.ResponseWriter, r *http.Request, username string, email string) error {
	user, err := controller.database.LoadUserFromUsername(username)
	if err != nil {
		return err
	}

	user, err = controller.SetUser(w, r, user)
	if err != nil {
		return err
	}

	verification := misc.GenerateRandomString(32)

	err = controller.mail.SendEmailVerification(user.Username, email, verification)
	if err != nil {
		return err
	}

	loginSession, _ := controller.store.Get(r, "eveauthLogin")

	loginSession.Values["emailVerification"] = verification

	return sessions.Save(r, w)
}

// VerifyEmail checks the given code and verifies the presented email address is correct, currently not implemented
func (controller *Controller) VerifyEmail(w http.ResponseWriter, r *http.Request, email string, verification string) error {
	loginSession, _ := controller.store.Get(r, "eveauthLogin")

	emailVerification, ok := loginSession.Values["emailVerification"].(string)
	if !ok {
		return fmt.Errorf("Failed to retrieve verification code from login session")
	}

	if !strings.EqualFold(emailVerification, verification) {
		return fmt.Errorf("Failed to verify email address")
	}

	user, err := controller.GetUser(r)
	if err != nil {
		return err
	}

	user.VerifiedEmail = true

	user, err = controller.SetUser(w, r, user)

	loginSession.Values["timestamp"] = time.Now().Unix()
	loginSession.Values["verifiedEmail"] = user.VerifiedEmail

	return sessions.Save(r, w)
}

// SendPasswordReset sends an email with a verification link to reset a user's password to the given address
func (controller *Controller) SendPasswordReset(w http.ResponseWriter, r *http.Request, username string, email string) error {
	user, err := controller.database.LoadUserFromUsername(username)
	if err != nil {
		return err
	}

	if !strings.EqualFold(email, user.Email) {
		return fmt.Errorf("Email addresses do not match")
	}

	verification := misc.GenerateRandomString(32)

	err = controller.mail.SendPasswordReset(username, email, verification)
	if err != nil {
		return err
	}

	loginSession, _ := controller.store.Get(r, "eveauthLogin")

	loginSession.Values["passwordReset"] = verification

	return sessions.Save(r, w)
}

func (controller *Controller) VerifyPasswordReset(w http.ResponseWriter, r *http.Request, email string, username string, verification string, password string) error {
	user, err := controller.database.LoadUserFromUsername(username)
	if err != nil {
		return err
	}

	if !strings.EqualFold(email, user.Email) {
		return fmt.Errorf("Email addresses do not match")
	}

	loginSession, _ := controller.store.Get(r, "eveauthLogin")

	passwordReset, ok := loginSession.Values["passwordReset"].(string)
	if !ok {
		return fmt.Errorf("Failed to retrieve password reset code from login session")
	}

	if !strings.EqualFold(passwordReset, verification) {
		return fmt.Errorf("Failed to verify password reset code")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	user, err = controller.SetUser(w, r, user)
	if err != nil {
		return err
	}

	return nil
}

// SaveAPIKey saves the given API key ID and verification code to the database and updated the user-object in the data session
func (controller *Controller) SaveAPIKey(w http.ResponseWriter, r *http.Request, apiKeyID string, apivCode string) error {
	dataSession, _ := controller.store.Get(r, "eveauthData")

	user, ok := dataSession.Values["user"].(*models.User)
	if !ok {
		return fmt.Errorf("Failed to retrieve user from data session")
	}

	keyID, err := strconv.ParseInt(apiKeyID, 10, 64)
	if err != nil {
		return err
	}

	account := models.NewAccount(user.ID, keyID, apivCode, 0, false, true)

	apiClient := misc.CreateAPIClient(account)

	apiInfo, err := apiClient.Info()
	if err != nil {
		return err
	}

	accessMask, err := strconv.ParseInt(apiInfo.AccessMask, 10, 64)
	if err != nil {
		return err
	}

	account.APIAccessMask = int(accessMask)

	accountCharacters, err := apiClient.AccountCharacters()
	if err != nil {
		return nil
	}

	for _, accountChar := range accountCharacters {
		corporationID, err := strconv.ParseInt(accountChar.CorporationID, 10, 64)
		if err != nil {
			return err
		}

		accountCharID, err := strconv.ParseInt(accountChar.ID, 10, 64)
		if err != nil {
			return err
		}

		// TODO handle non-existent corporation gracefully by fetching and creating it
		corporation, err := controller.database.LoadCorporationFromEVECorporationID(corporationID)
		if err != nil {
			return err
		}

		character := models.NewCharacter(account.ID, corporation.ID, accountChar.Name, accountCharID, false, true)

		account.Characters = append(account.Characters, character)
	}

	user.Accounts = append(user.Accounts, account)

	user, err = controller.database.SaveUser(user)
	if err != nil {
		return err
	}

	dataSession.Values["user"] = user

	return sessions.Save(r, w)
}

// GetUser returns the user-object stored in the data session
func (controller *Controller) GetUser(r *http.Request) (*models.User, error) {
	dataSession, _ := controller.store.Get(r, "eveauthData")

	user, ok := dataSession.Values["user"].(*models.User)
	if !ok {
		return nil, fmt.Errorf("Failed to retrieve user from data session")
	}

	return user, nil
}

// SetUser saves the given user object to the database and updates the data session reference
func (controller *Controller) SetUser(w http.ResponseWriter, r *http.Request, user *models.User) (*models.User, error) {
	user, err := controller.database.SaveUser(user)
	if err != nil {
		return nil, err
	}

	dataSession, _ := controller.store.Get(r, "eveauthData")

	dataSession.Values["user"] = user

	return user, sessions.Save(r, w)
}

// VerifyApplication verifies the application to be authorized to perform requests to the auth backend
func (controller *Controller) VerifyApplication(appID string, callback string, auth string) (*models.Application, error) {
	applicationID, err := strconv.ParseInt(appID, 10, 64)
	if err != nil {
		return nil, err
	}

	application, err := controller.database.LoadApplication(applicationID)
	if err != nil {
		return nil, err
	}

	verified := misc.VerifyMessageHMACSHA256(fmt.Sprintf("%d:%s", application.ID, application.Callback), auth, application.Secret)

	if !verified {
		return nil, fmt.Errorf("Failed to verify HMAC")
	}

	return application, nil
}

// EncodeUserPermissions encodes the user's current permissions in a JSON struct and returns the encrypted payload
func (controller *Controller) EncodeUserPermissions(r *http.Request, application *models.Application) (string, error) {
	dataSession, _ := controller.store.Get(r, "eveauthData")

	user, ok := dataSession.Values["user"].(*models.User)
	if !ok {
		return "", fmt.Errorf("Failed to retrieve user from data session")
	}

	authUser := user.ToAuthUser()

	payload, err := json.Marshal(authUser)
	if err != nil {
		return "", err
	}

	encryptedPayload, err := misc.EncryptAESCFB(string(payload), application.Secret)
	if err != nil {
		return "", err
	}

	return encryptedPayload, nil
}
