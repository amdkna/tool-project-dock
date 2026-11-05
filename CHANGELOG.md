# Changelog

Version policy: `x.y.z`; x remains 0 until owner changes it, y remains 1 until owner changes it, and z increments on every published update (maximum 999).

## [0.1.15] - 2026-11-05
- Updated update-now.bat to close the currently running Dock application after a successful pull and only when a new commit is present.
- Added scripts/stop-dock.ps1, targeting only this checkout's dist/dock.exe by full executable path to avoid terminating other Dock installations.
- Aborts rebuilding if the old Dock instance cannot be stopped and then rebuilds and launches as before.
- Updated documentation and application version.

## [0.1.14] - 2026-11-05
- PUSH now opens a small commit-message dialog with Commit & Push and Cancel.
- Commit & Push stages all non-ignored changes (git add -A), creates a commit if needed, and pushes the checked-out configured branch.
- Validates the repository root, remote URL and current branch before staging; reports stage, commit and push errors in the activity log.
- If no changes are staged, pushes existing local commits without creating an empty commit.
- Updated documentation and application version.

## [0.1.13] - 2026-11-05
- Added Info > Project Repo opening https://github.com/amdkna/www-freelanceops in the default browser.
- Added a Git Repo button to every repository row, opening its configured GitHub URL.
- Validated GitHub HTTPS URLs before launching the browser.
- Updated application version and documentation.

## [0.1.12] - 2026-11-05
- Added Test button to Add New dialog to check GitHub branch read access and local destination.
- Added checkbox to clone directly into chosen folder instead of chosen-folder/repo-name.
- Changed Add to clone first, save TOML only on success and refresh the list afterward.
- Removed disabling terminal authentication for clone; Git Credential Manager can now request interactive sign-in.
- Reused matching existing Git clones and supported empty selected destination folders.
- Updated documentation and version.

## [0.1.11] - 2026-11-05
- Added File > Install Tools menu with actions for required tools (Git + GCM), Git for Windows, Git Credential Manager, and optional GitHub CLI.
- Added scripts/install-tools.ps1 with prerequisite detection and Windows OS compatibility selection.
- On Windows 8/8.1, opens compatible legacy Git/GCM release pages rather than installing unsupported modern versions.
- On Windows 10/11, uses winget when available for Git and gh, otherwise official install pages.
- Updated documentation and application version.

## [0.1.10] - 2026-11-05
- Fixed Add New to persist repository settings in the ignored dist/config.toml before attempting Git clone.
- Refreshes repository cards immediately after saving, so network/authentication failures no longer lose the new entry.
- Displays configuration write failures explicitly and retains clone errors in the activity log.
- Updated README and development documentation.

## [0.1.9] - 2026-11-05
- Added root `update-now.bat`: clear screen, fast-forward pull, compare before/after Git commit IDs, and build/run only when updated.
- Renamed each repository's Folder button to Browse.
- Added a Create GitHub Token button in the token dialog opening GitHub's fine-grained token creation page.
- Explained least-privilege token setup: selected repositories, Contents read/write, Metadata read-only.
- Updated README and development documentation.

## [0.1.8] - 2026-11-05
- Added a GitHub token dialog using Git Credential Manager and Windows Credential Manager with an isolated Dock namespace.
- Added Explorer folder buttons for each configured repository.
- Added Add New repository dialog with HTTPS URL, branch and Windows folder chooser; clones to parent/repository folder.
- Added TOML repository entry append with validation, preserving existing entries.
- Added scripts/release.ps1 for versioned GitHub Releases, Windows builds and SHA-256 checksums.
- Updated documentation and app version.

## [0.1.7] - 2026-11-05
- Support native Windows paths with single backslashes in double-quoted directory fields.
- Preserve standard TOML literal and escaped path forms.
- Added path parsing tests and updated examples and documentation.

## [0.1.6] - 2026-11-05
- Fixed the Walk startup error `TTM_ADDTOOL failed` caused by a missing Windows Common Controls v6 application manifest.
- Added `assets/dock.manifest` and resource embedding to `scripts/build.ps1` using rsrc.
- Generated .syso resources are cleaned after builds; executable remains portable without a sidecar manifest.
- Updated documentation and version label.

## [0.1.5] - 2026-11-05
- Fixed build.ps1 failing to find Go when install.bat temporarily updated PATH but the calling PowerShell session still had an old PATH.
- Build script resolves go.exe from common installation locations and uses its absolute path.
- Installer adds Go to the user PATH for future terminals only when no existing user/machine PATH entry is present.
- Updated setup documentation and version label.

## [0.1.4] - 2026-11-03
- Added root `install.bat` to check and install Go (via winget) only if missing.
- Added effective PowerShell execution policy check; applies CurrentUser RemoteSigned only when required.
- Added safe fast-forward Git update with clean-working-tree, branch and remote checks.
- Updated application version, README and developer documentation.

## [0.1.3] - 2026-11-03
- Added support for multiple repositories via TOML `[[repositories]]`.
- Added thin, bordered rows with PULL/PUSH buttons to the left of repository name and branch.
- Added File > Refresh to reload config and redraw the repository list.
- Added per-repository URL, directory, branch, and remote settings and validation.
- Kept `config.toml` ignored and updated tracked `config.example.toml`.
- Updated README and developer documentation.

## [0.1.2] - 2026-11-03
- Added `scripts/build.ps1` for one-command Windows GUI builds and optional launch.
- Added prerequisite checks, explicit command exit-code handling and output validation.

## [0.1.1] - 2026-11-03
- Initial portable Windows GUI in Go using Walk.
- Added File > Config..., Info > ver. 0.1.1, PULL/PUSH and logging.

## [0.1.0] - 2026-11-03
- Documented the product requirements.
