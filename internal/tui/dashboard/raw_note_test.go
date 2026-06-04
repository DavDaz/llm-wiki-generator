package dashboard

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/charmbracelet/huh"
	"github.com/stretchr/testify/require"
)

func TestRawNoteViewRendersForm(t *testing.T) {
	m := NewRawNoteView(filepath.Join(t.TempDir(), "raw")).(*rawNoteModel)

	view := m.View()
	require.Contains(t, view, "New raw note")
	require.Contains(t, view, "raw:")
}

func TestRawNoteEscReturnsToRoot(t *testing.T) {
	m := NewRawNoteView(filepath.Join(t.TempDir(), "raw")).(*rawNoteModel)

	_, cmd := m.Update(keyMsg("esc"))
	require.NotNil(t, cmd)
	_, ok := runCmd(cmd).(BackToRootMsg)
	require.True(t, ok)
}

func TestRawNoteFormAllowsQAsTextInput(t *testing.T) {
	m := NewRawNoteView(filepath.Join(t.TempDir(), "raw")).(*rawNoteModel)

	_, cmd := m.Update(keyMsg("q"))
	if msg := runCmd(cmd); msg != nil {
		_, ok := msg.(BackToRootMsg)
		require.False(t, ok)
	}
}

func TestRawNoteCompletedFormWritesMarkdown(t *testing.T) {
	rawDir := filepath.Join(t.TempDir(), "raw")
	m := NewRawNoteView(rawDir).(*rawNoteModel)
	m.vals.title = "Creación de Usuario"
	m.vals.context = "La plataforma RENAB permite crear usuarios."
	m.vals.sectionTitle = "Pasos"
	m.vals.sectionBody = "1. Entrar en Administración."
	m.vals.notes = "Completar permisos."
	m.vals.save = true
	m.form.State = huh.StateCompleted

	next, cmd := m.Update(nil)
	require.Nil(t, cmd)
	m = next.(*rawNoteModel)
	require.NotEmpty(t, m.savedAt)

	path := filepath.Join(rawDir, "creacion-de-usuario.md")
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Contains(t, string(content), "# Creación de Usuario")
	require.Contains(t, string(content), "## Contexto\nLa plataforma RENAB permite crear usuarios.")
	require.Contains(t, string(content), "## Pasos\n1. Entrar en Administración.")
	require.Contains(t, string(content), "## Notas / Pendientes\nCompletar permisos.")
}

func TestRawNoteDoesNotOverwriteExistingFile(t *testing.T) {
	rawDir := filepath.Join(t.TempDir(), "raw")
	require.NoError(t, os.MkdirAll(rawDir, 0o755))
	path := filepath.Join(rawDir, "sistema-halfman.md")
	require.NoError(t, os.WriteFile(path, []byte("original"), 0o644))

	m := NewRawNoteView(rawDir).(*rawNoteModel)
	m.vals.title = "Sistema HALFMAN"
	m.vals.context = "Nuevo contexto."
	m.vals.save = true
	m.form.State = huh.StateCompleted

	next, cmd := m.Update(nil)
	require.NotNil(t, cmd)
	m = next.(*rawNoteModel)
	require.Empty(t, m.savedAt)
	require.Contains(t, m.errMsg, "already exists")

	content, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, "original", string(content))
}

func TestRawNoteRejectsUnsafeFileName(t *testing.T) {
	m := NewRawNoteView(filepath.Join(t.TempDir(), "raw")).(*rawNoteModel)
	m.vals.title = "Unsafe"
	m.vals.slug = "../unsafe"
	m.vals.context = "Contexto."
	m.vals.save = true
	m.form.State = huh.StateCompleted

	next, cmd := m.Update(nil)
	require.NotNil(t, cmd)
	m = next.(*rawNoteModel)
	require.Empty(t, m.savedAt)
	require.Contains(t, m.errMsg, "file name must not contain paths")
}

func TestRawNoteSavedStateCanStartAnotherNote(t *testing.T) {
	m := NewRawNoteView(filepath.Join(t.TempDir(), "raw")).(*rawNoteModel)
	m.savedAt = filepath.Join(m.rawDir, "one.md")

	next, cmd := m.Update(keyMsg("n"))
	require.NotNil(t, cmd)
	m = next.(*rawNoteModel)
	require.Empty(t, m.savedAt)
	require.Equal(t, "Detalles", m.vals.sectionTitle)
}
