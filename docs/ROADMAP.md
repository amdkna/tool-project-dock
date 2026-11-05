# Dock by Webigram — Product roadmap

**Current application baseline:** `1.0b` (beta).  
**Next planned major version:** `2.0`.

Dock is a lightweight, portable **Windows control center for projects in multiple programming languages**, not a Go-only application host.

## Version 1.0b — current baseline

Multi-repository TOML configuration, HTTPS GitHub credentials through GCM, Add New/Clone/Test, pull, commit/push, Explorer/GitHub shortcuts, dependency installation shortcuts, update/build scripts, activity log and a GitHub release script. The 1.0b label is not proof that a binary release or Windows validation has occurred.

## Version 2.0 — planned scope

1. Independent background agent with loopback-only management API.
2. Local HTTP/HTTPS gateway binding 127.0.0.1:80/443, reverse proxy and static site support for PHP, Python, JS/TS, Node and Go applications.
3. Self-signed development certificates and user-approved local CA with trust and rotation management.
4. Cross-language project process lifecycle management.
5. Minimal Docker Engine client (containers, images, status, logs, start/stop/restart, optional Compose).
6. Tiny local CI/CD executor using TOML pipelines, SQLite history, manual/scheduled/Git polling triggers, controlled jobs and logs.
7. Earlier Git/workspace visibility, safety and management improvements, retained as Version 2.0 goals.

**Design and acceptance details:** [VERSION-2.0-ARCHITECTURE.md](VERSION-2.0-ARCHITECTURE.md).  
**Feature backlog:** [NEXT-VERSION-2.0.md](NEXT-VERSION-2.0.md).  
**GitHub issue index:** [V2-ISSUES.md](V2-ISSUES.md).

## Release guardrails

- No automatic execution of untrusted project scripts.
- No Docker daemon bundled; control an existing local or remote engine.
- No plaintext token, private key or credential in source-controlled TOML.
- Management listener loopback-only by default; external exposure needs an explicit threat model.
- Explicit Windows compatibility testing (especially Windows 8/8.1) before claiming support.
- Preserve or safely migrate the 1.0b configuration.

See [CHANGELOG.md](../CHANGELOG.md) for historical implementation details.
