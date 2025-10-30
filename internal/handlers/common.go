package handlers

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

type TemplateRenderer struct {
	templates *template.Template
}

func NewRenderer() (*TemplateRenderer, error) {
	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"max": func(a, b int) int {
			if a > b {
				return a
			}
			return b
		},
		"divCeil": func(a, b int) int {
			if b <= 0 {
				return 1
			}
			if a == 0 {
				return 1
			}
			q := a / b
			if a%b != 0 {
				q++
			}
			if q < 1 {
				q = 1
			}
			return q
		},
		"defaultYear": func(val any) int {
			switch y := val.(type) {
			case int:
				if y > 0 {
					return y
				}
			case int64:
				if y > 0 {
					return int(y)
				}
			case float64:
				if y > 0 {
					return int(y)
				}
			}
			return time.Now().Year()
		},
	}
	t, err := template.New("").Funcs(funcMap).ParseGlob("internal/templates/*.gohtml")
	if err != nil {
		return nil, err
	}
	return &TemplateRenderer{templates: t}, nil
}

func (r *TemplateRenderer) Render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := r.templates.ExecuteTemplate(w, name, data); err != nil {
		log.Printf("template render error: %v", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
	}
}
