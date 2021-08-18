package services

import (
	"net/http"
	"time"

	"github.com/ehabterra/flash_api/internal/constants"

	"github.com/go-openapi/errors"

	"github.com/ehabterra/flash_api/internal/models"
)

type Bank interface {
	Connect(account *models.Account) error
	Upload(number string, amount int64) error
}

type Database interface {
	Connect(account *models.Account) error
	GetUserByUsernameOrEmail(usernameOrEmail string) (*models.User, error)
	AddAccountTransaction(transaction *models.AccountTransaction) error
	AddRecipientTransaction(transaction *models.RecipientTransaction) error
	UpdateBalance(id string, amount int64) error
	GetBalance(id string) (float64, error)
	CheckTransactionLimits(id string, transactionType constants.TransactionType, amount int64, limits map[time.Duration]int64) (bool, error)
	CheckAccountNumber(accountNumber string) (bool, error)
}

type Users struct {
	bank Bank
	db   Database
}

func NewUsers(bank Bank, db Database) *Users {
	return &Users{bank, db}
}

func (u Users) Connect(account *models.Account) error {
	if err := u.bank.Connect(account); err != nil {
		return errors.New(http.StatusBadRequest, "error happened while trying to connect to bank")
	}
	return u.db.Connect(account)

}

func (u Users) GetUserByUsernameOrEmail(usernameOrEmail string) (*models.User, error) {
	return u.db.GetUserByUsernameOrEmail(usernameOrEmail)
}

func (u Users) GetUserBalance(id string) (float64, error) {
	return u.db.GetBalance(id)
}

func (u Users) Send(id string, usernameOrEmail string, amount int64) error {
	// check sent money limits
	err := u.checkLimits(id, constants.SendTransactionType, amount)
	if err != nil {
		return err
	}

	// check balance
	balance, err := u.db.GetBalance(id)
	if err != nil {
		return err
	}

	if int64(balance) < amount {
		return errors.New(http.StatusInternalServerError, "balance not sufficient")
	}

	// get recipient by username or email
	recipient, err := u.db.GetUserByUsernameOrEmail(usernameOrEmail)
	if err != nil {
		return err
	}

	err = u.sendTransaction(id, recipient.ID, amount)
	if err != nil {
		return err
	}

	return nil
}

func (u Users) Upload(id string, number string, amount int64) error {
	exists, err := u.db.CheckAccountNumber(number)
	if err != nil || !exists {
		return err
	}

	err = u.checkLimits(id, constants.UploadTransactionType, amount)
	if err != nil {
		return err
	}

	if err := u.bank.Upload(number, amount); err != nil {
		return errors.New(http.StatusBadRequest, "error happened while trying to connect to bank")
	}

	err = u.uploadTransaction(id, number, amount)
	if err != nil {
		// TODO: queue the process to repeat again or refund the amount
		return err
	}

	return nil
}

func (u Users) uploadTransaction(id string, number string, amount int64) error {
	err := u.db.UpdateBalance(id, amount)
	if err != nil {
		return err
	}

	err = u.db.AddAccountTransaction(&models.AccountTransaction{
		UserID:        id,
		Type:          constants.UploadTransactionType,
		Amount:        float64(amount),
		AccountNumber: number,
	})
	if err != nil {
		return err
	}
	return nil
}

func (u Users) sendTransaction(id string, recipientID string, amount int64) error {
	// update sender balance
	err := u.db.UpdateBalance(id, -amount)
	if err != nil {
		return err
	}

	// update recipient balance
	err = u.db.UpdateBalance(recipientID, amount)
	if err != nil {
		return err
	}

	err = u.db.AddRecipientTransaction(&models.RecipientTransaction{
		UserID:      id,
		Type:        constants.SendTransactionType,
		Amount:      float64(amount),
		RecipientID: recipientID,
	})
	if err != nil {
		return err
	}
	return nil
}

func (u Users) checkLimits(id string, t constants.TransactionType, amount int64) error {
	limits, err := u.db.CheckTransactionLimits(
		id,
		t,
		amount,
		map[time.Duration]int64{
			time.Hour * 24:     10000,
			time.Hour * 24 * 7: 50000,
		},
	)
	if err != nil {
		return err
	}
	if limits {
		return errors.New(http.StatusInternalServerError, "transaction exceeded the daily or weekly limits")
	}
	return nil
}
