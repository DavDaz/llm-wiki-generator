// Package doctor provides bounded, read-only wiki health checks.
package doctor

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/DavDaz/llm-wiki-generator/internal/manifest"
	"github.com/DavDaz/llm-wiki-generator/internal/tools"
)

var (
	wikilinkRe = regexp.MustCompile(`\[\[([^\]]+)\]\]`)
	toolPaths  = map[string][]requiredPath{
		"claude-code": {
			{path: "CLAUDE.md"},
			{path: ".claude/skills", dir: true},
			{path: ".claude/skills/wiki-ingest/SKILL.md"},
			{path: ".claude/skills/wiki-query/SKILL.md"},
			{path: ".claude/skills/wiki-lint/SKILL.md"},
		},
		"opencode": {
			{path: "AGENTS.md"},
			{path: ".opencode/commands", dir: true},
			{path: ".opencode/commands/wiki-ingest.md"},
			{path: ".opencode/commands/wiki-query.md"},
			{path: ".opencode/commands/wiki-lint.md"},
		},
		"pi": {
			{path: "AGENTS.md"},
			{path: ".pi/prompts", dir: true},
			{path: ".pi/prompts/wiki-ingest.md"},
			{path: ".pi/prompts/wiki-query.md"},
			{path: ".pi/prompts/wiki-lint.md"},
		},
	}
)

// Severity describes the outcome level of a doctor finding.
type Severity string

const (
	SeverityPass  Severity = "pass"
	SeverityWarn  Severity = "warn"
	SeverityError Severity = "error"
)

// Finding is a single read-only diagnostic result.
type Finding struct {
	Check    string
	Severity Severity
	Path     string
	Message  string
}

// Report contains the bounded doctor findings for a wiki root.
type Report struct {
	WikiRoot string
	Findings []Finding
}

type mdFile struct{ name, path, content string }

type requiredPath struct {
	path string
	dir  bool
}

// Run executes bounded, read-only health checks for a wiki root.
func Run(root string) (Report, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return Report{}, fmt.Errorf("resolve wiki root: %w", err)
	}
	r := Report{WikiRoot: root}
	m := checkManifest(&r, root)
	checkCore(&r, root)
	checkTools(&r, root, m)
	if err := checkWiki(&r, root); err != nil {
		return Report{}, err
	}
	return r, nil
}

func checkManifest(r *Report, root string) *manifest.Manifest {
	path := filepath.Join(root, manifest.Filename)
	m, err := manifest.Load(root)
	if err != nil {
		r.add("manifest", SeverityError, path, fmt.Sprintf("cannot load manifest: %v", err))
		return nil
	}
	if err := m.Validate(); err != nil {
		r.add("manifest", SeverityError, path, fmt.Sprintf("manifest is invalid: %v", err))
		return m
	}
	r.add("manifest", SeverityPass, path, "manifest loads and validates")
	return m
}

func checkCore(r *Report, root string) {
	checks := []struct {
		path string
		dir  bool
	}{
		{filepath.Join(root, "raw"), true}, {filepath.Join(root, "wiki"), true},
		{filepath.Join(root, "wiki", "index.md"), false}, {filepath.Join(root, "wiki", "log.md"), false}, {filepath.Join(root, "wiki", "sources.json"), false},
	}
	missing := []string{}
	invalidPaths := []string{}
	for _, check := range checks {
		if err := validatePathKind(check.path, check.dir); err != nil {
			missing = append(missing, rel(root, check.path))
			invalidPaths = append(invalidPaths, check.path)
		}
	}
	if len(missing) == 0 {
		r.add("core-files", SeverityPass, filepath.Join(root, "wiki"), "required wiki directories and core files are present")
		return
	}
	sort.Strings(missing)
	path := root
	if len(invalidPaths) == 1 {
		path = invalidPaths[0]
	}
	r.add("core-files", SeverityError, path, "missing or invalid core paths: "+strings.Join(missing, ", "))
}

func checkTools(r *Report, root string, m *manifest.Manifest) {
	if m == nil {
		r.add("tool-outputs", SeverityWarn, root, "skipped tool output checks because the manifest could not be loaded")
		return
	}
	missing := []string{}
	invalidPaths := []string{}
	for _, name := range m.EnabledTools() {
		_, err := tools.Get(name)
		if err != nil {
			missing = append(missing, name+": unknown tool in manifest")
			continue
		}
		for _, req := range toolPaths[name] {
			if err := validatePathKind(filepath.Join(root, req.path), req.dir); err != nil {
				expected := "file"
				if req.dir {
					expected = "directory"
				}
				missing = append(missing, fmt.Sprintf("%s: missing or invalid %s (expected %s)", name, req.path, expected))
				invalidPaths = append(invalidPaths, filepath.Join(root, req.path))
			}
		}
	}
	if len(missing) == 0 {
		r.add("tool-outputs", SeverityPass, root, "enabled tool outputs are installed")
		return
	}
	sort.Strings(missing)
	path := root
	if len(invalidPaths) == 1 {
		path = invalidPaths[0]
	}
	r.add("tool-outputs", SeverityError, path, strings.Join(missing, "; "))
}

func checkWiki(r *Report, root string) error {
	files, pages, targets, err := readWiki(root)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		r.add("wikilinks", SeverityWarn, filepath.Join(root, "wiki"), "skipped wikilink checks because wiki/ is missing")
		r.add("index-coverage", SeverityWarn, filepath.Join(root, "wiki", "index.md"), "skipped index coverage because wiki/ is missing")
		return nil
	}
	broken := []string{}
	for _, file := range files {
		for _, match := range wikilinkRe.FindAllStringSubmatch(file.content, -1) {
			target := cleanTarget(match[1])
			if target != "" {
				if _, ok := targets[target]; !ok {
					broken = append(broken, rel(root, file.path)+" -> [["+target+"]]")
				}
			}
		}
	}
	if len(broken) == 0 {
		r.add("wikilinks", SeverityPass, filepath.Join(root, "wiki"), "wikilink targets resolve to wiki markdown pages")
	} else {
		sort.Strings(broken)
		r.add("wikilinks", SeverityError, filepath.Join(root, "wiki"), "broken wikilinks: "+strings.Join(broken, "; "))
	}
	index := ""
	for _, file := range files {
		if file.name == "index.md" {
			index = file.content
			break
		}
	}
	if index == "" {
		r.add("index-coverage", SeverityError, filepath.Join(root, "wiki", "index.md"), "wiki/index.md is missing or unreadable")
		return nil
	}
	missing := []string{}
	for _, page := range pages {
		if !strings.Contains(index, "[["+page+"]]") && !strings.Contains(index, page+".md") {
			missing = append(missing, page)
		}
	}
	if len(missing) == 0 {
		r.add("index-coverage", SeverityPass, filepath.Join(root, "wiki", "index.md"), "wiki/index.md covers discovered wiki pages")
		return nil
	}
	sort.Strings(missing)
	r.add("index-coverage", SeverityWarn, filepath.Join(root, "wiki", "index.md"), "pages missing from wiki/index.md: "+strings.Join(missing, ", "))
	return nil
}

func readWiki(root string) ([]mdFile, []string, map[string]struct{}, error) {
	entries, err := os.ReadDir(filepath.Join(root, "wiki"))
	if err != nil {
		return nil, nil, nil, err
	}
	files, pages, targets := []mdFile{}, []string{}, map[string]struct{}{}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}
		path := filepath.Join(root, "wiki", entry.Name())
		buf, err := os.ReadFile(path)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("read wiki file %s: %w", path, err)
		}
		files = append(files, mdFile{entry.Name(), path, string(buf)})
		slug := strings.TrimSuffix(entry.Name(), ".md")
		targets[slug] = struct{}{}
		if entry.Name() != "index.md" && entry.Name() != "log.md" {
			pages = append(pages, slug)
		}
	}
	return files, pages, targets, nil
}

func cleanTarget(raw string) string {
	target, _, _ := strings.Cut(raw, "|")
	target, _, _ = strings.Cut(target, "#")
	return strings.TrimSpace(target)
}

func rel(root string, path string) string {
	relPath, err := filepath.Rel(root, path)
	if err != nil {
		return path
	}
	return relPath
}

func validatePathKind(path string, wantDir bool) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if wantDir && !info.IsDir() {
		return fmt.Errorf("expected directory")
	}
	if !wantDir && info.IsDir() {
		return fmt.Errorf("expected file")
	}
	return nil
}

func (r *Report) add(check string, severity Severity, path string, message string) {
	r.Findings = append(r.Findings, Finding{Check: check, Severity: severity, Path: path, Message: message})
}
