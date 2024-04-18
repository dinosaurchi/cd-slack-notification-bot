package github

import "time"

type CDInfo struct {
	State    string      `json:"state"`
	Statuses []*cdStatus `json:"statuses"`
	SHA      string      `json:"sha"`
}

type cdStatus struct {
	URL         string    `json:"url"`
	ID          int       `json:"id"`
	NodeID      string    `json:"node_id"`
	State       string    `json:"state"`
	Description string    `json:"description"`
	TargetURL   string    `json:"target_url"`
	Context     string    `json:"context"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
