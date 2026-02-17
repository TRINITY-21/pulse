package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"pulse/internal/style"
)

type panelEntry struct {
	index   int
	content string
}

func (m Model) View() string {
	if m.width == 0 {
		return "  Loading..."
	}

	// Header
	title := style.TitleStyle.Render("⚡ pulse")
	clock := style.AccentStyle.Render(m.clock.Format("15:04:05"))
	gap := m.width - lipgloss.Width(title) - lipgloss.Width(clock) - 4
	if gap < 1 {
		gap = 1
	}
	header := style.HeaderStyle.Render(
		fmt.Sprintf("%s%s%s", title, strings.Repeat(" ", gap), clock),
	)

	// Collect visible panels
	count := m.visibleCount()
	if count == 0 {
		statusBar := style.StatusBarStyle.Render(
			"q quit  ·  w/c/n/g toggle panels",
		)
		return lipgloss.JoinVertical(lipgloss.Left, header, "\n  No panels visible. Press w/c/n/g to toggle.", statusBar)
	}

	// Calculate dimensions based on layout
	var grid string
	switch count {
	case 1:
		grid = m.layoutOne()
	case 2:
		grid = m.layoutTwo()
	case 3:
		grid = m.layoutThree()
	default:
		grid = m.layoutFour()
	}

	// Status bar — show toggle indicators
	var indicators []string
	labels := []string{"w:weather", "c:crypto", "n:news", "g:github"}
	for i, label := range labels {
		if m.Visible[i] {
			indicators = append(indicators, style.AccentStyle.Render(label))
		} else {
			indicators = append(indicators, style.SubtitleStyle.Render(label))
		}
	}

	statusBar := style.StatusBarStyle.Render(
		fmt.Sprintf("q quit  ·  r refresh  ·  tab focus  ·  %s", strings.Join(indicators, "  ")),
	)

	return lipgloss.JoinVertical(lipgloss.Left, header, grid, statusBar)
}

func (m Model) layoutOne() string {
	pw := m.width - 2
	ph := m.height - 3
	if ph < 5 {
		ph = 5
	}
	cw := pw - 4
	ch := ph - 3

	panels := m.getVisiblePanels(cw, ch)
	return m.renderPanel(panels[0].index, panels[0].content, pw, ph)
}

func (m Model) layoutTwo() string {
	pw := m.width/2 - 1
	ph := m.height - 3
	if ph < 5 {
		ph = 5
	}
	cw := pw - 4
	ch := ph - 3

	panels := m.getVisiblePanels(cw, ch)
	left := m.renderPanel(panels[0].index, panels[0].content, pw, ph)
	right := m.renderPanel(panels[1].index, panels[1].content, pw, ph)
	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}

func (m Model) layoutThree() string {
	fullW := m.width - 2
	halfW := m.width/2 - 1
	halfH := (m.height - 3) / 2
	if halfH < 5 {
		halfH = 5
	}

	panels := m.getVisiblePanels(halfW-4, halfH-3)

	top := m.renderPanel(panels[0].index, panels[0].content, fullW, halfH)
	bl := m.renderPanel(panels[1].index, panels[1].content, halfW, halfH)
	br := m.renderPanel(panels[2].index, panels[2].content, halfW, halfH)
	bottom := lipgloss.JoinHorizontal(lipgloss.Top, bl, br)
	return lipgloss.JoinVertical(lipgloss.Left, top, bottom)
}

func (m Model) layoutFour() string {
	pw := m.width/2 - 1
	ph := (m.height - 3) / 2
	if ph < 5 {
		ph = 5
	}
	cw := pw - 4
	ch := ph - 3
	if cw < 10 {
		cw = 10
	}
	if ch < 3 {
		ch = 3
	}

	panels := m.getVisiblePanels(cw, ch)

	tl := m.renderPanel(panels[0].index, panels[0].content, pw, ph)
	tr := m.renderPanel(panels[1].index, panels[1].content, pw, ph)
	bl := m.renderPanel(panels[2].index, panels[2].content, pw, ph)
	br := m.renderPanel(panels[3].index, panels[3].content, pw, ph)

	top := lipgloss.JoinHorizontal(lipgloss.Top, tl, tr)
	bottom := lipgloss.JoinHorizontal(lipgloss.Top, bl, br)
	return lipgloss.JoinVertical(lipgloss.Left, top, bottom)
}

func (m Model) getVisiblePanels(contentWidth, contentHeight int) []panelEntry {
	var panels []panelEntry

	type panelDef struct {
		index int
		view  func(int, int) string
	}

	all := []panelDef{
		{0, m.Weather.View},
		{1, m.Crypto.View},
		{2, m.News.View},
		{3, m.GitHub.View},
	}

	for _, p := range all {
		if m.Visible[p.index] {
			panels = append(panels, panelEntry{
				index:   p.index,
				content: p.view(contentWidth, contentHeight),
			})
		}
	}

	return panels
}

func (m Model) renderPanel(index int, content string, width, height int) string {
	s := style.PanelStyle
	if m.focused == index {
		s = style.PanelActiveStyle
	}
	return s.
		Width(width - 2).
		Height(height - 2).
		Render(content)
}
