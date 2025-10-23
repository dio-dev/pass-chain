# Pass Chain - Quick Start Script
# Run this script to start all services

Write-Host "🚀 Starting Pass Chain..." -ForegroundColor Cyan

# Navigate to docker directory
Set-Location -Path "$PSScriptRoot\infrastructure\docker"

Write-Host "📦 Building Docker images..." -ForegroundColor Yellow
docker-compose build

Write-Host "🔄 Starting services..." -ForegroundColor Yellow
docker-compose up -d

Write-Host ""
Write-Host "⏳ Waiting for services to be ready..." -ForegroundColor Yellow
Start-Sleep -Seconds 10

Write-Host ""
Write-Host "✅ Services are starting!" -ForegroundColor Green
Write-Host ""
Write-Host "📋 Service URLs:" -ForegroundColor Cyan
Write-Host "  Frontend:  http://localhost:3000" -ForegroundColor White
Write-Host "  Backend:   http://localhost:8080" -ForegroundColor White
Write-Host "  Vault UI:  http://localhost:8200" -ForegroundColor White
Write-Host ""
Write-Host "🔍 Check service status:" -ForegroundColor Cyan
docker-compose ps

Write-Host ""
Write-Host "📝 View logs with: docker-compose logs -f" -ForegroundColor Yellow
Write-Host "🛑 Stop services with: docker-compose down" -ForegroundColor Yellow
Write-Host ""
Write-Host "AUUUUFFFF! 🎉" -ForegroundColor Magenta

