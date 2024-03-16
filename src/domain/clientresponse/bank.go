package clientresponse

type BankAccount struct {
	ID            string  `json:"id"`
	BankID        int     `json:"bank_id"`
	AccountNumber string  `json:"account_number"`
	AccountName   string  `json:"account_name"`
	AccountAmount float64 `json:"account_amount"`
}
