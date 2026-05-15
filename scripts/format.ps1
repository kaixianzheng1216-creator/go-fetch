$ErrorActionPreference = "Stop"

$root = Resolve-Path (Join-Path $PSScriptRoot "..")
$frontend = Join-Path $root "frontend"
$trackingScript = Join-Path $root "internal\static\static\script.js"

function Assert-LastExitCode {
  if ($LASTEXITCODE -ne 0) {
    exit $LASTEXITCODE
  }
}

$goFiles = go list -f '{{range .GoFiles}}{{$.Dir}}/{{.}}{{println}}{{end}}{{range .TestGoFiles}}{{$.Dir}}/{{.}}{{println}}{{end}}' ./...
Assert-LastExitCode
$goFiles = @($goFiles | Where-Object { $_ })

if ($goFiles.Count -gt 0) {
  gofmt -w @goFiles
  Assert-LastExitCode
}

Push-Location $frontend
try {
  npm run format
  Assert-LastExitCode
}
finally {
  Pop-Location
}

npm --prefix $frontend exec prettier -- --write $trackingScript
Assert-LastExitCode
