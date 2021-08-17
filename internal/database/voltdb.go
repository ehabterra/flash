package database

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ehabterra/flash_api/internal/constants"

	"github.com/google/uuid"

	"github.com/VoltDB/voltdb-client-go/voltdbclient"
	"github.com/ehabterra/flash_api/internal/models"
)

type VoltDB struct {
	db *voltdbclient.Conn
}

func NewVoltDB(datasource string) *VoltDB {
	conn, err := voltdbclient.OpenConn(datasource)
	if err != nil {

	}

	return &VoltDB{db: conn}
}

func (v *VoltDB) Connect(account *models.Account) error {
	result, err := v.db.Exec("@AdHoc",
		[]driver.Value{
			`insert into USER_BANK_ACCOUNTS(
ACCOUNT_NUMBER,USER_ID,BANK_ID,BRANCH_NUMBER,HOLDER_NAME,REFERENCE) 
values(?,?,?,?,?,?);`,
			account.AccountNumber,
			account.UserID,
			account.BankID,
			account.BranchNumber,
			account.HolderName,
			account.Reference,
		})
	return executionErrorHandling(err, result)
}

func (v *VoltDB) AddAccountTransaction(transaction *models.AccountTransaction) error {
	id := uuid.New().String()

	result, err := v.db.Exec("@AdHoc",
		[]driver.Value{
			`insert into transactions(
ID, USER_ID,Date,Type,Amount) 
values(?,?,?,?,?);`,
			id,
			transaction.UserID,
			fmt.Sprint(time.Now().UnixNano() / int64(time.Microsecond)),
			fmt.Sprint(transaction.Type),
			fmt.Sprint(transaction.Amount),
		})

	result, err = v.db.Exec("@AdHoc",
		[]driver.Value{
			`insert into transaction_accounts(
transaction_id, ACCOUNT_NUMBER) 
values(?,?);`,
			id,
			transaction.AccountNumber,
		})

	return executionErrorHandling(err, result)
}

func (v *VoltDB) AddRecipientTransaction(transaction *models.RecipientTransaction) error {
	id := uuid.New().String()

	result, err := v.db.Exec("@AdHoc",
		[]driver.Value{
			`insert into transactions(
ID, USER_ID,Date,Type,Amount) 
values(?,?,?,?,?);`,
			id,
			transaction.UserID,
			fmt.Sprint(time.Now().UnixNano() / int64(time.Microsecond)),
			fmt.Sprint(transaction.Type),
			fmt.Sprint(transaction.Amount),
		})

	result, err = v.db.Exec("@AdHoc",
		[]driver.Value{
			`insert into transaction_recipients(
transaction_id, recipient_id) 
values(?,?);`,
			id,
			transaction.RecipientID,
		})
	return executionErrorHandling(err, result)
}

func (v *VoltDB) UpdateBalance(id string, amount int64) error {
	result, err := v.db.Exec("@AdHoc",
		[]driver.Value{
			`update users set balance = balance + ? where id = ?`,
			fmt.Sprint(amount),
			id,
		})
	return executionErrorHandling(err, result)
}

func (v *VoltDB) GetUserByUsernameOrEmail(usernameOrEmail string) (*models.User, error) {
	rows, err := v.db.Query("@AdHoc",
		[]driver.Value{
			`select ID,Username,Email,Password, Balance  From USERS where username = ? or email = ?;`,
			usernameOrEmail,
			usernameOrEmail,
		})
	if err != nil {
		return nil, err
	}

	return getUser(rows), nil
}

func (v *VoltDB) GetBalance(id string) (float64, error) {
	rows, err := v.db.Query("@AdHoc",
		[]driver.Value{
			`select Balance  From USERS where id = ?;`,
			id,
		})
	if err != nil {
		return 0, err
	}

	return printBalance(rows), nil
}

func (v *VoltDB) CheckAccountNumber(accountNumber string) (bool, error) {
	rows, err := v.db.Query("@AdHoc",
		[]driver.Value{
			`select account_number  From USER_BANK_ACCOUNTS where account_number = ?;`,
			accountNumber,
		})
	if err != nil {
		return false, err
	}

	return len(strings.TrimSpace(printAccountNumber(rows))) > 0, nil
}

func (v *VoltDB) CheckTransactionLimits(id string, transactionType constants.TransactionType, amount int64, limits map[time.Duration]int64) (bool, error) {
	if len(limits) == 0 {
		return true, nil
	}

	for k, value := range limits {
		t := time.Now().Add(-k)

		statement := "select count(*) transaction_count From Transactions where user_id = ? and type = ? and date >= ? group by user_id, type having sum(amount) > ?;"

		rows, err := v.db.Query("@AdHoc",
			[]driver.Value{
				statement,
				id,
				fmt.Sprint(transactionType),
				fmt.Sprint(t.UnixNano() / int64(time.Microsecond)),
				fmt.Sprint(value - amount),
			})
		if err != nil {
			return false, err
		}

		if printCount(rows) > 0 {
			return true, nil
		}
	}
	return false, nil
}

func printBalance(rows driver.Rows) float64 {
	voltRows := rows.(voltdbclient.VoltRows)
	if voltRows.AdvanceRow() {
		return getFloat64(voltRows, "Balance")
	}
	return 0
}

func printAccountNumber(rows driver.Rows) string {
	voltRows := rows.(voltdbclient.VoltRows)
	if voltRows.AdvanceRow() {
		return getString(voltRows, "account_number")
	}
	return ""
}

func printCount(rows driver.Rows) int64 {
	voltRows := rows.(voltdbclient.VoltRows)
	if voltRows.AdvanceRow() {
		return getInt64(voltRows, "transaction_count")
	}
	return 0
}

func getUser(rows driver.Rows) *models.User {
	voltRows := rows.(voltdbclient.VoltRows)
	if voltRows.AdvanceRow() {
		return &models.User{
			ID:       getString(voltRows, "ID"),
			Username: getString(voltRows, "Username"),
			Email:    getString(voltRows, "Email"),
			Password: getString(voltRows, "Password"),
			Balance:  getFloat64(voltRows, "Balance"),
		}
	}
	return nil
}

func getString(voltRows voltdbclient.VoltRows, col string) string {
	result, err := voltRows.GetStringByName(col)
	if err != nil {
		log.Println(err)
	}
	return result.(string)
}

func getFloat64(voltRows voltdbclient.VoltRows, col string) float64 {
	result, err := voltRows.GetDecimalByName(col)
	if err != nil {
		log.Println(err)
	}
	if result != nil {
		f, _ := result.(*big.Float).Float64()
		return f
	}
	return 0
}

func getInt64(voltRows voltdbclient.VoltRows, col string) int64 {
	result, err := voltRows.GetBigIntByName(col)
	if err != nil {
		log.Println(err)
	}
	if result != nil {
		n := result.(int64)
		return n
	}
	return 0
}

func executionErrorHandling(err error, result driver.Result) error {
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// if not record affected then return error
	if affected == 0 {
		return errors.New("no record inserted")
	}
	return nil
}
