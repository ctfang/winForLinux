# install.ps1
param(
    [string]$InstallDir = "$env:LOCALAPPDATA\win-tools",
    [string]$BackupDir = "$env:LOCALAPPDATA\win-tools-backup",
    [switch]$AddToPath,
    [switch]$Force,
    [switch]$Restore
)

$ErrorActionPreference = "Stop"

# Restore mode
if ($Restore) {
    if (-not (Test-Path $BackupDir)) {
        Write-Host "No backup directory found at: $BackupDir" -ForegroundColor Red
        exit 1
    }

    Write-Host "Restoring from backup..." -ForegroundColor Yellow

    Get-ChildItem $BackupDir -Directory | ForEach-Object {
        $cmdName = $_.Name
        $metadataPath = Join-Path $_.FullName "metadata.json"
        $backupFile = Get-ChildItem $_.FullName -Filter "*.original" | Select-Object -First 1

        if ($backupFile -and (Test-Path $metadataPath)) {
            $metadata = Get-Content $metadataPath | ConvertFrom-Json
            $originalPath = $metadata.OriginalPath

            Write-Host "Restoring $cmdName to: $originalPath" -ForegroundColor Cyan
            Copy-Item $backupFile.FullName $originalPath -Force
        }
    }

    # Remove install directory
    Remove-Item -Recurse -Force $InstallDir -ErrorAction SilentlyContinue
    Write-Host "Restore complete!" -ForegroundColor Green
    exit 0
}

# Install mode
if (-not (Test-Path ".\dist")) {
    Write-Host "dist directory not found. Run build.ps1 first." -ForegroundColor Red
    exit 1
}

Write-Host "Installing win-tools..." -ForegroundColor Cyan

# Create directories
New-Item -ItemType Directory -Force $InstallDir | Out-Null
New-Item -ItemType Directory -Force $BackupDir | Out-Null

# Backup existing commands
$commands = Get-ChildItem ".\dist\*.exe" | ForEach-Object { $_.BaseName }

foreach ($cmd in $commands) {
    $existing = Get-Command $cmd -ErrorAction SilentlyContinue
    if ($existing) {
        $backupPath = Join-Path $BackupDir $cmd
        New-Item -ItemType Directory -Force $backupPath | Out-Null

        # Backup original file
        Copy-Item $existing.Source "$backupPath\$cmd.exe.original" -Force

        # Save metadata
        @{
            OriginalPath = $existing.Source
            BackupDate = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
            ToolVersion = (Get-Item $existing.Source).VersionInfo.FileVersion
        } | ConvertTo-Json | Set-Content "$backupPath\metadata.json"

        Write-Host "Backed up: $($existing.Source)" -ForegroundColor Yellow
    }
}

# Install new commands
Copy-Item ".\dist\*.exe" $InstallDir -Force:$Force

# Add to PATH
if ($AddToPath) {
    $currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($currentPath -notlike "*$InstallDir*") {
        [Environment]::SetEnvironmentVariable("Path", "$currentPath;$InstallDir", "User")
        Write-Host "Added to PATH: $InstallDir" -ForegroundColor Green
        Write-Host "Restart your terminal to apply changes." -ForegroundColor Yellow
    }
}

Write-Host "`nInstallation complete!" -ForegroundColor Green
Write-Host "Installed to: $InstallDir" -ForegroundColor Green
Write-Host "Backups at: $BackupDir" -ForegroundColor Green
Write-Host "`nTo restore original commands: .\install.ps1 -Restore" -ForegroundColor Gray
