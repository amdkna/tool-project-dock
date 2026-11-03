# Dock by Webigram

**Pull. Run. Preview.**

Lightweight portable Windows Git project utility. **Current version: ver. 0.1.3.**

## Current features

- Native Go + Walk GUI.
- File > Config..., File > Refresh, Info > ver. 0.1.3.
- Multiple repositories defined using TOML `[[repositories]]` blocks.
- Compact bordered repository rows: PULL/PUSH buttons on the left, name and branch to the right.
- File > Refresh reloads the TOML and updates rows, without restarting.
- Non-blocking Git operations with an activity log.
- The private `config.toml` is ignored by Git; tracked sample: [config.example.toml](config.example.toml).

## Build / launch (PowerShell)

Requires Go 1.22+, Git for Windows and an internet connection for first-time Go dependencies.

```powershell
cd 'C:\devs\Tools\dev\project-dock'
git pull origin main
.\scripts\build.ps1 -Run
```

The executable is written to `dist/dock.exe`. Its local configuration is `dist/config.toml`, generated if missing. Edit through **File > Config...**, then select **File > Refresh**.

## Example config

```toml
[[repositories]]
name = "Website"
url = "https://github.com/example/website.git"
directory = 'C:\devs\website'
branch = "main"
remote = "origin"

[[repositories]]
name = "API"
url = "https://github.com/example/api.git"
directory = 'C:\devs\api'
branch = "dev"
remote = "origin"
```

Every repository needs name, URL, directory, branch; remote defaults to origin. Relative paths are relative to config.toml. Local clones must already exist. The configured URL must match the local remote exactly. TOML contains no credentials.

PULL uses `git pull --ff-only <remote> <branch>`. PUSH uses `git push <remote> <branch>` (already committed changes only). Neither force-pushes or automatically commits.

**Not yet included:** cloning, local web servers, project previews and robust concurrency across different Dock processes. Existing configurations from v0.1.2 must be migrated from `[project]`/`[git]` to `[[repositories]]`.

See [development guide](docs/DEVELOPMENT.md), [roadmap](docs/ROADMAP.md) and [changelog](CHANGELOG.md). The Windows build has not yet been verified in the authoring environment.
