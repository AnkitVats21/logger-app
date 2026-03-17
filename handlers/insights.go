package handlers

import (
	"logger-app/service"
	"net/http"
)

func InsightsPage(w http.ResponseWriter, r *http.Request) {
	// Simple render, the JS will fetch the data for the timeline/table.
	service.RenderTemplate(w, "insights.html", map[string]interface{}{
		"ActivePage": "insights",
	})
}
