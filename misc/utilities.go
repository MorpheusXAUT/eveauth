package misc

import (
	"crypto/rand"
	"math/big"

	"github.com/morpheusxaut/eveauth/models"

	"github.com/morpheusxaut/eveapi"
)

// GenerateRandomString returns a random alphanumerical string with the given length
func GenerateRandomString(length int) string {
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, length)

	for i := 0; i < len(b); i++ {
		r, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			i--
			continue
		}

		b[i] = chars[r.Int64()]
	}

	return string(b)
}

// CreateAPIClient creates an EVE Online API client used for querying API data
func CreateAPIClient(account *models.Account) *eveapi.API {
	api := eveapi.Simple(eveapi.Key{
		ID:    account.APIKeyID,
		VCode: account.APIvCode,
	})

	return api
}

// AuthStatus represents the result of an authentication attempt by the user
type AuthStatus int

const (
	// AuthStatusUnknown indicates an unknown status during the authentication
	AuthStatusUnknown AuthStatus = iota
	// AuthStatusError indicates an (internal) error during the authentication
	AuthStatusError
	// AuthStatusCredentialMismatch indicates a mismatch of the given passwords and stored bcrypt hash
	AuthStatusCredentialMismatch
	// AuthStatusUnverifiedEmail indicates the user's email address has not been verified yet
	AuthStatusUnverifiedEmail
	// AuthStatusSuccess indicates a successful authentication attempt
	AuthStatusSuccess
)

// String returns an easily readable (string) representation of the AuthStatus
func (authStatus AuthStatus) String() string {
	switch authStatus {
	case AuthStatusUnknown:
		return "unknown"
	case AuthStatusError:
		return "error"
	case AuthStatusCredentialMismatch:
		return "mismatch"
	case AuthStatusUnverifiedEmail:
		return "unverified"
	case AuthStatusSuccess:
		return "success"
	}

	return "unknown"
}
