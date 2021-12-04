package entity

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	UserID         int    `json:"user_id"`
	Username       int    `json:"username"`
	HashedPassword string `json:"hash_password"`
	AvatarURL      string `json:"avatar_url"`
}

type Transaction struct {
	TransactionID int       `json:"transaction_id,omitempty"`
	WalletID      int       `json:"wallet_id" validate:"required"`
	CategoryID    int       `json:"category_id" validate:"required"`
	Amount        int       `json:"amount" validate:"required,gt=0"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	Note          string    `json:"note"`
}

type Category struct {
	CategoryID int    `json:"category_id"`
	Name       string `json:"name"`
	IsExpense  bool   `json:"is_expense"`
	IconName   string `json:"icon_name"`
}

type Wallet struct {
	WalletID   int `json:"wallet_id"`
	CurrencyID int `json:"currency_id"`
	UserID     int `json:"user_id"`
}

type Currency struct {
	CurrencyID int    `json:"currency_id"`
	Code       string `json:"code"`
	Name       string `json:"name"`
}

type Session struct {
	SessionID string `json:"session_id"`
	UserID    int    `json:"user_id"`
}

type CustomValidator struct {
	validator *validator.Validate
}
