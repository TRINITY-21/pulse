package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"pulse/internal/config"
	"pulse/internal/ui"
)

func main() {
	showWeather := flag.Bool("weather", false, "show weather panel")
	showCrypto := flag.Bool("crypto", false, "show crypto panel")
	showNews := flag.Bool("news", false, "show news panel")
	showGitHub := flag.Bool("github", false, "show github panel")
	flag.Parse()

	cfg := config.Load()

	// If no flags specified, show all panels
	anySet := *showWeather || *showCrypto || *showNews || *showGitHub
	visible := [4]bool{
		!anySet || *showWeather,
		!anySet || *showCrypto,
		!anySet || *showNews,
		!anySet || *showGitHub,
	}

	m := ui.NewModel(cfg, visible)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
