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
  $unformatted = gofmt -l @goFiles
  Assert-LastExitCode
  if ($unformatted) {
    Write-Error "The following Go files need gofmt:`n$($unformatted -join "`n")"
  }
}

Push-Location $frontend
try {
  npm run format:check
  Assert-LastExitCode
}
finally {
  Pop-Location
}

npm --prefix $frontend exec prettier -- --check $trackingScript
Assert-LastExitCode
