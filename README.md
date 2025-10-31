# 🔐 Pass Chain

**Decentralized Password Manager with Zero-Knowledge Architecture**

[![Test Suite](https://github.com/yourusername/pass-chain/workflows/Test%20Suite/badge.svg)](https://github.com/yourusername/pass-chain/actions)
[![Coverage](https://codecov.io/gh/yourusername/pass-chain/branch/main/graph/badge.svg)](https://codecov.io/gh/yourusername/pass-chain)
[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/pass-chain)](https://goreportcard.com/report/github.com/yourusername/pass-chain)

Pass Chain is an enterprise-grade secrets management platform combining blockchain immutability, client-side encryption, and multi-tenant organization support. Nobody—not even Pass Chain—can decrypt your passwords.

## ✨ Features

### Core Security
- **Client-side encryption** - Passwords encrypted in browser using XChaCha20-Poly1305
- **Split-key security** - 2-of-3 Shamir Secret Sharing (Vault + Blockchain + User backup)
- **Wallet authentication** - MetaMask/WalletConnect, no passwords needed
- **Zero-knowledge** - Backend never sees plaintext credentials
- **Blockchain audit trail** - Immutable access logs (Hyperledger Fabric)

### Enterprise Features
- **Multi-tenant organizations** - Teams, departments, projects
- **Role-Based Access Control (RBAC)** - 6 default roles (Owner, Admin, Security Officer, Member, Auditor, Guest)
- **Project-based vaults** - Organize credentials by projects
- **Granular permissions** - Control who can create, read, update, delete, share, rotate
- **Invitation system** - Invite members via wallet address
- **Audit logs** - Full history of all actions (DB + Vault + Blockchain)

## 🚀 Quick Start

### Prerequisites
- Docker Desktop
- Minikube (for local Kubernetes)
- Node.js 18+
- Go 1.23+

### 1. Start Infrastructure (Kubernetes)

```powershell
# Windows
.\start-minikube.ps1

# Or manually
minikube start --memory=8192 --cpus=4
kubectl apply -f k8s/namespace.yaml
helm install passchain-infra k8s/charts/passchain-infra -n passchain
```

### 2. Port Forward Services

```powershell
# Backend
kubectl port-forward svc/passchain-backend 8080:8080 -n passchain

# PostgreSQL (optional, for debugging)
kubectl port-forward svc/passchain-postgres 5432:5432 -n passchain

# Vault (optional)
kubectl port-forward svc/passchain-vault 8200:8200 -n passchain
```

### 3. Start Frontend

```bash
cd frontend
npm install
npm run dev
```

Access at **http://localhost:3000**

## 🧪 Testing

### Quick Test

```bash
# Using Makefile (recommended)
make test-all

# Or manually
docker-compose -f docker-compose.test.yml up -d
cd backend && go test ./... -v
```

### Test Types

```bash
# Unit tests (fast, no infrastructure)
make test-unit
# or: cd backend && go test ./internal/... -v -short

# Integration tests (DB + Vault + Redis)
make test-integration
# or: go test ./test -v -run TestIntegration

# E2E flow tests (complete user journeys)
make test-e2e
# or: go test ./test -v -run TestEndToEndFlow

# Coverage report
make test-coverage
# Opens HTML report in browser
```

### Smoke Test (against running backend)

```bash
make smoke
# or: cd backend && go run cmd/smoke_test/main.go
```

### CI/CD

Tests run automatically on every push and PR:
- ✅ Unit tests with race detector
- ✅ Integration tests
- ✅ E2E flow tests
- ✅ Linting (golangci-lint)
- ✅ Security scanning (Trivy)
- ✅ Coverage reporting (Codecov)

See [`.github/workflows/test.yml`](.github/workflows/test.yml)

## 🏗️ Architecture

```
User (Web3 Wallet)
    ↓ Sign request
Frontend (Next.js + Web3)
    ↓ Encrypted payload
Backend API (Go + Gin)
    ├─→ PostgreSQL (metadata + encrypted credentials)
    ├─→ HashiCorp Vault (Share1 + audit backup)
    ├─→ Hyperledger Fabric (Share2 + immutable logs)
    └─→ Redis (cache)
```

### Data Flow

1. **Credential Storage**
   - Client generates random DEK (Data Encryption Key)
   - Client encrypts credential with DEK using XChaCha20-Poly1305
   - DEK split into 3 shares via Shamir Secret Sharing (2-of-3)
     - Share1 → HashiCorp Vault
     - Share2 → Hyperledger Fabric PDC
     - Share3 → User's browser localStorage (backup)
   - Encrypted credential → PostgreSQL
   - Audit log → Database + Vault + Fabric

2. **Credential Retrieval**
   - User signs access request with wallet
   - Backend validates signature
   - Backend retrieves Share1 (Vault) + Share2 (Fabric)
   - Client reconstructs DEK from shares
   - Client decrypts credential locally
   - Access logged to blockchain

## 📚 Tech Stack

| Layer | Technology | Purpose |
|-------|-----------|---------|
| Frontend | React + Next.js 14 + TypeScript | UI framework |
| Web3 | RainbowKit + Wagmi + viem | Wallet connection |
| UI Library | shadcn/ui + Tailwind CSS | Components & styling |
| Backend | Go 1.23 + Gin + GORM | API server |
| Database | PostgreSQL 15 | Metadata & encrypted data |
| Cache | Redis 7 | Session & query cache |
| Secrets | HashiCorp Vault | Key storage (Share1) |
| Blockchain | Hyperledger Fabric 2.5 | Immutable audit logs (Share2) |
| Orchestration | Kubernetes + Helm | Container management |
| IaC | Terraform | Cloud infrastructure |
| CI/CD | GitHub Actions | Testing & deployment |
| Monitoring | Prometheus + Grafana | Metrics & alerts |

## 📁 Project Structure

```
pass-chain/
├── backend/                  # Go backend
│   ├── cmd/
│   │   ├── server/          # Main API server
│   │   └── smoke_test/      # Smoke test binary
│   ├── internal/
│   │   ├── api/             # Gin routes & handlers
│   │   │   ├── handlers/    # Request handlers
│   │   │   ├── middleware/  # Auth, CORS, logging
│   │   │   └── router.go    # Route definitions
│   │   ├── database/        # GORM models & migrations
│   │   ├── models/          # Data models
│   │   └── services/        # Business logic (RBAC, Org, Vault, Fabric)
│   ├── test/                # Integration & E2E tests
│   └── pkg/                 # Shared utilities
├── frontend/                # Next.js frontend
│   ├── src/
│   │   ├── app/            # App router pages
│   │   ├── components/     # React components
│   │   ├── lib/            # Utils, API client, crypto
│   │   └── hooks/          # Custom React hooks
│   └── public/             # Static assets
├── k8s/                    # Kubernetes manifests
│   ├── charts/             # Helm charts
│   │   ├── passchain-infra/  # PostgreSQL, Redis, Vault
│   │   └── passchain-app/    # Backend, Frontend
│   └── *.yaml              # Base manifests
├── blockchain/             # Hyperledger Fabric
│   ├── chaincode/          # Smart contracts (Go)
│   ├── network/            # Fabric network config
│   └── scripts/            # Setup scripts
├── contracts/              # Solidity (ERC-20 token, NFT backup)
│   ├── src/
│   ├── test/
│   └── scripts/            # Deployment scripts
├── docs/                   # Documentation
│   ├── ENTERPRISE_SCHEMA.md
│   ├── ENTERPRISE_ROADMAP.md
│   ├── USER_STORIES.md
│   └── README.md
├── scripts/                # Helper scripts
│   ├── run-tests.ps1
│   └── run-tests.sh
├── .github/workflows/      # CI/CD pipelines
├── docker-compose.yml      # Local dev environment
├── docker-compose.test.yml # Test environment
└── Makefile               # Test & build commands
```

## 🛠️ Development

### Local Development (Docker Compose)

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f backend

# Rebuild after changes
docker-compose up -d --build backend
```

### Database Migrations

Migrations run automatically on backend startup. To run manually:

```bash
cd backend
go run cmd/server/main.go migrate
```

### Linting

```bash
# Backend
cd backend
golangci-lint run

# Frontend
cd frontend
npm run lint
```

### Building

```bash
# Backend binary
cd backend
go build -o bin/server cmd/server/main.go

# Frontend production
cd frontend
npm run build
npm start
```

## 🔐 Security

- **Client-side encryption**: All credentials encrypted before leaving browser
- **Split-key architecture**: No single point can decrypt data
- **Wallet signatures**: Every action signed by user's wallet
- **Audit trail**: Immutable blockchain logs
- **RBAC**: Granular permission system
- **Vault integration**: Enterprise-grade secret storage
- **Security scanning**: Trivy scans in CI/CD
- **Dependency auditing**: Automatic vulnerability scanning

## 📖 Documentation

- [Enterprise Schema](docs/ENTERPRISE_SCHEMA.md) - Database design
- [Implementation Plan](IMPLEMENTATION_PLAN.md) - Development roadmap
- [User Stories](docs/USER_STORIES.md) - Feature scenarios
- [Testing Guide](backend/test/README.md) - How to run tests
- [API Documentation](docs/API.md) - API endpoints (Swagger)
- [Product Documentation](pass-chain-product-documentation.md) - Why we built this

## 🚢 Deployment

### GKE (Google Kubernetes Engine)

```bash
# Apply Terraform
cd terraform/gke
terraform init
terraform apply

# Deploy with Helm
helm upgrade --install passchain k8s/charts/passchain-app \
  --namespace passchain \
  --set image.tag=v1.0.0
```

See [`.github/workflows/deploy-gke.yml`](.github/workflows/deploy-gke.yml)

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing`)
5. Open Pull Request

**Note**: All PRs must pass tests and linting before merge.

## 📝 License

MIT License - see [LICENSE](LICENSE)

## 🎯 Roadmap

- [x] Core password management
- [x] Wallet authentication
- [x] Split-key architecture (Vault + DB)
- [x] Enterprise organizations & RBAC
- [x] Project-based vaults
- [x] Audit logs (DB + Vault)
- [ ] Hyperledger Fabric integration (Share2 + blockchain logs)
- [ ] ERC-20 payment token
- [ ] ERC-721 NFT shard backup (IPFS)
- [ ] Mobile app (React Native)
- [ ] Browser extension
- [ ] Hardware wallet support (Ledger, Trezor)
- [ ] Self-hosted option
- [ ] SSO integration (SAML, OAuth)

## 💬 Support

- GitHub Issues: [Report bugs or request features](https://github.com/yourusername/pass-chain/issues)
- Documentation: [Read the docs](docs/)
- Discord: [Join our community](#)

---

**AUUUUFFFF! 🔥**
