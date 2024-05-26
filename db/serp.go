package db

import (
	"github.com/karust/openserp/core"
	"time"
)

type SERP struct {
	URL         string
	Title       string
	Description string
	ContactInfo core.ContactInfo
	Keywords    []string
	IsRead      bool
	CreatedAt   time.Time
}

type SearchQuery struct {
	Id         int
	Query      string
	Language   string
	Location   string
	IsCanceled bool
	CreatedAt  time.Time
}
