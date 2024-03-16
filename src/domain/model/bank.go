package model

var (
	GetBankAccountKey string = "gbAccountKey:%s"
)

type GetBankAccount struct {
	AccountNumber string `json:"account_number" query:"account_number"`
	AccountName   string `json:"account_name" query:"account_name"`
	BankID        int64  `json:"bank_id" query:"bank_id"`
}
