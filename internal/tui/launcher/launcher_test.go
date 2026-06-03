package launcher

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/stretchr/testify/require"
)

func launcherKeyMsg(s string) tea.KeyMsg {
	switch s {
	case "esc":
		return tea.KeyMsg{Type: tea.KeyEsc}
	case "ctrl+c":
		return tea.KeyMsg{Type: tea.KeyCtrlC}
	case "q":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
	default:
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
	}
}

func runLauncherCmd(cmd tea.Cmd) tea.Msg {
	if cmd == nil {
		return nil
	}
	return cmd()
}

func TestLauncherExitChoiceReturnsExitAction(t *testing.T) {
	m := New()
	m.vals.choice = "exit"
	m.form.State = huh.StateCompleted

	updated, cmd := m.Update(nil)
	lm, ok := updated.(Model)
	require.True(t, ok)
	require.NotNil(t, cmd)
	require.IsType(t, tea.QuitMsg{}, runLauncherCmd(cmd))
	require.Equal(t, ActionExit, lm.Result())
}

func TestLauncherEscAndQQuitCleanly(t *testing.T) {
	for _, key := range []string{"esc", "q"} {
		t.Run(key, func(t *testing.T) {
			m := New()
			updated, cmd := m.Update(launcherKeyMsg(key))
			lm, ok := updated.(Model)
			require.True(t, ok)
			require.NotNil(t, cmd)
			require.IsType(t, tea.QuitMsg{}, runLauncherCmd(cmd))
			require.Equal(t, ActionExit, lm.Result())
		})
	}
}

func TestLauncherCtrlCAborts(t *testing.T) {
	m := New()
	updated, cmd := m.Update(launcherKeyMsg("ctrl+c"))
	lm, ok := updated.(Model)
	require.True(t, ok)
	require.NotNil(t, cmd)
	require.IsType(t, tea.QuitMsg{}, runLauncherCmd(cmd))
	require.Equal(t, ActionAborted, lm.Result())
}
