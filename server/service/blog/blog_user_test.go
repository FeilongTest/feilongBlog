package blog

import (
	"testing"

	"feilongBlog/utils"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordMatches(t *testing.T) {
	const password = "correct horse battery staple"
	if !passwordMatches(utils.MD5V([]byte(password)), password) {
		t.Fatal("legacy MD5 password should match")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	if !passwordMatches(string(hash), password) {
		t.Fatal("bcrypt password should match")
	}
	if passwordMatches(string(hash), "wrong") {
		t.Fatal("wrong password should not match")
	}
}
