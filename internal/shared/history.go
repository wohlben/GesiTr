package shared

import (
	"encoding/json"
	"time"
)

// VersionEntry is the API response DTO for a single version in the history.
type VersionEntry struct {
	Version   int             `json:"version"`
	Snapshot  json.RawMessage `json:"snapshot"`
	ChangedAt time.Time       `json:"changedAt"`
	ChangedBy string          `json:"changedBy"`
}

// SnapshotJSON marshals a DTO to JSON for history storage.
func SnapshotJSON(dto any) string {
	data, _ := json.Marshal(dto)
	return string(data)
}
