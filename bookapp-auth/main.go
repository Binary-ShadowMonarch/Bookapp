// main.go
package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Login struct {
	HashedPassword string
	SessionToken   string
	CSRFToken      string
}

// Key is mail
var users = map[string]Login{}

// sessionToken → mail
var sessions = map[string]string{}

func main() {
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/protected", protected)

	fmt.Println("Listening on :8080…")
	http.ListenAndServe(":8080", nil)
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	mail := r.FormValue("mail")
	password := r.FormValue("password")
	if mail == "" || len(password) < 8 {
		http.Error(w, "invalid mail or password too short", http.StatusBadRequest)
		return
	}
	if _, exists := users[mail]; exists {
		http.Error(w, "user already exists", http.StatusConflict)
		return
	}
	hashed, err := hashPassword(password)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	users[mail] = Login{HashedPassword: hashed}
	fmt.Fprintln(w, "registered")
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	mail := r.FormValue("mail")
	password := r.FormValue("password")

	user, ok := users[mail]
	if !ok || !checkPassword(user.HashedPassword, password) {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// create session token
	sessToken := randomToken(32)
	sessions[sessToken] = mail

	// create CSRF token
	csrfToken := randomToken(32)

	// update user record
	u := users[mail]
	u.SessionToken = sessToken
	u.CSRFToken = csrfToken
	users[mail] = u

	// set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessToken,
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		HttpOnly: false,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	fmt.Fprintf(w, "logged in; your CSRF token is: %s\n", csrfToken)
}

func logout(w http.ResponseWriter, r *http.Request) {
	if err := Authorize(r); err != nil {
		er := http.StatusUnauthorized
		http.Error(w, "Unauthorized", er)
		return
	}

	// clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: false,
	})

	// clear tokens from db
	mail := r.FormValue("mail")
	user, _ := users[mail]
	user.SessionToken = ""
	user.CSRFToken = ""
	users[mail] = user

	fmt.Fprintln(w, "logged out")
}

func protected(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		er := http.StatusMethodNotAllowed
		http.Error(w, "Invalid request method", er)
		return
	}

	if err := Authorize(r); err != nil {
		er := http.StatusUnauthorized
		http.Error(w, "Unauthorized", er)
		return
	}

	mail := r.FormValue("mail")
	fmt.Fprintln(w, "CSRF validation succesful "+mail)
}

func randomToken(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)
}
