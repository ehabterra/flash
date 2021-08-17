package models

import (
	"github.com/ehabterra/flash_api/internal/constants"
)

type RecipientTransaction struct {
	UserID      string                    `json:"user_id"`
	Type        constants.TransactionType `json:"type"`
	Amount      float64                   `json:"amount"`
	RecipientID string                    `json:"recipient_id"`
}
