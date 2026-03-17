package handlers

import (
	"encoding/json"
	"logger-app/models"
	"logger-app/service"
	"logger-app/storage"
	"net/http"
	"strconv"
	"time"
)

func GetRangeSummary(w http.ResponseWriter, r *http.Request) {
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	if start == "" || end == "" {
		http.Error(w, "missing start or end", http.StatusBadRequest)
		return
	}

	events, err := storage.GetEventsInRange(start, end)
	if err != nil {
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
	// For accurate simulation of elapsed time up to the exact minute
	now := time.Now()

	events, err := storage.GetEventsForToday(now)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Midnight State Inheritance:
	// If the day starts without an event at 00:00:00, find the state before midnight
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if len(events) == 0 || events[0].Timestamp.After(startOfToday) {
		lastEvent, err := storage.GetLatestEvent(startOfToday)
		if err == nil && lastEvent != nil {
			// Synthetic event at midnight
			synthetic := models.Event{
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
