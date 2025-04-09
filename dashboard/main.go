package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"redblue-sim/models"
	"redblue-sim/shared"
)

func main() {
	db, err := shared.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/api/alerts", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w, r)
		if r.Method == http.MethodOptions {
			return
		}
		alerts := getLatestAlerts(db)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(alerts)
	})

	http.HandleFunc("/api/logs", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w, r)
		if r.Method == http.MethodOptions {
			return
		}
		logs := getLatestLogs(db)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(logs)
	})

	log.Println("Dashboard API server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// CORS for cross-origin requests from frontend
func enableCORS(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*") // Or restrict to http://localhost:3000
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		(*w).WriteHeader(http.StatusOK)
	}
}

func getLatestAlerts(db *sql.DB) []models.Alert {
	var alerts []models.Alert
	rows, err := db.Query(`SELECT id, timestamp, severity, message, source_ip, status FROM alerts ORDER BY timestamp DESC`)
	if err != nil {
		log.Println("Query alerts error:", err)
		return alerts
	}
	defer rows.Close()

	for rows.Next() {
		var a models.Alert
		if err := rows.Scan(&a.ID, &a.Timestamp, &a.Severity, &a.Message, &a.SourceIP, &a.Status); err != nil {
			log.Println("Scan alert error:", err)
			continue
		}
		alerts = append(alerts, a)
	}
	return alerts
}

func getLatestLogs(db *sql.DB) []models.Log {
	var logs []models.Log
	rows, err := db.Query(`SELECT id, timestamp, source_ip, action, description FROM logs ORDER BY timestamp DESC`)
	if err != nil {
		log.Println("Query logs error:", err)
		return logs
	}
	defer rows.Close()

	for rows.Next() {
		var l models.Log
		if err := rows.Scan(&l.ID, &l.Timestamp, &l.SourceIP, &l.Action, &l.Description); err != nil {
			log.Println("Scan log error:", err)
			continue
		}
		logs = append(logs, l)
	}
	return logs
}
