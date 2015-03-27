package misc

import "testing"

func TestEncryptDecrypt(t *testing.T) {
	message := "hello world"
	secret := "press_f_to_pay_respects_12345678"

	encrypted, err := EncryptAndAuthenticate(message, secret)
	if err != nil {
		t.Fatalf("Error encrypting: %s", err.Error())
	}

	decrypted, err := DecryptAndAuthenticate(encrypted, secret)
	if err != nil {
		t.Fatalf("Error decrypting: %s", err.Error())
	}

	if decrypted != message {
		t.Fatalf("Error decrypted: \"%s\" != \"%s\"", decrypted, message)
	}
}
