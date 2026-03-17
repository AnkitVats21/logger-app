package main

import (
	"log"
	"net/http"

	"logger-app/db"
	"logger-app/handlers"
	"logger-app/service"
	"time"
)

func init() {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Printf("Failed to load Asia/Kolkata: %v. Falling back to UTC.", err)
	} else {
		time.Local = loc
	}
}

func main() {
	service.TemplateFS = TemplateFS
	db.InitDB()
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	})
	mux.HandleFunc("/events", handlers.GetEvents)
	mux.HandleFunc("/log", handlers.LogEventFromQuery)
	mux.HandleFunc("/summary/range", handlers.GetRangeSummary)
	mux.HandleFunc("/summary/today", handlers.GetTodaySummary)
	mux.HandleFunc("/dashboard", handlers.DashboardPage)
	mux.HandleFunc("/insights", handlers.InsightsPage)
	loggedMux := handlers.LoggingMiddleware(mux)
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", loggedMux))
}
