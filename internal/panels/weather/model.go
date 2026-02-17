package weather

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
	data        Data
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
	return tea.Tick(10*time.Minute, func(t time.Time) tea.Msg {
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
			m.data = msg.Data
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

func (m Model) View(width, height int) string {
	if m.loading && m.lastUpdated.IsZero() {
		title := style.TitleStyle.Render("🌤 Weather")
		return fmt.Sprintf("%s\n\n  %s Loading weather...", title, m.spinner.View())
	}

	if m.err != nil && m.lastUpdated.IsZero() {
		title := style.TitleStyle.Render("🌤 Weather")
		return fmt.Sprintf("%s\n\n%s", title, style.ErrorStyle.Render("  "+m.err.Error()))
	}

	d := m.data
	icon := weatherIcon(d.Condition)
	title := style.TitleStyle.Render(fmt.Sprintf("%s Weather", icon))

	sunrise := time.Unix(d.Sunrise, 0).Format("15:04")
	sunset := time.Unix(d.Sunset, 0).Format("15:04")
	updated := m.lastUpdated.Format("15:04")

	var lines []string
	lines = append(lines, title)
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("  %s %s",
		style.BoldWhite.Render(d.City), loadingDot(m.loading)))
	lines = append(lines, fmt.Sprintf("  %s  %s",
		style.BoldWhite.Render(fmt.Sprintf("%.0f°C", d.Temp)),
		style.SubtitleStyle.Render(d.Description)))
	lines = append(lines, fmt.Sprintf("  %s",
		style.SubtitleStyle.Render(fmt.Sprintf("Feels like %.0f°C", d.FeelsLike))))
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("  %s  %s  %s  %s",
		style.SubtitleStyle.Render(fmt.Sprintf("H: %.0f°", d.TempMax)),
		style.SubtitleStyle.Render(fmt.Sprintf("L: %.0f°", d.TempMin)),
		style.SubtitleStyle.Render(fmt.Sprintf("💧 %d%%", d.Humidity)),
		style.SubtitleStyle.Render(fmt.Sprintf("💨 %.1fm/s", d.WindSpeed))))
	lines = append(lines, fmt.Sprintf("  %s  %s  %s",
		style.SubtitleStyle.Render(fmt.Sprintf("☁ %d%%", d.Clouds)),
		style.SubtitleStyle.Render(fmt.Sprintf("👁 %.0fkm", d.Visibility)),
		style.SubtitleStyle.Render(fmt.Sprintf("%dhPa", d.Pressure))))
	lines = append(lines, fmt.Sprintf("  %s  %s",
		style.SubtitleStyle.Render(fmt.Sprintf("🌅 %s", sunrise)),
		style.SubtitleStyle.Render(fmt.Sprintf("🌇 %s", sunset))))
	lines = append(lines, "")
	lines = append(lines, style.SubtitleStyle.Render(fmt.Sprintf("  Updated %s", updated)))

	return strings.Join(lines, "\n")
}

func loadingDot(loading bool) string {
	if loading {
		return style.AccentStyle.Render("*")
	}
	return ""
}

func weatherIcon(condition string) string {
	switch condition {
	case "Clear":
		return "☀"
	case "Clouds":
		return "☁"
	case "Rain", "Drizzle":
		return "🌧"
	case "Snow":
		return "❄"
	case "Thunderstorm":
		return "⛈"
	default:
		return "🌤"
	}
}
