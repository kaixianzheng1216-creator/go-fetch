$ErrorActionPreference = "Stop"

$root = Resolve-Path (Join-Path $PSScriptRoot "..")
$dashboard = Join-Path $root "web\dashboard"
$trackingScript = Join-Path $root "web\static\tracker.js"

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

Push-Location $dashboard
try {
  npm run format
  Assert-LastExitCode
}
finally {
  Pop-Location
}

npm --prefix $dashboard exec prettier -- --write $trackingScript
Assert-LastExitCode
