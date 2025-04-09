package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"redblue-sim/shared"
	"strconv"
	"strings"
	"time"
)

var db *sql.DB

func main() {
	rand.Seed(time.Now().UnixNano())

	var err error
	db, err = shared.InitDB()
	if err != nil {
		log.Fatal("Database initialization failed:", err)
	}
	defer db.Close()

	http.HandleFunc("/sqlinjection", handleSQLInjection)
	http.HandleFunc("/portscan", handlePortScan)
	http.HandleFunc("/maliciousupload", handleMaliciousUpload)
	http.HandleFunc("/ddos", handleDDoS)
	http.HandleFunc("/steganography", handleSteganography)
	http.HandleFunc("/api/alerts/red", getAlerts)
	log.Println("Red Team service running on port 8082...")
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func handleSQLInjection(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	type SQLPayload struct {
		Message  string
		Severity string
	}

	sqlPayloads := []SQLPayload{
		{Message: "SQL Injection: ' OR '1'='1'; --", Severity: "low"},
		{Message: "SQL Injection: admin' --", Severity: "low"},
		{Message: "SQL Injection: ' UNION SELECT NULL, version(); --", Severity: "medium"},
		{Message: "SQL Injection: ' AND 1=0 UNION SELECT username, password FROM users --", Severity: "high"},
		{Message: "SQL Injection: ' OR EXISTS(SELECT * FROM users WHERE username = 'admin') --", Severity: "medium"},
		{Message: "SQL Injection: ' OR SLEEP(5)--", Severity: "high"},
	}

	selected := sqlPayloads[rand.Intn(len(sqlPayloads))]
	sourceIP := randomSourceIP()

	_, err := db.Exec(`
		INSERT INTO alerts (severity, message, source_ip)
		VALUES (?, ?, ?);
	`, selected.Severity, selected.Message, sourceIP)
	if err != nil {
		http.Error(w, "Failed to insert SQL Injection alert", http.StatusInternalServerError)
		return
	}

	insertLog("SQL Injection", fmt.Sprintf("%s attempt.", selected.Severity), sourceIP)
	broadcastAlert(selected.Message, sourceIP)

	w.Write([]byte(fmt.Sprintf("SQL Injection (%s) attack logged from %s", selected.Severity, sourceIP)))
}

func handlePortScan(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w, r)
	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Dynamic options
	portSets := []struct {
		Ports    string
		Severity string
	}{
		{"22, 80, 443", "medium"},
		{"21, 23, 8080", "low"},
		{"3306, 5432, 1521", "high"}, // DB ports
		{"25, 110, 143", "medium"},   // Mail ports
		{"135, 139, 445", "high"},    // Windows SMB ports
	}

	selected := portSets[rand.Intn(len(portSets))]
	message := fmt.Sprintf("Port Scan detected on ports: %s", selected.Ports)
	sourceIP := randomSourceIP()

	_, err := db.Exec(`
		INSERT INTO alerts (severity, message, source_ip)
		VALUES (?, ?, ?);
	`, selected.Severity, message, sourceIP)

	if err != nil {
		http.Error(w, "Failed to insert Port Scan alert", http.StatusInternalServerError)
		return
	}

	insertLog("Port Scan", fmt.Sprintf("Detected scan on ports: %s", selected.Ports), sourceIP)
	broadcastAlert(message, sourceIP)

	w.Write([]byte(fmt.Sprintf("Port Scan (%s) attack logged from %s", selected.Severity, sourceIP)))
}

func getAlerts(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w, r)
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

func handleMaliciousUpload(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w, r)
	log.Printf("Incoming Content-Type: %s", r.Header.Get("Content-Type"))
	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		log.Printf("Invalid content type: %s", contentType)
		http.Error(w, "Expected multipart/form-data", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("FormFile error:", err)
		http.Error(w, "Need file upload", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 2. Read content
	fileData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Can't read file", http.StatusInternalServerError)
		return
	}

	// 3. Generate IP and store
	sourceIP := fmt.Sprintf("192.168.1.%d", rand.Intn(254)+1)

	// Insert alert (using your existing schema)
	res, err := db.Exec(`
        INSERT INTO alerts (severity, message, source_ip, status)
        VALUES (?, ?, ?, ?)`,
		"high",
		"Malicious file: "+handler.Filename,
		sourceIP,
		"not_defended")
	if err != nil {
		http.Error(w, "DB alert error", http.StatusInternalServerError)
		return
	}

	alertID, _ := res.LastInsertId()

	// Store file
	_, err = db.Exec(`
        INSERT INTO malicious_files (alert_id, filename, file_data)
        VALUES (?, ?, ?)`,
		alertID, handler.Filename, fileData)
	if err != nil {
		http.Error(w, "DB file error", http.StatusInternalServerError)
		return
	}

	// Log using your existing logs table
	db.Exec(`
        INSERT INTO logs (source_ip, action, description)
        VALUES (?, ?, ?)`,
		sourceIP,
		"malicious_upload",
		"Uploaded: "+handler.Filename)

	filename := handler.Filename

	insertLog("MaliciousUpload", fmt.Sprintf("Uploaded %s", filename), sourceIP)

	message := fmt.Sprintf("Malicious file uploaded: %s", filename)

	broadcastAlert(message, sourceIP)

	w.Write([]byte("Malicious upload logged: " + filename))

}

func handleDDoS(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sourceIP := randomSourceIP()
	message := "DDoS attempt: high request volume"

	for i := 0; i < 10; i++ {
		_, err := db.Exec(`
			INSERT INTO alerts (severity, message, source_ip, status)
			VALUES (?, ?, ?, 'not_defended');
		`, "high", message, sourceIP)
		if err != nil {
			http.Error(w, "Failed to insert DDoS alert", http.StatusInternalServerError)
			return
		}

		insertLog("DDoS", "High-severity DDoS traffic detected", sourceIP)
		broadcastAlert(message, sourceIP)
	}

	log.Printf("Simulated  DDoS alerts from IP %s", sourceIP)
	w.Write([]byte(fmt.Sprintf("Simulated  DDoS alerts from IP %s", sourceIP)))
}

func handleSteganography(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w, r)
	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form with size limit
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "File too large (max 10MB)", http.StatusBadRequest)
		return
	}

	// Get uploaded file
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "No image uploaded", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate image type
	buf := make([]byte, 512)
	if _, err := file.Read(buf); err != nil {
		http.Error(w, "Invalid image file", http.StatusBadRequest)
		return
	}

	if fileType := http.DetectContentType(buf); !strings.HasPrefix(fileType, "image/") {
		http.Error(w, "Only image files are allowed", http.StatusBadRequest)
		return
	}

	// Reset and read full image
	if _, err := file.Seek(0, 0); err != nil {
		http.Error(w, "Failed to read image", http.StatusInternalServerError)
		return
	}
	imgData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read image data", http.StatusInternalServerError)
		return
	}

	// Generate and embed payload
	payload := fmt.Sprintf("SECRET-%d", rand.Intn(10000))
	stegoImage := append(imgData, []byte("\nSTEGO:"+payload)...)

	// Sanitize filename
	filename := filepath.Base(handler.Filename)
	if filename == "." {
		filename = "uploaded_image"
	}

	// Database transaction
	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Insert alert
	res, err := tx.Exec(`
        INSERT INTO alerts (severity, message, source_ip)
        VALUES (?, ?, ?)`,
		"low",
		fmt.Sprintf("Steganography: hidden payload in image %s", filename),
		randomSourceIP())
	if err != nil {
		http.Error(w, "Failed to insert alert", http.StatusInternalServerError)
		return
	}

	alertID, err := res.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to get alert ID", http.StatusInternalServerError)
		return
	}

	// Store the modified image (with payload)
	if _, err := tx.Exec(`
        INSERT INTO stego_images (alert_id, filename, image_data)
        VALUES (?, ?, ?)`,
		alertID, filename, stegoImage); err != nil {
		http.Error(w, "Failed to store image", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Database commit failed", http.StatusInternalServerError)
		return
	}

	insertLog("Steganography", fmt.Sprintf("Uploaded %s with payload", filename), randomSourceIP())
	message := fmt.Sprintf("New stego image: %s", filename)
	sourceIP := randomSourceIP()
	broadcastAlert(message, sourceIP)

	// Return modified image
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"filename": handler.Filename,
		"message":  "Steganography attack successful",
		"status":   "uploaded",
	})
}

func randomSourceIP() string {
	lastOctet := rand.Intn(254) + 1
	return "192.168.1." + strconv.Itoa(lastOctet)
}

func enableCORS(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		(*w).WriteHeader(http.StatusOK)
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

func broadcastAlert(message string, sourceIP string) {
	blueTeamURL := os.Getenv("BLUE_TEAM_URL")
	if blueTeamURL == "" {
		blueTeamURL = "http://localhost:8081"
	}

	// Create a structured alert payload
	alert := struct {
		Message  string `json:"message"`
		SourceIP string `json:"source_ip"`
	}{
		Message:  message,
		SourceIP: sourceIP,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(alert)
	if err != nil {
		log.Println("Failed to marshal alert:", err)
		return
	}

	// Create request with JSON body
	req, err := http.NewRequest("POST", blueTeamURL+"/notify", bytes.NewReader(jsonData))
	if err != nil {
		log.Println("Failed to create request:", err)
		return
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Forwarded-For", sourceIP)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to notify Blue Team:", err)
		return
	}
	defer resp.Body.Close()

	// Ensure response is fully read
	if _, err := io.ReadAll(resp.Body); err != nil {
		log.Println("Failed to read response body:", err)
	}
}
