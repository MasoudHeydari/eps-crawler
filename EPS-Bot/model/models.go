package model

import (
	"time"
)

type SearchQueryRequest struct {
	Query    string `json:"q"`
	Location string `json:"loc"`
	Language string `json:"lang"`
}

type SERPResponse struct {
	URL         string
	Title       string
	Description string
	ContactInfo ContactInfo
	Keywords    []string
	IsRead      bool
	CreatedAt   time.Time
}

type SearchQueryResponse struct {
	Id         int
	Query      string
	Language   string
	Location   string
	IsCanceled bool
	CreatedAt  time.Time
}

type ContactInfo struct {
	Emails []string `json:"emails"`
	Phones []string `json:"phones"`
}
