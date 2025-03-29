package shared

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB initializes and returns a database connection with schema setup
func InitDB() (*sql.DB, error) {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "sim.db" // Default fallback for local dev
	}

	absPath, _ := filepath.Abs(dbPath)
	log.Println("🛠  InitDB() using path:", absPath)

	dsn := fmt.Sprintf("%s?_foreign_keys=1&_journal_mode=WAL", dbPath)
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure PRAGMA settings
	if _, err = db.Exec(`PRAGMA busy_timeout = 5000;`); err != nil {
		return nil, fmt.Errorf("failed to set busy_timeout: %w", err)
	}

	// Table creation SQLs
	tables := []string{
		`CREATE TABLE IF NOT EXISTS alerts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			severity TEXT NOT NULL CHECK(severity IN ('low', 'medium', 'high', 'critical')),
			message TEXT NOT NULL,
			source_ip TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			source_ip TEXT NOT NULL,
			action TEXT NOT NULL,
			description TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			role TEXT DEFAULT 'viewer' CHECK(role IN ('admin', 'editor', 'viewer'))
		);`,
	}

	for _, stmt := range tables {
		if _, err := db.Exec(stmt); err != nil {
			return nil, fmt.Errorf("failed to create table: %w", err)
		}
	}

	log.Println("Tables created successfully.")

	// List tables found in the DB for verification
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		log.Println("Could not query tables:", err)
		return db, nil
	}
	defer rows.Close()

	log.Println("Tables in the database:")
	for rows.Next() {
		var name string
		rows.Scan(&name)
		log.Println("  -", name)
	}

	log.Println("Database initialized at:", absPath)
	return db, nil
}
