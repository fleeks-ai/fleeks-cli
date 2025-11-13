# Fleeks CLI Windows Installer
# Usage: iwr -useb https://raw.githubusercontent.com/fleeks-ai/fleeks-cli/main/install.ps1 | iex

$ErrorActionPreference = "Stop"

Write-Host "Fleeks CLI Installer" -ForegroundColor Green
Write-Host ""

# Get latest version
Write-Host "Fetching latest version..."
try {
    $response = Invoke-RestMethod -Uri "https://api.github.com/repos/fleeks-ai/fleeks-cli/releases/latest"
    $version = $response.tag_name
    Write-Host "Latest version: $version" -ForegroundColor Green
    Write-Host ""
} catch {
    Write-Host "Failed to fetch latest version" -ForegroundColor Red
    exit 1
}

# Download URL
$downloadUrl = "https://github.com/fleeks-ai/fleeks-cli/releases/download/$version/fleeks-windows-amd64.exe"
$tempFile = "$env:TEMP\fleeks.exe"

Write-Host "Downloading Fleeks CLI..."
try {
    Invoke-WebRequest -Uri $downloadUrl -OutFile $tempFile
} catch {
    Write-Host "Download failed: $_" -ForegroundColor Red
    exit 1
}

# Install to user bin directory
$installDir = "$env:USERPROFILE\.fleeks\bin"
$installPath = "$installDir\fleeks.exe"

Write-Host "Installing to $installDir..."
New-Item -ItemType Directory -Force -Path $installDir | Out-Null
Move-Item -Force $tempFile $installPath

# Add to PATH if not already there
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -notlike "*$installDir*") {
    Write-Host "Adding to PATH..."
    [Environment]::SetEnvironmentVariable("Path", "$userPath;$installDir", "User")
    $env:Path = "$env:Path;$installDir"
}

Write-Host ""
Write-Host "✓ Fleeks CLI installed successfully!" -ForegroundColor Green
Write-Host ""
Write-Host "Installation location: $installPath" -ForegroundColor Cyan
Write-Host ""
Write-Host "Get started with:" -ForegroundColor Yellow
Write-Host "  fleeks --version" -ForegroundColor Gray
Write-Host "  fleeks auth login" -ForegroundColor Gray
Write-Host ""
Write-Host "Note: You may need to restart your terminal for PATH changes to take effect." -ForegroundColor Yellow
