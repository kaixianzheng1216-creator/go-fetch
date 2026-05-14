$ErrorActionPreference = "Stop"

$root = Resolve-Path (Join-Path $PSScriptRoot "..")
$frontend = Join-Path $root "frontend"
$trackingScript = Join-Path $root "internal\web\static\script.js"

$goFiles = Get-ChildItem -Path $root -Recurse -Filter "*.go" -File |
  Where-Object {
    $_.FullName -notlike "*\frontend\node_modules\*"
  } |
  ForEach-Object { $_.FullName }

if ($goFiles.Count -gt 0) {
  gofmt -w @goFiles
}

Push-Location $frontend
try {
  npm run format
}
finally {
  Pop-Location
}

npm --prefix $frontend exec prettier -- --write $trackingScript
