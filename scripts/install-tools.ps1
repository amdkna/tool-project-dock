# Dock prerequisite installer. Windows 8/8.1 use official legacy download pages.
[CmdletBinding()]
param([ValidateSet('Required','Git','GCM','GH')][string]$Tool='Required')
$ErrorActionPreference='Stop'
$os=[Environment]::OSVersion.Version
$legacy = ($os.Major -eq 6 -and $os.Minor -le 3)
$modern = ($os.Major -ge 10)
Write-Host ("Dock tools setup | Windows {0} | Selection: {1}" -f $os,$Tool) -ForegroundColor Cyan
Write-Host "Current Git and GCM releases require Windows 10+; legacy versions are provided for Windows 8.x."
Write-Host "GitHub CLI (gh) is optional. It is needed by scripts/release.ps1, NOT by Dock's Git operations."

function Open-Official([string]$url,[string]$message) {
    Write-Host $message -ForegroundColor Yellow
    Write-Host $url
    Start-Process -FilePath $url
}
function Have-Tool([string]$name) {
    return ($null -ne (Get-Command $name -ErrorAction SilentlyContinue))
}
function Have-GCM {
    if (-not (Have-Tool 'git')) {return $false}
    & git credential-manager --version 2>$null | Out-Null
    return ($LASTEXITCODE -eq 0)
}
function Install-Git {
    if (Have-Tool 'git') { Write-Host '[OK] Git is installed:' -ForegroundColor Green; & git --version; return }
    if ($legacy) {
        if ($os.Minor -eq 3) {
            Open-Official 'https://github.com/git-for-windows/git/releases/tag/v2.55.0.windows.5' 'Windows 8.1: select the last compatible Git for Windows 2.55.0(5) installer.'
        } else {
            Open-Official 'https://github.com/git-for-windows/git/releases/tag/v2.46.2.windows.1' 'Windows 8: select a legacy Git for Windows 2.46.2 installer compatible with your architecture.'
        }
        return
    }
    if (-not $modern) {throw 'Unsupported Windows version; install a supported Git version manually.'}
    if (Have-Tool 'winget') {
        Write-Host 'Installing Git for Windows using winget...'
        & winget install --id Git.Git -e --source winget --accept-source-agreements --accept-package-agreements
        if ($LASTEXITCODE -ne 0) {throw 'Git installation failed.'}
        Write-Host '[OK] Git installer completed. Restart Dock to refresh PATH.' -ForegroundColor Green
    } else {
        Open-Official 'https://git-scm.com/install/windows' 'winget is not available. Download official Git for Windows, enabling Git Credential Manager.'
    }
}
function Install-GCM {
    if (Have-GCM) {Write-Host '[OK] Git Credential Manager is already installed.' -ForegroundColor Green; return}
    if ($legacy) {
        Open-Official 'https://github.com/git-ecosystem/git-credential-manager/releases/tag/v2.9.1' 'Windows 8/8.1: use legacy GCM 2.9.1 (NOT GCM 3.x). Select the Windows installer.'
        return
    }
    if (-not $modern) {throw 'Unsupported Windows version for current Git Credential Manager.'}
    Open-Official 'https://github.com/git-ecosystem/git-credential-manager/releases/latest' 'Install official GCM for Windows, or reinstall Git for Windows with GCM selected.'
}
function Install-GH {
    if (Have-Tool 'gh') {Write-Host '[OK] GitHub CLI is installed.' -ForegroundColor Green; & gh --version; return}
    if ($legacy) {
        Open-Official 'https://github.com/cli/cli/releases' 'GitHub CLI is optional; current releases are not guaranteed to run on Windows 8.x. Find a compatible older release manually.'
        return
    }
    if (-not $modern) {throw 'Unsupported Windows version for automatic GitHub CLI install.'}
    if (Have-Tool 'winget') {
        Write-Host 'Installing GitHub CLI using winget...'
        & winget install --id GitHub.cli -e --source winget --accept-source-agreements --accept-package-agreements
        if ($LASTEXITCODE -ne 0) {throw 'GitHub CLI installation failed.'}
        Write-Host '[OK] GitHub CLI installed. Open a new terminal before running gh auth login.' -ForegroundColor Green
    } else {
        Open-Official 'https://cli.github.com/' 'winget is not available; download GitHub CLI from its official website.'
    }
}
switch ($Tool) {
 'Git' {Install-Git}
 'GCM' {Install-GCM}
 'GH' {Install-GH}
 'Required' {
    Install-Git
    # If Git just installed, GCM may already be present but this session may
    # not detect it yet. Avoid falsely claiming GCM was installed.
    if (Have-Tool 'git') {Install-GCM}
    else {Write-Host 'After Git installation finishes, rerun Install Required Tools to check GCM.' -ForegroundColor Yellow}
 }
}
Write-Host 'Done. Installations/changes may require a new Dock session.' -ForegroundColor Green
