# Development guide

## Windows setup and update

Run `install.bat` from the repository root (or by absolute path). It checks Git, installs Go via winget if missing, allows local PowerShell scripts with `CurrentUser RemoteSigned` if required, and fast-forward pulls `main` if the working tree is clean.

If Go was found only by the installer, its installation directory is added to the persistent user PATH only when not already present in user or machine PATH. A running PowerShell session may have a stale PATH; the build script independently locates Go in common install folders and works without restarting PowerShell.\n\nIt deliberately does not bypass machine or organizational Group Policy. It does not force-pull, reset or discard uncommitted files. Go installation through winget may prompt for elevated privileges. It does not run the application build.

```powershell
cd 'C:\devs\Tools\dev\project-dock'
.\install.bat
```

## Build

Windows 10/11, Go 1.22+ and Git for Windows required.

```powershell
cd 'C:\devs\Tools\dev\project-dock'
.\scripts\build.ps1 -Run
```

The build script uses `go run github.com/akavel/rsrc@latest` to compile `assets/dock.manifest` into a temporary `cmd/dock/rsrc_windows_amd64.syso`, which Go links into the Windows application. This is required by Walk's Common Controls v6 integration. The resource is removed in a finally block after the build. The first build needs network access for rsrc.\n\nCompiles `dist/dock.exe`. First run creates `dist/config.toml`, which is gitignored.

## Configuration

The sample `config.example.toml` is tracked. Use TOML `[[repositories]]` blocks, with `name`, `url`, `directory`, `branch`, and optional `remote` (default `origin`). The app validates required fields and rejects duplicate local directories and unknown config keys. Dock also accepts directory = "C:\devs\test\" with single backslashes in double quotes. Standard TOML literal paths and properly escaped strings remain valid.

**Migration from 0.1.2:** replace the old `[project]` and `[git]` sections with one or more `[[repositories]]` blocks. Set the URL to match `git remote get-url origin`. Config is read on startup and on **File > Refresh**; invalid config leaves previous rows unchanged.

## Git safety / limits

Each row's buttons operate on the configured local Git clone. The app validates that the remote URL matches configuration, then runs `git pull --ff-only` or `git push`. It never commits, forces a push or runs clone. Only one Git operation can run per Dock instance. No protection against different instances using the same local working directory.

Live output is not streamed; the completed command's output is shown in logs. No web preview, managed development server or project cloning yet. Large lists have no dedicated scrolling area in this version. Windows compilation could not be tested in the authoring environment.

## Versioning

The current application baseline is `1.0b` (beta), as requested by the product owner. The next planned major version is `2.0`. Keep the source constant, README and changelog aligned; the release script recognizes the beta label when creating a `v1.0b` tag. No GitHub release is created just by changing source files. Version 2.0 feature planning is maintained in [NEXT-VERSION-2.0.md](NEXT-VERSION-2.0.md).


## GitHub token security (v0.1.8)

File > GitHub Token... asks for the account username and a fine-grained PAT (masked entry).
The token is passed only via standard input to Git Credential Manager and stored in Windows Credential Manager using `GCM_CREDENTIAL_STORE=wincredman` and `GCM_NAMESPACE=webigram-dock`. The same namespace and credential helper are used for Git network operations. This keeps Dock separate from other Git client credentials. No token is saved to TOML, stdout or command arguments. Git Credential Manager must be present (normally supplied with Git for Windows). HTTPS GitHub remotes are supported, not SSH token authentication. Windows account security still matters; grant only necessary repository permissions and rotate compromised PATs.

## Add New and Explorer

Choose Add New, enter GitHub HTTPS URL and branch, and browse an existing parent folder. Dock writes the TOML repository block first and refreshes the GUI, then creates/reuses a subfolder named after the repository and clones when absent. If clone fails, the entry stays in config.toml for recovery; the Activity log shows the error. Use File > Config to correct the entry and File > Refresh to reload it. Browse buttons open the configured path in Explorer.

## Publishing releases

Install GitHub CLI: `winget install GitHub.cli` and authenticate with `gh auth login`. Run `.scriptselease.ps1` on main with a clean working tree synced to origin/main. The script creates a new version tag release, uploads the built Windows EXE and SHA-256 checksum, and refuses existing releases. It never silently replaces a release asset. The executable build and runtime features require Windows testing.

## Update Now (v0.1.9)

Run `update-now.bat` from the repository root (or absolute path). It uses `git rev-parse HEAD` before and after a safe `git pull --ff-only origin main`. If the revision changed, it runs `scripts/build.ps1 -Run`. If HEAD is unchanged, it skips compilation and does not start a new application instance. It stops on pull or build errors. This does not repair a missing EXE if there is no new commit; for a manual rebuild use `scripts/build.ps1 -Run`.

## Fine-grained GitHub token permissions

Within File > GitHub Token, choose **Create GitHub Token** to open https://github.com/settings/personal-access-tokens/new. Select the intended GitHub resource owner and limit access to only the desired repositories. For private repository pull and push, use **Repository permissions > Contents: Read and write**; GitHub automatically grants **Metadata: Read-only**. No administration, workflows or other permissions are required for ordinary Git push and pull (workflow changes may need additional permissions). Use a limited expiration and protect the Windows user profile.

## Install Tools menu (v0.1.11)

GCM is not built into Windows. It normally comes with Git for Windows when selected during setup. The optional `gh` utility is separately installed and is needed only for `scripts/release.ps1`.

Open **File > Install Tools** and select Install Required, Git for Windows, Git Credential Manager, or GitHub CLI. The menu runs `scripts/install-tools.ps1` in a visible PowerShell window and skips detected installations. On Windows 10/11, Git and gh use winget if available; Git Credential Manager opens official standalone installer releases if missing. On Windows 8.1, Git 2.55.0(5) and GCM 2.9.x legacy releases are offered; on Windows 8, Git 2.46.2 and GCM 2.9.x are offered. GitHub CLI compatibility on Windows 8.x is not guaranteed and only its release archive is opened for manual selection. Windows 8.x versions are out of vendor support and should not be used for security-sensitive workloads. The Dock executable itself has not been tested on Windows 8.x and a binary built with a modern Go toolchain may require a newer Windows version.

The application expects `scripts/install-tools.ps1` in the source checkout near the executable's `dist` directory. When distributing only `dock.exe`, include the `scripts` folder or install dependencies separately.

## Add New / Test / exact folder (v0.1.12)

Add New has a Test button next to Add; it checks the GitHub branch via git ls-remote and examines local destination safety. On private repositories, Git Credential Manager can request browser sign-in. Read access and branch existence are checked; write/push access cannot be verified without making an actual write. Uncheck "Clone directly into selected folder" to clone into selected-folder/repository-name. Check it to clone straight into the chosen folder. Existing folders must be empty or a matching Git checkout, never unrelated nonempty data. Clicking Add clones first, saves config.toml only on success and refreshes the UI. If clone fails, errors are logged and no invalid configuration entry is created.

If an older Dock version already wrote an entry for a failed clone, use File > Config to remove or correct the old [[repositories]] block before adding again. The newer Add flow does not delete those existing entries automatically.

## Repository browser shortcuts (v0.1.13)

**Info > Project Repo** opens https://github.com/amdkna/www-freelanceops. Each configured repository row includes **Git Repo** next to Browse, opening the `url` field of that repository in the default browser. Only HTTPS GitHub URLs are accepted, and invalid links produce a warning.

## PUSH commit dialog (v0.1.14)

Click PUSH on a repository row to open the commit-message dialog. Cancel exits without staging, committing, or pushing. **Commit & Push** checks the configured folder is the Git root and that the URL/branch match the local checkout, then runs `git add -A`, `git diff --cached --quiet`, `git commit -m "<message>"` if changes exist, and finally `git push <remote> HEAD:refs/heads/<branch>`. Git respects `.gitignore` for untracked files; already tracked files are still staged even if later added to `.gitignore`. This includes deletions and any existing staged changes. An empty working tree does not generate an empty commit but still pushes outgoing commits. Git identity (`user.name` and `user.email`) must be configured; if a commit or push fails, prior local staging/commits remain in place. No force push. Git Credential Manager may prompt for authentication.

## Update Now stops the old app (v0.1.15)

After a successful fast-forward pull, `update-now.bat` compares Git HEAD before/after. Only for a changed commit, `scripts/stop-dock.ps1` locates `dock.exe` processes using the Win32_Process executable path, stops only instances from this source checkout's `dist\dock.exe`, and waits for exit. If unable to inspect or stop matching processes, updating aborts before rebuilding. When no new commit exists, it neither kills a running Dock nor builds. Stopping an in-progress Dock Git operation interrupts it; ensure ongoing operations are complete before updating.

## Version 2.0 architecture planning (proposed)

Read [VERSION-2.0-ARCHITECTURE.md](VERSION-2.0-ARCHITECTURE.md) for cross-language gateway, 127.0.0.1:80/443 reverse proxy, local certificate authority, app process manager, Docker API client, TOML/SQLite-backed local CI/CD, security constraints and sample pipeline schema. See [V2-ISSUES.md](V2-ISSUES.md) for implementation work items. This is future design, not 1.0b functionality.
