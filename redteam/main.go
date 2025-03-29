package main

import (
	"log"
	"math/rand"
	"net/http"
	"redblue-sim/shared"
	"strconv"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	db, err := shared.InitDB()
	if err != nil {
		log.Fatal("Database initialization failed:", err)
	}
	defer db.Close()

	http.HandleFunc("/attack", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w, r)
		if r.Method == http.MethodOptions {
			return
		}

		sourceIP := "192.168.1." + strconv.Itoa(rand.Intn(255))

		_, err := db.Exec(`
			INSERT INTO alerts (severity, message, source_ip) 
			VALUES (?, ?, ?)`,
			"high",
			"Malicious payload detected",
			sourceIP,
		)

		if err != nil {
			http.Error(w, "Failed to insert alert", http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Attack logged from " + sourceIP))
	})

	log.Println("Red Team service running on port 8082...")
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func enableCORS(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		(*w).WriteHeader(http.StatusOK)
	}
}
