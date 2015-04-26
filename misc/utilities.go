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
