Set-Location bin

Write-Host "Generating SHA256 checksums..." -ForegroundColor Green

$files = Get-ChildItem fleeks-* -File
$checksums = @()

foreach ($file in $files) {
    Write-Host "  Hashing $($file.Name)..." -ForegroundColor Gray
    $hash = (Get-FileHash $file -Algorithm SHA256).Hash.ToLower()
    $checksums += "$hash  $($file.Name)"
}

$checksums | Out-File -FilePath SHA256SUMS -Encoding ASCII

Write-Host "`nChecksums generated:" -ForegroundColor Green
Get-Content SHA256SUMS

Write-Host "`n[SUCCESS] Checksums complete!" -ForegroundColor Green
