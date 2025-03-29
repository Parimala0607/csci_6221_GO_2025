package models

type Alert struct {
	ID        int
	Timestamp string
	Severity  string
	Message   string
	SourceIP  string
}

type Log struct {
	ID          int
	Timestamp   string
	SourceIP    string
	Action      string
	Description string
}

type User struct {
	ID           int
	Username     string
	PasswordHash string
	Role         string
}
