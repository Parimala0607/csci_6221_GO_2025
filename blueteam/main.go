package main

import (
	"log"
	"net"
	"net/http"
	"redblue-sim/shared"
)

func main() {
	db, err := shared.InitDB()
	if err != nil {
		log.Fatal("Database initialization failed:", err)
	}
	defer db.Close()

	http.HandleFunc("/defend", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w, r)
		if r.Method == http.MethodOptions {
			return
		}

		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Invalid remote address", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec(`
			INSERT INTO logs (source_ip, action, description) 
			VALUES (?, ?, ?)`,
			host,
			"firewall_block",
			"Blocked suspicious incoming connection",
		)

		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Defense action logged for IP: " + host))
	})

	log.Println("Blue Team service running on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func enableCORS(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		(*w).WriteHeader(http.StatusOK)
	}
}
