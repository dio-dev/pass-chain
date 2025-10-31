# Testing Guide

## Quick Start

```bash
# Start test infrastructure
docker-compose -f docker-compose.test.yml up -d

# Run all tests
cd backend && go test ./... -v

# Cleanup
docker-compose -f docker-compose.test.yml down -v
```

## Using Makefile (Recommended)

```bash
make test-all      # Run all tests with infrastructure
make test-unit     # Unit tests only
make test-e2e      # E2E flow tests
make test-coverage # Generate coverage report
make smoke         # Quick smoke test
```

## Manual Testing

### 1. Unit Tests (no infrastructure needed)
```bash
cd backend
go test ./internal/api/handlers -v
go test ./internal/services -v
```

### 2. Integration Tests (requires DB, Vault, Redis)
```bash
# Start infrastructure
docker-compose -f docker-compose.test.yml up -d

# Run tests
cd backend
go test ./test -v -run TestIntegration

# Cleanup
docker-compose -f docker-compose.test.yml down -v
```

### 3. E2E Flow Tests
```bash
cd backend
go test ./test -v -run TestEndToEndFlow
```

This tests the complete flow:
- User registration
- Organization creation
- Project creation
- Member invitation & acceptance
- Vault creation
- Credential CRUD
- RBAC enforcement
- Member removal
- Audit logs

### 4. Smoke Tests (against running backend)
```bash
cd backend
go run cmd/smoke_test/main.go
```

## Test Database

The test database runs on port `5433` (not `5432`) to avoid conflicts:

```yaml
Host: localhost
Port: 5433
Database: passchain_test
User: passchain
Password: passchain123
```

## Coverage

```bash
cd backend
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
# Open coverage.html in browser
```

## CI/CD

Tests run automatically on:
- Pull requests
- Pushes to main

See `.github/workflows/test.yml`

## Debugging

```bash
# Run specific test
go test -v -run TestCreateOrganization

# With race detector
go test -race ./...

# Verbose with logs
go test -v ./... 2>&1 | tee test.log
```

## Common Issues

### "invalid input syntax for type uuid"
- Fixed: All UUIDs now explicitly generated in services
- If you see this, check that user.ID is not empty before creating org

### "connection refused"
- Make sure test infrastructure is running
- Check: `docker-compose -f docker-compose.test.yml ps`

### "table does not exist"
- Run migrations: `make test-setup` starts infra and runs migrations automatically

## AUUUUFFFF! 🔥
