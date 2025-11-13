# Fleeks CLI Windows Installer - CDN Distribution
# Usage: iwr -useb http://downloads.fleeks.ai/install.ps1 | iex

$ErrorActionPreference = "Stop"

# Configuration
$CDN_BASE_URL = "http://downloads.fleeks.ai"
$BINARY_NAME = "fleeks-windows-amd64.exe"

Write-Host "╔════════════════════════════════════════╗" -ForegroundColor Blue
Write-Host "║   Fleeks CLI Installation Script       ║" -ForegroundColor Blue
Write-Host "╚════════════════════════════════════════╝" -ForegroundColor Blue
Write-Host ""

# Get latest version
Write-Host "[INFO] Fetching latest version..." -ForegroundColor Green
try {
    $response = Invoke-RestMethod -Uri "https://api.github.com/repos/fleeks-ai/fleeks-cli/releases/latest" -ErrorAction SilentlyContinue
    $version = $response.tag_name -replace '^v', ''
    Write-Host "[INFO] Latest version: v$version" -ForegroundColor Green
    $downloadUrl = "$CDN_BASE_URL/v$version/$BINARY_NAME"
} catch {
    Write-Host "[WARN] Could not determine latest version, using 'latest' folder" -ForegroundColor Yellow
    $version = "latest"
    $downloadUrl = "$CDN_BASE_URL/latest/$BINARY_NAME"
}
Write-Host ""

# Download binary
$tempFile = "$env:TEMP\fleeks.exe"

Write-Host "[INFO] Downloading from: $downloadUrl" -ForegroundColor Green
try {
    Invoke-WebRequest -Uri $downloadUrl -OutFile $tempFile -ErrorAction Stop
    Write-Host "[SUCCESS] Downloaded successfully" -ForegroundColor Green
} catch {
    Write-Host "[ERROR] Download failed: $_" -ForegroundColor Red
    exit 1
}

# Verify checksum (optional)
if ($version -ne "latest") {
    Write-Host "[INFO] Verifying checksum..." -ForegroundColor Green
    try {
        $checksumUrl = "$CDN_BASE_URL/v$version/SHA256SUMS.txt"
        $checksums = (Invoke-WebRequest -Uri $checksumUrl -ErrorAction SilentlyContinue).Content
        $expectedHash = ($checksums -split "`n" | Select-String $BINARY_NAME).ToString().Split()[0]
        $actualHash = (Get-FileHash $tempFile -Algorithm SHA256).Hash.ToLower()
        
        if ($expectedHash -eq $actualHash) {
            Write-Host "[SUCCESS] Checksum verified" -ForegroundColor Green
        } else {
            Write-Host "[WARN] Checksum verification failed" -ForegroundColor Yellow
        }
    } catch {
        Write-Host "[WARN] Checksum verification not available" -ForegroundColor Yellow
    }
}

# Install to user bin directory
$installDir = "$env:USERPROFILE\.fleeks\bin"
$installPath = "$installDir\fleeks.exe"

Write-Host "[INFO] Installing to $installDir..." -ForegroundColor Green
New-Item -ItemType Directory -Force -Path $installDir | Out-Null
Move-Item -Force $tempFile $installPath
Write-Host "[SUCCESS] Installed to $installPath" -ForegroundColor Green

# Add to PATH if not already there
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -notlike "*$installDir*") {
    Write-Host "[INFO] Adding to PATH..." -ForegroundColor Green
    [Environment]::SetEnvironmentVariable("Path", "$userPath;$installDir", "User")
    $env:Path = "$env:Path;$installDir"
    Write-Host "[SUCCESS] Added to PATH" -ForegroundColor Green
} else {
    Write-Host "[INFO] Already in PATH" -ForegroundColor Green
}

Write-Host ""
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Blue
Write-Host "Installation complete!" -ForegroundColor Green
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Blue
Write-Host ""
Write-Host "Get started with:" -ForegroundColor Yellow
Write-Host "  fleeks --version" -ForegroundColor Gray
Write-Host "  fleeks auth login" -ForegroundColor Gray
Write-Host ""
Write-Host "Note: You may need to restart your terminal for PATH changes to take effect." -ForegroundColor Yellow
Write-Host ""
