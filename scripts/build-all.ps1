param(
    [string]$Version = "dev"
)

$ErrorActionPreference = "Stop"
$Commit = git rev-parse --short HEAD
$BuildTime = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")
$LdFlags = "-s -w -X 'main.Version=$Version' -X 'main.BuildDate=$BuildTime'"

Write-Host "Building Fleeks CLI v$Version..." -ForegroundColor Green
Write-Host "Commit: $Commit" -ForegroundColor Gray
Write-Host "Build Time: $BuildTime" -ForegroundColor Gray
Write-Host ""

# Create bin directory
New-Item -ItemType Directory -Force -Path bin | Out-Null

# Windows AMD64
Write-Host "Building Windows AMD64..." -ForegroundColor Yellow
$env:GOOS = "windows"; $env:GOARCH = "amd64"; $env:CGO_ENABLED = "0"
go build -trimpath -ldflags $LdFlags -o bin/fleeks-windows-amd64.exe .

# macOS Intel
Write-Host "Building macOS Intel..." -ForegroundColor Yellow
$env:GOOS = "darwin"; $env:GOARCH = "amd64"
go build -trimpath -ldflags $LdFlags -o bin/fleeks-darwin-amd64 .

# macOS Apple Silicon
Write-Host "Building macOS Apple Silicon..." -ForegroundColor Yellow
$env:GOOS = "darwin"; $env:GOARCH = "arm64"
go build -trimpath -ldflags $LdFlags -o bin/fleeks-darwin-arm64 .

# Linux AMD64
Write-Host "Building Linux AMD64..." -ForegroundColor Yellow
$env:GOOS = "linux"; $env:GOARCH = "amd64"
go build -trimpath -ldflags $LdFlags -o bin/fleeks-linux-amd64 .

# Linux ARM64
Write-Host "Building Linux ARM64..." -ForegroundColor Yellow
$env:GOOS = "linux"; $env:GOARCH = "arm64"
go build -trimpath -ldflags $LdFlags -o bin/fleeks-linux-arm64 .

Write-Host ""
Write-Host "[SUCCESS] Build complete! Binaries in bin/" -ForegroundColor Green
Write-Host ""
Get-ChildItem bin/ | Format-Table Name, Length
