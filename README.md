# Dock by Webigram

**Pull. Run. Preview.**

Lightweight portable Windows Git project utility. **Current version: ver. 0.1.15.**

## Current features

- Native Go + Walk GUI.
- File > Config..., File > Refresh, Info > ver. 0.1.15.
- Multiple repositories defined using TOML `[[repositories]]` blocks.
- Compact bordered repository rows: PULL/PUSH buttons on the left, name and branch to the right.
- File > Refresh reloads the TOML and updates rows, without restarting.
- Non-blocking Git operations with an activity log.
- **Browse** button on each repository row opens Windows Explorer; **Git Repo** opens that repository's configured GitHub URL in the default browser.\n- **Info > Project Repo** opens https://github.com/amdkna/www-freelanceops in the default browser.
- **Add New...** includes Test, Add, and an optional "Clone directly into selected folder" checkbox. Add clones (or verifies an existing clone), saves its entry in ignored TOML **only after success**, and refreshes the GUI. Test checks remote branch/read access and local destination.
- **File > Install Tools** checks dependencies and offers Git for Windows, GCM and optional GitHub CLI installations.
- **File > GitHub Token...** includes a direct Create GitHub Token link and required permissions; it saves a masked token via Git Credential Manager in Windows Credential Manager (dedicated Dock namespace, never TOML).
- The private `config.toml` is ignored by Git; tracked sample: [config.example.toml](config.example.toml).

## One-click update

Run `update-now.bat` from the repository root. It clears the console, fast-forward pulls `origin/main`, compares Git HEAD before and after pulling, and **stops only this checkout's running `dist/dock.exe` process** before running `scripts/build.ps1 -Run` **only if a new commit was pulled**. No changes means no rebuild and no launch. A failed pull or failed process stop cancels the rebuild. Any active operations in the old Dock instance will be interrupted when an update requires rebuilding. This is based on Git commit IDs, not an artifact checksum; manually deleted builds are not regenerated when the commit is unchanged.

```powershell
.\update-now.bat
```

## Initial setup / update (Command Prompt or PowerShell)

From the cloned repository, run the root-level `install.bat`:

```powershell
cd 'C:\devs\Tools\dev\project-dock'
.\install.bat
```

The installer checks Git, installs Go with winget only when missing, makes Go discoverable in future terminals where necessary, checks the effective PowerShell policy and uses `CurrentUser RemoteSigned` only when necessary, and safely fast-forward pulls `origin/main`. It stops if the working tree is dirty or the current branch is not main. An internet connection is needed to install Go and pull updates.

## Build / launch (PowerShell)

Requires Go 1.22+, Git for Windows and an internet connection for first-time Go dependencies. The build script now also finds Go in common installation directories even when the current PowerShell PATH is stale.

```powershell
cd 'C:\devs\Tools\dev\project-dock'
.\scripts\build.ps1 -Run
```

The build script generates a temporary Windows resource from `assets/dock.manifest` using `rsrc` and embeds Common Controls v6 into the executable to prevent Walk's `TTM_ADDTOOL failed` startup error. The temporary `.syso` is removed after compilation; the manifest itself remains tracked.\n\nThe executable is written to `dist/dock.exe`. Its local configuration is `dist/config.toml`, generated if missing. Edit through **File > Config...**, then select **File > Refresh**.

## Example config

```toml
[[repositories]]
name = "Website"
url = "https://github.com/example/website.git"
directory = "C:\devs\website\"
branch = "main"
remote = "origin"

[[repositories]]
name = "API"
url = "https://github.com/example/api.git"
directory = "C:\devs\api\"
branch = "dev"
remote = "origin"
```

Dock accepts directory = "C:\devs\test\" with single backslashes and a trailing slash. This is a Dock-specific convenience; literal TOML paths and properly escaped paths remain supported. All other fields follow standard TOML.

Every repository needs name, URL, directory, branch; remote defaults to origin. Relative paths are relative to config.toml. Add New clones (or verifies the existing clone) before saving a repository entry; if cloning fails, no broken entry is added, while rows already in TOML still expect their local directories to exist. The configured URL must match the local remote exactly. TOML contains no credentials.

PULL uses `git pull --ff-only <remote> <branch>`. PUSH opens a commit-message dialog. Clicking **Commit & Push** stages all non-ignored changes using `git add -A`, commits them using your message (if there are changes), then pushes the configured branch. If no changes exist, it pushes existing commits. Cancel makes no changes. Pull remains fast-forward only, and Dock never force-pushes.

**Not yet included:** web preview, local web servers, project previews and robust concurrency across different Dock processes. Existing configurations from v0.1.2 must be migrated from `[project]`/`[git]` to `[[repositories]]`.

See [development guide](docs/DEVELOPMENT.md), [roadmap](docs/ROADMAP.md) and [changelog](CHANGELOG.md). The Windows build has not yet been verified in the authoring environment.

## Publishing Windows releases

GitHub CLI (gh) must be installed and authenticated using `gh auth login`. Publish the next application version from a clean, synchronized main branch:

```powershell
.\scripts\release.ps1
```

The release script builds `distdock.exe`, generates a SHA-256 checksum and publishes both artifacts as a NEW GitHub release tagged with the app version. It refuses to overwrite an existing version. Use `-SkipBuild` only when the existing EXE was built from exactly the current source.

## Install Tools compatibility

**Git Credential Manager is not part of Windows itself.** It is normally bundled with Git for Windows, but `gh` is a separate optional program. Choose **File > Install Tools > Install Required** to check Git + GCM. The three individual actions are also available. On Windows 10/11, supported dependencies install with winget where available; otherwise official download pages open. On Windows 8.1/8, Dock directs you to compatible older Git and GCM releases instead of attempting current versions. No silent installation of obsolete packages is attempted. Windows 8.x operation of the Dock EXE itself has not been verified. See [tool support details](docs/DEVELOPMENT.md).
