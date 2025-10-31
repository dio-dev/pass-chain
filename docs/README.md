<<<<<<< HEAD
# 📚 Pass Chain - Documentation Index

## 🎯 Getting Started
- **[README.md](../README.md)** - Project overview
- **[QUICKSTART.md](../QUICKSTART.md)** - Quick start guide
- **[GETTING_STARTED.md](GETTING_STARTED.md)** - Detailed setup guide
- **[ProjectSpec.md](../ProjectSpec.md)** - Original project specification
- **[ProjectPlan.md](../ProjectPlan.md)** - Development roadmap
- **[ProjectStack.md](../ProjectStack.md)** - Technology stack

---

## 🏗️ Architecture
- **[ARCHITECTURE.md](ARCHITECTURE.md)** - System architecture overview
- **[SECURITY_ARCHITECTURE.md](SECURITY_ARCHITECTURE.md)** - Security design & threat model
- **[ENTERPRISE_SCHEMA.md](ENTERPRISE_SCHEMA.md)** - Database schema for multi-tenant orgs
- **[architecture/overview.md](architecture/overview.md)** - Detailed architecture diagrams
=======
---
hidden: true
---

# Pass Chain Documentation

![Pass Chain](https://img.shields.io/badge/Pass%20Chain-v1.0-purple?style=for-the-badge) [![License](https://img.shields.io/badge/license-MIT-blue?style=for-the-badge)](../LICENSE/) [![GitHub](https://img.shields.io/github/stars/yourusername/pass-chain?style=for-the-badge\&logo=github)](https://github.com/yourusername/pass-chain)

**Secure, decentralized password management with blockchain audit trail**

[Quick Start](./#quick-start) • [Architecture](./#architecture) • [Security](./#security) • [API](./#api) • [Deploy](./#deployment)

***

## 🚀 Quick Start

Get Pass Chain running in 5 minutes:

```bash
# Clone and start
git clone https://github.com/yourusername/pass-chain
cd pass-chain
./start-minikube.ps1

# Port forward backend
kubectl port-forward svc/passchain-backend 8080:8080 -n passchain

# Start frontend
cd frontend && npm run dev
```

Visit http://localhost:3000

[**Full Quick Start Guide →**](getting-started/quickstart.md)
>>>>>>> a80db2be5f85dd2e804ef25dc51340dc37e09d5c

***

<<<<<<< HEAD
## 🚀 Implementation
- **[IMPLEMENTATION_PLAN.md](../IMPLEMENTATION_PLAN.md)** - 6-week implementation plan (current)
- **[ENTERPRISE_ROADMAP.md](ENTERPRISE_ROADMAP.md)** - Enterprise features roadmap
- **[USER_STORIES.md](USER_STORIES.md)** - 24 user stories across 10 epics
=======
## 📚 Documentation

### Getting Started

* 📖 [Quick Start](getting-started/quickstart.md)
* 🔧 [Installation](getting-started/installation.md)
* 🎯 [First Credential](getting-started/first-credential.md)

### Architecture

* 🏗️ [System Overview](architecture/overview.md)
* 🔐 [Encryption Deep Dive](architecture/encryption.md)
* 🔑 [Split-Key Security](architecture/split-key.md)
* 👛 [Wallet Authentication](architecture/wallet-auth.md)
* ⛓️ [Blockchain Integration](architecture/blockchain.md)

### Deployment

* ☸️ [Kubernetes Setup](deployment/kubernetes.md)
* ☁️ [AWS Deployment](deployment/aws.md)
* 🌐 [GKE Deployment](deployment/gke.md)
* 🔐 [Vault Configuration](deployment/vault.md)
* ⛓️ [Fabric Setup](deployment/fabric.md)

### Security

* 🛡️ [Security Model](security/model.md)
* ⚠️ [Threat Analysis](security/threats.md)
* ✅ [Best Practices](security/best-practices.md)
* 📋 [Compliance (SOC2/GDPR)](security/compliance.md)

### API Reference

* 🔌 [Backend REST API](api/backend.md)
* ⛓️ [Fabric Chaincode](api/chaincode.md)
* 💻 [Frontend SDK](api/frontend.md)

### Development

* 💻 [Dev Setup](development/setup.md)
* 🤝 [Contributing](development/contributing.md)
* 🧪 [Testing](development/testing.md)
* 🐛 [Debugging](development/debugging.md)

### FAQ
>>>>>>> a80db2be5f85dd2e804ef25dc51340dc37e09d5c

* ❓ [General Questions](faq/general.md)
* 🔒 [Security Questions](faq/security.md)
* 🔧 [Troubleshooting](faq/troubleshooting.md)

***

## 🔐 Security & Audit
- **[AUDIT_VAULT_INTEGRATION.md](AUDIT_VAULT_INTEGRATION.md)** - Dual audit logging (DB + Vault)
- **[SMART_CONTRACTS.md](SMART_CONTRACTS.md)** - ERC-20 & ERC-721 contracts

---

<<<<<<< HEAD
## 📡 API Reference
- **[API.md](API.md)** - REST API documentation
- **Backend Handlers:**
  - Credentials: `backend/internal/api/handlers/credentials.go`
  - Audit: `backend/internal/api/handlers/audit.go`
  - Organizations: `backend/internal/api/handlers/organization_handler.go` (to implement)
  - Members: `backend/internal/api/handlers/member_handler.go` (to implement)
  - Projects: `backend/internal/api/handlers/project_handler.go` (to implement)
  - Vaults: `backend/internal/api/handlers/vault_handler.go` (to implement)
=======
* 🔐 **Zero-Knowledge Encryption** - We literally can't decrypt your passwords
* 🔑 **Split-Key Security** - Key split across Vault, blockchain, and your device
* 👛 **Wallet Authentication** - Use MetaMask instead of master passwords
* ⛓️ **Blockchain Audit Trail** - Immutable logs on Hyperledger Fabric
* ☸️ **Enterprise Ready** - Kubernetes deployment, SOC2 compliant
>>>>>>> a80db2be5f85dd2e804ef25dc51340dc37e09d5c

---

## 🛠️ Development
- **[SETUP.md](SETUP.md)** - Local development setup
- **[CLOUDFLARE_TUNNEL.md](CLOUDFLARE_TUNNEL.md)** - Expose local dev to internet
- **Backend:**
  - **[backend/README.md](../backend/README.md)** - Backend setup
  - **[backend/TESTING.md](../backend/TESTING.md)** - Testing guide
- **Frontend:**
  - **[frontend/README.md](../frontend/README.md)** - Frontend setup
- **Contracts:**
  - **[contracts/README.md](../contracts/README.md)** - Smart contracts
  - **[contracts/HARDHAT_DEPLOYMENT.md](../contracts/HARDHAT_DEPLOYMENT.md)** - Deploy to testnets

---

## ☸️ Infrastructure & Deployment
- **[infrastructure/README.md](../infrastructure/README.md)** - Infrastructure overview
- **Kubernetes:**
  - `infrastructure/k8s/backend/` - Backend manifests
  - `infrastructure/k8s/frontend/` - Frontend manifests
  - `infrastructure/k8s/vault/` - Vault manifests
  - `infrastructure/k8s/helm/` - Helm charts
- **Docker:**
  - `infrastructure/docker/docker-compose.yml` - Local Docker setup
- **Terraform:**
  - `infrastructure/terraform/` - IaC for cloud resources
- **Scripts:**
  - `start-minikube.ps1` / `start-minikube.sh` - Start local K8s
  - `start.ps1` / `start.sh` - Start all services
  - `update.ps1` / `update.sh` - Update deployments

---

## ⛓️ Blockchain
- **[blockchain/README.md](../blockchain/README.md)** - Blockchain overview
- **[blockchain/fabric/FABRIC_SETUP.md](../blockchain/fabric/FABRIC_SETUP.md)** - Hyperledger Fabric setup
- **Chaincode:**
  - `blockchain/chaincode/credentials/chaincode.go` - Credentials chaincode

---

## 📖 Product Documentation
- **[product/INDEX.md](product/INDEX.md)** - Product docs index
- **[product/why-pass-chain.md](product/why-pass-chain.md)** - Why we built this
- **[product/how-it-works.md](product/how-it-works.md)** - How it works

---

## 🗂️ Documentation Structure

```
<<<<<<< HEAD
docs/
├── README.md                           # This file
├── GETTING_STARTED.md                  # Setup guide
├── API.md                              # API reference
├── ARCHITECTURE.md                     # System architecture
├── SECURITY_ARCHITECTURE.md            # Security design
├── ENTERPRISE_SCHEMA.md                # Database schema
├── ENTERPRISE_ROADMAP.md               # Enterprise roadmap
├── USER_STORIES.md                     # User stories
├── AUDIT_VAULT_INTEGRATION.md          # Audit logging
├── SMART_CONTRACTS.md                  # Smart contracts
├── SETUP.md                            # Dev setup
├── CLOUDFLARE_TUNNEL.md                # Tunnel setup
├── architecture/
│   └── overview.md                     # Architecture diagrams
├── getting-started/
│   ├── quickstart.md                   # Quick start
│   └── README.md
└── product/
    ├── INDEX.md                        # Product docs index
    ├── why-pass-chain.md               # Why
    └── how-it-works.md                 # How it works
=======
1. Enter password → 2. Encrypt in browser → 3. Split key into 3 parts
                    ↓
    Part 1 → Vault | Part 2 → Blockchain | Part 3 → Your device
                    ↓
4. To decrypt: Wallet signature + any 2 of 3 parts
```

[**Deep Dive into Architecture →**](architecture/overview.md)

***

## 🛡️ Security Guarantees

| What We **CAN'T** Do      | What You **GET**         |
| ------------------------- | ------------------------ |
| ❌ Decrypt your passwords  | ✅ Client-side encryption |
| ❌ Access your credentials | ✅ Split-key architecture |
| ❌ See plaintext data      | ✅ Blockchain audit trail |
| ❌ Recover without wallet  | ✅ Complete control       |

[**Security Model →**](security/model.md)

***

## 🏗️ Tech Stack

**Frontend**: React, Next.js 14, TailwindCSS, Web3.js, Wagmi\
**Backend**: Go 1.21, Gin, GORM, Fabric SDK\
**Infrastructure**: Kubernetes, Vault, PostgreSQL, Redis, Fabric\
**Blockchain**: Hyperledger Fabric 2.5

[**Full Stack Details →**](../ProjectStack.md)

***

## 📦 Project Structure

```
pass-chain/
├── frontend/          # React/Next.js app
├── backend/           # Go API server
├── blockchain/        # Fabric chaincode
├── infrastructure/    # K8s manifests
└── docs/             # This documentation
>>>>>>> a80db2be5f85dd2e804ef25dc51340dc37e09d5c
```

***

## 🎯 For Different Audiences

### 👨‍💻 Developers:
1. [GETTING_STARTED.md](GETTING_STARTED.md)
2. [ARCHITECTURE.md](ARCHITECTURE.md)
3. [API.md](API.md)
4. [IMPLEMENTATION_PLAN.md](../IMPLEMENTATION_PLAN.md)

### 🔐 Security Engineers:
1. [SECURITY_ARCHITECTURE.md](SECURITY_ARCHITECTURE.md)
2. [AUDIT_VAULT_INTEGRATION.md](AUDIT_VAULT_INTEGRATION.md)
3. [SMART_CONTRACTS.md](SMART_CONTRACTS.md)

### 🏢 Product Managers:
1. [product/why-pass-chain.md](product/why-pass-chain.md)
2. [USER_STORIES.md](USER_STORIES.md)
3. [ENTERPRISE_ROADMAP.md](ENTERPRISE_ROADMAP.md)

### 🚀 DevOps:
1. [infrastructure/README.md](../infrastructure/README.md)
2. [SETUP.md](SETUP.md)
3. [blockchain/fabric/FABRIC_SETUP.md](../blockchain/fabric/FABRIC_SETUP.md)

***

<<<<<<< HEAD
## 🔄 Recently Updated
- **2025-10-30:** IMPLEMENTATION_PLAN.md (6-week plan)
- **2025-10-30:** USER_STORIES.md (24 stories)
- **2025-10-30:** AUDIT_VAULT_INTEGRATION.md (Vault backup)
- **2025-10-30:** ENTERPRISE_SCHEMA.md (Multi-tenant DB)
- **2025-10-30:** ENTERPRISE_ROADMAP.md (Phase breakdown)
=======
## 📄 License

MIT License - see [LICENSE](../LICENSE/)
>>>>>>> a80db2be5f85dd2e804ef25dc51340dc37e09d5c

***

## 📞 Need Help?

<<<<<<< HEAD
- **Implementation Questions:** See [IMPLEMENTATION_PLAN.md](../IMPLEMENTATION_PLAN.md)
- **Setup Issues:** See [GETTING_STARTED.md](GETTING_STARTED.md)
- **API Questions:** See [API.md](API.md)
- **Security Questions:** See [SECURITY_ARCHITECTURE.md](SECURITY_ARCHITECTURE.md)

---

**AUUUUFFFF!** 🔥
=======
* 🏠 [Main Site](https://passchain.io)
* 💻 [GitHub](https://github.com/yourusername/pass-chain)
* 📖 [Documentation](https://docs.passchain.io)
* 💬 [Discord](https://discord.gg/passchain)
* 🐦 [Twitter](https://twitter.com/passchain)

***

**Built with ❤️ by developers who care about security**

_AUUUUFFFF!_ 🔥
>>>>>>> a80db2be5f85dd2e804ef25dc51340dc37e09d5c
