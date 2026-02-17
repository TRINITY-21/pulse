package news

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"pulse/internal/style"
)

type Model struct {
	stories     []Story
	selected    int
	lastUpdated time.Time
	loading     bool
	err         error
	spinner     spinner.Model
}

type OpenURLMsg struct {
	URL string
}

func New() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = style.AccentStyle
	return Model{
		loading: true,
		spinner: s,
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(5*time.Minute, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(FetchCmd(), m.spinner.Tick)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ResponseMsg:
		m.loading = false
		if msg.Error != nil {
			m.err = msg.Error
		} else {
			m.stories = msg.Stories
			m.err = nil
			m.lastUpdated = time.Now()
		}
		return m, tickCmd()

	case TickMsg:
		m.loading = true
		return m, tea.Batch(FetchCmd(), m.spinner.Tick)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *Model) SelectNext() {
	if len(m.stories) > 0 {
		m.selected = (m.selected + 1) % len(m.stories)
	}
}

func (m *Model) SelectPrev() {
	if len(m.stories) > 0 {
		m.selected = (m.selected - 1 + len(m.stories)) % len(m.stories)
	}
}

func (m Model) SelectedURL() string {
	if len(m.stories) == 0 {
		return ""
	}
	s := m.stories[m.selected]
	if s.URL != "" {
		return s.URL
	}
	return fmt.Sprintf("https://news.ycombinator.com/item?id=%d", s.ID)
}

func (m Model) View(width, height int) string {
	title := style.TitleStyle.Render("📰 News")

	if m.loading && m.lastUpdated.IsZero() {
		return fmt.Sprintf("%s\n\n  %s Loading stories...", title, m.spinner.View())
	}

	if m.err != nil && m.lastUpdated.IsZero() {
		return fmt.Sprintf("%s\n\n%s", title, style.ErrorStyle.Render("  "+m.err.Error()))
	}

	var lines []string
	lines = append(lines, title)
	lines = append(lines, "")

	maxTitle := width - 8
	if maxTitle < 20 {
		maxTitle = 20
	}

	maxStories := height / 3
	if maxStories < 3 {
		maxStories = 3
	}
	if maxStories > len(m.stories) {
		maxStories = len(m.stories)
	}

	for i := 0; i < maxStories; i++ {
		story := m.stories[i]
		t := story.Title
		if len(t) > maxTitle {
			t = t[:maxTitle-3] + "..."
		}

		cursor := "  "
		if i == m.selected {
			cursor = style.AccentStyle.Render("▸ ")
		}

		lines = append(lines, fmt.Sprintf("%s%s %s",
			cursor,
			style.AccentStyle.Render(fmt.Sprintf("%d.", i+1)),
			t,
		))

		// Detail line: domain/snippet + metadata
		detail := domainFrom(story.URL)
		if story.Text != "" {
			snippet := stripHTML(story.Text)
			maxSnippet := maxTitle - 5
			if maxSnippet > 80 {
				maxSnippet = 80
			}
			if len(snippet) > maxSnippet {
				snippet = snippet[:maxSnippet-3] + "..."
			}
			detail = snippet
		}

		meta := fmt.Sprintf("%d pts · %s · %d comments · %s",
			story.Score, story.By, story.Comments, timeAgo(story.Time))

		if detail != "" {
			lines = append(lines, fmt.Sprintf("     %s  %s",
				style.SubtitleStyle.Render(detail),
				style.SubtitleStyle.Render("· "+meta),
			))
		} else {
			lines = append(lines, fmt.Sprintf("     %s",
				style.SubtitleStyle.Render(meta),
			))
		}
	}

	lines = append(lines, "")
	lines = append(lines, style.SubtitleStyle.Render(
		fmt.Sprintf("  ↑↓ navigate · o open · Updated %s", m.lastUpdated.Format("15:04")),
	))

	return strings.Join(lines, "\n")
}

func domainFrom(rawURL string) string {
	if rawURL == "" {
		return ""
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	host := u.Hostname()
	host = strings.TrimPrefix(host, "www.")
	return host
}

func stripHTML(s string) string {
	var out strings.Builder
	inTag := false
	for _, r := range s {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			continue
		}
		if !inTag {
			if r == '\n' {
				out.WriteRune(' ')
			} else {
				out.WriteRune(r)
			}
		}
	}
	return strings.TrimSpace(out.String())
}

func timeAgo(unix int64) string {
	d := time.Since(time.Unix(unix, 0))
	switch {
	case d.Hours() >= 24:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	case d.Hours() >= 1:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	case d.Minutes() >= 1:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	default:
		return "just now"
	}
}
