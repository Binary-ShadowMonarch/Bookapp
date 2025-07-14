// internal/models/user.go
package models

type User struct {
	ID             int
	Email          string
	HashedPassword string
	Provider       string // 'local', 'google', 'github'
	ProviderID     string // e.g. OAuth subject or userID
}
