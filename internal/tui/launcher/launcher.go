// Package launcher provides the root TUI menu shown when llm-wiki is run
// outside of an existing wiki directory.
package launcher

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"

	"github.com/DavDaz/llm-wiki-generator/internal/tui/styles"
)

// Action represents the user's choice from the launcher menu.
type Action int

const (
	ActionNone   Action = iota
	ActionNew           // open the init wizard
	ActionGuide         // print the guide
	ActionAborted       // user cancelled
)

type Model struct {
	form   *huh.Form
	choice string
	done   bool
}

func New() Model {
	m := Model{choice: "new"}
	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("llm-wiki").
				Description("What would you like to do?").
				Options(
					huh.NewOption("Create a new wiki", "new"),
					huh.NewOption("Read the guide", "guide"),
				).
				Value(&m.choice),
		),
	).WithTheme(huh.ThemeCatppuccin())
	return m
}

func (m Model) Init() tea.Cmd {
	return m.form.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "ctrl+c" {
		m.done = true
		m.choice = "abort"
		return m, tea.Quit
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	if m.form.State == huh.StateCompleted || m.form.State == huh.StateAborted {
		m.done = true
		return m, tea.Quit
	}

	return m, cmd
}

func (m Model) View() string {
	if m.done {
		return ""
	}
	return styles.Title.Render("") + m.form.View()
}

// Result returns the action the user selected.
func (m Model) Result() Action {
	if m.form.State == huh.StateAborted || m.choice == "abort" {
		return ActionAborted
	}
	switch m.choice {
	case "new":
		return ActionNew
	case "guide":
		return ActionGuide
	}
	return ActionNone
}
