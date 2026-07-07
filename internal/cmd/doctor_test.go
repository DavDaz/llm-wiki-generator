package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/DavDaz/llm-wiki-generator/internal/generator"
)

func TestDoctorCommandIsRegistered(t *testing.T) {
	command, _, err := rootCmd.Find([]string{"doctor"})
	require.NoError(t, err)
	require.Same(t, doctorCmd, command)
}

func TestRunDoctorInsideWikiPrintsReport(t *testing.T) {
	originalWD, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, os.Chdir(originalWD))
	})

	wikiRoot := initDoctorWiki(t)
	require.NoError(t, os.Chdir(wikiRoot))

	command := &cobra.Command{}
	var out bytes.Buffer
	command.SetOut(&out)

	err = runDoctor(command, nil)
	require.NoError(t, err)
	require.Contains(t, out.String(), "Doctor report for ")
	require.Contains(t, out.String(), "PASS manifest [wiki.toml]: manifest loads and validates")
	require.Contains(t, out.String(), "PASS core-files [wiki]: required wiki directories and core files are present")
	require.Contains(t, out.String(), "Summary: 5 passed, 0 warning(s), 0 error(s)")
}

func TestRunDoctorPrintsAffectedPathForFailures(t *testing.T) {
	originalWD, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, os.Chdir(originalWD))
	})

	wikiRoot := initDoctorWiki(t)
	require.NoError(t, os.Remove(filepath.Join(wikiRoot, "wiki", "log.md")))
	require.NoError(t, os.Chdir(wikiRoot))

	command := &cobra.Command{}
	var out bytes.Buffer
	command.SetOut(&out)

	err = runDoctor(command, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "doctor found 1 error(s)")
	require.Contains(t, out.String(), "ERROR core-files [")
	require.Contains(t, out.String(), "wiki/log.md")
}

func TestRunDoctorOutsideWikiReturnsClearError(t *testing.T) {
	originalWD, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, os.Chdir(originalWD))
	})

	require.NoError(t, os.Chdir(t.TempDir()))

	command := &cobra.Command{}
	var out bytes.Buffer
	command.SetOut(&out)

	err = runDoctor(command, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no wiki found in current directory")
	require.Empty(t, out.String())
}

func initDoctorWiki(t *testing.T) string {
	t.Helper()

	root, err := generator.Init(generator.InitConfig{
		ParentDir:  t.TempDir(),
		Name:       "Legal Wiki",
		Slug:       "legal-wiki",
		Language:   "es",
		ClaudeCode: true,
	})
	require.NoError(t, err)

	writeTestFile(t, filepath.Join(root, "wiki", "alpha.md"), "# Alpha\n")
	appendTestFile(t, filepath.Join(root, "wiki", "index.md"), "| [[alpha]] | Alpha | referencia | vigente | 2026-07-07 |\n")

	return root
}

func writeTestFile(t *testing.T, path string, content string) {
	t.Helper()
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
}

func appendTestFile(t *testing.T, path string, content string) {
	t.Helper()
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
	require.NoError(t, err)
	defer file.Close()
	_, err = file.WriteString(content)
	require.NoError(t, err)
}
