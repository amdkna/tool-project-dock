# Development and release guide

## Requirements
Windows 10/11, Go 1.22+, Git for Windows and configured Git credentials. An internet connection is needed to fetch Go modules on first build.

## Build from PowerShell
```powershell
cd 'C:\devs\Tools\dev\project-dock'
git pull origin main
go mod tidy
New-Item -ItemType Directory -Force -Path 'dist' | Out-Null
go build -ldflags="-H windowsgui" -o "dist/dock.exe" ./cmd/dock
.\dist\dock.exe
```

## Usage
1. Launch Dock.
2. Select **File > Config...** (opens the per-executable `config.toml` in Notepad).
3. Change `project.directory` to an **existing Git clone**; use TOML literal strings for Windows paths: `directory = 'C:\devs\my-project'`.
4. Set `git.remote` (typically `origin`) and `git.branch` (typically `main`).
5. Save the file, click PULL or PUSH and inspect the activity log.

PULL executes `git -C <directory> pull --ff-only <remote> <branch>` and refuses non-fast-forward merges. PUSH executes `git -C <directory> push <remote> <branch>` and transfers existing commits only; it does **not** stage or commit files. No force-push is provided. Git terminal prompting is disabled; authenticate using Git for Windows first.

## Limitations of 0.1.1
- A local cloned repository must already exist. Clone and project preview are not implemented.
- Git output is displayed **after the command completes**, not streamed live.
- One command at a time per instance; separate instances should have separate executable directories and should not operate simultaneously on one Git working tree.
- The app creates `config.toml` beside its executable. That directory must be writable.
- Go module dependencies have not been built in the authoring environment; test the Windows build before treating this as a release.

## Versioning
- `x=0` until an explicit main-version change.
- `y=1` until the owner explicitly requests another upgrade series.
- Increment `z` for each new published app change (0–999).
- Update the source version in `cmd/dock/main_windows.go`, README and changelog together **before pushing**. No automatic bump currently exists.

## Why TOML?
TOML is readable and supports comments and organized sections. It is less indentation-sensitive than YAML and more comfortable to hand-edit than JSON. Validate fields on load and avoid putting credentials in it. Relative paths are resolved against the config directory.
