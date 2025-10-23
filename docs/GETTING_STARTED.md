# Pass Chain - Getting Started Guide

Welcome to Pass Chain! This guide will help you get the project up and running.

## Quick Start (5 minutes)

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/pass-chain.git
cd pass-chain
```

### 2. Start Local Services

```bash
cd infrastructure/docker
docker-compose up -d
```

This starts:
- PostgreSQL (port 5432)
- Redis (port 6379)
- Vault (port 8200)
- Backend API (port 8080)
- Frontend (port 3000)

### 3. Access the Application

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **API Docs**: http://localhost:8080/swagger
- **Vault UI**: http://localhost:8200

### 4. Connect Your Wallet

1. Install MetaMask browser extension
2. Open http://localhost:3000
3. Click "Connect Wallet"
4. Approve the connection

You're ready to go! 🎉

## Development Setup

### Prerequisites

- **Node.js** 18+ and pnpm
- **Go** 1.21+
- **Docker** & Docker Compose
- **Git**

Optional (for full stack):
- **Terraform** 1.6+
- **kubectl** & **Helm**
- **Hyperledger Fabric** binaries

### Backend Development

```bash
cd backend

# Install dependencies
go mod download

# Copy environment file
cp .env.example .env

# Run database migrations (via docker-compose)
docker-compose up -d postgres

# Run the server
go run cmd/server/main.go
```

Backend runs on http://localhost:8080

### Frontend Development

```bash
cd frontend

# Install dependencies
pnpm install

# Copy environment file
cp .env.example .env.local

# Run development server
pnpm dev
```

Frontend runs on http://localhost:3000

### Blockchain Development

```bash
cd blockchain

# Start Fabric network (requires Fabric binaries)
./network.sh up createChannel -ca

# Deploy chaincode
./network.sh deployCC -ccn credentials -ccp ./chaincode/credentials -ccl go

# Test chaincode
./scripts/test-chaincode.sh
```

## Project Structure

```
pass-chain/
├── backend/            # Go API server
├── frontend/           # Next.js web app
├── blockchain/         # Hyperledger Fabric
├── infrastructure/     # Terraform & K8s
├── docs/              # Documentation
├── ProjectSpec.md     # Project specification
├── ProjectStack.md    # Technology stack
└── ProjectPlan.md     # Execution plan
```

## Common Tasks

### Run Tests

```bash
# Backend tests
cd backend && go test ./...

# Frontend tests
cd frontend && pnpm test

# E2E tests
cd frontend && pnpm test:e2e
```

### Build for Production

```bash
# Backend
cd backend
go build -o bin/server cmd/server/main.go

# Frontend
cd frontend
pnpm build
```

### Database Migrations

Migrations run automatically when the backend starts. To run manually:

```bash
cd backend
go run cmd/server/main.go
```

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend
```

## Configuration

### Backend (.env)

```bash
ENVIRONMENT=development
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=passchain
VAULT_ADDR=http://localhost:8200
VAULT_TOKEN=dev-only-token
REDIS_ADDR=localhost:6379
JWT_SECRET=your-secret-key
```

### Frontend (.env.local)

```bash
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_WALLET_CONNECT_PROJECT_ID=your_project_id
```

Get a WalletConnect Project ID from: https://cloud.walletconnect.com

## Troubleshooting

### Port Already in Use

```bash
# Stop all services
docker-compose down

# Check what's using the port
# Windows PowerShell:
netstat -ano | findstr :8080

# Kill the process
taskkill /PID <process_id> /F
```

### Database Connection Failed

```bash
# Restart PostgreSQL
docker-compose restart postgres

# Check logs
docker-compose logs postgres
```

### Vault Connection Failed

```bash
# Check Vault status
docker-compose logs vault

# Restart Vault
docker-compose restart vault
```

### Frontend Build Errors

```bash
# Clear cache and reinstall
cd frontend
rm -rf node_modules .next
pnpm install
pnpm dev
```

## Next Steps

1. **Read the Documentation**
   - [Project Specification](../ProjectSpec.md)
   - [Technology Stack](../ProjectStack.md)
   - [Project Plan](../ProjectPlan.md)

2. **Explore the Code**
   - [Backend README](../backend/README.md)
   - [Frontend README](../frontend/README.md)
   - [Blockchain README](../blockchain/README.md)
   - [Infrastructure README](../infrastructure/README.md)

3. **Set Up for Production**
   - Configure AWS credentials
   - Set up Terraform backend
   - Deploy infrastructure
   - Configure monitoring

## Getting Help

- **Documentation**: Check the `/docs` folder
- **Issues**: Open an issue on GitHub
- **Discussions**: Join our Discord
- **Email**: support@passchain.io

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](../CONTRIBUTING.md) for guidelines.

## License

MIT License - see [LICENSE](../LICENSE) for details.

---

**Happy coding!** 🚀




