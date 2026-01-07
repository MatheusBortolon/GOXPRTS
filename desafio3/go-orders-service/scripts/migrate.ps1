$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent $PSCommandPath
$ProjectDir = Split-Path -Parent $ScriptDir
$MigrationsDir = Join-Path $ProjectDir "migrations"
Write-Host "Running database migrations..."
docker run --rm -v "${MigrationsDir}:/migrations" --network host migrate/migrate `
  -path=/migrations/ `
  -database "postgres://user:password@localhost:5432/orders?sslmode=disable" up
Write-Host "Migrations completed."
