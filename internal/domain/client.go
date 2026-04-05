package domain

import "time"

type Client struct {
	ID        int64
	Name      string
	Email     string
	ShopURL   string
	Status    string
	Plan      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
