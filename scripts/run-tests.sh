#!/bin/bash

set -e

echo "🧪 Pass Chain E2E Test Suite"
echo "=============================="

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Step 1: Start test infrastructure
echo ""
echo "📦 Step 1: Starting test infrastructure..."
docker-compose -f docker-compose.test.yml up -d

# Wait for services to be healthy
echo "⏳ Waiting for services to be ready..."
sleep 10

# Check PostgreSQL
echo -n "   Checking PostgreSQL... "
until docker exec passchain-test-db pg_isready -U passchain > /dev/null 2>&1; do
    sleep 1
done
echo -e "${GREEN}✓${NC}"

# Check Redis
echo -n "   Checking Redis... "
until docker exec passchain-test-redis redis-cli ping > /dev/null 2>&1; do
    sleep 1
done
echo -e "${GREEN}✓${NC}"

# Check Vault
echo -n "   Checking Vault... "
until docker exec passchain-test-vault vault status > /dev/null 2>&1; do
    sleep 1
done
echo -e "${GREEN}✓${NC}"

echo -e "${GREEN}✅ Infrastructure ready${NC}"

# Step 2: Run unit tests
echo ""
echo "🔬 Step 2: Running unit tests..."
cd backend
go test ./internal/api/handlers -v -short || {
    echo -e "${RED}❌ Unit tests failed${NC}"
    exit 1
}
echo -e "${GREEN}✅ Unit tests passed${NC}"

# Step 3: Run integration tests
echo ""
echo "🔗 Step 3: Running integration tests..."
go test ./test -v -run TestIntegration || {
    echo -e "${RED}❌ Integration tests failed${NC}"
    docker-compose -f ../docker-compose.test.yml logs
    exit 1
}
echo -e "${GREEN}✅ Integration tests passed${NC}"

# Step 4: Run E2E flow tests
echo ""
echo "🎬 Step 4: Running E2E flow tests..."
go test ./test -v -run TestEndToEndFlow || {
    echo -e "${RED}❌ E2E tests failed${NC}"
    docker-compose -f ../docker-compose.test.yml logs
    exit 1
}
echo -e "${GREEN}✅ E2E tests passed${NC}"

# Step 5: Run smoke tests
echo ""
echo "💨 Step 5: Running smoke tests..."
cd ..

# Build smoke test
go build -o bin/smoke_test ./backend/cmd/smoke_test

# Start backend in test mode (background)
echo "   Starting backend..."
cd backend
DATABASE_HOST=localhost \
DATABASE_PORT=5433 \
DATABASE_NAME=passchain_test \
VAULT_ADDR=http://localhost:8201 \
VAULT_TOKEN=test-token \
PORT=8080 \
go run cmd/server/main.go &
BACKEND_PID=$!

# Wait for backend to start
sleep 5

# Run smoke tests
cd ..
./bin/smoke_test || {
    echo -e "${RED}❌ Smoke tests failed${NC}"
    kill $BACKEND_PID 2>/dev/null || true
    exit 1
}

# Kill backend
kill $BACKEND_PID 2>/dev/null || true

echo -e "${GREEN}✅ Smoke tests passed${NC}"

# Step 6: Generate coverage report
echo ""
echo "📊 Step 6: Generating coverage report..."
cd backend
go test ./... -coverprofile=coverage.out -covermode=atomic
go tool cover -html=coverage.out -o coverage.html
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
echo -e "   Coverage: ${GREEN}$COVERAGE${NC}"

# Step 7: Cleanup
echo ""
echo "🧹 Step 7: Cleanup..."
cd ..
docker-compose -f docker-compose.test.yml down -v

echo ""
echo "=============================="
echo -e "${GREEN}✅ All tests passed!${NC}"
echo "=============================="

