package model

import (
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	Coins     int       `json:"coins"`
	CreatedAt time.Time `json:"created_at"`
}

type MerchItem struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type UserInventory struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	MerchID   uint      `json:"merch_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}

type CoinTransaction struct {
	ID        uint      `json:"id"`
	FromID    uint      `json:"from_id"`
	ToID      uint      `json:"to_id"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

var MerchCatalog = map[string]int{
	"t-shirt":    80,
	"cup":        20,
	"book":       50,
	"pen":        10,
	"powerbank":  200,
	"hoody":      300,
	"umbrella":   200,
	"socks":      10,
	"wallet":     50,
	"pink-hoody": 500,
}
