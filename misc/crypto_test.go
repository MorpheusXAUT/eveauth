package misc

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCryptoEncryptDecrypt(t *testing.T) {
	Convey("Trying to encrypt and decrypt a message", t, func() {
		message := "hello world"
		secret := "press_f_to_pay_respects_12345678"

		Convey("Encrypting the test message", func() {
			encrypted, err := EncryptAndAuthenticate(message, secret)

			Convey("The returned error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Decrypting the encrypted message", func() {
				decrypted, err := DecryptAndAuthenticate(encrypted, secret)

				Convey("The returned error should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The returned messages should match", func() {
					So(decrypted, ShouldEqual, message)
				})
			})
		})
	})
}
