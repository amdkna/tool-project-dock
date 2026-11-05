# Dock by Webigram — Version 2.0 GitHub issue tracker

**Status:** Planning, none of these issue descriptions imply implementation.  
**Baseline:** 1.0b. **Target:** 2.0.  
**Issue list:** https://github.com/amdkna/tool-project-dock/issues  
**Technical architecture:** [VERSION-2.0-ARCHITECTURE.md](VERSION-2.0-ARCHITECTURE.md)  
**Feature roadmap:** [NEXT-VERSION-2.0.md](NEXT-VERSION-2.0.md)

## Purpose

This document is the canonical issue index. Every GitHub issue contains a user story, technical specification, acceptance criteria, validation tests, dependency notes and definition-of-done checklist. Priorities: **P0** required platform foundation; **P1** important usability/safety; **P2** later enhancement.

## Foundation

| Issue | Priority | Work item |
|---|---|---|
| [#1](https://github.com/amdkna/tool-project-dock/issues/1) | P0 | Dock Agent and GUI control API |
| [#2](https://github.com/amdkna/tool-project-dock/issues/2) | P0 | Config migration and secrets |

## Web Gateway

| Issue | Priority | Work item |
|---|---|---|
| [#3](https://github.com/amdkna/tool-project-dock/issues/3) | P0 | HTTP/HTTPS 80/443 |
| [#4](https://github.com/amdkna/tool-project-dock/issues/4) | P0 | Host routing + WebSocket |
| [#5](https://github.com/amdkna/tool-project-dock/issues/5) | P1 | Static website hosting |
| [#6](https://github.com/amdkna/tool-project-dock/issues/6) | P0 | Gateway route editor |

## TLS

| Issue | Priority | Work item |
|---|---|---|
| [#7](https://github.com/amdkna/tool-project-dock/issues/7) | P0 | Self-signed certificate issuer |
| [#8](https://github.com/amdkna/tool-project-dock/issues/8) | P0 | Local CA, trust, rotation |

## App Manager

| Issue | Priority | Work item |
|---|---|---|
| [#9](https://github.com/amdkna/tool-project-dock/issues/9) | P0 | Cross-language process manager |
| [#10](https://github.com/amdkna/tool-project-dock/issues/10) | P1 | Health, logs, preview |

## Docker

| Issue | Priority | Work item |
|---|---|---|
| [#11](https://github.com/amdkna/tool-project-dock/issues/11) | P0 | Docker Engine connections |
| [#12](https://github.com/amdkna/tool-project-dock/issues/12) | P0 | Containers UI + lifecycle |
| [#13](https://github.com/amdkna/tool-project-dock/issues/13) | P2 | Docker Compose integration |

## CI/CD

| Issue | Priority | Work item |
|---|---|---|
| [#14](https://github.com/amdkna/tool-project-dock/issues/14) | P0 | TOML pipelines and approval |
| [#15](https://github.com/amdkna/tool-project-dock/issues/15) | P0 | Local pipeline worker |
| [#16](https://github.com/amdkna/tool-project-dock/issues/16) | P0 | SQLite history and logs |
| [#17](https://github.com/amdkna/tool-project-dock/issues/17) | P1 | Schedules and Git polling |
| [#18](https://github.com/amdkna/tool-project-dock/issues/18) | P2 | Dashboard, artifacts, recovery |

## Git

| Issue | Priority | Work item |
|---|---|---|
| [#19](https://github.com/amdkna/tool-project-dock/issues/19) | P1 | Status + ahead/behind |
| [#20](https://github.com/amdkna/tool-project-dock/issues/20) | P1 | Diff preview and history |

## Workspace

| Issue | Priority | Work item |
|---|---|---|
| [#21](https://github.com/amdkna/tool-project-dock/issues/21) | P1 | Project editor, VS Code, search |
| [#22](https://github.com/amdkna/tool-project-dock/issues/22) | P2 | Branch, stash, bulk pull |
| [#23](https://github.com/amdkna/tool-project-dock/issues/23) | P1 | Profiles, backups, templates |

## Release

| Issue | Priority | Work item |
|---|---|---|
| [#24](https://github.com/amdkna/tool-project-dock/issues/24) | P0 | Security and Windows QA |

## Recommended delivery order

1. **Foundation:** #1, #2, #24 threat modeling / compatibility gate.
2. **Gateway:** #3, #6, #4, #5.
3. **Certificates:** #7, #8, with explicit Windows trust consent.
4. **Multi-language application runner:** #9, #10.
5. **Docker client:** #11, #12, then #13.
6. **Local CI/CD:** #14, #16, #15, then #17 and #18.
7. **Git and workspace improvements:** #19–#23, in dependency order.
8. **Final release gate:** finish #24 with security and Windows test matrix.

## Common release acceptance criteria

- Previous `dist/config.toml` remains usable after upgrade or is backed up and reversibly migrated.
- Language runtimes remain external; Go powers the Dock control plane, not every hosted app.
- Management/control API defaults to loopback only; routes on 127.0.0.1:80/443 are configurable and port conflicts are diagnosed.
- No unapproved local CA trust installation, daemon exposure, secret logging, untrusted script execution or force-push.
- Docker manager requires an existing Engine, not a Dock-provided daemon.
- SQLite history persists and pipeline runs continue independently of the GUI while the agent is running.
- Supported Windows versions are validated, not inferred; legacy Windows 8/8.1 remain a feasibility decision.

## Maintenance notes

Keep the issue bodies and this index synchronized when priorities or scope change. Do not label features as shipped merely because an issue exists. For each completed issue, update docs, tests, changelog and its GitHub status. The existing app remains 1.0b until the owner chooses a new release label.
