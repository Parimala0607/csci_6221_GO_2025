package models

type Alert struct {
	ID        int
	Timestamp string
	Severity  string
	Message   string
	SourceIP  string
	Status    string
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

type BlockedIP struct {
	IPAddress string `json:"ip_address"`
	Reason    string `json:"reason"`
	Timestamp string `json:"blocked_at"`
}
