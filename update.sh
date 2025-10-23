#!/bin/bash
# Pass Chain - Quick Update Script (Linux/Mac)
# Updates frontend or backend with new code

set -e

COMPONENT=$1

if [ -z "$COMPONENT" ]; then
    echo "Usage: ./update.sh [frontend|backend|both]"
    exit 1
fi

echo "🔄 Pass Chain - Quick Update"
echo ""

# Point to Minikube Docker
echo "📦 Configuring Docker for Minikube..."
eval $(minikube docker-env)

update_component() {
    NAME=$1
    echo ""
    echo "🔨 Building $NAME..."
    docker build -t "passchain/$NAME:2.0.0" "./$NAME" -q
    
    if [ $? -eq 0 ]; then
        echo "✅ $NAME built successfully!"
        
        echo "🔄 Restarting deployment..."
        kubectl rollout restart "deployment/passchain-$NAME" -n passchain
        
        echo "⏳ Waiting for rollout..."
        kubectl rollout status "deployment/passchain-$NAME" -n passchain --timeout=60s
        
        if [ $? -eq 0 ]; then
            echo "✅ $NAME updated successfully!"
        else
            echo "❌ $NAME rollout failed!"
        fi
    else
        echo "❌ $NAME build failed!"
    fi
}

# Update components
if [ "$COMPONENT" = "both" ]; then
    update_component "backend"
    update_component "frontend"
else
    update_component "$COMPONENT"
fi

echo ""
echo "📊 Current pod status:"
kubectl get pods -n passchain

echo ""
echo "✅ Update complete!"
echo ""
echo "🌐 Access your app:"
MINIKUBE_IP=$(minikube ip)
echo "   Frontend: http://$MINIKUBE_IP:30300"
echo "   Backend:  http://$MINIKUBE_IP:30800"
echo ""

