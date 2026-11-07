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

New-Item -ItemType Directory -Force $InstallDir | Out-Null
New-Item -ItemType Directory -Force $BackupDir | Out-Null

$commands = Get-ChildItem ".\dist\*.exe" | ForEach-Object { $_.BaseName }

foreach ($cmd in $commands) {
    # Only consider Application commands with a real file Source (skip aliases/functions)
    $existing = Get-Command $cmd -CommandType Application -ErrorAction SilentlyContinue |
        Where-Object { $_.Source -and (Test-Path $_.Source) } |
        Select-Object -First 1

    if ($existing -and $existing.Source -notlike "$InstallDir*") {
        $backupPath = Join-Path $BackupDir $cmd
        New-Item -ItemType Directory -Force $backupPath | Out-Null

        Copy-Item $existing.Source "$backupPath\$cmd.exe.original" -Force

        @{
            OriginalPath = $existing.Source
            BackupDate   = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
            ToolVersion  = (Get-Item $existing.Source).VersionInfo.FileVersion
        } | ConvertTo-Json | Set-Content "$backupPath\metadata.json"

        Write-Host "Backed up: $($existing.Source)" -ForegroundColor Yellow
    }
}

Copy-Item ".\dist\*.exe" $InstallDir -Force

# Ensure PATH: put InstallDir first for current session and user env
$env:Path = $InstallDir + [IO.Path]::PathSeparator + ($env:Path -replace [regex]::Escape($InstallDir + [IO.Path]::PathSeparator), '')

if ($AddToPath) {
    $currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if (-not $currentPath) { $currentPath = "" }
    $parts = $currentPath -split ';' | Where-Object { $_ -and $_ -ne $InstallDir }
    $newPath = (@($InstallDir) + $parts) -join ';'
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    Write-Host "Added/refreshed PATH: $InstallDir" -ForegroundColor Green
}

# Remove PowerShell aliases that shadow common tools
$aliasesToRemove = @("ls", "cat", "cp", "mv", "rm", "mkdir", "sort", "echo")
foreach ($a in $aliasesToRemove) {
    if (Get-Alias $a -ErrorAction SilentlyContinue) {
        Remove-Item "Alias:$a" -Force -ErrorAction SilentlyContinue
        Write-Host "Removed PowerShell alias: $a" -ForegroundColor Gray
    }
}

Write-Host "`nInstallation complete!" -ForegroundColor Green
Write-Host "Installed to: $InstallDir" -ForegroundColor Green
Write-Host "Backups at: $BackupDir" -ForegroundColor Green
Write-Host "`nTo restore original commands: .\install.ps1 -Restore" -ForegroundColor Gray
