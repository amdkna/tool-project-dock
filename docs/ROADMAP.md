# Dock by Webigram — Product roadmap

**Current application baseline:** `1.0b` (beta).  
**Next planned major version:** `2.0`.

Dock is a lightweight portable Windows GUI for managing local Git repositories and eventually launching and previewing development projects.

## Version 1.0b — current baseline

The beta includes multi-repository TOML configuration, GitHub HTTPS authentication through Git Credential Manager, Add New with clone and access test, configurable clone destination, pull, stage/commit/push, Explorer and GitHub shortcuts, installation helpers, update/build scripts, activity logging and GitHub release publishing script.

The Version 1.0b label reflects the owner's requested version name. It does not imply a GitHub Release has been published or that Windows compatibility testing is complete.

## Version 2.0 — planned

All feature proposals and acceptance criteria are tracked in **[NEXT-VERSION-2.0.md](NEXT-VERSION-2.0.md)**. They are not implemented in 1.0b.

### Focus
1. Safe Git status, changes preview and synchronized operations.
2. Developer workspace usability and project management.
3. Project commands, processes and preview capabilities.
4. Reliable configuration backup, diagnostics and multiple-instance safety.

### Design constraints
- Maintain a fast, lightweight, primarily native Windows experience.
- No automatic execution of unknown project scripts.
- Never silently discard user changes, force-push, or disclose tokens.
- Explicitly document and test supported Windows versions rather than promising compatibility without verification.
- Existing 1.0b configuration should remain readable or be migrated safely.

See [CHANGELOG.md](../CHANGELOG.md) for previous implementation history.
