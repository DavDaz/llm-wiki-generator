package rawnotes

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSlugify(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "simple title", in: "Creación de Usuario", want: "creacion-de-usuario"},
		{name: "punctuation", in: "Bug RENAB_IONIC 9 de enero", want: "bug-renab-ionic-9-de-enero"},
		{name: "trim dashes", in: " ¿Sistema HALFMAN? ", want: "sistema-halfman"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, Slugify(tt.in))
		})
	}
}

func TestValidateSlugRejectsUnsafeNames(t *testing.T) {
	tests := []string{"", "../note", "folder/note", "note.md", "Note", "bad_name"}
	for _, slug := range tests {
		t.Run(slug, func(t *testing.T) {
			require.Error(t, ValidateSlug(slug))
		})
	}
}

func TestRenderIncludesManualRawSections(t *testing.T) {
	got := Render(Note{
		Title:        "Creacion de Usuarios",
		Context:      "La plataforma RENAB permite crear usuarios.",
		SectionTitle: "Pasos",
		SectionBody:  "1. Abrir Administración.",
		Notes:        "Completar permisos.",
	})

	require.Contains(t, got, "# Creacion de Usuarios\n\n")
	require.Contains(t, got, "## Contexto\nLa plataforma RENAB permite crear usuarios.\n")
	require.Contains(t, got, "## Pasos\n1. Abrir Administración.\n")
	require.Contains(t, got, "## Notas / Pendientes\nCompletar permisos.\n")
}

func TestRenderOmitsEmptyOptionalSections(t *testing.T) {
	got := Render(Note{Title: "Sistema HALFMAN", Context: "Log intermediario."})

	require.Contains(t, got, "# Sistema HALFMAN")
	require.Contains(t, got, "## Contexto")
	require.NotContains(t, got, "## Detalles")
	require.NotContains(t, got, "## Notas / Pendientes")
}

func TestWriteCreatesNoteWithoutOverwriting(t *testing.T) {
	rawDir := filepath.Join(t.TempDir(), "raw")

	path, err := Write(rawDir, Note{
		Title:       "Creación de Usuario",
		Context:     "Contexto inicial.",
		SectionBody: "Pasos iniciales.",
	})
	require.NoError(t, err)
	require.Equal(t, filepath.Join(rawDir, "creacion-de-usuario.md"), path)

	content, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Contains(t, string(content), "# Creación de Usuario")
	require.Contains(t, string(content), "## Detalles\nPasos iniciales.")

	_, err = Write(rawDir, Note{Title: "Creación de Usuario", Context: "Nuevo contenido."})
	require.Error(t, err)

	contentAfter, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, content, contentAfter)
}
