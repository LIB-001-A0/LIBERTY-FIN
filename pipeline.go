package fin

import (
	"errors"
	"sort"
)

// Analyze implements SIGNAL -> FUSE -> MAP -> ATTACK MAP (defensive risk
// review) -> RECONSTRUCT -> EXPAND (report only) on offline event data.
func Analyze(events []Event) (Report, error) {
	if len(events) == 0 {
		return Report{}, errors.New("no events supplied")
	}
	seen := map[string]bool{}
	copyEvents := append([]Event(nil), events...)
	sort.SliceStable(copyEvents, func(i, j int) bool { return copyEvents[i].Time.Before(copyEvents[j].Time) })
	for _, e := range copyEvents {
		if e.ID == "" || seen[e.ID] {
			return Report{}, errors.New("missing or duplicate event id")
		}
		seen[e.ID] = true
	}
	nodes, edges := BuildGraph(copyEvents)
	decisions := make([]Decision, 0, len(copyEvents))
	for _, e := range copyEvents {
		decisions = append(decisions, Guard(e))
	}
	return Report{Events: len(copyEvents), Nodes: nodes, Edges: edges, Sessions: TrackSessions(copyEvents), Decisions: decisions}, nil
}
