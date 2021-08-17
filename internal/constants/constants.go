package constants

const JWTSecretKey = "flashSecret"

type TransactionType int

const (
	UploadTransactionType TransactionType = iota + 1
	SendTransactionType
)
