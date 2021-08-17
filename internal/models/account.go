package models

type Account struct {
	AccountNumber string `json:"account_number"`
	UserID        string `json:"user_id"`
	BankID        string `json:"bank_id"`
	BranchNumber  string `json:"branch_number"`
	HolderName    string `json:"holder_name"`
	Reference     string `json:"reference"`
}
