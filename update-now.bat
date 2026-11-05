@echo off
setlocal EnableExtensions DisableDelayedExpansion
cls
cd /d "%~dp0"
if errorlevel 1 (
    echo [ERROR] Cannot access the Project Dock directory.
    exit /b 1
)

echo ========================================
echo   Dock by Webigram - Update Now
echo ========================================

where git >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Git for Windows was not found.
    exit /b 1
)

for /f "delims=" %%B in ('git branch --show-current 2^>nul') do set "CURRENT_BRANCH=%%B"
if /i not "%CURRENT_BRANCH%"=="main" (
    echo [ERROR] Expected branch main, found "%CURRENT_BRANCH%".
    exit /b 1
)

for /f "delims=" %%H in ('git rev-parse HEAD 2^>nul') do set "BEFORE=%%H"
if not defined BEFORE (
    echo [ERROR] Cannot determine the current Git revision.
    exit /b 1
)

echo [1/3] Current revision: %BEFORE%
echo [2/3] Checking for updates...
git pull --ff-only origin main
if errorlevel 1 (
    echo [ERROR] Git pull failed. Build cancelled.
    exit /b 1
)

for /f "delims=" %%H in ('git rev-parse HEAD 2^>nul') do set "AFTER=%%H"
if not defined AFTER (
    echo [ERROR] Cannot determine the updated Git revision.
    exit /b 1
)

if /i "%BEFORE%"=="%AFTER%" (
    echo [OK] Already up to date. No new commit; skipping build and launch.
    exit /b 0
)

echo [3/3] New revision detected: %AFTER%
echo Building and launching Dock...
powershell.exe -NoProfile -File "%~dp0scripts\build.ps1" -Run
if errorlevel 1 (
    echo [ERROR] Build or launch failed. See the PowerShell output above.
    exit /b 1
)

echo [OK] Update completed and Dock launched.
exit /b 0
