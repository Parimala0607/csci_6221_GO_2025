package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"redblue-sim/models"
	"redblue-sim/shared"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var db *sql.DB

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)
var clientsMutex sync.Mutex

func main() {
	var err error
	db, err = shared.InitDB()
	if err != nil {
		log.Fatal("Database initialization failed:", err)
	}
	defer db.Close()

	http.HandleFunc("/defend-sqlinjection", defendSQLInjection)
	http.HandleFunc("/defend-portscan", defendPortScan)
	http.HandleFunc("/api/blueteam/ddos/defend", defendDDoS)
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/notify", handleNotification)
	http.HandleFunc("/api/alerts/blue", getAlerts)
	http.HandleFunc("/defend-maliciousupload", defendMaliciousUpload)
	http.HandleFunc("/defend-steganography", defendSteganography)
	http.HandleFunc("/api/blocked", handleBlockedIPs)

	log.Println("Blue Team service running on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

// CORS handler
func enableCORS(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	clientsMutex.Lock()
	clients[conn] = true
	clientsMutex.Unlock()

	log.Println("New WebSocket client connected")

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			clientsMutex.Lock()
			delete(clients, conn)
			clientsMutex.Unlock()
			conn.Close()
			log.Println("WebSocket client disconnected")
			break
		}
	}
}

func handleNotification(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	msg := string(body)

	clientsMutex.Lock()
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			log.Println("WebSocket send error:", err)
			client.Close()
			delete(clients, client)
		}
	}
	clientsMutex.Unlock()

	w.WriteHeader(http.StatusOK)
}

func getAlerts(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	rows, err := db.Query("SELECT id, severity, message, source_ip, status FROM alerts ORDER BY id DESC")
	if err != nil {
		http.Error(w, "Failed to fetch alerts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var alerts []map[string]interface{}
	for rows.Next() {
		var id int
		var severity, message, sourceIP string
		var status sql.NullString

		if err := rows.Scan(&id, &severity, &message, &sourceIP, &status); err != nil {
			continue
		}

		alerts = append(alerts, map[string]interface{}{
			"id":        id,
			"severity":  severity,
			"message":   message,
			"source_ip": sourceIP,
			"status":    status.String,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alerts)
}

func defendSQLInjection(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pattern := "%sql injection%"
	_, err := db.Exec(`
		UPDATE alerts
		SET status = 'defended'
		WHERE LOWER(message) LIKE ?
		  AND (status IS NULL OR LOWER(status) <> 'defended');
	`, pattern)

	if err != nil {
		http.Error(w, "Failed to update SQL Injection defense", http.StatusInternalServerError)
		return
	}

	insertLog("SQL Injection Defense", "Blocked suspicious SQL query attempt from untrusted source", "blue-team")
	w.Write([]byte("SQL Injection attempt successfully blocked and logged"))
}

func defendPortScan(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pattern := "%port scan%"
	_, err := db.Exec(`
		UPDATE alerts
		SET status = 'defended'
		WHERE LOWER(message) LIKE ?
		  AND (status IS NULL OR LOWER(status) <> 'defended');
	`, pattern)

	if err != nil {
		http.Error(w, "Failed to update Port Scan defense", http.StatusInternalServerError)
		return
	}

	insertLog("Port Scan Defense", "Multiple unsolicited connection attempts blocked — port scan defense activated", "blue-team")
	w.Write([]byte("Port scanning activity detected and blocked"))
}

// defendSteganography defends against a simulated Steganography attack.
func defendSteganography(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var (
		alertID   int
		filename  string
		message   string
		sourceIP  string
		imageData []byte
	)

	err := db.QueryRow(`
	SELECT a.id, a.message, a.source_ip, si.filename, si.image_data
	FROM alerts a
	JOIN stego_images si ON a.id = si.alert_id
	WHERE a.message LIKE '%Steganography%'
	ORDER BY a.id DESC
	LIMIT 1
`).Scan(&alertID, &message, &sourceIP, &filename, &imageData)
	if err != nil {
		http.Error(w, "No steganography attacks to defend against", http.StatusNotFound)
		return
	}

	hiddenData := ""
	if idx := bytes.LastIndex(imageData, []byte("STEGO:")); idx != -1 {
		hiddenData = string(imageData[idx+len("STEGO:"):])
	}

	blockReason := fmt.Sprintf("Blocked IP %s due to hidden data in %s", sourceIP, filename)
	if err := BlockIP(sourceIP, blockReason); err != nil {
		log.Printf("Failed to block IP %s: %v", sourceIP, err)
	}

	_, err = db.Exec(`
        UPDATE alerts
        SET status = 'defended'
        WHERE id = ? AND (status IS NULL OR status <> 'defended')
    `, alertID)
	if err != nil {
		http.Error(w, "Failed to update alert status", http.StatusInternalServerError)
		return
	}

	insertLog("Steganography Defense",
		fmt.Sprintf("Analyzed image %s, found: %s", filename, hiddenData),
		"blue-team")

	result := map[string]interface{}{
		"status":      "defended",
		"filename":    filename,
		"hidden_data": hiddenData,
		"source_ip":   sourceIP,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// defendMaliciousUpload defends against a simulated Malicious File Upload.
func defendMaliciousUpload(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var (
		alertID  int
		filename string
		sourceIP string
		fileData []byte
	)

	err := db.QueryRow(`
		SELECT a.id, mf.filename, a.source_ip, mf.file_data
		FROM alerts a
		JOIN malicious_files mf ON a.id = mf.alert_id
		WHERE a.status = 'not_defended'
		AND a.message LIKE '%Malicious file%'
		ORDER BY a.id DESC
		LIMIT 1
	`).Scan(&alertID, &filename, &sourceIP, &fileData)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "No files to defend", http.StatusNotFound)
		} else {
			log.Printf("QueryRow failed: %v", err)
			http.Error(w, "DB query failed", http.StatusInternalServerError)
		}
		return
	}

	keywords := []string{"virus", "malware", "exploit", "payload", "attack"}
	var foundKeywords []string
	content := strings.ToLower(string(fileData))
	for _, kw := range keywords {
		if strings.Contains(content, kw) {
			foundKeywords = append(foundKeywords, kw)
		}
	}

	var blocked bool
	var blockMessage string
	if len(foundKeywords) > 0 {
		blockMessage = fmt.Sprintf("Blocked IP %s - File %s contained: %s",
			sourceIP, filename, strings.Join(foundKeywords, ", "))
		if err := BlockIP(sourceIP, blockMessage); err != nil {
			log.Printf("BlockIP failed: %v", err)
			blockMessage += " (but IP blocking failed)"
		} else {
			blocked = true
		}
	} else {
		blockMessage = fmt.Sprintf("IP %s not blocked - no malicious content found", sourceIP)
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "DB transaction start failed", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	result, err := tx.Exec(`
	UPDATE alerts 
	SET status = 'defended' 
	WHERE id = ? 
	  AND status = 'not_defended'
	  AND message LIKE '%Malicious file%'`, alertID)
	if err != nil {
		log.Printf("Failed to update alert status for id %d: %v", alertID, err)
		http.Error(w, "Status update failed", http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		log.Printf("No rows updated for id %d (maybe already defended)", alertID)
		http.Error(w, "Status update failed", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Transaction commit failed: %v", err)
		http.Error(w, "Transaction commit failed", http.StatusInternalServerError)
		return
	}

	insertLog("Malicious Upload Defense", "Blocked IP due to malicious file", sourceIP)

	response := map[string]interface{}{
		"filename":       filename,
		"defended":       true,
		"blocked":        blocked,
		"block_message":  blockMessage,
		"keywords_found": foundKeywords,
		"source_ip":      sourceIP,
		"file_size":      len(fileData),
		"file_content":   string(fileData),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Response encode failed: %v", err)
	}
}

func BlockIP(ip string, reason string) error {
	if ip == "" {
		return fmt.Errorf("empty IP address")
	}
	_, err := db.Exec(`
		INSERT OR REPLACE INTO blocked_ips (ip_address, reason)
		VALUES (?, ?)`, ip, reason)
	return err
}

var (
	lastRequestTimes = make(map[string]time.Time)
	requestCounts    = make(map[string]int)
	rateLimitWindow  = 2 * time.Second
)

func defendDDoS(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query(`
		SELECT DISTINCT source_ip FROM alerts
		WHERE message LIKE '%DDoS%'
		AND (status IS NULL OR status = 'not_defended');
	`)
	if err != nil {
		http.Error(w, "Failed to fetch IPs", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var ips []string
	for rows.Next() {
		var ip string
		if err := rows.Scan(&ip); err == nil {
			ips = append(ips, ip)
		}
	}

	if len(ips) == 0 {
		w.Write([]byte("No DDoS alerts left to defend"))
		return
	}

	var (
		totalUpdated int
		blockedIPs   []string
		defendedIPs  []string
		mu           sync.Mutex
		wg           sync.WaitGroup
	)

	for _, ip := range ips {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()

			lastTime, seen := lastRequestTimes[ip]
			if seen && time.Since(lastTime) < rateLimitWindow {
				log.Printf("Skipping %s: rate limited", ip)
				return
			}
			lastRequestTimes[ip] = time.Now()
			requestCounts[ip]++
			isAbusive := requestCounts[ip] > 1

			// Perform update
			result, err := db.Exec(`
				UPDATE alerts
				SET status = 'defended'
				WHERE source_ip = ?
				AND message LIKE '%DDoS%'
				AND (status IS NULL OR status = 'not_defended')
			`, ip)
			if err != nil {
				log.Printf("Failed to update alerts for %s: %v", ip, err)
				return
			}

			rowsAffected, _ := result.RowsAffected()
			if rowsAffected == 0 {
				return
			}

			mu.Lock()
			totalUpdated += int(rowsAffected)
			mu.Unlock()

			if isAbusive {
				blockReason := fmt.Sprintf("Blocked IP %s due to repeated DDoS attempts", ip)
				if err := BlockIP(ip, blockReason); err == nil {
					mu.Lock()
					blockedIPs = append(blockedIPs, ip)
					mu.Unlock()
					insertLog("DDoS Defense", blockReason, ip)
				} else {
					log.Printf("BlockIP failed for %s: %v", ip, err)
				}
			} else {
				mu.Lock()
				defendedIPs = append(defendedIPs, ip)
				mu.Unlock()
				insertLog("DDoS Defense", fmt.Sprintf("DDoS alert defended for IP %s", ip), ip)
			}
		}(ip)
	}

	wg.Wait() // wait for all goroutines

	response := map[string]interface{}{
		"total_alerts_defended": totalUpdated,
		"ips_blocked":           blockedIPs,
		"ips_defended":          defendedIPs,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func insertLog(action, description, sourceIP string) {
	_, err := db.Exec(`
		INSERT INTO logs (source_ip, action, description)
		VALUES (?, ?, ?);
	`, sourceIP, action, description)
	if err != nil {
		log.Println("Failed to insert into logs:", err)
	}
}

func getIP(r *http.Request) string {
	ip := r.RemoteAddr
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ip = forwarded
	}
	return ip
}

func handleBlockedIPs(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	rows, err := db.Query(`SELECT ip_address, reason, blocked_at FROM blocked_ips`)
	if err != nil {
		log.Println("Error querying blocked IPs:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []models.BlockedIP
	for rows.Next() {
		var b models.BlockedIP
		if err := rows.Scan(&b.IPAddress, &b.Reason, &b.Timestamp); err != nil {
			continue
		}
		results = append(results, b)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
