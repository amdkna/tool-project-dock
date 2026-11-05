# Stop only Dock instances whose executable is this checkout's dist\dock.exe.
# Avoid terminating unrelated applications with the same executable name.
[CmdletBinding()]
param()
$ErrorActionPreference = 'Stop'
$root = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$target = [System.IO.Path]::GetFullPath((Join-Path $root 'dist\dock.exe'))
$comparison = [StringComparison]::OrdinalIgnoreCase

try {
    $candidates = @(Get-CimInstance Win32_Process -Filter "Name = 'dock.exe'" -ErrorAction Stop)
} catch {
    Write-Error "Could not inspect running Dock processes: $_"
    exit 1
}

$matches = @($candidates | Where-Object {
    $_.ExecutablePath -and [string]::Equals(
        [System.IO.Path]::GetFullPath($_.ExecutablePath), $target, $comparison
    )
})

if ($matches.Count -eq 0) {
    Write-Host '[OK] No running Dock instance from this checkout.' -ForegroundColor Green
    exit 0
}

foreach ($instance in $matches) {
    Write-Host ("Stopping Dock process PID {0}..." -f $instance.ProcessId) -ForegroundColor Yellow
    try {
        Stop-Process -Id $instance.ProcessId -Force -ErrorAction Stop
    } catch {
        Write-Error ("Could not stop Dock process PID {0}: {1}" -f $instance.ProcessId, $_)
        exit 1
    }
}

foreach ($instance in $matches) {
    try {
        Wait-Process -Id $instance.ProcessId -Timeout 10 -ErrorAction Stop
    } catch {
        # A process that already exited is fine; otherwise stop the rebuild.
        if (Get-Process -Id $instance.ProcessId -ErrorAction SilentlyContinue) {
            Write-Error ("Dock PID {0} did not exit. Build cancelled." -f $instance.ProcessId)
            exit 1
        }
    }
}

Write-Host '[OK] Previous Dock instance(s) stopped.' -ForegroundColor Green
exit 0
