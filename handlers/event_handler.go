package handlers

import (
	"logger-app/service"
	"logger-app/storage"
	"net/http"
	"strconv"
)

const pageSize = 50

func LogEventFromQuery(w http.ResponseWriter, r *http.Request) {
	place := r.URL.Query().Get("place")

	if place == "" {
		http.Error(w, "place missing", http.StatusBadRequest)
		return
	}

	err := service.LogEvent(place)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Logged: " + place))
}

func GetEvents(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	limit := 50
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	offset := (page - 1) * limit

	events, err := storage.GetEventsPaginated(offset, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalEvents, err := storage.GetTotalEventsCount()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := (totalEvents + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}

	data := struct {
		Events      interface{}
		CurrentPage int
		TotalPages  int
		Limit       int
		TotalEvents int
		HasMore     bool
	}{
		Events:      events,
		CurrentPage: page,
		TotalPages:  totalPages,
		Limit:       limit,
		TotalEvents: totalEvents,
		HasMore:     page < totalPages,
	}

	service.RenderTemplate(w, "events.html", data)
}
