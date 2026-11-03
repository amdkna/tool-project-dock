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

if ($env:OS -ne 'Windows_NT') {
    throw 'This build script requires Windows.'
}

$go = Get-Command go -ErrorAction SilentlyContinue
if (-not $go) {
    throw 'Go was not found in PATH. Install Go 1.22+ from https://go.dev/dl/ and reopen PowerShell.'
}

Push-Location $root
try {
    Write-Host '[1/4] Go version:' -ForegroundColor Cyan
    & go version
    if ($LASTEXITCODE -ne 0) { throw 'Unable to execute Go.' }

    Write-Host '[2/4] Resolving Go dependencies...' -ForegroundColor Cyan
    & go mod tidy
    if ($LASTEXITCODE -ne 0) { throw 'Go dependency resolution failed.' }

    Write-Host '[3/4] Preparing output directory...' -ForegroundColor Cyan
    New-Item -ItemType Directory -Path $dist -Force | Out-Null

    Write-Host '[4/4] Building Windows GUI...' -ForegroundColor Cyan
    & go build -trimpath -ldflags '-H windowsgui' -o $target ./cmd/dock
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
    Pop-Location
}
