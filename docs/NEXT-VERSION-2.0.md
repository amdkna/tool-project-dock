# Dock by Webigram — Version 2.0 feature plan

**Status:** Proposed / not implemented  
**Current release baseline:** 1.0b (beta)  
**Target:** 2.0

## Product goal

Move Dock from a basic Git client to a lightweight Windows development workspace manager: understand the state of projects, commit safely, run them locally, and open previews without bundling a large development stack.

This is a **planning document**, not a statement that these features already exist.

## Priority 1 — Git visibility and safe daily workflow

| Feature | Behavior in 2.0 | Acceptance criteria |
|---|---|---|
| Git Status Indicators | Show Clean, Modified, Staged, Ahead, Behind, Diverged, and Conflicts on each project row. | User can identify local/remote state at a glance; fetch does not mutate working files. |
| Changes Preview | Show added, modified, removed and untracked files, with a readable diff before staging/committing. | User can review the actual files to be committed; ignored files remain excluded unless tracked. |
| Git Sync | Separate Fetch / Check Updates from Pull and Push; show incoming/outgoing commit counts. | Check operation is read-only apart from remote-tracking refs; any non-fast-forward condition is clearly reported. |
| Open in VS Code | Open the selected project in VS Code using a configured or discovered editor. | Missing VS Code displays actionable guidance, without blocking other Dock features. |
| Edit / Remove Project | Edit name, URL, branch and path; remove Dock registration without deleting the actual project folder. | Explicit save/cancel, validation of changes, no accidental file deletion. |

## Priority 2 — Workspace management and developer productivity

| Feature | Behavior in 2.0 | Acceptance criteria |
|---|---|---|
| Branch Manager | List, switch, create and delete local branches with safeguards. | Uncommitted changes/conflicts block destructive operations until handled; no forced branch deletion by default. |
| Multi-project Pull | Fetch/pull selected repositories in a managed queue. | Shows individual progress/results and skips unsafe repositories rather than overwriting work. |
| Git History | View recent commits, authors, dates and messages; optionally open selected commits on GitHub. | Works without requiring GitHub CLI; history is fetched locally. |
| Stash Manager | Stash changes with a label, list stashes and restore safely. | Restoring a stash does not silently overwrite conflicting changes. |
| Run Commands | Configure per-project commands (e.g., npm run dev, python app.py, docker compose up). | User must explicitly approve commands; per-project processes can be started/stopped and logs inspected. |
| Application Preview | Open the project's configured localhost URL and show connectivity status. | Port/URL is editable; an unavailable service produces a useful explanation, not a false green status. |
| Favorites and Search | Pin, search and filter repositories or group frequently used projects. | Large lists remain responsive and searchable. |

## Priority 3 — Expansion and reliability

| Feature | Behavior in 2.0 | Acceptance criteria |
|---|---|---|
| Project Health Checks | Detect missing Git/runtime dependencies, port conflicts and unreachable services. | Non-destructive checks with suggested fixes; no automatic installation without confirmation. |
| Project Templates | Create or clone projects from approved starter templates. | New paths are validated; templates never execute scripts implicitly. |
| Profiles / Workspaces | Group and switch between Work, Personal and Client project sets. | Profile switching retains independent settings and never leaks credentials. |
| Backup / Export | Export/import configuration and restore it on another Windows computer. | Tokens are never exported in plaintext; validate paths and duplicates before import. |

## Recommended implementation order

1. **2.0-A — Git Safety:** status indicators, changes preview, Git Sync.
2. **2.0-B — Workspace:** Edit/Remove, VS Code integration, search/favorites, branch manager, history and stash.
3. **2.0-C — Launch & Preview:** configurable commands, process management, local previews.
4. **2.0-D — Reliability:** multi-project pull, checks, templates, profiles, backups and multi-instance hardening.

## Engineering and security requirements

- Keep small binary size and fast startup as a design constraint; minimize additional runtime dependencies.
- Support a documented Windows compatibility matrix. Windows 8 and 8.1 are unsupported by current OS vendors, so require explicit feasibility and security testing before advertising compatibility.
- Credentials should stay in an OS-protected store; never put PATs in TOML, logs or release packages.
- Preserve local commits and uncommitted work on failures; clearly distinguish fetching, pulling, committing and pushing.
- Provide confirmation before actions that overwrite local state, delete a branch or stop a running project.
- Keep Git operations responsive with cancellable progress and per-repository locking.
- Preserve the user's existing `config.toml` and offer a reversible migration if the schema changes.
- Verify Windows compilation, GUI behavior and update/install workflows before publishing 2.0.

## Out of scope for now

- Built-in Git hosting or replacing GitHub.
- Automatic force-push, silent hard reset, or unattended execution of unknown scripts.
- Bundling large runtimes such as Docker, Node.js or Python inside Dock itself.

## Release definition

Version 2.0 should be called delivered only after selected scope is implemented, tested on the supported Windows versions, documented and released. Features may be split across interim builds, but this file remains the planning source of truth.
