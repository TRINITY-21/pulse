package style

import "github.com/charmbracelet/lipgloss"

var (
	Purple = lipgloss.Color("#A855F7")
	Cyan   = lipgloss.Color("#22D3EE")
	Green  = lipgloss.Color("#22C55E")
	Red    = lipgloss.Color("#EF4444")
	Yellow = lipgloss.Color("#EAB308")
	Dim    = lipgloss.Color("#666666")
	White  = lipgloss.Color("#FFFFFF")

	BorderNormal = lipgloss.Color("#333333")
	BorderActive = Purple

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Purple)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(Dim)

	BoldWhite = lipgloss.NewStyle().
			Bold(true).
			Foreground(White)

	AccentStyle = lipgloss.NewStyle().
			Foreground(Cyan)

	ErrorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Red)

	PositiveStyle = lipgloss.NewStyle().
			Foreground(Green)

	NegativeStyle = lipgloss.NewStyle().
			Foreground(Red)

	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderNormal).
			Padding(0, 1)

	PanelActiveStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(BorderActive).
				Padding(0, 1)

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1)

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(Dim).
			Padding(0, 1)
)
