param(
    [Parameter(Mandatory=$true)]
    [string]$Version
)

$ErrorActionPreference = "Stop"

Write-Host "Creating release $Version..." -ForegroundColor Green
Write-Host ""

# Check if gh CLI is installed
if (!(Get-Command gh -ErrorAction SilentlyContinue)) {
    Write-Host "Error: GitHub CLI (gh) is not installed" -ForegroundColor Red
    Write-Host "Install from: https://cli.github.com/"
    exit 1
}

# Check if logged in
try {
    gh auth status 2>&1 | Out-Null
} catch {
    Write-Host "Error: Not logged in to GitHub CLI" -ForegroundColor Red
    Write-Host "Run: gh auth login"
    exit 1
}

# 1. Build binaries
Write-Host "Step 1: Building binaries..." -ForegroundColor Yellow
& .\scripts\build-all.ps1 -Version $Version.TrimStart('v')
Write-Host ""

# 2. Generate checksums
Write-Host "Step 2: Generating checksums..." -ForegroundColor Yellow
& .\scripts\generate-checksums.ps1
Write-Host ""

# 3. Create GitHub release
Write-Host "Step 3: Creating GitHub release..." -ForegroundColor Yellow

# Check if release already exists
$releaseExists = $false
try {
    gh release view $Version 2>&1 | Out-Null
    $releaseExists = $true
} catch {}

if ($releaseExists) {
    Write-Host "Release $Version already exists!" -ForegroundColor Yellow
    $response = Read-Host "Delete and recreate? (y/n)"
    if ($response -eq 'y') {
        gh release delete $Version --yes
    } else {
        exit 1
    }
}

# Create release notes
$releaseNotes = @"
## Fleeks CLI $Version

### Installation

#### macOS
``````bash
# Intel Macs
curl -L https://github.com/fleeks-ai/fleeks-cli/releases/download/$Version/fleeks-darwin-amd64 -o fleeks
chmod +x fleeks
sudo mv fleeks /usr/local/bin/

# Apple Silicon (M1/M2/M3)
curl -L https://github.com/fleeks-ai/fleeks-cli/releases/download/$Version/fleeks-darwin-arm64 -o fleeks
chmod +x fleeks
sudo mv fleeks /usr/local/bin/
``````

#### Linux
``````bash
# AMD64
curl -L https://github.com/fleeks-ai/fleeks-cli/releases/download/$Version/fleeks-linux-amd64 -o fleeks
chmod +x fleeks
sudo mv fleeks /usr/local/bin/

# ARM64
curl -L https://github.com/fleeks-ai/fleeks-cli/releases/download/$Version/fleeks-linux-arm64 -o fleeks
chmod +x fleeks
sudo mv fleeks /usr/local/bin/
``````

#### Windows (PowerShell)
``````powershell
Invoke-WebRequest -Uri "https://github.com/fleeks-ai/fleeks-cli/releases/download/$Version/fleeks-windows-amd64.exe" -OutFile fleeks.exe
# Add to PATH or move to a directory in your PATH
``````

### Package Managers

``````bash
# Homebrew (macOS/Linux)
brew install fleeks-ai/fleeks/fleeks

# NPM (Cross-platform)
npm install -g @fleeks/cli
``````

### What's Changed

See the [CHANGELOG](https://github.com/fleeks-ai/fleeks-cli/blob/main/CHANGELOG.md) for full details.

### Verification

Verify your download with SHA256 checksums in ``SHA256SUMS``.
"@

# Save to temp file
$tempFile = New-TemporaryFile
$releaseNotes | Out-File -FilePath $tempFile.FullName -Encoding UTF8

# Create the release
gh release create $Version `
    bin/fleeks-windows-amd64.exe `
    bin/fleeks-darwin-amd64 `
    bin/fleeks-darwin-arm64 `
    bin/fleeks-linux-amd64 `
    bin/fleeks-linux-arm64 `
    bin/SHA256SUMS `
    --title "Release $Version" `
    --notes-file $tempFile.FullName

# Clean up
Remove-Item $tempFile.FullName

Write-Host ""
Write-Host "✓ Release $Version created successfully!" -ForegroundColor Green
Write-Host "View at: https://github.com/fleeks-ai/fleeks-cli/releases/tag/$Version" -ForegroundColor Cyan
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "  1. Update Homebrew formula" -ForegroundColor Gray
Write-Host "  2. Update NPM package" -ForegroundColor Gray
Write-Host "  3. Announce on social media and documentation" -ForegroundColor Gray
