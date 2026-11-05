# Build and publish a NEW versioned GitHub Release.
# Run from anywhere: .\scripts\release.ps1
[CmdletBinding()]
param([switch]$SkipBuild)
$ErrorActionPreference='Stop'
$root=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$repo='amdkna/tool-project-dock'
$source=Join-Path $root 'cmd\dock\main_windows.go'
$exe=Join-Path $root 'dist\dock.exe'

foreach($tool in @('git','gh')){
 if(-not (Get-Command $tool -ErrorAction SilentlyContinue)){
  throw "Required command '$tool' is missing. Install Git for Windows and GitHub CLI (winget install GitHub.cli)."
 }
}
Push-Location $root
try {
 $match=[regex]::Match((Get-Content -LiteralPath $source -Raw),'const\s+version\s*=\s*"(\d+\.\d+\.\d+)"')
 if(-not $match.Success){throw 'Cannot determine app version from source.'}
 $version=$match.Groups[1].Value
 $tag="v$version"
 Write-Host "Publishing $tag to $repo" -ForegroundColor Cyan

 $branch=(& git branch --show-current).Trim()
 if($LASTEXITCODE -ne 0 -or $branch -ne 'main'){throw 'Releases must be made from the main branch.'}
 $dirty=@(& git status --porcelain --untracked-files=normal)
 if($LASTEXITCODE -ne 0 -or $dirty.Count -gt 0){throw 'Working tree is not clean. Commit or stash changes first.'}
 & git fetch origin main
 if($LASTEXITCODE -ne 0){throw 'Unable to fetch origin/main.'}
 $local=(& git rev-parse HEAD).Trim()
 $remote=(& git rev-parse origin/main).Trim()
 if($local -ne $remote){throw 'Local HEAD differs from origin/main. Pull and push changes first.'}

 & gh auth status
 if($LASTEXITCODE -ne 0){throw 'GitHub CLI authentication required: run gh auth login.'}
 & gh release view $tag --repo $repo 2>$null | Out-Null
 if($LASTEXITCODE -eq 0){throw "Release $tag already exists. Bump the version rather than overwriting a release."}

 if(-not $SkipBuild){
  & (Join-Path $PSScriptRoot 'build.ps1')
  if(-not $?){throw 'Build failed.'}
 }
 if(-not (Test-Path -LiteralPath $exe -PathType Leaf)){throw 'Build output is missing: dist\dock.exe'}
 if((Get-Item -LiteralPath $exe).Length -eq 0){throw 'Build output is empty.'}

 $sumPath=Join-Path $root 'dist\dock.exe.sha256'
 $hash=(Get-FileHash -LiteralPath $exe -Algorithm SHA256).Hash.ToLowerInvariant()
 Set-Content -LiteralPath $sumPath -Value ("$hash  dock.exe") -Encoding ascii
 Write-Host "Creating release $tag..." -ForegroundColor Cyan
 & gh release create $tag $exe $sumPath --repo $repo --target $local --title "Dock by Webigram $tag" --generate-notes
 if($LASTEXITCODE -ne 0){throw 'GitHub release creation failed.'}
 Write-Host "Published $tag successfully." -ForegroundColor Green
}
finally {Pop-Location}
