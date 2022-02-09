package encryption

import (
	"golang.org/x/crypto/bcrypt"
)

func Encrypt(s string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

func Compare(s string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(s))
	if err != nil {
		return false
	}
	return true
}
