package news

import "time"

type Story struct {
	ID       int
	Title    string
	URL      string
	Text     string
	Score    int
	By       string
	Time     int64
	Comments int
}

type apiStory struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Text        string `json:"text"`
	Score       int    `json:"score"`
	By          string `json:"by"`
	Time        int64  `json:"time"`
	Type        string `json:"type"`
	Descendants int    `json:"descendants"`
}

type ResponseMsg struct {
	Stories []Story
	Error   error
}

type TickMsg time.Time
