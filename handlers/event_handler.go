package handlers

import (
	"log"
	"logger-app/models"
	"logger-app/service"
	"logger-app/storage"
	"net/http"
	"strconv"
)

const pageSize = 50

func LogEventFromQuery(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userid")
	place := r.URL.Query().Get("place")

	if userID == "" {
		http.Error(w, "userid missing", http.StatusBadRequest)
		return
	}

	if place == "" {
		http.Error(w, "place missing", http.StatusBadRequest)
		return
	}

	exists, _ := storage.UserExists(userID)
	if !exists {
		http.Error(w, "invalid userid", http.StatusUnauthorized)
		return
	}

	err := service.LogEvent(userID, place)
	if err != nil {
		log.Printf("LogEvent Error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Logged: " + place))
}

func GetEvents(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userid")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	if userID == "" {
		if cookie, err := r.Cookie("userid"); err == nil {
			userID = cookie.Value
		}
	}

	if userID == "" {
		service.RenderTemplate(w, "events.html", struct {
			Events      []models.Event
			CurrentPage int
			TotalPages  int
			TotalEvents int
			Limit       int
			HasMore     bool
		}{Limit: pageSize, CurrentPage: 1})
		return
	}

	exists, _ := storage.UserExists(userID)
	if !exists {
		// If ID is invalid, we still render the page but the checkUser() on frontend 
		// will detect the 401 from subsequent API calls or we can clear it here.
		// For HTML pages, better to let the frontend modal handle it via API validation.
	}

	page := 1
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	limit := 50
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	offset := (page - 1) * limit

	events, err := storage.GetEventsPaginated(userID, offset, limit)
	if err != nil {
		log.Printf("GetEvents Error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalEvents, err := storage.GetTotalEventsCount(userID)
	if err != nil {
		log.Printf("GetTotalEventsCount Error: %v", err)
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
		ActivePage  string
		UserID      string
	}{
		Events:      events,
		CurrentPage: page,
		TotalPages:  totalPages,
		Limit:       limit,
		TotalEvents: totalEvents,
		HasMore:     page < totalPages,
		ActivePage:  "events",
		UserID:      userID,
	}

	service.RenderTemplate(w, "events.html", data)
}
