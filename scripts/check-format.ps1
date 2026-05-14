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
  $unformatted = gofmt -l @goFiles
  if ($unformatted) {
    Write-Error "以下 Go 文件需要执行 gofmt：`n$($unformatted -join "`n")"
  }
}

Push-Location $frontend
try {
  npm run format:check
}
finally {
  Pop-Location
}

npm --prefix $frontend exec prettier -- --check $trackingScript
