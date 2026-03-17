package service

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"time"
)

// TemplateFS is set by main at startup so RenderTemplate uses embedded files.
var TemplateFS fs.FS

func FormatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	return fmt.Sprintf("%dh %dm", h, m)
}

func EndOfDay(date string) time.Time {
	t, _ := time.Parse("2006-01-02", date)
	return t.Add(24 * time.Hour)
}

var funcMap = template.FuncMap{
	"add": func(a, b int) int { return a + b },
	"sub": func(a, b int) int { return a - b },
}

// RenderTemplate renders a page using the shared layout and nav, sourcing files
// from TemplateFS (embedded) if set, or falling back to disk.
func RenderTemplate(w http.ResponseWriter, page string, data interface{}) {
	patterns := []string{
		"templates/layout.html",
		"templates/components/nav.html",
		"templates/pages/" + page,
	}

	var tmpl *template.Template
	var err error

	if TemplateFS != nil {
		tmpl, err = template.New("layout").Funcs(funcMap).ParseFS(TemplateFS, patterns...)
	} else {
		tmpl, err = template.New("layout").Funcs(funcMap).ParseFiles(patterns...)
	}

	if err != nil {
		http.Error(w, "Template Parse Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, "Template Execution Error: "+err.Error(), http.StatusInternalServerError)
	}
}
