# Fleeks CLI Build Script
# Usage: .\build.ps1 [version]

param(
    [string]$Version = "dev",
    [switch]$Release
)

Write-Host "🚀 Building Fleeks CLI..." -ForegroundColor Cyan

# Get build info
$GitCommit = git rev-parse --short HEAD 2>$null
if (-not $GitCommit) { $GitCommit = "unknown" }

$BuildTime = Get-Date -Format "yyyy-MM-dd HH:mm:ss"

# Build flags
$ldflags = "-X 'github.com/fleeks/cmd.Version=$Version' " +
           "-X 'github.com/fleeks/cmd.GitCommit=$GitCommit' " +
           "-X 'github.com/fleeks/cmd.BuildTime=$BuildTime'"

if ($Release) {
    # Production build with optimization
    $ldflags += " -s -w"
    Write-Host "Building RELEASE version $Version..." -ForegroundColor Green
} else {
    Write-Host "Building DEV version $Version..." -ForegroundColor Yellow
}

# Create bin directory
if (-not (Test-Path "bin")) {
    New-Item -ItemType Directory -Path "bin" | Out-Null
}

# Build for current platform
Write-Host "Building for Windows (amd64)..." -ForegroundColor Cyan
go build -ldflags="$ldflags" -o "bin\fleeks.exe" main.go

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Build successful!" -ForegroundColor Green
    Write-Host "   Location: bin\fleeks.exe" -ForegroundColor White
    Write-Host "   Version:  $Version" -ForegroundColor White
    Write-Host "   Commit:   $GitCommit" -ForegroundColor White
    Write-Host "   Built:    $BuildTime" -ForegroundColor White
    
    # Show file size
    $size = (Get-Item "bin\fleeks.exe").Length / 1MB
    Write-Host "   Size:     $([math]::Round($size, 2)) MB" -ForegroundColor White
    
    Write-Host "`n📖 Test it with:" -ForegroundColor Cyan
    Write-Host "   .\bin\fleeks.exe --environment development" -ForegroundColor Yellow
} else {
    Write-Host "❌ Build failed!" -ForegroundColor Red
    exit 1
}

# Optionally build for other platforms
if ($Release) {
    Write-Host "`n🌍 Building for other platforms..." -ForegroundColor Cyan
    
    # macOS Intel
    Write-Host "Building for macOS (amd64)..." -ForegroundColor Cyan
    $env:GOOS = "darwin"
    $env:GOARCH = "amd64"
    go build -ldflags="$ldflags" -o "bin\fleeks-darwin-amd64" main.go
    
    # macOS Apple Silicon
    Write-Host "Building for macOS (arm64)..." -ForegroundColor Cyan
    $env:GOOS = "darwin"
    $env:GOARCH = "arm64"
    go build -ldflags="$ldflags" -o "bin\fleeks-darwin-arm64" main.go
    
    # Linux
    Write-Host "Building for Linux (amd64)..." -ForegroundColor Cyan
    $env:GOOS = "linux"
    $env:GOARCH = "amd64"
    go build -ldflags="$ldflags" -o "bin\fleeks-linux-amd64" main.go
    
    # Reset environment
    $env:GOOS = ""
    $env:GOARCH = ""
    
    Write-Host "`n✅ All builds completed!" -ForegroundColor Green
    Get-ChildItem "bin" | ForEach-Object {
        $fileSize = $_.Length / 1MB
        $sizeMB = [math]::Round($fileSize, 2)
        $fileName = $_.Name
        $output = "   " + $fileName + ": " + $sizeMB + " MB"
        Write-Host $output -ForegroundColor White
    }
}
