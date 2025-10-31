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

---

## 🚀 Implementation
- **[IMPLEMENTATION_PLAN.md](../IMPLEMENTATION_PLAN.md)** - 6-week implementation plan (current)
- **[ENTERPRISE_ROADMAP.md](ENTERPRISE_ROADMAP.md)** - Enterprise features roadmap
- **[USER_STORIES.md](USER_STORIES.md)** - 24 user stories across 10 epics

---

## 🔐 Security & Audit
- **[AUDIT_VAULT_INTEGRATION.md](AUDIT_VAULT_INTEGRATION.md)** - Dual audit logging (DB + Vault)
- **[SMART_CONTRACTS.md](SMART_CONTRACTS.md)** - ERC-20 & ERC-721 contracts

---

## 📡 API Reference
- **[API.md](API.md)** - REST API documentation
- **Backend Handlers:**
  - Credentials: `backend/internal/api/handlers/credentials.go`
  - Audit: `backend/internal/api/handlers/audit.go`
  - Organizations: `backend/internal/api/handlers/organization_handler.go` (to implement)
  - Members: `backend/internal/api/handlers/member_handler.go` (to implement)
  - Projects: `backend/internal/api/handlers/project_handler.go` (to implement)
  - Vaults: `backend/internal/api/handlers/vault_handler.go` (to implement)

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
```

---

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

---

## 🔄 Recently Updated
- **2025-10-30:** IMPLEMENTATION_PLAN.md (6-week plan)
- **2025-10-30:** USER_STORIES.md (24 stories)
- **2025-10-30:** AUDIT_VAULT_INTEGRATION.md (Vault backup)
- **2025-10-30:** ENTERPRISE_SCHEMA.md (Multi-tenant DB)
- **2025-10-30:** ENTERPRISE_ROADMAP.md (Phase breakdown)

---

## 📞 Need Help?

- **Implementation Questions:** See [IMPLEMENTATION_PLAN.md](../IMPLEMENTATION_PLAN.md)
- **Setup Issues:** See [GETTING_STARTED.md](GETTING_STARTED.md)
- **API Questions:** See [API.md](API.md)
- **Security Questions:** See [SECURITY_ARCHITECTURE.md](SECURITY_ARCHITECTURE.md)

---

**AUUUUFFFF!** 🔥
