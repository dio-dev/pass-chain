# Pass Chain - Minikube Quick Start
Write-Host "Pass Chain v2.0 - Minikube Deployment" -ForegroundColor Cyan

# Check prerequisites
if (!(Get-Command minikube -ErrorAction SilentlyContinue)) {
    Write-Host "Minikube not found. Install from: https://minikube.sigs.k8s.io/docs/start/" -ForegroundColor Red
    exit 1
}

if (!(Get-Command kubectl -ErrorAction SilentlyContinue)) {
    Write-Host "kubectl not found. Install from: https://kubernetes.io/docs/tasks/tools/" -ForegroundColor Red
    exit 1
}

if (!(Get-Command helm -ErrorAction SilentlyContinue)) {
    Write-Host "Helm not found. Install from: https://helm.sh/docs/intro/install/" -ForegroundColor Red
    exit 1
}

Write-Host "Prerequisites found!" -ForegroundColor Green

# Start Minikube
Write-Host "Starting Minikube..." -ForegroundColor Yellow
minikube start --cpus=4 --memory=8192 --disk-size=40g --driver=docker

# Enable addons
minikube addons enable storage-provisioner
minikube addons enable default-storageclass
minikube addons enable metrics-server

# Add Helm repos
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add hashicorp https://helm.releases.hashicorp.com
helm repo update

# Build images
Write-Host "Building Docker images (5-10 minutes)..." -ForegroundColor Yellow
& minikube docker-env | Invoke-Expression
docker build -t passchain/backend:2.0.0 .\backend -q
docker build -t passchain/frontend:2.0.0 .\frontend -q

# Create namespace
kubectl create namespace passchain --dry-run=client -o yaml | kubectl apply -f -

# Deploy Pass Chain
Write-Host "Deploying Pass Chain (5-10 minutes)..." -ForegroundColor Yellow
helm install passchain .\infrastructure\k8s\helm\passchain --namespace passchain --values .\infrastructure\k8s\helm\passchain\values-minikube.yaml --timeout 15m

# Wait for ready
Write-Host "Waiting for services..." -ForegroundColor Yellow
kubectl wait --for=condition=ready pod --all -n passchain --timeout=600s

Write-Host "Pass Chain deployed successfully!" -ForegroundColor Green
$MINIKUBE_IP = minikube ip
Write-Host "Frontend: http://${MINIKUBE_IP}:30300" -ForegroundColor Cyan
Write-Host "Backend:  http://${MINIKUBE_IP}:30800" -ForegroundColor Cyan
Write-Host "Vault:    http://${MINIKUBE_IP}:30200" -ForegroundColor Cyan

Write-Host "AUUUUFFFF! Ready!" -ForegroundColor Magenta

# Open frontend
$response = Read-Host "Open frontend? (Y/n)"
if ($response -eq "" -or $response -eq "Y" -or $response -eq "y") {
    minikube service passchain-frontend -n passchain
}
