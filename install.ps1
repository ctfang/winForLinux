# install.ps1
param(
    [string]$InstallDir = "$env:LOCALAPPDATA\win-tools",
    [string]$BackupDir = "$env:LOCALAPPDATA\win-tools-backup",
    [switch]$AddToPath,
    [switch]$Force,
    [switch]$Restore
)

$ErrorActionPreference = "Stop"

function Get-WinForLinuxProfileBlock {
    return @'

# BEGIN winForLinux / win-tools
# Prefer win-tools binaries over PowerShell aliases/functions
$winToolsDir = Join-Path $env:LOCALAPPDATA 'win-tools'
if (Test-Path $winToolsDir) {
    $parts = @($winToolsDir)
    foreach ($p in ($env:Path -split [IO.Path]::PathSeparator)) {
        if ($p -and $p -ne $winToolsDir) { $parts += $p }
    }
    $env:Path = $parts -join [IO.Path]::PathSeparator
}
foreach ($a in @('ls','cat','cp','mv','rm','mkdir','sort','md','echo','pwd','clear','diff','tee')) {
    if (Get-Alias $a -ErrorAction SilentlyContinue) {
        Remove-Item "Alias:$a" -Force -ErrorAction SilentlyContinue
    }
}
# mkdir is a function in Windows PowerShell; remove so mkdir.exe is used
if (Test-Path Function:mkdir) {
    Remove-Item Function:mkdir -Force -ErrorAction SilentlyContinue
}
# clear is typically an alias for Clear-Host
if (Get-Alias clear -ErrorAction SilentlyContinue) {
    Remove-Item Alias:clear -Force -ErrorAction SilentlyContinue
}
if (Get-Alias pwd -ErrorAction SilentlyContinue) {
    Remove-Item Alias:pwd -Force -ErrorAction SilentlyContinue
}
# UTF-8 for console + native exe pipes (PS 5.1 defaults $OutputEncoding to ASCII,
# which turns non-ASCII into '?' in pipes like: cat README.md | grep -i win).
# Idempotent; skip chcp when stdout is redirected; never fail profile load.
if ($Host.Name -eq 'ConsoleHost') {
    try {
        $utf8 = New-Object System.Text.UTF8Encoding $false
        [Console]::OutputEncoding = $utf8
        [Console]::InputEncoding  = $utf8
        $global:OutputEncoding    = $utf8
        if (-not [Console]::IsOutputRedirected) {
            chcp 65001 > $null
        }
    } catch { }
}
# END winForLinux / win-tools
'@
}

function Update-WinForLinuxProfile {
    $profilePath = Join-Path ([Environment]::GetFolderPath('MyDocuments')) 'WindowsPowerShell\Microsoft.PowerShell_profile.ps1'
    $profileDir = Split-Path $profilePath -Parent
    if (-not (Test-Path $profileDir)) {
        New-Item -ItemType Directory -Force $profileDir | Out-Null
    }

    $block = Get-WinForLinuxProfileBlock
    $existing = ''
    if (Test-Path $profilePath) {
        $existing = Get-Content $profilePath -Raw -ErrorAction SilentlyContinue
        if ($null -eq $existing) { $existing = '' }
    }

    $begin = '# BEGIN winForLinux / win-tools'
    if ($existing -match [regex]::Escape($begin)) {
        $pattern = '(?s)# BEGIN winForLinux / win-tools.*?# END winForLinux / win-tools\r?\n?'
        $updated = [regex]::Replace($existing, $pattern, $block.TrimStart() + "`r`n")
        Set-Content -Path $profilePath -Value $updated.TrimEnd() -Encoding UTF8
    } else {
        $sep = if ($existing.Trim().Length -gt 0) { "`r`n`r`n" } else { '' }
        Set-Content -Path $profilePath -Value ($existing.TrimEnd() + $sep + $block.TrimStart()) -Encoding UTF8
    }
    Write-Host "Updated PowerShell profile: $profilePath" -ForegroundColor Green
}

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
$aliasesToRemove = @("ls", "cat", "cp", "mv", "rm", "mkdir", "sort", "echo", "pwd", "clear", "diff", "tee")
foreach ($a in $aliasesToRemove) {
    if (Get-Alias $a -ErrorAction SilentlyContinue) {
        Remove-Item "Alias:$a" -Force -ErrorAction SilentlyContinue
        Write-Host "Removed PowerShell alias: $a" -ForegroundColor Gray
    }
}

# mkdir is a function in Windows PowerShell
if (Test-Path Function:mkdir) {
    Remove-Item Function:mkdir -Force -ErrorAction SilentlyContinue
}
# clear is typically an alias for Clear-Host
if (Get-Alias clear -ErrorAction SilentlyContinue) {
    Remove-Item Alias:clear -Force -ErrorAction SilentlyContinue
}
if (Get-Alias pwd -ErrorAction SilentlyContinue) {
    Remove-Item Alias:pwd -Force -ErrorAction SilentlyContinue
}

Update-WinForLinuxProfile

Write-Host ""
Write-Host "Installation complete!" -ForegroundColor Green
Write-Host "Installed to: $InstallDir" -ForegroundColor Green
Write-Host "Backups at: $BackupDir" -ForegroundColor Green
Write-Host ""
Write-Host "To restore original commands: .\install.ps1 -Restore" -ForegroundColor Gray
