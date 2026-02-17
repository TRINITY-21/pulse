package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"pulse/internal/config"
)

const eventCount = 8

func FetchCmd(cfg config.Config) tea.Cmd {
	return func() tea.Msg {
		url := fmt.Sprintf("https://api.github.com/users/%s/events?per_page=%d",
			cfg.GitHubUser, eventCount)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return ResponseMsg{Error: fmt.Errorf("request failed: %w", err)}
		}

		req.Header.Set("Accept", "application/vnd.github.v3+json")
		if cfg.GitHubToken != "" {
			req.Header.Set("Authorization", "Bearer "+cfg.GitHubToken)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return ResponseMsg{Error: fmt.Errorf("request failed: %w", err)}
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return ResponseMsg{Error: fmt.Errorf("API returned %d", resp.StatusCode)}
		}

		var apiEvents []apiEvent
		if err := json.NewDecoder(resp.Body).Decode(&apiEvents); err != nil {
			return ResponseMsg{Error: fmt.Errorf("decode failed: %w", err)}
		}

		var events []Event
		for _, e := range apiEvents {
			created, _ := time.Parse(time.RFC3339, e.CreatedAt)
			action, detail, url := formatEvent(e)
			events = append(events, Event{
				Type:    e.Type,
				Repo:    e.Repo.Name,
				Action:  action,
				Detail:  detail,
				URL:     url,
				Created: created,
			})
		}

		return ResponseMsg{Events: events}
	}
}

func formatEvent(e apiEvent) (action, detail, url string) {
	switch e.Type {
	case "PushEvent":
		action = "Pushed to"
		if len(e.Payload.Commits) == 1 {
			detail = firstLine(e.Payload.Commits[0].Message)
		} else if len(e.Payload.Commits) > 1 {
			detail = fmt.Sprintf("%d commits · %s", len(e.Payload.Commits), firstLine(e.Payload.Commits[0].Message))
		}
		url = fmt.Sprintf("https://github.com/%s", e.Repo.Name)
	case "CreateEvent":
		if e.Payload.RefType == "repository" {
			action = "Created repo"
		} else {
			action = fmt.Sprintf("Created %s", e.Payload.RefType)
			detail = e.Payload.Ref
		}
		url = fmt.Sprintf("https://github.com/%s", e.Repo.Name)
	case "WatchEvent":
		action = "Starred"
		url = fmt.Sprintf("https://github.com/%s", e.Repo.Name)
	case "PullRequestEvent":
		pa := e.Payload.Action
		if pa == "closed" && e.Payload.PullRequest.Merged {
			pa = "merged"
		}
		action = fmt.Sprintf("PR %s on", pa)
		detail = e.Payload.PullRequest.Title
		url = e.Payload.PullRequest.HTMLURL
	case "IssuesEvent":
		action = fmt.Sprintf("Issue %s on", e.Payload.Action)
		detail = e.Payload.Issue.Title
		url = e.Payload.Issue.HTMLURL
	case "IssueCommentEvent":
		action = "Commented on"
		detail = e.Payload.Issue.Title
		url = e.Payload.Comment.HTMLURL
	case "ForkEvent":
		action = "Forked"
		detail = "→ " + e.Payload.Forkee.FullName
		url = fmt.Sprintf("https://github.com/%s", e.Payload.Forkee.FullName)
	case "DeleteEvent":
		action = fmt.Sprintf("Deleted %s in", e.Payload.RefType)
		detail = e.Payload.Ref
		url = fmt.Sprintf("https://github.com/%s", e.Repo.Name)
	default:
		action = e.Type
		url = fmt.Sprintf("https://github.com/%s", e.Repo.Name)
	}
	return
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}
