package data

import (
	"database/sql"
	"time"

	"nathejk.dk/internal/validator"
)

type Payment struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Version   uint64    `json:"version"`

	Timestamp       time.Time `json:"ts"`
	ShopNumber      string    `json:"shopNumber"`
	Amount          int       `json:"amount"`
	Currency        string    `json:"currency"`
	Message         string    `json:"message"`
	UserPhoneNumber string    `json:"userPhoneNumber"`
	UserName        string    `json:"userName"`
}

func (p *Payment) Validate(v validator.Validator) {
	v.Check(p.Timestamp.IsZero(), "timestamp", "must be provided")
}

type PaymentModel struct {
	DB *sql.DB
}

func (m PaymentModel) Insert(p *Payment) error {
	query := `
		INSERT INTO payment (ts, shopNumber, amount, currency, message, userPhoneNumber, userName)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, version`
	args := []any{p.Timestamp, p.ShopNumber, p.Amount, p.Currency, p.Message, p.UserPhoneNumber, p.UserName}
	return m.DB.QueryRow(query, args...).Scan(&p.ID, &p.CreatedAt, &p.Version)
}

func (m PaymentModel) Latest() time.Time {
	query := `SELECT ts FROM payment ORDER BY id DESC LIMIT 1`

	var t string
	if err := m.DB.QueryRow(query).Scan(&t); err != nil {
		return time.Time{}
	}
	ts, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return time.Time{}
	}
	return ts
}
