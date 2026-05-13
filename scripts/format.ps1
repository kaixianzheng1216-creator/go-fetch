$ErrorActionPreference = "Stop"

$root = Resolve-Path (Join-Path $PSScriptRoot "..")
$goFiles = Get-ChildItem -Path $root -Recurse -Filter "*.go" -File |
  Where-Object {
    $_.FullName -notlike "*\frontend\node_modules\*"
  } |
  ForEach-Object { $_.FullName }

if ($goFiles.Count -gt 0) {
  gofmt -w @goFiles
}

Push-Location (Join-Path $root "frontend")
try {
  npm run format
}
finally {
  Pop-Location
}
