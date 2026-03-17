package models

import "time"

type Event struct {
	ID        int
	UserID    string
	Place     string
	Timestamp time.Time
}
