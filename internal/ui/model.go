package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"pulse/internal/config"
	"pulse/internal/panels/crypto"
	"pulse/internal/panels/github"
	"pulse/internal/panels/news"
	"pulse/internal/panels/weather"
)

type clockTickMsg time.Time

type Model struct {
	Config  config.Config
	width   int
	height  int
	focused int
	clock   time.Time
	Visible [4]bool // 0=weather, 1=crypto, 2=news, 3=github
	Weather weather.Model
	Crypto  crypto.Model
	News    news.Model
	GitHub  github.Model
}

func NewModel(cfg config.Config, visible [4]bool) Model {
	// Focus the first visible panel
	focused := 0
	for i, v := range visible {
		if v {
			focused = i
			break
		}
	}
	return Model{
		Config:  cfg,
		focused: focused,
		clock:   time.Now(),
		Visible: visible,
		Weather: weather.New(cfg),
		Crypto:  crypto.New(cfg),
		News:    news.New(),
		GitHub:  github.New(cfg),
	}
}

func (m *Model) focusNext() {
	for i := 1; i <= 4; i++ {
		next := (m.focused + i) % 4
		if m.Visible[next] {
			m.focused = next
			return
		}
	}
}

func clockTickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return clockTickMsg(t)
	})
}

func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{clockTickCmd()}
	if m.Visible[0] {
		cmds = append(cmds, m.Weather.Init())
	}
	if m.Visible[1] {
		cmds = append(cmds, m.Crypto.Init())
	}
	if m.Visible[2] {
		cmds = append(cmds, m.News.Init())
	}
	if m.Visible[3] {
		cmds = append(cmds, m.GitHub.Init())
	}
	return tea.Batch(cmds...)
}

func (m Model) visibleCount() int {
	count := 0
	for _, v := range m.Visible {
		if v {
			count++
		}
	}
	return count
}
