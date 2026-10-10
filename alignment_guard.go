package fin

// Guard evaluates policy metadata in supplied events. It does not execute
// privileged controls. A human-approved integration is required for enforcement.
func Guard(e Event) Decision {
	d := Decision{EventID: e.ID, Action: "observe", Reason: "telemetry only"}
	if !e.Allowed {
		d.Action = "review_block"
		d.Reason = "policy denied event"
	}
	if e.Kind == "prompt_injection" {
		d.Action = "quarantine_review"
		d.Reason = "untrusted instruction detected"
	}
	if e.Kind == "session_revoke" {
		d.Action = "verify_revoke"
		d.Reason = "confirm session revocation in authorized audit logs"
	}
	return d
}
