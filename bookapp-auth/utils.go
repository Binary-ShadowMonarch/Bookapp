// utils.go
package main

import "golang.org/x/crypto/bcrypt"

// hashPassword returns the bcrypt hash of the password at cost 10.
func hashPassword(pw string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(bytes), err
}

// checkPassword compares a bcrypt hashed password with its possible plaintext equivalent.
func checkPassword(hashed, pw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pw))
	return err == nil
}
