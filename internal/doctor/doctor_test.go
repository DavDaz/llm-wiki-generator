package doctor_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/DavDaz/llm-wiki-generator/internal/doctor"
	"github.com/DavDaz/llm-wiki-generator/internal/generator"
	"github.com/DavDaz/llm-wiki-generator/internal/manifest"
)

func TestRunReportsHealthyWiki(t *testing.T) {
	root := initWiki(t)
	writeFile(t, filepath.Join(root, "wiki", "alpha.md"), "# Alpha\n\nVer [[beta]].\n")
	writeFile(t, filepath.Join(root, "wiki", "beta.md"), "# Beta\n")
	appendFile(t, filepath.Join(root, "wiki", "index.md"), "| [[alpha]] | Alpha | referencia | vigente | 2026-07-07 |\n| [[beta]] | Beta | referencia | vigente | 2026-07-07 |\n")

	report, err := doctor.Run(root)
	require.NoError(t, err)

	assert.Equal(t, doctor.SeverityPass, severityFor(t, report, "manifest"))
	assert.Equal(t, doctor.SeverityPass, severityFor(t, report, "core-files"))
	assert.Equal(t, doctor.SeverityPass, severityFor(t, report, "tool-outputs"))
	assert.Equal(t, doctor.SeverityPass, severityFor(t, report, "wikilinks"))
	assert.Equal(t, doctor.SeverityPass, severityFor(t, report, "index-coverage"))
}

func TestRunReportsFailuresWithoutWriting(t *testing.T) {
	root := initWiki(t)
	writeFile(t, filepath.Join(root, "wiki", "orphan.md"), "# Orphan\n\nBroken [[missing-page]].\n")
	writeFile(t, filepath.Join(root, "wiki", "covered.md"), "# Covered\n")
	appendFile(t, filepath.Join(root, "wiki", "index.md"), "| [[covered]] | Covered | referencia | vigente | 2026-07-07 |\n")
	require.NoError(t, os.Remove(filepath.Join(root, "wiki", "log.md")))
	require.NoError(t, os.Remove(filepath.Join(root, ".claude", "skills", "wiki-lint", "SKILL.md")))
	invalidateManifest(t, root)

	before := snapshotTree(t, root)
	report, err := doctor.Run(root)
	after := snapshotTree(t, root)
	require.NoError(t, err)
	assert.Equal(t, before, after)

	assert.Equal(t, doctor.SeverityError, severityFor(t, report, "manifest"))
	assert.Equal(t, doctor.SeverityError, severityFor(t, report, "core-files"))
	assert.Equal(t, doctor.SeverityError, severityFor(t, report, "tool-outputs"))
	assert.Equal(t, doctor.SeverityError, severityFor(t, report, "wikilinks"))
	assert.Equal(t, doctor.SeverityWarn, severityFor(t, report, "index-coverage"))
	assert.Contains(t, messageFor(t, report, "core-files"), "wiki/log.md")
	assert.Contains(t, messageFor(t, report, "tool-outputs"), ".claude/skills/wiki-lint/SKILL.md")
	assert.Contains(t, messageFor(t, report, "wikilinks"), "missing-page")
	assert.Contains(t, messageFor(t, report, "index-coverage"), "orphan")
}

func TestRunReportsDirectoryPlacedAtRequiredCoreFilePath(t *testing.T) {
	requiredFiles := []string{
		filepath.Join("wiki", "index.md"),
		filepath.Join("wiki", "log.md"),
		filepath.Join("wiki", "sources.json"),
	}

	for _, relPath := range requiredFiles {
		t.Run(relPath, func(t *testing.T) {
			root := initWiki(t)
			absPath := filepath.Join(root, relPath)

			require.NoError(t, os.Remove(absPath))
			require.NoError(t, os.Mkdir(absPath, 0o755))

			report, err := doctor.Run(root)
			require.NoError(t, err)

			assert.Equal(t, doctor.SeverityError, severityFor(t, report, "core-files"))
			assert.Contains(t, messageFor(t, report, "core-files"), relPath)
		})
	}
}

func TestRunReportsDirectoryPlacedAtToolOutputFilePath(t *testing.T) {
	toolFiles := []string{
		"CLAUDE.md",
		filepath.Join(".claude", "skills", "wiki-lint", "SKILL.md"),
	}

	for _, relPath := range toolFiles {
		t.Run(relPath, func(t *testing.T) {
			root := initWiki(t)
			absPath := filepath.Join(root, relPath)

			require.NoError(t, os.Remove(absPath))
			require.NoError(t, os.Mkdir(absPath, 0o755))

			report, err := doctor.Run(root)
			require.NoError(t, err)

			assert.Equal(t, doctor.SeverityError, severityFor(t, report, "tool-outputs"))
			assert.Contains(t, messageFor(t, report, "tool-outputs"), relPath)
			assert.Contains(t, messageFor(t, report, "tool-outputs"), "expected file")
		})
	}
}

func initWiki(t *testing.T) string {
	t.Helper()
	root, err := generator.Init(generator.InitConfig{
		ParentDir:  t.TempDir(),
		Name:       "Legal Wiki",
		Slug:       "legal-wiki",
		Language:   "es",
		ClaudeCode: true,
	})
	require.NoError(t, err)
	return root
}

func invalidateManifest(t *testing.T, root string) {
	t.Helper()
	m, err := manifest.Load(root)
	require.NoError(t, err)
	m.Wiki.Slug = "not valid"
	require.NoError(t, m.Save(root))
}

func snapshotTree(t *testing.T, root string) map[string]string {
	t.Helper()
	snapshot := map[string]string{}
	require.NoError(t, filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		require.NoError(t, err)
		rel, relErr := filepath.Rel(root, path)
		require.NoError(t, relErr)
		if d.IsDir() {
			snapshot[rel+"/"] = "dir"
			return nil
		}
		buf, readErr := os.ReadFile(path)
		require.NoError(t, readErr)
		snapshot[rel] = string(buf)
		return nil
	}))
	return snapshot
}

func writeFile(t *testing.T, path string, content string) {
	t.Helper()
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
}

func appendFile(t *testing.T, path string, content string) {
	t.Helper()
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
	require.NoError(t, err)
	defer file.Close()
	_, err = file.WriteString(content)
	require.NoError(t, err)
}

func severityFor(t *testing.T, report doctor.Report, check string) doctor.Severity {
	t.Helper()
	for _, finding := range report.Findings {
		if finding.Check == check {
			return finding.Severity
		}
	}
	t.Fatalf("missing finding for %s", check)
	return ""
}

func messageFor(t *testing.T, report doctor.Report, check string) string {
	t.Helper()
	for _, finding := range report.Findings {
		if finding.Check == check {
			return finding.Message
		}
	}
	t.Fatalf("missing finding for %s", check)
	return ""
}
