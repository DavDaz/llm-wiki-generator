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

// values holds mutable form state behind a pointer so it survives
// Bubbletea's copy semantics when the Model is passed by value.
type values struct {
	choice string
}

type Model struct {
	form    *huh.Form
	vals    *values // pointer — valid across Bubbletea copies
	aborted bool
	done    bool
}

func New() Model {
	v := &values{choice: "new"}
	return Model{
		vals: v,
		form: buildForm(v),
	}
}

func buildForm(v *values) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("llm-wiki").
				Description("What would you like to do?").
				Options(
					huh.NewOption("Create a new wiki", "new"),
					huh.NewOption("Read the guide", "guide"),
				).
				Value(&v.choice),
		),
	).WithTheme(huh.ThemeCatppuccin())
}

func (m Model) Init() tea.Cmd {
	return m.form.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "ctrl+c" {
		m.aborted = true
		m.done = true
		return m, tea.Quit
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	if m.form.State == huh.StateCompleted {
		m.done = true
		return m, tea.Quit
	}
	if m.form.State == huh.StateAborted {
		m.aborted = true
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
	if m.aborted {
		return ActionAborted
	}
	switch m.vals.choice {
	case "new":
		return ActionNew
	case "guide":
		return ActionGuide
	}
	return ActionNone
}
