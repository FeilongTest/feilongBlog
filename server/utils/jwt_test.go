package utils

import (
	"testing"
	"time"

	model "feilongBlog/model/blog"
	jwt "github.com/golang-jwt/jwt/v4"
)

func TestJWTRejectsInvalidSignaturesAlgorithmsAndClaims(t *testing.T) {
	key := []byte("test-signing-key-not-for-a-deployed-site")
	verifier := &JWT{SigningKey: key}
	for _, test := range []struct {
		name               string
		method             jwt.SigningMethod
		key                interface{}
		expires, notBefore int64
		valid              bool
	}{
		{"valid administrator", jwt.SigningMethodHS256, key, time.Now().Add(time.Hour).Unix(), 0, true},
		{"wrong key", jwt.SigningMethodHS256, []byte("another-test-key"), time.Now().Add(time.Hour).Unix(), 0, false},
		{"wrong HMAC algorithm", jwt.SigningMethodHS384, key, time.Now().Add(time.Hour).Unix(), 0, false},
		{"unsigned token", jwt.SigningMethodNone, jwt.UnsafeAllowNoneSignatureType, time.Now().Add(time.Hour).Unix(), 0, false},
		{"expired", jwt.SigningMethodHS256, key, time.Now().Add(-time.Hour).Unix(), 0, false},
		{"expired and wrong signature", jwt.SigningMethodHS256, []byte("another-test-key"), time.Now().Add(-time.Hour).Unix(), 0, false},
		{"not active yet", jwt.SigningMethodHS256, key, time.Now().Add(2 * time.Hour).Unix(), time.Now().Add(time.Hour).Unix(), false},
	} {
		t.Run(test.name, func(t *testing.T) {
			claims := model.CustomClaims{BaseClaims: model.BaseClaims{ID: 1, Admin: 1}, StandardClaims: jwt.StandardClaims{ExpiresAt: test.expires, NotBefore: test.notBefore}}
			signed, err := jwt.NewWithClaims(test.method, claims).SignedString(test.key)
			if err != nil {
				t.Fatal(err)
			}
			parsed, err := verifier.ParseToken(signed)
			if test.valid {
				if err != nil || parsed == nil || parsed.ID != 1 || parsed.Admin != 1 {
					t.Fatalf("valid token rejected: %v", err)
				}
			} else if err == nil || parsed != nil {
				t.Fatal("invalid token exposed trusted claims")
			}
		})
	}
	if parsed, err := verifier.ParseToken("invalid.token"); err == nil || parsed != nil {
		t.Fatal("malformed token exposed trusted claims")
	}
}
