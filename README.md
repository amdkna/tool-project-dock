# Dock by Webigram

**Pull. Run. Preview.**

A lightweight portable Windows developer tool for Git-based projects. The first release is a native GUI foundation; local web preview comes in a later iteration.

## Current version: ver. 0.1.1

Implemented:
- **File > Config...** to open an auto-generated `config.toml` in Notepad.
- **Info > ver. 0.1.1** menu label.
- Two main buttons: **PULL** and **PUSH**, operating on a configured existing local Git repository.
- Activity log and non-blocking Git execution.

Not yet implemented: Clone, start/stop web servers, project previews, automatic project detection and full multi-instance isolation.

## Windows build
Requires Go 1.22+ and Git for Windows:

```powershell
cd 'C:\devs\Tools\dev\project-dock'
git pull origin main
go mod tidy
New-Item -ItemType Directory -Force -Path 'dist' | Out-Null
go build -ldflags="-H windowsgui" -o 'dist/dock.exe' ./cmd/dock
.\dist\dock.exe
```

The executable generates `config.toml` beside itself on first start. Edit that file through **File > Config...** and set `project.directory` to an existing Git clone.

See [development instructions](docs/DEVELOPMENT.md), [roadmap](docs/ROADMAP.md), [example configuration](config.example.toml) and [changelog](CHANGELOG.md).

**Note:** Source has been committed, but the Windows build still needs to be validated on a machine with Go module access.
