package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"smartassistant/internal/config"
	sdb "smartassistant/internal/db"
	"smartassistant/internal/handlers"
	"smartassistant/internal/models"
)

func main() {
	cfg := config.Load()

	pool, err := sdb.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	defer pool.Close()

	// ensure migrations (best-effort apply of 001)
	if b, err := os.ReadFile("migrations/001_init.sql"); err == nil {
		if _, err := pool.Exec(context.Background(), string(b)); err != nil {
			log.Printf("warning: applying migration failed: %v", err)
		}
	}

	renderer, err := handlers.NewRenderer()
	if err != nil {
		log.Fatalf("templates error: %v", err)
	}

	qaModel := &models.QAModel{DB: pool}
	qa := &handlers.QARoutes{Renderer: renderer, Model: qaModel}
	ask := &handlers.AskRoutes{Renderer: renderer, FlaskURL: cfg.FlaskURL}

	mux := http.NewServeMux()

	// static
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// UI
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		renderer.Render(w, "layout", map[string]any{
			"ContentTemplate": "content_home",
			"Year":            time.Now().Year(),
		})
	})
	mux.HandleFunc("/ask", ask.AskPage)
	mux.HandleFunc("/qa", qa.List)
	mux.HandleFunc("/qa/create", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			qa.CreateForm(w, r)
		case http.MethodPost:
			qa.Create(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/qa/edit", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			qa.EditForm(w, r)
		case http.MethodPost:
			qa.Edit(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/qa/delete", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			qa.DeleteConfirm(w, r)
		case http.MethodPost:
			qa.Delete(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// API
	mux.HandleFunc("/api/ask", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		ask.AskAPI(w, r)
	})

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           logging(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
	log.Printf("Go server listening on :%s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		dur := time.Since(start)
		q := r.URL.RawQuery
		if q != "" {
			q = "?" + q
		}
		log.Printf("%s %s%s %s", r.Method, r.URL.Path, q, dur)
	})
}

// small helpers for templates if needed
func atoi(s string, def int) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}

// end
