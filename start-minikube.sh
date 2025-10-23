#!/bin/bash
# Pass Chain - Minikube Quick Start (Linux/macOS)

echo "🚀 Pass Chain v2.0 - Minikube Deployment"
echo ""

# Check prerequisites
echo "📋 Checking prerequisites..."

if ! command -v minikube &> /dev/null; then
    echo "❌ Minikube not found. Install from: https://minikube.sigs.k8s.io/docs/start/"
    exit 1
fi

if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl not found. Install from: https://kubernetes.io/docs/tasks/tools/"
    exit 1
fi

if ! command -v helm &> /dev/null; then
    echo "❌ Helm not found. Install from: https://helm.sh/docs/intro/install/"
    exit 1
fi

echo "✅ All prerequisites found!"
echo ""

# Start Minikube
echo "🎯 Starting Minikube..."
minikube start --cpus=4 --memory=8192 --disk-size=40g --driver=docker

if [ $? -ne 0 ]; then
    echo "❌ Failed to start Minikube"
    exit 1
fi

# Enable addons
echo "🔧 Enabling Minikube addons..."
minikube addons enable storage-provisioner
minikube addons enable default-storageclass
minikube addons enable metrics-server

# Add Helm repos
echo "📦 Adding Helm repositories..."
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add hashicorp https://helm.releases.hashicorp.com
helm repo update

# Build images
echo "🐳 Building Docker images..."
echo "  (This may take 5-10 minutes)"

# Point to Minikube's Docker
eval $(minikube docker-env)

# Build backend
echo "  Building backend..."
docker build -t passchain/backend:2.0.0 ./backend

# Build frontend
echo "  Building frontend..."
docker build -t passchain/frontend:2.0.0 ./frontend

echo "✅ Images built successfully!"
echo ""

# Create namespace
echo "📁 Creating Kubernetes namespace..."
kubectl create namespace passchain --dry-run=client -o yaml | kubectl apply -f -

# Deploy Pass Chain
echo "🚀 Deploying Pass Chain to Minikube..."
echo "  (This may take 5-10 minutes)"

helm install passchain ./infrastructure/k8s/helm/passchain \
  --namespace passchain \
  --values ./infrastructure/k8s/helm/passchain/values-minikube.yaml \
  --timeout 15m

if [ $? -ne 0 ]; then
    echo "❌ Failed to deploy Pass Chain"
    exit 1
fi

echo ""
echo "⏳ Waiting for services to be ready..."
kubectl wait --for=condition=ready pod --all -n passchain --timeout=600s

echo ""
echo "✅ Pass Chain deployed successfully!"
echo ""

# Get service URLs
echo "📋 Service URLs:"
MINIKUBE_IP=$(minikube ip)

echo "  Frontend:  http://${MINIKUBE_IP}:30300"
echo "  Backend:   http://${MINIKUBE_IP}:30800"
echo "  Vault:     http://${MINIKUBE_IP}:30200"
echo ""

echo "🔍 Quick Commands:"
echo "  Open Frontend:  minikube service passchain-frontend -n passchain"
echo "  View Pods:      kubectl get pods -n passchain"
echo "  View Logs:      kubectl logs -f deployment/passchain-backend -n passchain"
echo "  Stop Minikube:  minikube stop"
echo "  Delete All:     minikube delete"
echo ""

echo "📚 Documentation:"
echo "  Complete Guide: MINIKUBE_SETUP.md"
echo "  GKE Migration:  GKE_MIGRATION.md"
echo ""

echo "🎉 AUUUUFFFF! Pass Chain is running on Minikube!"
echo ""

# Ask to open frontend
read -p "Open frontend in browser? (Y/n): " response
if [ -z "$response" ] || [ "$response" = "Y" ] || [ "$response" = "y" ]; then
    minikube service passchain-frontend -n passchain
fi

