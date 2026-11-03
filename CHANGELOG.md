# Changelog

Version policy: `x.y.z`. Main version `x` is currently 0; upgrade series `y` stays at 1 until the owner changes it; change counter `z` increments on each published application update (maximum 999). This is a project-specific convention, not necessarily Semantic Versioning.

## [0.1.2] - 2026-11-03
- Added `scripts/build.ps1` for one-command Windows GUI builds and optional launch.
- Added prerequisite checks, explicit command exit-code handling and output validation.
- Updated README and development guide with the new build workflow.

## [0.1.1] - 2026-11-03
- Initial portable Windows GUI in Go using Walk.
- Added File > Config... to create and edit TOML settings.
- Added Info > ver. 0.1.1.
- Added asynchronous Git PULL and PUSH commands with result logs.
- Added build instructions, configuration documentation and roadmap.

## [0.1.0] - 2026-11-03
- Requirements documented before implementation.
