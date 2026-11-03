# Changelog

Version policy: `x.y.z`; x remains 0 until owner changes it, y remains 1 until owner changes it, and z increments on every published update (maximum 999).

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
