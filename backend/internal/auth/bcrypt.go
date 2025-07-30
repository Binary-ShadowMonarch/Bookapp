package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword creates a secure hash of a password using bcrypt
// I use this when someone registers to store their password safely
// bcrypt automatically adds salt and is designed to be slow (which is good for security)
func HashPassword(pw string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword compares a plain text password with a hashed password
// I use this during login to check if the password is correct
// returns true if the password matches, false otherwise
func CheckPassword(hashed, pw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pw))
	return err == nil
}
