package domain

import "errors"

var (
	ErrClientNotFound      = errors.New("client not found")
	ErrClientEmailExists   = errors.New("client email already exists")
	ErrInvalidClientID     = errors.New("invalid client id")
	ErrInvalidClientStatus = errors.New("status must be either active or inactive")
)
