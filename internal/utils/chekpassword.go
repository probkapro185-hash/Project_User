package utils

import "golang.org/x/crypto/bcrypt"

func CheckPassword(passhash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passhash), []byte(password))
	return err == nil
}
