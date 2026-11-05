# Build Dock by Webigram on Windows.
# Usage (from repository root): .\scripts\build.ps1
# Optional: .\scripts\build.ps1 -Run
[CmdletBinding()]
param(
    [switch]$Run
)

$ErrorActionPreference = 'Stop'
$root = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$dist = Join-Path $root 'dist'
$target = Join-Path $dist 'dock.exe'
$manifest = Join-Path $root 'assets\dock.manifest'
$resource = Join-Path $root 'cmd\dock\rsrc_windows_amd64.syso'

if ($env:OS -ne 'Windows_NT') {
    throw 'This build script requires Windows.'
}

# Batch files cannot update the PATH of their parent PowerShell session.
# Resolve Go independently, including common Windows installation locations.
$goCommand = Get-Command go.exe -ErrorAction SilentlyContinue
$goExe = if ($goCommand) { $goCommand.Source } else { $null }

if (-not $goExe) {
    $candidates = @(
        (Join-Path $env:ProgramFiles 'Go\bin\go.exe'),
        (Join-Path $env:LOCALAPPDATA 'Programs\Go\bin\go.exe')
    )
    foreach ($candidate in $candidates) {
        if (Test-Path -LiteralPath $candidate -PathType Leaf) {
            $goExe = $candidate
            $env:PATH = (Split-Path -Path $goExe -Parent) + ';' + $env:PATH
            break
        }
    }
}
if (-not $goExe) {
    throw 'Go was not found. Run install.bat or install Go from https://go.dev/dl/.'
}

Push-Location $root
try {
    Write-Host '[1/5] Go version:' -ForegroundColor Cyan
    & $goExe version
    if ($LASTEXITCODE -ne 0) { throw 'Unable to execute Go.' }

    Write-Host '[2/5] Resolving Go dependencies...' -ForegroundColor Cyan
    & $goExe mod tidy
    if ($LASTEXITCODE -ne 0) { throw 'Go dependency resolution failed.' }

    Write-Host '[3/5] Preparing output directory...' -ForegroundColor Cyan
    New-Item -ItemType Directory -Path $dist -Force | Out-Null

    if (-not (Test-Path -LiteralPath $manifest -PathType Leaf)) {
        throw "Application manifest missing: $manifest"
    }
    if (Test-Path -LiteralPath $resource) {
        throw "Unexpected existing resource file: $resource. Remove it after confirming it is not needed."
    }

    Write-Host '[4/5] Embedding Windows Common Controls v6 manifest...' -ForegroundColor Cyan
    & $goExe run github.com/akavel/rsrc@latest -arch amd64 -manifest $manifest -o $resource
    if ($LASTEXITCODE -ne 0) { throw 'Windows manifest resource generation failed.' }
    if (-not (Test-Path -LiteralPath $resource -PathType Leaf)) {
        throw 'Manifest resource generator returned success but did not create .syso.'
    }

    Write-Host '[5/5] Building Windows GUI...' -ForegroundColor Cyan
    & $goExe build -trimpath -ldflags '-H windowsgui' -o $target ./cmd/dock
    if ($LASTEXITCODE -ne 0) { throw 'Go build failed. Check compiler output above.' }
    if (-not (Test-Path -LiteralPath $target -PathType Leaf)) {
        throw 'Build reported success but dock.exe was not found.'
    }

    $file = Get-Item -LiteralPath $target
    Write-Host ('Build successful: {0} ({1:N0} bytes)' -f $file.FullName, $file.Length) -ForegroundColor Green

    if ($Run) {
        Write-Host 'Launching Dock...' -ForegroundColor Cyan
        Start-Process -FilePath $target -WorkingDirectory $dist
    }
}
finally {
    # The resource is generated inside the Go package only for compilation.
    # Keep the repository clean and make every build regenerate the manifest.
    if (Test-Path -LiteralPath $resource -PathType Leaf) {
        Remove-Item -LiteralPath $resource -Force
    }
    Pop-Location
}
