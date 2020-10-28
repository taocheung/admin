package model

import (
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestAddUser(t *testing.T) {
	password, _ := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)
	t.Log(string(password))
}
