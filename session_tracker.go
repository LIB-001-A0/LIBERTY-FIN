package fin

// TrackSessions computes status from redacted, pre-existing event records.
// It never extracts tokens or contacts a production identity provider.
func TrackSessions(events []Event) map[string]string {
	status := make(map[string]string)
	for _, e := range events {
		if e.SessionID == "" {
			continue
		}
		switch e.Kind {
		case "session_open":
			status[e.SessionID] = "active"
		case "session_revoke":
			status[e.SessionID] = "revoked"
		case "session_expire":
			status[e.SessionID] = "expired"
		}
	}
	return status
}
