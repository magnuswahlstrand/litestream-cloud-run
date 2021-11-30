package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func newRouter(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("request received", r.Method)
		//// Store page view.
		if _, err := db.Exec(`INSERT INTO page_views (timestamp) VALUES (?);`, time.Now().Format(time.RFC3339)); err != nil {
			log.Println("insert", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Read total page views.
		var n int
		if err := db.QueryRow(`SELECT COUNT(1) FROM page_views;`).Scan(&n); err != nil {
			log.Println("count", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Print total page views.
		fmt.Fprintf(w, "This server has been visited %d times.\n", n)
	})

	return r
}

func main() {
	log.Printf("start!")

	if err := run(); err != nil {
		log.Printf("error: %v", err)
		fmt.Fprintf(os.Stderr, "run:", err)
		os.Exit(1)
	}
}

func configureDatabase(db *sql.DB) error {
	// Create table for storing page views.
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS page_views (id INTEGER PRIMARY KEY, timestamp TEXT);`); err != nil {
		return err
	}

	if _, err := db.Exec(`PRAGMA busy_timeout = 5000;`); err != nil {
		return fmt.Errorf("configure busy_timeout failed: %s", err)
	}
	return nil
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer stop()

	dsn := os.Getenv("DB_PATH")
	if dsn == "" {
		return fmt.Errorf("required: DB_PATH")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Open database file.
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return fmt.Errorf("sql.Open: %v", err)
	}
	defer db.Close()

	if err := configureDatabase(db); err != nil {
		return fmt.Errorf("failed to configure database: %s", err)
	}

	// Run web server.
	log.Printf("listening on :%s\n", port)
	go http.ListenAndServe(":"+port, newRouter(db))

	<-ctx.Done()
	log.Print("received shutdown signal, shutting down")

	return nil
}
