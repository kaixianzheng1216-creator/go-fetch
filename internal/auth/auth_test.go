package auth

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestCheckPassword(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}

	if !CheckPassword(string(hash), "secret") {
		t.Fatal("expected password to match")
	}
	if CheckPassword(string(hash), "wrong") {
		t.Fatal("expected password mismatch")
	}
}
