package storage

import (
	"logger-app/db"
	"logger-app/models"
	"time"
)

func InsertEvent(event models.Event) error {
	query := `INSERT INTO events(user_id, place, timestamp) VALUES(?, ?, ?)`

	_, err := db.DB.Exec(query, event.UserID, event.Place, event.Timestamp)
	return err
}

func GetLatestEvent(userID string, now time.Time) (*models.Event, error) {
	row := db.DB.QueryRow("SELECT id, user_id, place, timestamp FROM events WHERE user_id = ? AND timestamp <= ? ORDER BY timestamp DESC LIMIT 1", userID, now)
	var e models.Event
	err := row.Scan(&e.ID, &e.UserID, &e.Place, &e.Timestamp)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func GetAllEvents(userID string) ([]models.Event, error) {
	rows, err := db.DB.Query("SELECT id, user_id, place, timestamp FROM events WHERE user_id = ? ORDER BY timestamp DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event

	for rows.Next() {
		var e models.Event
		err := rows.Scan(&e.ID, &e.UserID, &e.Place, &e.Timestamp)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func GetEventsPaginated(userID string, offset, limit int) ([]models.Event, error) {
	rows, err := db.DB.Query(
		"SELECT id, user_id, place, timestamp FROM events WHERE user_id = ? ORDER BY timestamp DESC LIMIT ? OFFSET ?",
		userID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event

	for rows.Next() {
		var e models.Event
		err := rows.Scan(&e.ID, &e.UserID, &e.Place, &e.Timestamp)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func GetEventsInRange(userID string, start, end string) ([]models.Event, error) {
	query := `
	SELECT id, user_id, place, timestamp
	FROM events
	WHERE user_id = ? AND DATE(timestamp) BETWEEN ? AND ?
	ORDER BY timestamp ASC
	`

	rows, err := db.DB.Query(query, userID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event

	for rows.Next() {
		var e models.Event
		err := rows.Scan(&e.ID, &e.UserID, &e.Place, &e.Timestamp)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func GetEventsForToday(userID string, now time.Time) ([]models.Event, error) {
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Format("2006-01-02 15:04:05")
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location()).Format("2006-01-02 15:04:05")

	query := `
	SELECT id, user_id, place, timestamp
	FROM events
	WHERE user_id = ? AND timestamp BETWEEN ? AND ?
	ORDER BY timestamp ASC
	`

	rows, err := db.DB.Query(query, userID, startOfDay, endOfDay)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var e models.Event
		err := rows.Scan(&e.ID, &e.UserID, &e.Place, &e.Timestamp)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func GetTotalEventsCount(userID string) (int, error) {
	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM events WHERE user_id = ?", userID).Scan(&count)
	return count, err
}
