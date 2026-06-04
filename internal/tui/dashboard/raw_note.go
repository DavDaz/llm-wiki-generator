package dashboard

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"

	"github.com/DavDaz/llm-wiki-generator/internal/rawnotes"
	"github.com/DavDaz/llm-wiki-generator/internal/tui/styles"
)

type rawNoteValues struct {
	title        string
	slug         string
	context      string
	sectionTitle string
	sectionBody  string
	notes        string
	save         bool
}

type rawNoteModel struct {
	rawDir  string
	form    *huh.Form
	vals    *rawNoteValues
	errMsg  string
	savedAt string
}

// NewRawNoteView creates a form for writing a manual markdown source note into raw/.
func NewRawNoteView(rawDir string) tea.Model {
	vals := &rawNoteValues{sectionTitle: "Detalles", save: true}
	return &rawNoteModel{
		rawDir: rawDir,
		form:   newRawNoteForm(vals),
		vals:   vals,
	}
}

func newRawNoteForm(v *rawNoteValues) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Title").
				Description("Used as the # heading.").
				Placeholder("Sistema HALFMAN").
				Value(&v.title).
				Validate(requiredText("title is required")),
			huh.NewInput().
				Title("File name (optional)").
				Description("Without .md. Leave empty to generate it from the title.").
				Placeholder("sistema-halfman").
				Value(&v.slug).
				Validate(optionalSlug),
			huh.NewText().
				Title("Contexto").
				Description("What this note covers and why it exists.").
				Placeholder("Describe the topic...").
				Value(&v.context).
				Validate(requiredText("context is required")).
				WithHeight(5),
			huh.NewInput().
				Title("Main section title").
				Description("Example: Pasos, Hallazgos, Proceso.").
				Value(&v.sectionTitle),
			huh.NewText().
				Title("Main section body").
				Description("Write the first content section. You can edit the file later.").
				Placeholder("1. Primer paso...").
				Value(&v.sectionBody).
				WithHeight(7),
			huh.NewText().
				Title("Notas / Pendientes").
				Description("Optional uncertainties, reminders, or incomplete details.").
				Value(&v.notes).
				WithHeight(4),
			huh.NewConfirm().
				Title("Save raw note?").
				Affirmative("Save").
				Negative("Cancel").
				Value(&v.save),
		),
	).WithTheme(huh.ThemeCatppuccin())
}

func requiredText(message string) func(string) error {
	return func(value string) error {
		if strings.TrimSpace(value) == "" {
			return errors.New(message)
		}
		return nil
	}
}

func optionalSlug(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return rawnotes.ValidateSlug(value)
}

func (m *rawNoteModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *rawNoteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.savedAt != "" {
		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case "esc", "q":
				return m, BackToRoot()
			case "n":
				m.vals = &rawNoteValues{sectionTitle: "Detalles", save: true}
				m.form = newRawNoteForm(m.vals)
				m.errMsg = ""
				m.savedAt = ""
				return m, m.form.Init()
			}
		}
		return m, nil
	}

	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "esc" {
		return m, BackToRoot()
	}

	f, cmd := m.form.Update(msg)
	if updated, ok := f.(*huh.Form); ok {
		m.form = updated
	}

	if m.form.State == huh.StateCompleted {
		if !m.vals.save {
			return m, BackToRoot()
		}

		path, err := rawnotes.Write(m.rawDir, rawnotes.Note{
			Title:        m.vals.title,
			Slug:         m.vals.slug,
			Context:      m.vals.context,
			SectionTitle: m.vals.sectionTitle,
			SectionBody:  m.vals.sectionBody,
			Notes:        m.vals.notes,
		})
		if err != nil {
			m.errMsg = err.Error()
			m.form = newRawNoteForm(m.vals)
			return m, m.form.Init()
		}
		m.savedAt = path
		m.errMsg = ""
		return m, nil
	}

	if m.form.State == huh.StateAborted {
		return m, BackToRoot()
	}

	return m, cmd
}

func (m *rawNoteModel) View() string {
	var b strings.Builder
	b.WriteString(styles.Title.Render("New raw note"))
	b.WriteString("\n")
	b.WriteString(styles.Muted.Render(fmt.Sprintf("raw: %s", m.rawDir)))
	b.WriteString("\n\n")

	if m.errMsg != "" {
		b.WriteString(styles.Warning.Render("  ✗ " + m.errMsg))
		b.WriteString("\n\n")
	}

	if m.savedAt != "" {
		rel, err := filepath.Rel(m.rawDir, m.savedAt)
		if err != nil {
			rel = m.savedAt
		}
		b.WriteString(styles.Bold.Render("  Saved raw note: "))
		b.WriteString(rel)
		b.WriteString("\n")
		b.WriteString(styles.Muted.Render("  It will be processed later by wiki ingest."))
		b.WriteString("\n")
		b.WriteString(styles.KeyHint.Render("\n  [n] new note  [esc] back  [ctrl+c] quit"))
		b.WriteString("\n")
		return b.String()
	}

	b.WriteString(m.form.View())
	b.WriteString(styles.KeyHint.Render("\n  [esc] back  [ctrl+c] quit"))
	b.WriteString("\n")
	return b.String()
}
