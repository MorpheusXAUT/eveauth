package misc

import (
	"math/rand"

	"github.com/morpheusxaut/eveauth/models"

	//"github.com/nixwaro/eveapi"
	"github.com/morpheusxaut/eveapi"
)

// GenerateRandomString returns a random alphanumerical string with the given length
func GenerateRandomString(length int) string {
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, length)

	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
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
