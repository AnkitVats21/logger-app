package storage

import (
	"logger-app/db"
	"logger-app/models"
	"time"
)

func InsertEvent(event models.Event) error {
	query := `INSERT INTO events(place, timestamp) VALUES(?, ?)`

	_, err := db.DB.Exec(query, event.Place, event.Timestamp)
	return err
}

func GetLatestEvent(now time.Time) (*models.Event, error) {
	row := db.DB.QueryRow("SELECT id, place, timestamp FROM events WHERE timestamp <= ? ORDER BY timestamp DESC LIMIT 1", now)
	var e models.Event
	err := row.Scan(&e.ID, &e.Place, &e.Timestamp)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func GetAllEvents() ([]models.Event, error) {
	rows, err := db.DB.Query("SELECT id, place, timestamp FROM events ORDER BY timestamp DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event

	for rows.Next() {
		var e models.Event
		err := rows.Scan(&e.ID, &e.Place, &e.Timestamp)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func GetEventsPaginated(offset, limit int) ([]models.Event, error) {
	rows, err := db.DB.Query(
		"SELECT id, place, timestamp FROM events ORDER BY timestamp DESC LIMIT ? OFFSET ?",
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event

	for rows.Next() {
		var e models.Event
		err := rows.Scan(&e.ID, &e.Place, &e.Timestamp)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func GetEventsInRange(start, end string) ([]models.Event, error) {
	query := `
	SELECT id, place, timestamp
	FROM events
	WHERE DATE(timestamp) BETWEEN ? AND ?
	ORDER BY timestamp ASC
	`

	rows, err := db.DB.Query(query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event

	for rows.Next() {
		var e models.Event
		err := rows.Scan(&e.ID, &e.Place, &e.Timestamp)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func GetEventsForToday(now time.Time) ([]models.Event, error) {
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Format("2006-01-02 15:04:05")
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location()).Format("2006-01-02 15:04:05")

	query := `
	SELECT id, place, timestamp
	FROM events
	WHERE timestamp BETWEEN ? AND ?
	ORDER BY timestamp ASC
	`

	rows, err := db.DB.Query(query, startOfDay, endOfDay)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var e models.Event
		err := rows.Scan(&e.ID, &e.Place, &e.Timestamp)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func GetTotalEventsCount() (int, error) {
	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM events").Scan(&count)
	return count, err
}
