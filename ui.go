package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const maxWidth = 20

var (
	highlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	selectedStyle  = lipgloss.NewStyle().Bold(true)
)

type model struct {
	initialized bool
	width       int
	height      int

	choices   configs
	highlight int
	selected  int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.initialized = true
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.selected != -1 {
				// Disable current VPN connection if one is active.
				disable(m.choices[m.selected])
			}
			if m.selected == m.highlight {
				// Highlighted VPN is already active.
				// Only disconnect.
				m.selected = -1
				break
			}

			enable(m.choices[m.highlight])
			m.selected = m.highlight
		case "j":
			m.highlight += 1
			if m.highlight >= len(m.choices) {
				m.highlight = 0
			}
		case "k":
			m.highlight -= 1
			if m.highlight < 0 {
				m.highlight = len(m.choices) - 1
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	if !m.initialized {
		return ""
	}

	height := m.height - 2

	page := m.highlight / height
	start := page * height
	end := page*height + height
	if end > len(m.choices) {
		end = len(m.choices)
	}

	header := "Disconnected"
	if m.selected != -1 {
		header = "Connected"
	}

	lines := []string{
		header,
		"",
	}
	visible := m.choices[start:end]
	for i, c := range visible {
		if i == height {
			break
		}
		line := c.name
		if i+page*height == m.selected {
			line = selectedStyle.Render(line)
		}
		if i+page*height == m.highlight {
			line = highlightStyle.Render(line)
		}
		lines = append(lines, line)
	}
	return lipgloss.JoinVertical(lipgloss.Center, lines...)
}
