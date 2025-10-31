# Pass Chain E2E Test Suite (PowerShell)

Write-Host "🧪 Pass Chain E2E Test Suite" -ForegroundColor Cyan
Write-Host "==============================" -ForegroundColor Cyan

$ErrorActionPreference = "Stop"

# Step 1: Start test infrastructure
Write-Host ""
Write-Host "📦 Step 1: Starting test infrastructure..." -ForegroundColor Yellow
docker-compose -f docker-compose.test.yml up -d

# Wait for services
Write-Host "⏳ Waiting for services to be ready..." -ForegroundColor Yellow
Start-Sleep -Seconds 10

# Check PostgreSQL
Write-Host "   Checking PostgreSQL... " -NoNewline
$retries = 30
$ready = $false
for ($i = 0; $i -lt $retries; $i++) {
    try {
        docker exec passchain-test-db pg_isready -U passchain 2>$null
        if ($LASTEXITCODE -eq 0) {
            $ready = $true
            break
        }
    } catch {}
    Start-Sleep -Seconds 1
}
if ($ready) {
    Write-Host "✓" -ForegroundColor Green
} else {
    Write-Host "✗" -ForegroundColor Red
    exit 1
}

# Check Redis
Write-Host "   Checking Redis... " -NoNewline
$retries = 30
$ready = $false
for ($i = 0; $i -lt $retries; $i++) {
    try {
        docker exec passchain-test-redis redis-cli ping 2>$null
        if ($LASTEXITCODE -eq 0) {
            $ready = $true
            break
        }
    } catch {}
    Start-Sleep -Seconds 1
}
if ($ready) {
    Write-Host "✓" -ForegroundColor Green
} else {
    Write-Host "✗" -ForegroundColor Red
    exit 1
}

# Check Vault
Write-Host "   Checking Vault... " -NoNewline
$retries = 30
$ready = $false
for ($i = 0; $i -lt $retries; $i++) {
    try {
        docker exec passchain-test-vault vault status 2>$null
        if ($LASTEXITCODE -ne 127) {  # Vault status returns non-zero when sealed/ready
            $ready = $true
            break
        }
    } catch {}
    Start-Sleep -Seconds 1
}
if ($ready) {
    Write-Host "✓" -ForegroundColor Green
} else {
    Write-Host "✗" -ForegroundColor Red
    exit 1
}

Write-Host "✅ Infrastructure ready" -ForegroundColor Green

# Step 2: Run unit tests
Write-Host ""
Write-Host "🔬 Step 2: Running unit tests..." -ForegroundColor Yellow
Set-Location backend
go test ./internal/api/handlers -v -short
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Unit tests failed" -ForegroundColor Red
    Set-Location ..
    exit 1
}
Write-Host "✅ Unit tests passed" -ForegroundColor Green

# Step 3: Run integration tests
Write-Host ""
Write-Host "🔗 Step 3: Running integration tests..." -ForegroundColor Yellow
go test ./test -v -run TestIntegration
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Integration tests failed" -ForegroundColor Red
    Set-Location ..
    docker-compose -f docker-compose.test.yml logs
    exit 1
}
Write-Host "✅ Integration tests passed" -ForegroundColor Green

# Step 4: Run E2E flow tests
Write-Host ""
Write-Host "🎬 Step 4: Running E2E flow tests..." -ForegroundColor Yellow
go test ./test -v -run TestEndToEndFlow
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ E2E tests failed" -ForegroundColor Red
    Set-Location ..
    docker-compose -f docker-compose.test.yml logs
    exit 1
}
Write-Host "✅ E2E tests passed" -ForegroundColor Green

# Step 5: Generate coverage report
Write-Host ""
Write-Host "📊 Step 5: Generating coverage report..." -ForegroundColor Yellow
go test ./... -coverprofile=coverage.out -covermode=atomic
go tool cover -html=coverage.out -o coverage.html
$coverage = (go tool cover -func=coverage.out | Select-String "total" | ForEach-Object { $_.Line.Split()[2] })
Write-Host "   Coverage: $coverage" -ForegroundColor Green

# Step 6: Cleanup
Write-Host ""
Write-Host "🧹 Step 6: Cleanup..." -ForegroundColor Yellow
Set-Location ..
docker-compose -f docker-compose.test.yml down -v

Write-Host ""
Write-Host "==============================" -ForegroundColor Cyan
Write-Host "✅ All tests passed!" -ForegroundColor Green
Write-Host "==============================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Coverage report: backend/coverage.html" -ForegroundColor Cyan

