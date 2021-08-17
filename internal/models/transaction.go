package models

import (
	"time"

	"github.com/ehabterra/flash_api/internal/constants"
)

type Transaction struct {
	ID            string                    `json:"id"`
	UserID        string                    `json:"user_id"`
	Type          constants.TransactionType `json:"type"`
	Amount        float64                   `json:"amount"`
	Date          time.Time                 `json:"date"`
	AccountNumber *string                   `json:"account_number"`
	RecipientID   string                    `json:"recipient_id"`
}
