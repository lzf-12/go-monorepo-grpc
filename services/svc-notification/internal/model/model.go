package model

import (
	"time"
)

type EmailLog struct {
	ID        string    `json:"id" db:"id"`
	To        string    `json:"to" db:"to_email"`
	From      string    `json:"from" db:"from_email"`
	Subject   string    `json:"subject" db:"subject"`
	Body      string    `json:"body" db:"body"`
	IsHTML    bool      `json:"is_html" db:"is_html"`
	Status    string    `json:"status" db:"status"`
	Error     *string   `json:"error" db:"error"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}