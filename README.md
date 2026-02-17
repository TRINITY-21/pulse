<p align="center">
  <h1 align="center">⚡ pulse</h1>
  <p align="center">Live terminal dashboard — weather, crypto, news, and GitHub activity in one view.</p>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/BubbleTea-FF75B5?style=flat&logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/Lip%20Gloss-A855F7?style=flat" />
</p>

---

## What It Does

Pulse is a 4-panel TUI dashboard that streams real-time data into your terminal:

| Panel | Source | Refresh |
|-------|--------|---------|
| **Weather** | OpenWeatherMap API | 10 min |
| **Crypto** | CoinGecko API | 30 sec |
| **News** | Hacker News (Firebase) | 5 min |
| **GitHub** | GitHub Events API | 5 min |

Each panel auto-refreshes independently. The clock in the header ticks every second.

## Layout

```
┌──────────────────────────────────────────────────┐
│  ⚡ pulse                              15:04:05  │
├───────────────────────┬──────────────────────────┤
│  ☀ Weather            │  📈 Crypto               │
│  Istanbul             │  BTC  $97,234  +2.3%     │
│  22°C  Partly Cloudy  │     MCap $1.3T · Vol $45B│
│  Feels like 24°C      │  ETH  $3,456   -0.5%     │
│  H: 25° L: 18°       │     MCap $415B · Vol $18B │
│  💧 55%  💨 3.2m/s     │  SOL  $145     +5.1%     │
│  🌅 06:42  🌇 18:15   │     MCap $67B · Vol $3B   │
├───────────────────────┼──────────────────────────┤
│  📰 News              │  🐙 GitHub               │
│  ▸ 1. GrapheneOS...   │  ▸ ★ Pushed to auxcord   │
│     grapheneos.org    │     feat: add shuffle     │
│    2. Four Column...  │    ★ PR merged on popkorn │
│     garbagecollec...  │     Fix edge case in...   │
├──────────────────────────────────────────────────┤
│  q quit · r refresh · tab focus · w c n g toggle │
└──────────────────────────────────────────────────┘
```

Panels adapt dynamically — show 1, 2, 3, or all 4 with automatic layout reflow.

## Install

```bash
git clone https://github.com/TRINITY-21/pulse.git
cd pulse
go mod download
```

## Setup

Copy the example env and add your keys:

```bash
cp .env.example .env
```

```env
OPENWEATHER_API_KEY=your_key_here    # https://openweathermap.org/api
WEATHER_CITY=Istanbul
GITHUB_USERNAME=TRINITY-21
GITHUB_TOKEN=                         # optional, for higher rate limits
CRYPTO_COINS=bitcoin,ethereum,solana
```

Only the OpenWeatherMap key is required. Crypto and News work without any keys. GitHub works without a token but with lower rate limits.

## Usage

```bash
# all panels
go run main.go

# specific panels only
go run main.go --weather --crypto
go run main.go --news --github
go run main.go --crypto
```

## Keybindings

| Key | Action |
|-----|--------|
| `q` / `Ctrl+C` | Quit |
| `r` | Refresh all panels |
| `Tab` | Cycle focus between panels |
| `1` `2` `3` `4` | Jump to specific panel |
| `w` `c` `n` `g` | Toggle weather/crypto/news/github |
| `↑` `↓` / `j` `k` | Navigate items (news/github) |
| `o` / `Enter` | Open selected item in browser |

## Architecture

Built with the Elm Architecture via [BubbleTea](https://github.com/charmbracelet/bubbletea):

```
main.go                    → tea.NewProgram entry point
internal/
  config/config.go         → .env loader → Config struct
  style/style.go           → Lip Gloss styles (purple/cyan theme)
  ui/
    model.go               → Root Model composing 4 sub-panels
    update.go              → Message routing + key handling
    view.go                → Dynamic grid layout (1-4 panels)
    keys.go                → All keybindings
  panels/
    weather/               → OpenWeatherMap (types, fetch, model)
    crypto/                → CoinGecko (types, fetch, model)
    news/                  → Hacker News Firebase (types, fetch, model)
    github/                → GitHub Events API (types, fetch, model)
```

Each panel owns its own refresh cycle, loading state, and error handling. Messages flow through the root `Update()` and get routed to the relevant sub-panel.

## Built With

- [Go](https://go.dev)
- [BubbleTea](https://github.com/charmbracelet/bubbletea) — TUI framework (Elm Architecture)
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — Terminal styling & layout
- [Bubbles](https://github.com/charmbracelet/bubbles) — Spinner component
