// Package styles defines shared Lipgloss styles for all TUI components.
package styles

import "github.com/charmbracelet/lipgloss"

var (
	Primary   = lipgloss.NewStyle().Foreground(lipgloss.Color("99"))  // purple
	Muted     = lipgloss.NewStyle().Foreground(lipgloss.Color("240")) // gray
	Success   = lipgloss.NewStyle().Foreground(lipgloss.Color("76"))  // green
	Warning   = lipgloss.NewStyle().Foreground(lipgloss.Color("214")) // orange
	Bold      = lipgloss.NewStyle().Bold(true)
	Faint     = lipgloss.NewStyle().Faint(true)

	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("99")).
		MarginBottom(1)

	Box = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("99")).
		Padding(1, 2)

	KeyHint = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		MarginTop(1)
)
