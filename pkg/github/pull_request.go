package github

import "time"

type PullRequest struct {
	URL        string          `json:"url"`
	ID         int             `json:"id"`
	NodeID     string          `json:"node_id"`
	Number     int             `json:"number"`
	State      string          `json:"state"`
	Locked     bool            `json:"locked"`
	Title      string          `json:"title"`
	Author     *UserCommitInfo `json:"user"`
	MergedUser *UserCommitInfo `json:"merged_by"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
	ClosedAt   time.Time       `json:"closed_at"`
	MergedAt   time.Time       `json:"merged_at"`
	MergeSHA   string          `json:"merge_commit_sha"`
	Merged     bool            `json:"merged"`
}
