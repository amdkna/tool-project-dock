# Changelog

Version policy: `x.y.z`; x remains 0 until owner changes it, y remains 1 until owner changes it, and z increments on every published update (maximum 999).

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
