package things

// ActionResult represents a normalized Things callback response.
type ActionResult struct {
	Action              string            `json:"action"`
	ThingsID            string            `json:"things_id,omitempty"`
	ThingsIDs           []string          `json:"things_ids,omitempty"`
	ThingsSchemeVersion string            `json:"things_scheme_version,omitempty"`
	ThingsClientVersion string            `json:"things_client_version,omitempty"`
	Callback            map[string]string `json:"callback,omitempty"`
}
