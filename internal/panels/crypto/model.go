package crypto

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
	coins       []CoinData
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
	return tea.Tick(30*time.Second, func(t time.Time) tea.Msg {
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
			m.coins = msg.Coins
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
	title := style.TitleStyle.Render("📈 Crypto")

	if m.loading && m.lastUpdated.IsZero() {
		return fmt.Sprintf("%s\n\n  %s Loading prices...", title, m.spinner.View())
	}

	if m.err != nil && m.lastUpdated.IsZero() {
		return fmt.Sprintf("%s\n\n%s", title, style.ErrorStyle.Render("  "+m.err.Error()))
	}

	var lines []string
	lines = append(lines, title)
	lines = append(lines, "")

	for _, coin := range m.coins {
		price := formatPrice(coin.Price)
		change := fmt.Sprintf("%+.1f%%", coin.Change24h)

		changeStyle := style.PositiveStyle
		if coin.Change24h < 0 {
			changeStyle = style.NegativeStyle
		}

		lines = append(lines, fmt.Sprintf("  %s  %s  %s",
			style.BoldWhite.Render(coin.Symbol),
			price,
			changeStyle.Render(change),
		))
		lines = append(lines, fmt.Sprintf("     %s",
			style.SubtitleStyle.Render(fmt.Sprintf("MCap %s · Vol %s",
				formatCompact(coin.MarketCap),
				formatCompact(coin.Volume24h))),
		))
	}

	lines = append(lines, "")
	updated := m.lastUpdated.Format("15:04:05")
	lines = append(lines, style.SubtitleStyle.Render(fmt.Sprintf("  Updated %s", updated)))

	return strings.Join(lines, "\n")
}

func formatCompact(n float64) string {
	switch {
	case n >= 1e12:
		return fmt.Sprintf("$%.1fT", n/1e12)
	case n >= 1e9:
		return fmt.Sprintf("$%.1fB", n/1e9)
	case n >= 1e6:
		return fmt.Sprintf("$%.1fM", n/1e6)
	case n >= 1e3:
		return fmt.Sprintf("$%.0fK", n/1e3)
	default:
		return fmt.Sprintf("$%.0f", n)
	}
}

func formatPrice(price float64) string {
	if price >= 1 {
		whole := int64(price + 0.5)
		s := fmt.Sprintf("%d", whole)
		// Add commas
		if len(s) > 3 {
			var result []byte
			for i, c := range s {
				if i > 0 && (len(s)-i)%3 == 0 {
					result = append(result, ',')
				}
				result = append(result, byte(c))
			}
			return "$" + string(result)
		}
		return "$" + s
	}
	return fmt.Sprintf("$%.4f", price)
}
