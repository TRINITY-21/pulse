package ui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Quit          key.Binding
	Refresh       key.Binding
	Tab           key.Binding
	Panel1        key.Binding
	Panel2        key.Binding
	Panel3        key.Binding
	Panel4        key.Binding
	ToggleWeather key.Binding
	ToggleCrypto  key.Binding
	ToggleNews    key.Binding
	ToggleGitHub  key.Binding
	Up            key.Binding
	Down          key.Binding
	Open          key.Binding
}

var Keys = keyMap{
	Quit:          key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	Refresh:       key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "refresh")),
	Tab:           key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "focus")),
	Panel1:        key.NewBinding(key.WithKeys("1"), key.WithHelp("1", "weather")),
	Panel2:        key.NewBinding(key.WithKeys("2"), key.WithHelp("2", "crypto")),
	Panel3:        key.NewBinding(key.WithKeys("3"), key.WithHelp("3", "news")),
	Panel4:        key.NewBinding(key.WithKeys("4"), key.WithHelp("4", "github")),
	ToggleWeather: key.NewBinding(key.WithKeys("w"), key.WithHelp("w", "toggle weather")),
	ToggleCrypto:  key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "toggle crypto")),
	ToggleNews:    key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "toggle news")),
	ToggleGitHub:  key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "toggle github")),
	Up:            key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
	Down:          key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
	Open:          key.NewBinding(key.WithKeys("o", "enter"), key.WithHelp("o", "open")),
}
