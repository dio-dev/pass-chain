# Pass Chain - Quick Update Script
# Updates frontend or backend with new code

param(
    [Parameter(Mandatory=$true)]
    [ValidateSet("frontend", "backend", "both")]
    [string]$Component
)

Write-Host "Pass Chain - Quick Update" -ForegroundColor Cyan
Write-Host ""

# Point to Minikube Docker
Write-Host "Configuring Docker for Minikube..." -ForegroundColor Yellow
& minikube docker-env | Invoke-Expression

function Update-Component {
    param([string]$Name)
    
    Write-Host ""
    Write-Host "Building $Name..." -ForegroundColor Yellow
    docker build -t "passchain/${Name}:2.0.0" ".\$Name" -q
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "$Name built successfully!" -ForegroundColor Green
        
        Write-Host "Restarting deployment..." -ForegroundColor Yellow
        kubectl rollout restart "deployment/passchain-$Name" -n passchain
        
        Write-Host "Waiting for rollout..." -ForegroundColor Yellow
        kubectl rollout status "deployment/passchain-$Name" -n passchain --timeout=60s
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "$Name updated successfully!" -ForegroundColor Green
        } else {
            Write-Host "$Name rollout failed!" -ForegroundColor Red
        }
    } else {
        Write-Host "$Name build failed!" -ForegroundColor Red
    }
}

# Update components
if ($Component -eq "both") {
    Update-Component "backend"
    Update-Component "frontend"
} else {
    Update-Component $Component
}

Write-Host ""
Write-Host "Current pod status:" -ForegroundColor Cyan
kubectl get pods -n passchain

Write-Host ""
Write-Host "Update complete!" -ForegroundColor Green
Write-Host ""
Write-Host "Access your app:" -ForegroundColor Cyan
$MINIKUBE_IP = minikube ip
Write-Host "   Frontend: http://${MINIKUBE_IP}:30300" -ForegroundColor White
Write-Host "   Backend:  http://${MINIKUBE_IP}:30800" -ForegroundColor White
Write-Host ""
