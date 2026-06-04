package rawnotes

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// Note is a manually-authored source note saved under a wiki raw/ directory.
type Note struct {
	Title        string
	Slug         string
	Context      string
	SectionTitle string
	SectionBody  string
	Notes        string
}

// Render returns the markdown representation for a manual raw note.
func Render(n Note) string {
	var b strings.Builder
	title := strings.TrimSpace(n.Title)
	sectionTitle := strings.TrimSpace(n.SectionTitle)
	if sectionTitle == "" {
		sectionTitle = "Detalles"
	}

	b.WriteString("# ")
	b.WriteString(title)
	b.WriteString("\n\n")

	b.WriteString("## Contexto\n")
	b.WriteString(strings.TrimSpace(n.Context))
	b.WriteString("\n")

	if body := strings.TrimSpace(n.SectionBody); body != "" {
		b.WriteString("\n## ")
		b.WriteString(sectionTitle)
		b.WriteString("\n")
		b.WriteString(body)
		b.WriteString("\n")
	}

	if notes := strings.TrimSpace(n.Notes); notes != "" {
		b.WriteString("\n## Notas / Pendientes\n")
		b.WriteString(notes)
		b.WriteString("\n")
	}

	return b.String()
}

// Write creates a markdown file for n under rawDir without overwriting an existing file.
func Write(rawDir string, n Note) (string, error) {
	if strings.TrimSpace(n.Title) == "" {
		return "", errors.New("title is required")
	}

	slug := strings.TrimSpace(n.Slug)
	if slug == "" {
		slug = Slugify(n.Title)
	}
	if err := ValidateSlug(slug); err != nil {
		return "", err
	}

	rawAbs, err := filepath.Abs(rawDir)
	if err != nil {
		return "", fmt.Errorf("resolve raw directory: %w", err)
	}
	if err := os.MkdirAll(rawAbs, 0o755); err != nil {
		return "", fmt.Errorf("create raw directory: %w", err)
	}

	path := filepath.Join(rawAbs, slug+".md")
	if filepath.Dir(path) != rawAbs {
		return "", errors.New("note path must stay inside raw directory")
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return "", fmt.Errorf("raw note already exists: %s", filepath.Base(path))
		}
		return "", fmt.Errorf("create raw note: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString(Render(n)); err != nil {
		return "", fmt.Errorf("write raw note: %w", err)
	}

	return path, nil
}

// Slugify converts a title into a safe raw note file slug.
func Slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(replaceSpanishAccents(s)))
	var b strings.Builder
	lastDash := false
	for _, r := range s {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
			lastDash = false
		case r == '-' || r == '_' || unicode.IsSpace(r):
			if !lastDash && b.Len() > 0 {
				b.WriteRune('-')
				lastDash = true
			}
		default:
			if !lastDash && b.Len() > 0 {
				b.WriteRune('-')
				lastDash = true
			}
		}
	}
	return strings.Trim(b.String(), "-")
}

// ValidateSlug rejects path-like or empty note slugs.
func ValidateSlug(slug string) error {
	if slug == "" {
		return errors.New("file name is required")
	}
	if slug != filepath.Base(slug) || strings.Contains(slug, "..") || strings.ContainsAny(slug, `/\\`) {
		return errors.New("file name must not contain paths")
	}
	if strings.HasSuffix(slug, ".md") {
		return errors.New("file name should omit .md")
	}
	for _, r := range slug {
		if !(unicode.IsLower(r) || unicode.IsDigit(r) || r == '-') {
			return errors.New("file name must use lowercase letters, numbers, and hyphens")
		}
	}
	return nil
}

func replaceSpanishAccents(s string) string {
	replacer := strings.NewReplacer(
		"á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u", "ü", "u", "ñ", "n",
		"Á", "a", "É", "e", "Í", "i", "Ó", "o", "Ú", "u", "Ü", "u", "Ñ", "n",
	)
	return replacer.Replace(s)
}
