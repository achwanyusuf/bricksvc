package response

import "time"

type TransactionInfo struct {
	RequestURI    string    `json:"request_uri"`
	RequestMethod string    `json:"request_method"`
	RequestID     string    `json:"request_id"`
	Timestamp     time.Time `json:"timestamp"`
	ErrorCode     int64     `json:"error_code,omitempty"`
	Cause         string    `json:"cause,omitempty"`
}

type Response struct {
	TransactionInfo TransactionInfo `json:"transaction_info"`
	Code            int64           `json:"status_code"`
	Message         string          `json:"message,omitempty"`
	Translation     *Translation    `json:"translation,omitempty"`
}

type Translation struct {
	EN string `json:"en"`
}
