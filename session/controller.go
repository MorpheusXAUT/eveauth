package session

import (
	"database/sql"
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

	"github.com/boj/redistore"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/guregu/null.v2/zero"
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

	store, err := redistore.NewRediStoreWithDB(10, "tcp", controller.config.RedisHost, controller.config.RedisPassword, controller.config.RedisDB, securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))
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

// SetCSRFToken saves the given CSRF token for the current session
func (controller *Controller) SetCSRFToken(w http.ResponseWriter, r *http.Request, token string) error {
	loginSession, _ := controller.store.Get(r, "eveauthLogin")

	loginSession.Values["csrfToken"] = token

	return sessions.Save(r, w)
}

// GetCSRFToken retrieves the CSRF token for the current session or sets a new one if the session was newly created
func (controller *Controller) GetCSRFToken(w http.ResponseWriter, r *http.Request) string {
	loginSession, _ := controller.store.Get(r, "eveauthLogin")

	if loginSession.IsNew {
		err := controller.SetCSRFToken(w, r, misc.GenerateRandomString(32))
		if err != nil {
			misc.Logger.Warnf("Failed to set CSRF token: [%v]", err)
			return ""
		}

		return loginSession.Values["csrfToken"].(string)
	}

	token, ok := loginSession.Values["csrfToken"].(string)
	if !ok {
		return ""
	}

	return token
}

// VerifyCSRFToken checks whether the provided CSRF token of the request and the stored session token match
func (controller *Controller) VerifyCSRFToken(w http.ResponseWriter, r *http.Request) bool {
	loginSession, _ := controller.store.Get(r, "eveauthLogin")

	if loginSession.IsNew {
		err := controller.SetCSRFToken(w, r, misc.GenerateRandomString(32))
		if err != nil {
			misc.Logger.Warnf("Failed to set CSRF token: [%v]", err)
			return false
		}

		return true
	}

	token, ok := loginSession.Values["csrfToken"].(string)
	if !ok {
		return false
	}

	err := r.ParseForm()
	if err != nil {
		return false
	}

	csrfToken := r.FormValue("csrfToken")

	return strings.EqualFold(token, csrfToken)
}

// IsLoggedIn checks whether the user is currently logged in and has an appropriate timestamp set
func (controller *Controller) IsLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	loginSession, _ := controller.store.Get(r, "eveauthLogin")

	if loginSession.IsNew {
		err := controller.SetCSRFToken(w, r, misc.GenerateRandomString(32))
		if err != nil {
			misc.Logger.Warnf("Failed to set CSRF token: [%v]", err)
			return false
		}

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

	return sessions.Save(r, w)
}

// GetLoginRedirect retrieves the previously set path for redirection after login
func (controller *Controller) GetLoginRedirect(w http.ResponseWriter, r *http.Request) string {
	loginSession, _ := controller.store.Get(r, "eveauthLogin")

	if loginSession.IsNew {
		err := controller.SetCSRFToken(w, r, misc.GenerateRandomString(32))
		if err != nil {
			misc.Logger.Warnf("Failed to set CSRF token: [%v]", err)
			return ""
		}

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

	logErr := controller.database.SaveLoginAttempt(loginAttempt)
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

	csrfToken := misc.GenerateRandomString(32)

	err = controller.SetCSRFToken(w, r, csrfToken)
	if err != nil {
		return err
	}

	loginSession, _ := controller.store.Get(r, "eveauthLogin")

	loginSession.Values["username"] = user.Username
	loginSession.Values["userID"] = user.ID
	loginSession.Values["timestamp"] = time.Now().Unix()
	loginSession.Values["verifiedEmail"] = user.VerifiedEmail

	loginSession.Options.MaxAge = 604800

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

// VerifyEmail checks the given code and verifies the presented email address is correct
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

// VerifyPasswordReset checks the given code and changes the user's password if the request is valid
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

	account := models.NewAccount(user.ID, keyID, apivCode, 0, true)

	apiClient := misc.CreateAPIClient(account)

	apiInfo, err := apiClient.Info()
	if err != nil {
		return err
	}

	account.APIAccessMask = apiInfo.AccessMask

	accountCharacters, err := apiClient.AccountCharacters()
	if err != nil {
		return nil
	}

	for _, accountChar := range accountCharacters {
		corporation, err := controller.database.LoadCorporationFromEVECorporationID(accountChar.CorporationID)
		if err == sql.ErrNoRows {
			misc.Logger.Tracef("No corporation with ID %d found, fetching corporation sheet...", accountChar.CorporationID)

			corporationSheet, err := apiClient.CorpCorporationSheet(accountChar.CorporationID)
			if err != nil {
				return err
			}

			corporation = models.NewCorporation(corporationSheet.Name, corporationSheet.Ticker, accountChar.CorporationID, corporationSheet.CEOID, zero.NewInt(0, false), zero.NewString("", false), true)

			corporation, err = controller.database.SaveCorporation(corporation)
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}

		character := models.NewCharacter(account.ID, corporation.ID, accountChar.Name, accountChar.ID, false, true)

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

// DeleteAPIKey removes the given API key from the user and database
func (controller *Controller) DeleteAPIKey(w http.ResponseWriter, r *http.Request, apiKeyID string) error {
	dataSession, _ := controller.store.Get(r, "eveauthData")

	user, ok := dataSession.Values["user"].(*models.User)
	if !ok {
		return fmt.Errorf("Failed to retrieve user from data session")
	}

	keyID, err := strconv.ParseInt(apiKeyID, 10, 64)
	if err != nil {
		return err
	}

	user, err = controller.database.RemoveAPIKeyFromUser(user, keyID)
	if err != nil {
		return err
	}

	user, err = controller.SetUser(w, r, user)
	if err != nil {
		return err
	}

	return nil
}

// GetUser returns the user-object stored in the data session
func (controller *Controller) GetUser(r *http.Request) (*models.User, error) {
	dataSession, _ := controller.store.Get(r, "eveauthData")

	user, ok := dataSession.Values["user"].(*models.User)
	if !ok {
		misc.Logger.Tracef("Failed to retrieve user from data session, checking database...")

		loginSession, _ := controller.store.Get(r, "eveauthLogin")

		userID, ok := loginSession.Values["userID"].(int64)
		if !ok {
			return nil, fmt.Errorf("Failed to retrieve user from data session and database")
		}

		var err error
		user, err = controller.database.LoadUser(userID)
		if err != nil {
			return nil, err
		}
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

// UpdateUser updates the current user's settings with the new given values
func (controller *Controller) UpdateUser(w http.ResponseWriter, r *http.Request, email string, oldPassword string, newPassword string) (*models.User, error) {
	user, err := controller.GetUser(r)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return nil, err
	}

	user.Email = email

	if len(newPassword) > 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		user.Password = string(hashedPassword)
	}

	user, err = controller.SetUser(w, r, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// AddGroupToUser adds the group with the given ID to the user
func (controller *Controller) AddGroupToUser(userID int64, groupID int64) error {
	user, err := controller.database.LoadUser(userID)
	if err != nil {
		return err
	}

	group, err := controller.database.LoadGroup(groupID)
	if err != nil {
		return err
	}

	user.Groups = append(user.Groups, group)

	_, err = controller.database.SaveUser(user)
	if err != nil {
		return err
	}

	return nil
}

// AddUserRoleToUser adds the role with the given ID to the user
func (controller *Controller) AddUserRoleToUser(userID int64, roleID int64, roleGranted bool) error {
	user, err := controller.database.LoadUser(userID)
	if err != nil {
		return err
	}

	role, err := controller.database.LoadRole(roleID)
	if err != nil {
		return err
	}

	userRole := models.NewUserRole(user.ID, role, false, roleGranted)

	user.UserRoles = append(user.UserRoles, userRole)

	_, err = controller.database.SaveUser(user)
	if err != nil {
		return err
	}

	return nil
}

// AddGroupRoleToGroup adds the role with the given ID to the group
func (controller *Controller) AddGroupRoleToGroup(groupID int64, roleID int64, roleGranted bool) error {
	group, err := controller.database.LoadGroup(groupID)
	if err != nil {
		return err
	}

	role, err := controller.database.LoadRole(roleID)
	if err != nil {
		return err
	}

	groupRole := models.NewGroupRole(group.ID, role, false, roleGranted)

	group.GroupRoles = append(group.GroupRoles, groupRole)

	_, err = controller.database.SaveGroup(group)
	if err != nil {
		return err
	}

	return nil
}

// CreateNewGroup creates a new group, saves it to the database and returns the updated model
func (controller *Controller) CreateNewGroup(groupName string) (*models.Group, error) {
	group := models.NewGroup(groupName, true)

	var err error
	group, err = controller.database.SaveGroup(group)
	if err != nil {
		return nil, err
	}

	return group, nil
}

// CreateNewRole creates a new role, saves it to the database and returns the updated model
func (controller *Controller) CreateNewRole(roleName string, roleLocked bool) (*models.Role, error) {
	role := models.NewRole(roleName, true, roleLocked)

	var err error
	role, err = controller.database.SaveRole(role)
	if err != nil {
		return nil, err
	}

	return role, nil
}

// GetUserAccounts returns the accounts associated with the current user
func (controller *Controller) GetUserAccounts(r *http.Request) ([]*models.Account, error) {
	user, err := controller.GetUser(r)
	if err != nil {
		return nil, err
	}

	return user.Accounts, nil
}

// GetUserCharacters returns the characters associated with all accounts of the current user
func (controller *Controller) GetUserCharacters(r *http.Request) ([]*models.Character, error) {
	user, err := controller.GetUser(r)
	if err != nil {
		return nil, err
	}

	var characters []*models.Character

	for _, account := range user.Accounts {
		characters = append(characters, account.Characters...)
	}

	return characters, err
}

// SetDefaultCharacter sets the default character for the current user
func (controller *Controller) SetDefaultCharacter(w http.ResponseWriter, r *http.Request, characterID string) error {
	charID, err := strconv.ParseInt(characterID, 10, 64)
	if err != nil {
		return err
	}

	user, err := controller.GetUser(r)
	if err != nil {
		return err
	}

	for _, account := range user.Accounts {
		for _, character := range account.Characters {
			if character.DefaultCharacter {
				character.DefaultCharacter = false
			}

			if character.ID == charID {
				character.DefaultCharacter = true
			}
		}
	}

	_, err = controller.SetUser(w, r, user)
	if err != nil {
		return err
	}

	return nil
}

// HasUserRole checks whether the current user has a role with the given name granted
func (controller *Controller) HasUserRole(r *http.Request, role string) bool {
	user, err := controller.GetUser(r)
	if err != nil {
		return false
	}

	roleStatus := user.HasRole(role)

	return roleStatus == models.RoleStatusGranted
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

	encryptedPayload, err := misc.EncryptAndAuthenticate(string(payload), application.Secret)
	if err != nil {
		return "", err
	}

	return encryptedPayload, nil
}

// LoadAllUsers retrieves all currently registered users
func (controller *Controller) LoadAllUsers() ([]*models.User, error) {
	users, err := controller.database.LoadAllUsers()
	if err != nil {
		return nil, err
	}

	return users, nil
}

// LoadUserFromUserID retrieves a user with the given user ID
func (controller *Controller) LoadUserFromUserID(userID int64) (*models.User, error) {
	user, err := controller.database.LoadUser(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// LoadAllGroups retrieves all currently existing groups
func (controller *Controller) LoadAllGroups() ([]*models.Group, error) {
	groups, err := controller.database.LoadAllGroups()
	if err != nil {
		return nil, err
	}

	return groups, nil
}

// LoadGroupFromGroupID retrieves a group with the given group ID
func (controller *Controller) LoadGroupFromGroupID(groupID int64) (*models.Group, error) {
	group, err := controller.database.LoadGroup(groupID)
	if err != nil {
		return nil, err
	}

	return group, nil
}

// LoadAvailableGroupsForUser retrieves all groups the user can be added to
func (controller *Controller) LoadAvailableGroupsForUser(userID int64) ([]*models.Group, error) {
	availableGroups, err := controller.database.LoadAvailableGroupsForUser(userID)
	if err != nil {
		return nil, err
	}

	return availableGroups, nil
}

// LoadAllRoles retrieves all currently existing roles
func (controller *Controller) LoadAllRoles() ([]*models.Role, error) {
	roles, err := controller.database.LoadAllRoles()
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// LoadAvailableUserRolesForUser retrieves all roles the user can be assigned
func (controller *Controller) LoadAvailableUserRolesForUser(userID int64) ([]*models.Role, error) {
	availableUserRoles, err := controller.database.LoadAvailableUserRolesForUser(userID)
	if err != nil {
		return nil, err
	}

	return availableUserRoles, nil
}

// LoadAvailableGroupRolesForGroup retrieves all roles the group can be assigned
func (controller *Controller) LoadAvailableGroupRolesForGroup(groupID int64) ([]*models.Role, error) {
	availableGroupRoles, err := controller.database.LoadAvailableGroupRolesForGroup(groupID)
	if err != nil {
		return nil, err
	}

	return availableGroupRoles, nil
}

// QueryCorporationName queries the database for the name of the corporation with the given ID
func (controller *Controller) QueryCorporationName(corporationID int64) (string, error) {
	corporationName, err := controller.database.LoadCorporationNameFromID(corporationID)
	if err != nil {
		return "", err
	}

	return corporationName, nil
}
