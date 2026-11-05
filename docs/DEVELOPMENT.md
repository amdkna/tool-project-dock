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

`x.y.z`: x=0 until explicitly changed, y=1 until owner requests otherwise, z increments with each application update (up to 999). Keep source version, changelog and README synchronized.


## GitHub token security (v0.1.8)

File > GitHub Token... asks for the account username and a fine-grained PAT (masked entry).
The token is passed only via standard input to Git Credential Manager and stored in Windows Credential Manager using `GCM_CREDENTIAL_STORE=wincredman` and `GCM_NAMESPACE=webigram-dock`. The same namespace and credential helper are used for Git network operations. This keeps Dock separate from other Git client credentials. No token is saved to TOML, stdout or command arguments. Git Credential Manager must be present (normally supplied with Git for Windows). HTTPS GitHub remotes are supported, not SSH token authentication. Windows account security still matters; grant only necessary repository permissions and rotate compromised PATs.

## Add New and Explorer

Choose Add New, enter GitHub HTTPS URL and branch, and browse an existing parent folder. Dock creates/reuses a subfolder named after the repository, clones when absent, verifies an existing clone's remote and writes a TOML repository block. A failed clone is not added to the config. Folder buttons open the configured path in Explorer.

## Publishing releases

Install GitHub CLI: `winget install GitHub.cli` and authenticate with `gh auth login`. Run `.scriptselease.ps1` on main with a clean working tree synced to origin/main. The script creates a new version tag release, uploads the built Windows EXE and SHA-256 checksum, and refuses existing releases. It never silently replaces a release asset. The executable build and runtime features require Windows testing.
