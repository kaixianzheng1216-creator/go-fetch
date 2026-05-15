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
  $unformatted = gofmt -l @goFiles
  Assert-LastExitCode
  if ($unformatted) {
    Write-Error "The following Go files need gofmt:`n$($unformatted -join "`n")"
  }
}

Push-Location $dashboard
try {
  npm run format:check
  Assert-LastExitCode
}
finally {
  Pop-Location
}

npm --prefix $dashboard exec prettier -- --check $trackingScript
Assert-LastExitCode
