package model

import "time"

var (
	DefaultRedisExpiration time.Duration = 5 * time.Minute
	TokenTypeBearer                      = "Bearer"
)

type Header struct {
	CacheControl  string `reqHeader:"Cache-Control"`
	Authorization string `reqHeader:"Authorization"`
	APIKey        string `reqHeader:"API-Key"`
}

type BaseInformation struct {
	CreatedBy int64     `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedBy int64     `json:"updated_by"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedBy int64     `json:"deleted_by"`
	DeletedAt time.Time `json:"deleted_at"`
}

type Pagination struct {
	CurrentPage     int64  `json:"current_page"`
	CurrentElements int64  `json:"current_elements"`
	TotalPages      int64  `json:"total_pages"`
	TotalElements   int64  `json:"total_elements"`
	SortBy          string `json:"sort_by"`
}
