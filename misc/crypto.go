package misc

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrInvalidCipherText describes invalid ciphertext
	ErrInvalidCipherText = errors.New("Invalid CipherText")
)

func createGCM(secret string) (cipher.AEAD, error) {
	// Convert the key to bytes
	key := []byte(secret)

	// Create a new AES cipher with the key
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create the GCM instance
	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	return gcm, nil
}

// EncryptAndAuthenticate encrypts and authenticates a message using AES-GCM
func EncryptAndAuthenticate(message, secret string) (string, error) {
	// Create a GCM instance
	gcm, err := createGCM(secret)
	if err != nil {
		return "", err
	}

	// Create a nonce from crypto random data
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	// Encrypt and authenticate
	plainText := []byte(message)
	cipherText := gcm.Seal(nil, nonce, plainText, nil)

	// Encode and return
	encodedCipherText := base64.URLEncoding.EncodeToString(cipherText)
	encodedNonce := base64.URLEncoding.EncodeToString(nonce)

	return fmt.Sprintf("%s~%s", encodedNonce, encodedCipherText), nil
}

// DecryptAndAuthenticate decrypts and authenticates a message using AES-GCM
func DecryptAndAuthenticate(encodedCipherText, secret string) (string, error) {
	// Decode the pieces
	pieces := strings.SplitN(encodedCipherText, "~", 2)
	if len(pieces) != 2 {
		return "", ErrInvalidCipherText
	}

	nonceLen := base64.URLEncoding.DecodedLen(len(pieces[0]))
	nonce := make([]byte, nonceLen)
	_, err := base64.URLEncoding.Decode(nonce, []byte(pieces[0]))
	if err != nil {
		return "", err
	}

	cipherTextLen := base64.URLEncoding.DecodedLen(len(pieces[1]))
	cipherText := make([]byte, cipherTextLen)
	_, err = base64.URLEncoding.Decode(cipherText, []byte(pieces[1]))
	if err != nil {
		return "", err
	}

	// Create a GCM instance
	gcm, err := createGCM(secret)
	if err != nil {
		return "", err
	}

	// Attempt to open the encrypted data
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
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
