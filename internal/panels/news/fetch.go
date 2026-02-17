package news

import (
	"encoding/json"
	"fmt"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
)

const storyCount = 8

func FetchCmd() tea.Cmd {
	return func() tea.Msg {
		resp, err := http.Get("https://hacker-news.firebaseio.com/v0/topstories.json")
		if err != nil {
			return ResponseMsg{Error: fmt.Errorf("request failed: %w", err)}
		}
		defer resp.Body.Close()

		var ids []int
		if err := json.NewDecoder(resp.Body).Decode(&ids); err != nil {
			return ResponseMsg{Error: fmt.Errorf("decode failed: %w", err)}
		}

		if len(ids) > storyCount {
			ids = ids[:storyCount]
		}

		var stories []Story
		for _, id := range ids {
			url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", id)
			r, err := http.Get(url)
			if err != nil {
				continue
			}

			var s apiStory
			json.NewDecoder(r.Body).Decode(&s)
			r.Body.Close()

			if s.Title == "" {
				continue
			}

			stories = append(stories, Story{
				ID:       s.ID,
				Title:    s.Title,
				URL:      s.URL,
				Text:     s.Text,
				Score:    s.Score,
				By:       s.By,
				Time:     s.Time,
				Comments: s.Descendants,
			})
		}

		return ResponseMsg{Stories: stories}
	}
}
