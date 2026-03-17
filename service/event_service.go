package service

import (
	"database/sql"
	"errors"
	"fmt"
	"logger-app/models"
	"logger-app/storage"
	"time"
)

func LogEvent(place string) error {
	latest, err := storage.GetLatestEvent(time.Now())
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to fetch latest event: %w", err)
	}

	if latest != nil {
		// Enforce state machine rules
		if place == latest.Place {
			return fmt.Errorf("Invalid Transition: already at %s", latest.Place)
		}

		if (latest.Place == "home" && place == "office") ||
			(latest.Place == "office" && place == "home") {
			return errors.New("Invalid Transition: must go 'outside' before transitioning between Home and Office")
		}
	}

	event := models.Event{
		Place:     place,
		Timestamp: time.Now(),
	}

	return storage.InsertEvent(event)
}
