package fin

import "time"

// Event contains synthetic, redacted telemetry only. Never store credentials.
type Event struct {
	ID        string    `json:"id"`
	Time      time.Time `json:"time"`
	Kind      string    `json:"kind"`
	Actor     string    `json:"actor"`
	Source    string    `json:"source"`
	Target    string    `json:"target"`
	SessionID string    `json:"session_id,omitempty"`
	Allowed   bool      `json:"allowed"`
	Detail    string    `json:"detail,omitempty"`
}

type Decision struct {
	EventID string `json:"event_id"`
	Action  string `json:"action"`
	Reason  string `json:"reason"`
}

type Report struct {
	Events    int               `json:"events"`
	Nodes     []string          `json:"nodes"`
	Edges     []Edge            `json:"edges"`
	Sessions  map[string]string `json:"sessions"`
	Decisions []Decision        `json:"decisions"`
}
