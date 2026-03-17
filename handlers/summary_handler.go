package handlers

import (
	"encoding/json"
	"log"
	"logger-app/models"
	"logger-app/service"
	"logger-app/storage"
	"net/http"
	"strconv"
	"time"
)

func GetRangeSummary(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userid")
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	if userID == "" {
		http.Error(w, "userid missing", http.StatusBadRequest)
		return
	}

	if start == "" || end == "" {
		http.Error(w, "missing start or end", http.StatusBadRequest)
		return
	}

	exists, _ := storage.UserExists(userID)
	if !exists {
		http.Error(w, "invalid userid", http.StatusUnauthorized)
		return
	}

	events, err := storage.GetEventsInRange(userID, start, end)
	if err != nil {
		log.Printf("GetEventsInRange Error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := service.CalculateRangeSummary(events)

	page := 1
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	limit := 50 // Default limit
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	totalDays := len(result.Days)
	totalPages := (totalDays + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}

	startIdx := (page - 1) * limit
	endIdx := startIdx + limit
	if startIdx > totalDays {
		startIdx = totalDays
	}
	if endIdx > totalDays {
		endIdx = totalDays
	}

	paginatedDays := []models.DailySummary{} // Default to prevent null
	if len(result.Days) > 0 {
		paginatedDays = result.Days[startIdx:endIdx]
	}

	response := map[string]any{
		"total_home_time":    service.FormatDuration(result.TotalHome),
		"total_office_time":  service.FormatDuration(result.TotalOffice),
		"total_outside_time": service.FormatDuration(result.TotalOutside),
		"total_commute_time": service.FormatDuration(result.TotalCommute),
		"days":               paginatedDays,
		"current_page":       page,
		"total_pages":        totalPages,
		"total_items":        totalDays,
		"limit":              limit,
	}

	json.NewEncoder(w).Encode(response)
}

func GetTodaySummary(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userid")
	if userID == "" {
		http.Error(w, "userid missing", http.StatusBadRequest)
		return
	}

	exists, _ := storage.UserExists(userID)
	if !exists {
		http.Error(w, "invalid userid", http.StatusUnauthorized)
		return
	}

	// For accurate simulation of elapsed time up to the exact minute
	now := time.Now()

	events, err := storage.GetEventsForToday(userID, now)
	if err != nil {
		log.Printf("GetEventsForToday Error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Midnight State Inheritance:
	// If the day starts without an event at 00:00:00, find the state before midnight
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if len(events) == 0 || events[0].Timestamp.After(startOfToday) {
		lastEvent, err := storage.GetLatestEvent(userID, startOfToday)
		if err == nil && lastEvent != nil {
			// Synthetic event at midnight
			synthetic := models.Event{
				UserID:    userID,
				Place:     lastEvent.Place,
				Timestamp: startOfToday,
			}
			events = append([]models.Event{synthetic}, events...)
		}
	}

	result := service.CalculateTodaySummary(events, now)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
