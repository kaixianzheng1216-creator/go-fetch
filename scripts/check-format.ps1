$ErrorActionPreference = "Stop"

$root = Resolve-Path (Join-Path $PSScriptRoot "..")
$goFiles = Get-ChildItem -Path $root -Recurse -Filter "*.go" -File |
  Where-Object {
    $_.FullName -notlike "*\reference\umami\*" -and
    $_.FullName -notlike "*\frontend\node_modules\*"
  } |
  ForEach-Object { $_.FullName }

if ($goFiles.Count -gt 0) {
  $unformatted = gofmt -l @goFiles
  if ($unformatted) {
    Write-Error "Go files need gofmt:`n$($unformatted -join "`n")"
  }
}

Push-Location (Join-Path $root "frontend")
try {
  npm run format:check
}
finally {
  Pop-Location
}
