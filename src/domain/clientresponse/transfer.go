package clientresponse

type Transfer struct {
	ID                     string `json:"id"`
	Amount                 int    `json:"amount"`
	Status                 string `json:"status"`
	TransactionDate        string `json:"transaction_date"`
	SourceBankAccount      string `json:"source_bank_account"`
	DestinationBankAccount string `json:"destination_bank_account"`
	SourceBankID           int    `json:"source_bank_id"`
	DestinationBankID      int    `json:"destination_bank_id"`
}
