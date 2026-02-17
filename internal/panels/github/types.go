package github

import "time"

type Event struct {
	Type    string
	Repo    string
	Action  string
	Detail  string
	URL     string
	Created time.Time
}

type apiEvent struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
	Payload struct {
		Action  string `json:"action"`
		Ref     string `json:"ref"`
		RefType string `json:"ref_type"`
		Commits []struct {
			Message string `json:"message"`
		} `json:"commits"`
		PullRequest struct {
			Title   string `json:"title"`
			HTMLURL string `json:"html_url"`
			Merged  bool   `json:"merged"`
		} `json:"pull_request"`
		Issue struct {
			Title   string `json:"title"`
			HTMLURL string `json:"html_url"`
		} `json:"issue"`
		Comment struct {
			HTMLURL string `json:"html_url"`
		} `json:"comment"`
		Forkee struct {
			FullName string `json:"full_name"`
		} `json:"forkee"`
	} `json:"payload"`
	CreatedAt string `json:"created_at"`
}

type ResponseMsg struct {
	Events []Event
	Error  error
}

type TickMsg time.Time
