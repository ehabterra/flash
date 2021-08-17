package models

import (
	"github.com/ehabterra/flash_api/internal/constants"
)

type AccountTransaction struct {
	UserID        string                    `json:"user_id"`
	Type          constants.TransactionType `json:"type"`
	Amount        float64                   `json:"amount"`
	AccountNumber string                    `json:"account_number"`
}
