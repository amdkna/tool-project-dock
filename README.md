# Dock by Webigram

**Pull. Run. Preview.**

A lightweight portable Windows developer tool for Git-based projects. The first release is a native GUI foundation; local web preview comes in a later iteration.

## Current version: ver. 0.1.2

Implemented:
- **File > Config...** to open an auto-generated `config.toml` in Notepad.
- **Info > ver. 0.1.2** menu label.
- Two main buttons: **PULL** and **PUSH**, operating on a configured existing local Git repository.
- Activity log and non-blocking Git execution.

Not yet implemented: Clone, start/stop web servers, project previews, automatic project detection and full multi-instance isolation.

## Windows build
Requires Windows, Go 1.22+ and Git for Windows.

From the repository root in PowerShell:

```powershell
git pull origin main
.\scripts\build.ps1
```

To build and launch in one step:

```powershell
.\scripts\build.ps1 -Run
```

The script resolves Go dependencies, creates `dist` and compiles `dist/dock.exe`. It works regardless of the current working directory when invoked by path. It does **not** automatically pull Git changes.

The executable generates `config.toml` beside itself on first start. Edit that file through **File > Config...** and set `project.directory` to an existing Git clone.

See [development instructions](docs/DEVELOPMENT.md), [roadmap](docs/ROADMAP.md), [example configuration](config.example.toml) and [changelog](CHANGELOG.md).

**Note:** Source has been committed, but the Windows build still needs to be validated on a machine with Go module access.
