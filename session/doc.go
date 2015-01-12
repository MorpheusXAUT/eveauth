// Package session provides functionality to manage user sessions and temporary storages.
// The web handlers will mainly interact with this package to retrieve information.
// If the required information is not found within the session storage, the database backend is being queried.
package session
