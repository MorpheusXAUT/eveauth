package web

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/morpheusxaut/eveauth/database"
	"github.com/morpheusxaut/eveauth/misc"
	"github.com/morpheusxaut/eveauth/models"
	"github.com/morpheusxaut/eveauth/session"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
)

// Controller provides functionality for handling web requests and accessing session and backend data
type Controller struct {
	Config    *misc.Configuration
	Database  database.Connection
	Session   *session.Controller
	Templates *Templates
	Checksums *AssetChecksums
	RedisPool *redis.Pool

	router *mux.Router
}

// SetupController prepares the web controller and initialises the router and handled routes
func SetupController(config *misc.Configuration, db database.Connection, sessions *session.Controller, templates *Templates, checksums *AssetChecksums) *Controller {
	controller := &Controller{
		Config:    config,
		Database:  db,
		Session:   sessions,
		Templates: templates,
		Checksums: checksums,
		router:    mux.NewRouter().StrictSlash(true),
	}

	controller.RedisPool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.RedisHost)
			if err != nil {
				return nil, err
			}

			if len(config.RedisPassword) > 0 {
				_, err := c.Do("AUTH", config.RedisPassword)
				if err != nil {
					c.Close()
					return nil, err
				}
			}

			if len(config.RedisDB) > 0 {
				_, err = c.Do("SELECT", config.RedisDB)
				if err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	routes := SetupRoutes(controller)

	for _, route := range routes {
		controller.router.Methods(route.Methods...).Path(route.Pattern).Name(route.Name).Handler(controller.ServeHTTP(route.HandlerFunc, route.Name))
	}

	controller.router.PathPrefix("/").Handler(http.FileServer(http.Dir("app/assets")))

	return controller
}

// ServeHTTP acts as a middleware between parsed requests, logging the requests and replacing the remote address with the proxy-value if needed
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

		if (r.Method == "POST" || r.Method == "PUT") && !controller.Session.VerifyCSRFToken(w, r) {
			misc.Logger.Warnf("Failed to verify CSRF token")

			var userID int64

			user, err := controller.Session.GetUser(r)
			if err != nil {
				misc.Logger.Warnf("Failed to get user from session: [%v]", err)
				userID = -1
			} else {
				userID = user.ID
			}

			csrfFailure := models.NewCSRFFailure(userID, r)

			err = controller.Database.SaveCSRFFailure(csrfFailure)
			if err != nil {
				misc.Logger.Errorf("Failed to save CSRF failure: [%v]", err)
			}

			response := make(map[string]interface{})
			response["pageType"] = 1
			response["pageTitle"] = "Error"
			response["status"] = 1
			response["result"] = fmt.Errorf("An error occurred, please try again!")

			controller.SendResponse(w, r, "index", response)
		} else {
			inner.ServeHTTP(w, r)
		}

		misc.Logger.Debugf("ServeHTTP: [%s] %s %q {%s} - %s ", r.Method, r.RemoteAddr, r.RequestURI, name, time.Since(start))
	})
}

// HandleRequests starts the blocking call to handle web requests
func (controller *Controller) HandleRequests() {
	misc.Logger.Infof("Listening for HTTP requests on %q...", controller.Config.HTTPHost)

	http.Handle("/", controller.router)
	err := http.ListenAndServe(controller.Config.HTTPHost, nil)

	misc.Logger.Criticalf("Received error while listening for HTTP requests: [%v]", err)
}

// SetAuthorizationToken stores a temporary authorization token for the given user and app
func (controller *Controller) SetAuthorizationToken(userID int64, appID int64, token string) error {
	c := controller.RedisPool.Get()
	defer c.Close()

	err := c.Send("SET", fmt.Sprintf("authorization_token_%d_%d", appID, userID), token)
	if err != nil {
		return err
	}

	err = c.Send("EXPIRE", fmt.Sprintf("authorization_token_%d_%d", appID, userID), 300)
	if err != nil {
		return err
	}

	return nil
}

// GetAuthorizationToken tries to retrieve a temporary authorization for the given user and app
func (controller *Controller) GetAuthorizationToken(userID int64, appID int64) (string, error) {
	c := controller.RedisPool.Get()
	defer c.Close()

	token, err := redis.String(c.Do("GET", fmt.Sprintf("authorization_token_%d_%d", appID, userID)))
	if err != nil {
		return "", err
	}

	return token, nil
}

// EncryptUserPermissions retrieves the data for the given user and app and encrypted the user's permissions using the app secret
func (controller *Controller) EncryptUserPermissions(userID int64, appID int64) (string, error) {
	user, err := controller.Database.LoadUser(userID)
	if err != nil {
		return "", err
	}

	authUser := user.ToAuthUser()

	application, err := controller.Database.LoadApplication(appID)
	if err != nil {
		return "", err
	}

	payload, err := json.Marshal(authUser)
	if err != nil {
		return "", err
	}

	encryptedPayload, err := misc.EncryptAndAuthenticate(string(payload), application.Secret)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString([]byte(encryptedPayload)), nil
}

// AddGroupToUser adds the group with the given ID to the user
func (controller *Controller) AddGroupToUser(userID int64, groupID int64) error {
	user, err := controller.Database.LoadUser(userID)
	if err != nil {
		return err
	}

	group, err := controller.Database.LoadGroup(groupID)
	if err != nil {
		return err
	}

	user.Groups = append(user.Groups, group)

	_, err = controller.Database.SaveUser(user)
	if err != nil {
		return err
	}

	return nil
}

// AddUserRoleToUser adds the role with the given ID to the user
func (controller *Controller) AddUserRoleToUser(userID int64, roleID int64, roleGranted bool) error {
	user, err := controller.Database.LoadUser(userID)
	if err != nil {
		return err
	}

	role, err := controller.Database.LoadRole(roleID)
	if err != nil {
		return err
	}

	userRole := models.NewUserRole(user.ID, role, false, roleGranted)

	user.UserRoles = append(user.UserRoles, userRole)

	_, err = controller.Database.SaveUser(user)
	if err != nil {
		return err
	}

	return nil
}

// AddGroupRoleToGroup adds the role with the given ID to the group
func (controller *Controller) AddGroupRoleToGroup(groupID int64, roleID int64, roleGranted bool) error {
	group, err := controller.Database.LoadGroup(groupID)
	if err != nil {
		return err
	}

	role, err := controller.Database.LoadRole(roleID)
	if err != nil {
		return err
	}

	groupRole := models.NewGroupRole(group.ID, role, false, roleGranted)

	group.GroupRoles = append(group.GroupRoles, groupRole)

	_, err = controller.Database.SaveGroup(group)
	if err != nil {
		return err
	}

	return nil
}

// CreateNewGroup creates a new group, saves it to the database and returns the updated model
func (controller *Controller) CreateNewGroup(groupName string) (*models.Group, error) {
	group := models.NewGroup(groupName, true)

	var err error
	group, err = controller.Database.SaveGroup(group)
	if err != nil {
		return nil, err
	}

	return group, nil
}

// CreateNewRole creates a new role, saves it to the database and returns the updated model
func (controller *Controller) CreateNewRole(roleName string, roleLocked bool) (*models.Role, error) {
	role := models.NewRole(roleName, true, roleLocked)

	var err error
	role, err = controller.Database.SaveRole(role)
	if err != nil {
		return nil, err
	}

	return role, nil
}

// VerifyApplication verifies the application to be authorized to perform requests to the auth backend
func (controller *Controller) VerifyApplication(appID string, callback string, auth string) (*models.Application, error) {
	applicationID, err := strconv.ParseInt(appID, 10, 64)
	if err != nil {
		return nil, err
	}

	application, err := controller.Database.LoadApplication(applicationID)
	if err != nil {
		return nil, err
	}

	verified := misc.VerifyMessageHMACSHA256(fmt.Sprintf("%d:%s", application.ID, application.Callback), auth, application.Secret)

	if !verified {
		return nil, fmt.Errorf("Failed to verify HMAC")
	}

	return application, nil
}

// LoadAllUsers retrieves all currently registered users
func (controller *Controller) LoadAllUsers() ([]*models.User, error) {
	users, err := controller.Database.LoadAllUsers()
	if err != nil {
		return nil, err
	}

	return users, nil
}

// LoadUserFromUserID retrieves a user with the given user ID
func (controller *Controller) LoadUserFromUserID(userID int64) (*models.User, error) {
	user, err := controller.Database.LoadUser(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// LoadAllGroups retrieves all currently existing groups
func (controller *Controller) LoadAllGroups() ([]*models.Group, error) {
	groups, err := controller.Database.LoadAllGroups()
	if err != nil {
		return nil, err
	}

	return groups, nil
}

// LoadGroupFromGroupID retrieves a group with the given group ID
func (controller *Controller) LoadGroupFromGroupID(groupID int64) (*models.Group, error) {
	group, err := controller.Database.LoadGroup(groupID)
	if err != nil {
		return nil, err
	}

	return group, nil
}

// LoadAvailableGroupsForUser retrieves all groups the user can be added to
func (controller *Controller) LoadAvailableGroupsForUser(userID int64) ([]*models.Group, error) {
	availableGroups, err := controller.Database.LoadAvailableGroupsForUser(userID)
	if err != nil {
		return nil, err
	}

	return availableGroups, nil
}

// LoadAllRoles retrieves all currently existing roles
func (controller *Controller) LoadAllRoles() ([]*models.Role, error) {
	roles, err := controller.Database.LoadAllRoles()
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// LoadAvailableUserRolesForUser retrieves all roles the user can be assigned
func (controller *Controller) LoadAvailableUserRolesForUser(userID int64) ([]*models.Role, error) {
	availableUserRoles, err := controller.Database.LoadAvailableUserRolesForUser(userID)
	if err != nil {
		return nil, err
	}

	return availableUserRoles, nil
}

// LoadAvailableGroupRolesForGroup retrieves all roles the group can be assigned
func (controller *Controller) LoadAvailableGroupRolesForGroup(groupID int64) ([]*models.Role, error) {
	availableGroupRoles, err := controller.Database.LoadAvailableGroupRolesForGroup(groupID)
	if err != nil {
		return nil, err
	}

	return availableGroupRoles, nil
}

// QueryCorporationName queries the database for the name of the corporation with the given ID
func (controller *Controller) QueryCorporationName(corporationID int64) (string, error) {
	corporationName, err := controller.Database.LoadCorporationNameFromID(corporationID)
	if err != nil {
		return "", err
	}

	return corporationName, nil
}
