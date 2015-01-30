package misc

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"strconv"

	"github.com/morpheusxaut/eveauth/models"

	"github.com/nixwaro/eveapi"
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
		ID:    strconv.Itoa(int(account.APIKeyID)),
		VCode: account.APIvCode,
	})

	return api
}

// CalculateMessageHMACSHA256 calculates the HMAC of a message using the SHA256 algorithm and the given secret
func CalculateMessageHMACSHA256(message string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))

	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

// VerifyMessageHMACSHA256 verifies the received HMAC of a message using the SHA256 algorithm and the given secret
func VerifyMessageHMACSHA256(message string, calculated string, secret string) bool {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))

	expectedHMAC := h.Sum(nil)

	calculatedHMAC, err := base64.URLEncoding.DecodeString(calculated)
	if err != nil {
		Logger.Errorf("Failed to decode HMAC string: [%v]", err)
		return false
	}

	return hmac.Equal(calculatedHMAC, expectedHMAC)
}
