# uninstall.ps1
param(
    [string]$InstallDir = "$env:LOCALAPPDATA\win-tools",
    [switch]$Restore
)

$ErrorActionPreference = "Stop"

if ($Restore) {
    & .\install.ps1 -Restore
    exit 0
}

Write-Host "Uninstalling win-tools..." -ForegroundColor Yellow

# Remove from PATH
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -like "*$InstallDir*") {
    $newPath = ($currentPath -split ';' | Where-Object { $_ -ne $InstallDir }) -join ';'
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    Write-Host "Removed from PATH" -ForegroundColor Green
}

# Remove install directory
if (Test-Path $InstallDir) {
    Remove-Item -Recurse -Force $InstallDir
    Write-Host "Removed: $InstallDir" -ForegroundColor Green
}

Write-Host "`nUninstall complete!" -ForegroundColor Green
Write-Host "Backups preserved at: $env:LOCALAPPDATA\win-tools-backup" -ForegroundColor Gray
Write-Host "To restore backups: .\uninstall.ps1 -Restore" -ForegroundColor Gray
