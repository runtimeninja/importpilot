package domain

import "time"

type User struct {
	ID           int64
	ClientID     *int64
	Email        string
	PasswordHash string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
