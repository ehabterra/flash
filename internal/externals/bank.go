package externals

import "github.com/ehabterra/flash_api/internal/models"

type Bank struct {
}

func NewBank() *Bank {
	return &Bank{}
}

func (b *Bank) Connect(account *models.Account) error {
	return nil
}
func (b *Bank) Upload(number string, amount int64) error {
	return nil
}
