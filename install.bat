@echo off
setlocal EnableExtensions DisableDelayedExpansion
cd /d "%~dp0"
if errorlevel 1 (
  echo [ERROR] Cannot enter the project directory.
  exit /b 1
)

echo ======================================
echo   Dock by Webigram - Setup and Update
echo ======================================

where git >nul 2>&1
if errorlevel 1 (
  echo [ERROR] Git for Windows was not found in PATH.
  echo Install Git from https://git-scm.com/downloads/win
  exit /b 1
)

echo [1/3] Checking Go...
where go >nul 2>&1
if not errorlevel 1 goto go_ready

if exist "%ProgramFiles%\Go\bin\go.exe" (
  set "PATH=%ProgramFiles%\Go\bin;%PATH%"
  goto go_ready
)
if exist "%LOCALAPPDATA%\Programs\Go\bin\go.exe" (
  set "PATH=%LOCALAPPDATA%\Programs\Go\bin;%PATH%"
  goto go_ready
)

echo Go is missing. Installing through Windows Package Manager...
where winget >nul 2>&1
if errorlevel 1 (
  echo [ERROR] winget is unavailable. Install Go manually from https://go.dev/dl/
  exit /b 1
)
winget install --id GoLang.Go -e --source winget --accept-package-agreements --accept-source-agreements
if errorlevel 1 (
  echo [ERROR] Go installation failed. Check output above.
  exit /b 1
)

rem winget may update the user/machine PATH after this process has started.
if exist "%ProgramFiles%\Go\bin\go.exe" set "PATH=%ProgramFiles%\Go\bin;%PATH%"
if exist "%LOCALAPPDATA%\Programs\Go\bin\go.exe" set "PATH=%LOCALAPPDATA%\Programs\Go\bin;%PATH%"
where go >nul 2>&1
if errorlevel 1 (
  echo [ERROR] Go installed but is not visible in this session. Open a new terminal and retry.
  exit /b 1
)

:go_ready
go version
if errorlevel 1 (
  echo [ERROR] Go is installed but cannot run.
  exit /b 1
)

echo [2/3] Checking PowerShell execution policy...
rem Check effective policy, not only CurrentUser: GPO may override personal settings.
powershell.exe -NoProfile -NonInteractive -Command "$p=Get-ExecutionPolicy; if ($p -in @('Restricted','AllSigned')) {exit 1} else {exit 0}"
if not errorlevel 1 goto policy_ready
echo Setting CurrentUser policy to RemoteSigned...
powershell.exe -NoProfile -NonInteractive -Command "try { Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy RemoteSigned -Force -ErrorAction Stop; $p=Get-ExecutionPolicy; if ($p -in @('Restricted','AllSigned')) { Write-Error ('Effective policy is still ' + $p + '; check Group Policy.'); exit 1 } } catch { Write-Error $_; exit 1 }"
if errorlevel 1 (
  echo [ERROR] Could not allow local PowerShell scripts. Check organizational policy.
  exit /b 1
)
:policy_ready
echo PowerShell local script execution is permitted.

echo [3/3] Updating repository...
git rev-parse --is-inside-work-tree >nul 2>&1
if errorlevel 1 (
  echo [ERROR] The installer must run inside the cloned Git repository.
  exit /b 1
)
for /f "delims=" %%B in ('git branch --show-current') do set "CURRENT_BRANCH=%%B"
if /I not "%CURRENT_BRANCH%"=="main" (
  echo [ERROR] You are on branch "%CURRENT_BRANCH%". Switch to main before updating.
  exit /b 1
)
for /f "delims=" %%S in ('git status --porcelain --untracked-files=normal') do (
  echo [ERROR] Your working tree has local changes. Commit, stash or move them before updating.
  exit /b 1
)
git remote get-url origin >nul 2>&1
if errorlevel 1 (
  echo [ERROR] The origin remote is missing.
  exit /b 1
)
git pull --ff-only origin main
if errorlevel 1 (
  echo [ERROR] Git pull failed. Your local work was not forcibly overwritten.
  exit /b 1
)

echo.
echo [OK] Setup complete and repository updated.
echo Build with: scripts\build.ps1
echo Or build and run with: scripts\build.ps1 -Run
exit /b 0
