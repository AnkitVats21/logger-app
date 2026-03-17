package models

import "time"

type Event struct {
	ID        int
	Place     string
	Timestamp time.Time
}
