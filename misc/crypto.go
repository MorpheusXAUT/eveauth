package misc

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

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

// EncryptAESCFB encrypts a given string using AES-CFB and the given 32 bytes key
func EncryptAESCFB(message string, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	b := base64.URLEncoding.EncodeToString([]byte(message))

	ciphertext := make([]byte, aes.BlockSize+len(b))

	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)

	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// DecryptAESCFB decrypts a given (encrypted) string using AES-CFB and the given 32 bytes key
func DecryptAESCFB(encrypted string, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	text, err := base64.URLEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	if len(text) < aes.BlockSize {
		return "", errors.New("Ciphertext too short")
	}

	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)

	cfb.XORKeyStream(text, text)

	data, err := base64.URLEncoding.DecodeString(string(text))
	if err != nil {
		return "", err
	}
	return string(data), nil
}
