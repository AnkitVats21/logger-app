package handlers

import (
	"logger-app/service"
	"net/http"
)

func DashboardPage(w http.ResponseWriter, r *http.Request) {
	service.RenderTemplate(w, "dashboard.html", nil)
}
