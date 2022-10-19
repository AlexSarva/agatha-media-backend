package models

import (
	"time"

	"github.com/google/uuid"
)

type Status struct {
	Result string `json:"result"`
}

type User struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Username string    `json:"username" db:"username"`
	Email    string    `json:"email" db:"email"`
	Password string    `json:"password" db:"passwd"`
	Token    string    `json:"token" db:"token"`
	TokenExp time.Time `json:"token_expires" db:"token_expires"`
}

type UserLogin struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"passwd"`
}

type Token struct {
	Username string    `json:"username" db:"username"`
	Email    string    `json:"email" db:"email"`
	Type     string    `json:"type"`
	Token    string    `json:"token" db:"token"`
	TokenExp time.Time `json:"token_expires" db:"token_expires"`
}

type UserInfo struct {
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
}
