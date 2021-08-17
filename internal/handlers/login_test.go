package handlers

import (
	"log"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestGenerateHash(t *testing.T) {
	pass := "pass"
	password, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(string(password))

}
