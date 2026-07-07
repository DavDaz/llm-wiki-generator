package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/DavDaz/llm-wiki-generator/internal/doctor"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Run read-only structural wiki health checks",
	RunE:  runDoctor,
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}

func runDoctor(cmd *cobra.Command, _ []string) error {
	_, wikiRoot, err := loadManifestFromCwd()
	if err != nil {
		return err
	}

	report, err := doctor.Run(wikiRoot)
	if err != nil {
		return err
	}

	w := cmd.OutOrStdout()
	passed, warned, failed := summarizeDoctorFindings(report.Findings)

	fmt.Fprintf(w, "Doctor report for %s\n", report.WikiRoot)
	for _, finding := range report.Findings {
		path := formatFindingPath(report.WikiRoot, finding.Path)
		if path == "" {
			fmt.Fprintf(w, "%s %s: %s\n", strings.ToUpper(string(finding.Severity)), finding.Check, finding.Message)
			continue
		}
		fmt.Fprintf(w, "%s %s [%s]: %s\n", strings.ToUpper(string(finding.Severity)), finding.Check, path, finding.Message)
	}
	fmt.Fprintf(w, "Summary: %d passed, %d warning(s), %d error(s)\n", passed, warned, failed)

	if failed > 0 {
		return fmt.Errorf("doctor found %d error(s)", failed)
	}

	return nil
}

func summarizeDoctorFindings(findings []doctor.Finding) (passed int, warned int, failed int) {
	for _, finding := range findings {
		switch finding.Severity {
		case doctor.SeverityPass:
			passed++
		case doctor.SeverityWarn:
			warned++
		case doctor.SeverityError:
			failed++
		}
	}

	return passed, warned, failed
}

func formatFindingPath(root string, path string) string {
	if path == "" {
		return ""
	}
	if relPath, err := filepath.Rel(root, path); err == nil && relPath != "." && !strings.HasPrefix(relPath, "..") {
		return relPath
	}
	return path
}
