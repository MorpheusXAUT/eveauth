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

func TestCryptoHMACSHA256(t *testing.T) {
	Convey("Trying to calculate and verify the HMAC-SHA256 of a message", t, func() {
		message := "hello world"
		secret := "press_f_to_pay_respects_12345678"
		expected := "9koYaQtlVXY3p2SIZ896-F5E2NBE57Yuz5FfKBLbt6k="

		Convey("Calculating the HMAC-SHA256", func() {
			calculated := CalculateMessageHMACSHA256(message, secret)

			Convey("The calculated result should match the expected", func() {
				So(calculated, ShouldEqual, expected)
			})

			Convey("Verifying the calculated result", func() {
				verified := VerifyMessageHMACSHA256(message, expected, secret)

				Convey("The verification should indicate a match", func() {
					So(verified, ShouldBeTrue)
				})
			})
		})
	})
}
