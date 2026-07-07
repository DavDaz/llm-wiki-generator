# Code Context

## Files Retrieved
1. `internal/cmd/root.go` (lines 51-84) - no-arg startup routing; decides inside-wiki dashboard vs outside-wiki launcher.
2. `internal/cmd/manage.go` (lines 20-29) - explicit `llm-wiki manage` path; opens the manage root menu.
3. `internal/tui/dashboard/root.go` (lines 31-37, 49-59, 81-148) - manage root menu with Tools/Drafts/Published/Deprecated/Exit and subview handling.
4. `internal/tui/dashboard/dashboard.go` (lines 43-47, 78-108) - tools-only dashboard subview; `q`/`esc` emit `BackToRoot`, not `tea.Quit`.
5. `internal/tui/launcher/launcher.go` (lines 15-20, 35-55, 62-84, 94-106) - outside-wiki launcher; only has Create/Guide choices and only explicit key exit is `ctrl+c`.
6. `internal/cmd/guide.go` (lines 20-39) - guide command; outside-wiki root loops back to launcher after this exits.
7. `internal/cmd/root_test.go` (lines 14-43) - current test locks in no-arg inside-wiki routing to `*dashboard.Model` tools-only view.
8. `internal/tui/dashboard/root_test.go` (lines 11-24, 27-37, 104-112) - tests manage root menu has Exit and root-level `q`/`esc`/`ctrl+c` quit.
9. `internal/tui/dashboard/dashboard_test.go` (lines 11-32) - tests tools view emits BackToRoot on `q`/`esc`, but does not quit on `ctrl+c` when standalone.
10. `README.md` (lines 50-68) - docs show both `llm-wiki manage` and no-arg inside wiki, but note says no-arg routing was intentionally left unchanged.

## Key Code

`internal/cmd/root.go` currently routes no-arg inside a wiki to the tools-only subview:

```go
func runRoot(_ *cobra.Command, _ []string) error {
    m, wikiRoot, err := loadManifestFromCwd()
    if err == nil {
        d := dashboard.NewTools(m, wikiRoot)
        _, runErr := runProgram(d, tea.WithAltScreen())
        return runErr
    }
    // outside wiki launcher loop...
}
```

`internal/cmd/manage.go` routes explicit `manage` to the intended manage menu:

```go
func runManage(_ *cobra.Command, _ []string) error {
    _, wikiRoot, err := loadManifestFromCwd()
    if err != nil { return err }
    root := dashboard.NewRoot(wikiRoot)
    p := tea.NewProgram(root, tea.WithAltScreen())
    _, err = p.Run()
    return err
}
```

`internal/tui/dashboard/root.go` is the manage/dashboard root users expect:

```go
var rootMenuOptions = []rootMenuOption{
    {label: "Tools backends", value: "tools"},
    {label: "Drafts (status: borrador)", value: "drafts"},
    {label: "Published (status: vigente)", value: "published"},
    {label: "Deprecated (status: deprecado)", value: "deprecated"},
    {label: "Exit", value: "exit"},
}
```

`internal/tui/launcher/launcher.go` has no visible Exit option:

```go
Options(
    huh.NewOption("Create a new wiki", "new"),
    huh.NewOption("Read the guide", "guide"),
)
```

Outside a wiki, `runRoot` loops after guide closes:

```go
case launcher.ActionGuide:
    if err := runGuide(nil, nil); err != nil { return err }
    // guide closed → loop back to launcher
```

## Architecture

Startup flow:

- `cmd/llm-wiki/main.go` calls `cmd.Execute()`.
- Cobra root command in `internal/cmd/root.go` uses `RunE: runRoot` for plain `llm-wiki`.
- `runRoot` calls `loadManifestFromCwd()` to detect a wiki by `wiki.toml` in the current working directory.
- If detected, current code opens `dashboard.NewTools(...)`, the tools-only model from `internal/tui/dashboard/dashboard.go`.
- `llm-wiki manage` instead opens `dashboard.NewRoot(wikiRoot)`, which is the newer manage root menu with page buckets and Exit.
- If not detected, `runRoot` opens `launcher.New()` in a loop. Launcher can choose Create or Guide. Guide exits back into the loop, returning to launcher.

Likely root causes:

1. Inside-wiki no-arg problem: `runRoot` is wired to `dashboard.NewTools`, while `runManage` is wired to `dashboard.NewRoot`. This exactly explains why `llm-wiki manage` shows the manage/dashboard root menu but plain `llm-wiki` does not. The current test `TestRunRootInsideWikiKeepsNoArgToolsDashboardRouting` explicitly preserves the old behavior.
2. Cannot exit plain `llm-wiki`: outside-wiki launcher has no Exit menu item and only explicitly handles `ctrl+c`. If the user chooses Guide, pressing `q` exits the guide but `runRoot` loops back to launcher. The only visible choices are Create or Guide, so there is no obvious normal exit path. If plain `llm-wiki` is run inside a wiki with current `NewTools` routing, `q`/`esc` emit `BackToRootMsg`; because the tools model is running standalone, that message is ignored by Bubble Tea as a final message and does not represent a root quit. `ctrl+c` is also intentionally no-op in the standalone tools model test.

## Missing Tests

Recommended tests to add/update:

1. Replace `internal/cmd/root_test.go:14-43` with a test asserting no-arg inside wiki launches the manage root (`dashboard.NewRoot`) rather than `*dashboard.Model` tools-only. Since `rootModel` is unexported, options include checking `fmt.Sprintf("%T", model)` equals `*dashboard.rootModel`, or exporting/testing through behavior if preferred.
2. Add a cmd-level test for `runManage` and `runRoot` using the same `runProgram` seam. `runManage` currently constructs `tea.NewProgram` directly, so it is harder to test without refactoring it to use `runProgram`.
3. Add `internal/tui/launcher/launcher_test.go` covering:
   - visible `Exit` option exists;
   - selecting Exit returns an exit/none action that makes `runRoot` return;
   - `q` and `esc` abort/quit from launcher, not just `ctrl+c`.
4. Add cmd-level outside-wiki launcher-loop tests with stubbed `runProgram` results:
   - Guide then Exit returns cleanly;
   - Exit returns cleanly without starting wizard.
5. Add/adjust dashboard test if tools model is ever allowed to run standalone: `q`/`esc` should quit standalone or be run only under `NewRoot`. Best fix is to stop running it standalone.

## Recommended Fix

1. In `internal/cmd/root.go`, change inside-wiki no-arg routing from `dashboard.NewTools(m, wikiRoot)` to `dashboard.NewRoot(wikiRoot)`. The manifest value is only needed for detection; if unused, assign `_` or refactor detection.
2. In `internal/cmd/manage.go`, use the shared `runProgram` helper instead of constructing `tea.NewProgram` directly. This improves parity and makes startup routing testable.
3. In `internal/tui/launcher/launcher.go`, add an explicit `Exit` option and corresponding action/result. Also handle `q` and `esc` like `ctrl+c` for predictable keyboard exit.
4. Keep `dashboard.NewTools` as a subview under `NewRoot`; avoid launching it standalone because its `q`/`esc` semantics are “back to root”, not process quit.
5. Update README line 68 if behavior changes: remove the note saying no-arg routing is unchanged, and state that plain `llm-wiki` inside a wiki opens the manage menu.

## Start Here

Open `internal/cmd/root.go` first. The highest-impact issue is the single routing mismatch at lines 53-57: no-arg inside-wiki uses `dashboard.NewTools`, while `llm-wiki manage` uses `dashboard.NewRoot`.

## Supervisor coordination

No coordination needed. Engram memory tools were not available in this subagent toolset, so no memory save was performed.
