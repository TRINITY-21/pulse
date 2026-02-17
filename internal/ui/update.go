package ui

import (
	"os/exec"
	"runtime"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"pulse/internal/panels/crypto"
	"pulse/internal/panels/github"
	"pulse/internal/panels/news"
	"pulse/internal/panels/weather"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, Keys.Tab):
			m.focusNext()
			return m, nil
		case key.Matches(msg, Keys.Panel1):
			m.focused = 0
			return m, nil
		case key.Matches(msg, Keys.Panel2):
			m.focused = 1
			return m, nil
		case key.Matches(msg, Keys.Panel3):
			m.focused = 2
			return m, nil
		case key.Matches(msg, Keys.Panel4):
			m.focused = 3
			return m, nil
		case key.Matches(msg, Keys.Refresh):
			var refreshCmds []tea.Cmd
			if m.Visible[0] {
				refreshCmds = append(refreshCmds, weather.FetchCmd(m.Config))
			}
			if m.Visible[1] {
				refreshCmds = append(refreshCmds, crypto.FetchCmd(m.Config))
			}
			if m.Visible[2] {
				refreshCmds = append(refreshCmds, news.FetchCmd())
			}
			if m.Visible[3] {
				refreshCmds = append(refreshCmds, github.FetchCmd(m.Config))
			}
			return m, tea.Batch(refreshCmds...)

		// Toggle panels
		case key.Matches(msg, Keys.ToggleWeather):
			if m.visibleCount() > 1 || !m.Visible[0] {
				m.Visible[0] = !m.Visible[0]
				if m.Visible[0] {
					return m, m.Weather.Init()
				}
			}
			return m, nil
		case key.Matches(msg, Keys.ToggleCrypto):
			if m.visibleCount() > 1 || !m.Visible[1] {
				m.Visible[1] = !m.Visible[1]
				if m.Visible[1] {
					return m, m.Crypto.Init()
				}
			}
			return m, nil
		case key.Matches(msg, Keys.ToggleNews):
			if m.visibleCount() > 1 || !m.Visible[2] {
				m.Visible[2] = !m.Visible[2]
				if m.Visible[2] {
					return m, m.News.Init()
				}
			}
			return m, nil
		case key.Matches(msg, Keys.ToggleGitHub):
			if m.visibleCount() > 1 || !m.Visible[3] {
				m.Visible[3] = !m.Visible[3]
				if m.Visible[3] {
					return m, m.GitHub.Init()
				}
			}
			return m, nil

		// News panel navigation (when focused)
		case m.focused == 2 && key.Matches(msg, Keys.Down):
			m.News.SelectNext()
			return m, nil
		case m.focused == 2 && key.Matches(msg, Keys.Up):
			m.News.SelectPrev()
			return m, nil
		case m.focused == 2 && key.Matches(msg, Keys.Open):
			if u := m.News.SelectedURL(); u != "" {
				return m, openURL(u)
			}
			return m, nil

		// GitHub panel navigation (when focused)
		case m.focused == 3 && key.Matches(msg, Keys.Down):
			m.GitHub.SelectNext()
			return m, nil
		case m.focused == 3 && key.Matches(msg, Keys.Up):
			m.GitHub.SelectPrev()
			return m, nil
		case m.focused == 3 && key.Matches(msg, Keys.Open):
			if u := m.GitHub.SelectedURL(); u != "" {
				return m, openURL(u)
			}
			return m, nil
		}

	case clockTickMsg:
		m.clock = time.Time(msg)
		return m, clockTickCmd()
	}

	// Route all other messages to visible sub-panels
	var cmd tea.Cmd

	if m.Visible[0] {
		m.Weather, cmd = m.Weather.Update(msg)
		cmds = append(cmds, cmd)
	}
	if m.Visible[1] {
		m.Crypto, cmd = m.Crypto.Update(msg)
		cmds = append(cmds, cmd)
	}
	if m.Visible[2] {
		m.News, cmd = m.News.Update(msg)
		cmds = append(cmds, cmd)
	}
	if m.Visible[3] {
		m.GitHub, cmd = m.GitHub.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func openURL(url string) tea.Cmd {
	return func() tea.Msg {
		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "darwin":
			cmd = exec.Command("open", url)
		case "linux":
			cmd = exec.Command("xdg-open", url)
		default:
			cmd = exec.Command("open", url)
		}
		cmd.Run()
		return nil
	}
}
