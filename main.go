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

	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		icon, err := StaticFS.ReadFile("static/favicon.ico")
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "image/x-icon")
		w.Write(icon)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	})
	mux.HandleFunc("/events", handlers.GetEvents)
	mux.HandleFunc("/log", handlers.LogEventFromQuery)
	mux.HandleFunc("/api/range-summary", handlers.GetRangeSummary)
	mux.HandleFunc("/api/today-summary", handlers.GetTodaySummary)
	mux.HandleFunc("/dashboard", handlers.DashboardPage)
	mux.HandleFunc("/insights", handlers.InsightsPage)
	loggedMux := handlers.LoggingMiddleware(mux)
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", loggedMux))
}
