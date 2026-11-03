# Development guide

## Build

Windows 10/11, Go 1.22+ and Git for Windows required.

```powershell
cd 'C:\devs\Tools\dev\project-dock'
git pull origin main
.\scripts\build.ps1 -Run
```

Compiles `dist/dock.exe`. First run creates `dist/config.toml`, which is gitignored.

## Configuration

The sample `config.example.toml` is tracked. Use TOML `[[repositories]]` blocks, with `name`, `url`, `directory`, `branch`, and optional `remote` (default `origin`). The app validates required fields and rejects duplicate local directories and unknown config keys. Use single quotes for Windows paths to avoid escaping backslashes.

**Migration from 0.1.2:** replace the old `[project]` and `[git]` sections with one or more `[[repositories]]` blocks. Set the URL to match `git remote get-url origin`. Config is read on startup and on **File > Refresh**; invalid config leaves previous rows unchanged.

## Git safety / limits

Each row's buttons operate on the configured local Git clone. The app validates that the remote URL matches configuration, then runs `git pull --ff-only` or `git push`. It never commits, forces a push or runs clone. Only one Git operation can run per Dock instance. No protection against different instances using the same local working directory.

Live output is not streamed; the completed command's output is shown in logs. No web preview, managed development server or project cloning yet. Large lists have no dedicated scrolling area in this version. Windows compilation could not be tested in the authoring environment.

## Versioning

`x.y.z`: x=0 until explicitly changed, y=1 until owner requests otherwise, z increments with each application update (up to 999). Keep source version, changelog and README synchronized.
