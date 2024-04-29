package db

import "time"

type SERP struct {
	URL         string
	Title       string
	Description string
	Location    string
	ContactInfo []string
	Keywords    []string
	IsRead      bool
	CreatedAt   time.Time
}
