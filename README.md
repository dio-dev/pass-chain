# 🔐 Pass Chain

**Decentralized Password Manager with Zero-Knowledge Architecture**

Pass Chain is an enterprise-grade secrets management platform combining blockchain immutability, client-side encryption, and multi-tenant organization support. Nobody—not even Pass Chain—can decrypt your passwords.

---

## 🚀 Quick Start

### Prerequisites
- **Docker** & **Docker Compose**
- **Minikube** (for local Kubernetes)
- **kubectl**
- **Node.js 18+** (for frontend/contracts)
- **Go 1.23+** (for backend)

### Local Development

```bash
# 1. Start Minikube
minikube start --cpus=4 --memory=8192

# 2. Deploy infrastructure
cd infrastructure/k8s
kubectl apply -f namespace.yaml
helm install passchain-postgresql ./charts/postgresql -n passchain
helm install passchain-redis ./charts/redis -n passchain
helm install passchain-vault ./charts/vault -n passchain

# 3. Build and deploy backend
docker build -t passchain/backend:2.0.0 -f backend/Dockerfile backend
minikube image load passchain/backend:2.0.0
kubectl apply -f infrastructure/k8s/backend/

# 4. Build and deploy frontend
cd frontend
npm install
npm run build
docker build -t passchain/frontend:2.0.0 .
minikube image load passchain/frontend:2.0.0
kubectl apply -f ../infrastructure/k8s/frontend/

# 5. Port forward services
kubectl port-forward -n passchain svc/passchain-backend 8080:8080 &
kubectl port-forward -n passchain svc/passchain-frontend 3000:3000 &

# 6. Access the app
open http://localhost:3000
```

---

## 🏗️ Architecture

### Core Components

| Component | Technology | Purpose |
|-----------|-----------|---------|
| **Frontend** | Next.js 14 + RainbowKit | Wallet-based auth, client-side encryption |
| **Backend** | Go (Gin framework) | API coordination, RBAC, audit logging |
| **Vault** | HashiCorp Vault | Secure storage of encryption key shards |
| **Database** | PostgreSQL | Encrypted credentials & metadata |
| **Cache** | Redis | Session management & performance |
| **Blockchain** | Hyperledger Fabric (future) | Immutable audit logs & key sharding |

### Security Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    User (Web3 Wallet)                       │
│                  EIP-191 Signature Auth                     │
└────────────────────┬────────────────────────────────────────┘
                     │
         ┌───────────▼────────────┐
         │   Frontend (Browser)    │
         │  Client-Side Encryption │
         │   XChaCha20-Poly1305    │
         └───────────┬────────────┘
                     │
         ┌───────────▼────────────┐
         │   Backend API (Go)      │
         │  • RBAC Enforcement    │
         │  • Audit Logging       │
         │  • Shard Management    │
         └─┬─────────────────┬───┘
           │                 │
  ┌────────▼─────┐    ┌─────▼──────────┐
  │ HashiCorp    │    │  PostgreSQL    │
  │ Vault (KV)   │    │  (Encrypted)   │
  │ • Share1     │    │  • Ciphertext  │
  │ • Audit Logs │    │  • Share2      │
  └──────────────┘    └────────────────┘
```

### Encryption Flow

1. **Credential Storage:**
   - User generates 256-bit DEK (Data Encryption Key)
   - Splits DEK → Share1 (Vault) + Share2 (DB) + Share3 (Browser localStorage)
   - Encrypts credential with DEK using XChaCha20-Poly1305
   - Stores ciphertext + nonce in database

2. **Credential Retrieval:**
   - Backend fetches Share1 from Vault
   - Backend returns Share1 + Share2 to authenticated user
   - Frontend gets Share3 from localStorage
   - Reconstructs DEK: `Share1 XOR Share2 XOR Share3`
   - Decrypts credential locally (plaintext never leaves browser)

---

## 📂 Project Structure

```
pass-chain/
├── backend/              # Go API server
│   ├── cmd/server/       # Main entry point
│   ├── internal/
│   │   ├── api/          # HTTP handlers & routes
│   │   ├── database/     # GORM models & migrations
│   │   ├── middleware/   # Auth, logging, CORS
│   │   ├── models/       # Enterprise domain models
│   │   └── services/     # Business logic (RBAC, Vault, Fabric)
│   └── pkg/              # Shared utilities
├── frontend/             # Next.js 14 app
│   ├── src/
│   │   ├── app/          # App router pages
│   │   ├── components/   # React components (shadcn/ui)
│   │   ├── contexts/     # React contexts (Org, Auth)
│   │   └── lib/          # Crypto utils, API client
├── contracts/            # Solidity smart contracts
│   ├── contracts/        # ERC-20 token, ERC-721 NFT
│   ├── scripts/          # Hardhat deployment scripts
│   └── test/             # Contract tests
├── blockchain/           # Hyperledger Fabric (future)
│   └── chaincode/        # Fabric smart contracts
├── infrastructure/       # Deployment configs
│   ├── k8s/              # Kubernetes manifests
│   │   ├── backend/
│   │   ├── frontend/
│   │   └── charts/       # Helm charts
│   └── terraform/        # GKE/AWS provisioning (future)
└── docs/                 # Documentation
    ├── ENTERPRISE_SCHEMA.md      # Database schema
    ├── ENTERPRISE_ROADMAP.md     # Implementation roadmap
    ├── USER_STORIES.md           # User personas & stories
    ├── UX_DESIGN_PHILOSOPHY.md   # UI/UX principles
    └── README.md                 # Docs index
```

---

## 🎯 Features

### ✅ Implemented
- ✅ Wallet-based authentication (MetaMask, WalletConnect)
- ✅ Client-side encryption (XChaCha20-Poly1305)
- ✅ 3-shard encryption key splitting (Vault + DB + Browser)
- ✅ CRUD operations for credentials
- ✅ Audit logging (PostgreSQL + Vault backup)
- ✅ Multi-tenant organizations
- ✅ Role-Based Access Control (6 default roles)
- ✅ Personal & organization vaults
- ✅ Project-scoped credentials
- ✅ Kubernetes deployment (Minikube/GKE)
- ✅ Docker containerization
- ✅ ERC-20 payment token (PassChainToken)
- ✅ ERC-721 NFT for shard backup

### 🚧 In Progress
- 🚧 Hyperledger Fabric integration
- 🚧 2-of-3 Shamir Secret Sharing (proper implementation)
- 🚧 Key rotation mechanism
- 🚧 Credential sharing between users
- 🚧 Team invitations & onboarding

### 🔮 Planned
- 🔮 Browser extensions (Chrome, Firefox)
- 🔮 Mobile apps (iOS, Android)
- 🔮 CLI tool
- 🔮 Terraform GKE/AWS modules
- 🔮 Self-hosted enterprise edition
- 🔮 IPFS integration for shard storage
- 🔮 Zero-knowledge proof verification

---

## 🔐 Security Features

| Feature | Status | Description |
|---------|--------|-------------|
| **Client-Side Encryption** | ✅ | Plaintext never leaves browser |
| **Zero-Knowledge Auth** | ✅ | Wallet signatures (EIP-191/SIWE) |
| **Split-Key Storage** | ✅ | 3 shards across Vault/DB/Browser |
| **Immutable Audit Logs** | ✅ | PostgreSQL + Vault KV backup |
| **RBAC** | ✅ | Org-scoped permissions |
| **Key Rotation** | 🚧 | User-initiated key updates |
| **Blockchain Anchoring** | 🔮 | Fabric for audit trail |
| **Hardware Security** | 🔮 | HSM/TPM support |

---

## 📖 Documentation

See [`docs/README.md`](./docs/README.md) for:
- Enterprise architecture & database schema
- User stories & personas
- API documentation
- Deployment guides
- Security best practices

---

## 🛠️ Development

### Backend (Go)
```bash
cd backend
go mod tidy
go run cmd/server/main.go
# API: http://localhost:8080
```

### Frontend (Next.js)
```bash
cd frontend
npm install
npm run dev
# App: http://localhost:3000
```

### Smart Contracts (Hardhat)
```bash
cd contracts
npm install
npx hardhat compile
npx hardhat test
npx hardhat run scripts/deploy-local.js --network localhost
```

---

## 🚀 Deployment

### Minikube (Local)
```bash
cd infrastructure/k8s
./scripts/deploy-minikube.sh
```

### GKE (Production - Future)
```bash
cd infrastructure/terraform/gke
terraform init
terraform apply
# Then apply Kubernetes manifests
```

---

## 🤝 Contributing

We're not accepting external contributions yet, but star the repo to stay updated!

---

## 📄 License

[MIT License](./LICENSE)

---

## 🔗 Links

- **Documentation**: [./docs](./docs)
- **Implementation Plan**: [IMPLEMENTATION_PLAN.md](./IMPLEMENTATION_PLAN.md)
- **Product Docs**: [./docs/product](./docs/product)

---

**Built with ❤️ using Web3, Go, React, and a lot of coffee ☕**

**AUUUUFFFF!** 🔥
