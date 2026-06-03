# build.ps1
param(
    [string]$OutputDir = ".\dist",
    [string]$Version = "1.0.0",
    [switch]$Clean
)

$ErrorActionPreference = "Stop"

if ($Clean) {
    Write-Host "Cleaning dist directory..." -ForegroundColor Yellow
    Remove-Item -Recurse -Force $OutputDir -ErrorAction SilentlyContinue
}

New-Item -ItemType Directory -Force $OutputDir | Out-Null

$commands = @("ls", "grep", "cat", "which")

foreach ($cmd in $commands) {
    Write-Host "Building $cmd..." -ForegroundColor Cyan
    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
    go build -o "$OutputDir\$cmd.exe" -ldflags "-s -w -X main.version=$Version" ".\cmd\$cmd"
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Failed to build $cmd" -ForegroundColor Red
        exit 1
    }
}

Write-Host "`nBuild complete!" -ForegroundColor Green
Write-Host "Output: $OutputDir" -ForegroundColor Green
Get-ChildItem "$OutputDir\*.exe" | ForEach-Object {
    Write-Host "  - $($_.Name) ($([math]::Round($_.Length/1KB))KB)" -ForegroundColor Gray
}
