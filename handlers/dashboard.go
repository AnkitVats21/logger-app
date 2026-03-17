package handlers

import (
	"logger-app/service"
	"logger-app/storage"
	"net/http"
)

func DashboardPage(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userid")
	if userID == "" {
		if cookie, err := r.Cookie("userid"); err == nil {
			userID = cookie.Value
		}
	}

	userName := ""
	if userID != "" {
		if u, err := storage.GetUserByID(userID); err == nil {
			userName = u.Name
		}
	}

	service.RenderTemplate(w, "dashboard.html", map[string]interface{}{
		"ActivePage": "dashboard",
		"UserName":   userName,
	})
}
