package github

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"pulse/internal/config"
	"pulse/internal/style"
)

type Model struct {
	config      config.Config
	events      []Event
	selected    int
	lastUpdated time.Time
	loading     bool
	err         error
	spinner     spinner.Model
}

func New(cfg config.Config) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = style.AccentStyle
	return Model{
		config:  cfg,
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
	return tea.Batch(FetchCmd(m.config), m.spinner.Tick)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ResponseMsg:
		m.loading = false
		if msg.Error != nil {
			m.err = msg.Error
		} else {
			m.events = msg.Events
			m.err = nil
			m.lastUpdated = time.Now()
		}
		return m, tickCmd()

	case TickMsg:
		m.loading = true
		return m, tea.Batch(FetchCmd(m.config), m.spinner.Tick)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *Model) SelectNext() {
	if len(m.events) > 0 {
		m.selected = (m.selected + 1) % len(m.events)
	}
}

func (m *Model) SelectPrev() {
	if len(m.events) > 0 {
		m.selected = (m.selected - 1 + len(m.events)) % len(m.events)
	}
}

func (m Model) SelectedURL() string {
	if len(m.events) == 0 {
		return ""
	}
	return m.events[m.selected].URL
}

func (m Model) View(width, height int) string {
	title := style.TitleStyle.Render("🐙 GitHub")

	if m.loading && m.lastUpdated.IsZero() {
		return fmt.Sprintf("%s\n\n  %s Loading activity...", title, m.spinner.View())
	}

	if m.err != nil && m.lastUpdated.IsZero() {
		return fmt.Sprintf("%s\n\n%s", title, style.ErrorStyle.Render("  "+m.err.Error()))
	}

	var lines []string
	lines = append(lines, title)
	lines = append(lines, "")

	maxDetail := width - 8
	if maxDetail < 20 {
		maxDetail = 20
	}

	maxEvents := height / 3
	if maxEvents < 3 {
		maxEvents = 3
	}
	if maxEvents > len(m.events) {
		maxEvents = len(m.events)
	}

	for i := 0; i < maxEvents; i++ {
		event := m.events[i]
		repo := event.Repo
		if parts := strings.SplitN(repo, "/", 2); len(parts) == 2 {
			repo = parts[1]
		}

		cursor := "  "
		if i == m.selected {
			cursor = style.AccentStyle.Render("▸ ")
		}

		lines = append(lines, fmt.Sprintf("%s%s %s %s  %s",
			cursor,
			style.AccentStyle.Render("★"),
			event.Action,
			style.BoldWhite.Render(repo),
			style.SubtitleStyle.Render(timeAgo(event.Created)),
		))

		if event.Detail != "" {
			d := event.Detail
			if len(d) > maxDetail {
				d = d[:maxDetail-3] + "..."
			}
			lines = append(lines, fmt.Sprintf("     %s",
				style.SubtitleStyle.Render(d),
			))
		}
	}

	lines = append(lines, "")
	lines = append(lines, style.SubtitleStyle.Render(
		fmt.Sprintf("  ↑↓ navigate · o open · Updated %s", m.lastUpdated.Format("15:04")),
	))

	return strings.Join(lines, "\n")
}

func timeAgo(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "now"
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	}
}
