package models

// User holds credential + session data.
// Email is the unique key.
type User struct {
	Email          string
	HashedPassword string
}
