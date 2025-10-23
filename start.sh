#!/bin/bash
# Pass Chain - Quick Start Script (Linux/macOS)

echo "🚀 Starting Pass Chain..."

# Navigate to docker directory
cd "$(dirname "$0")/infrastructure/docker"

echo "📦 Building Docker images..."
docker-compose build

echo "🔄 Starting services..."
docker-compose up -d

echo ""
echo "⏳ Waiting for services to be ready..."
sleep 10

echo ""
echo "✅ Services are starting!"
echo ""
echo "📋 Service URLs:"
echo "  Frontend:  http://localhost:3000"
echo "  Backend:   http://localhost:8080"
echo "  Vault UI:  http://localhost:8200"
echo ""
echo "🔍 Check service status:"
docker-compose ps

echo ""
echo "📝 View logs with: docker-compose logs -f"
echo "🛑 Stop services with: docker-compose down"
echo ""
echo "AUUUUFFFF! 🎉"




