# Pass Chain - Quick Start

## What is Pass Chain?

Secure password manager using:
- **Client-side encryption** (XChaCha20-Poly1305)
- **Split-key storage** (Vault + Database + localStorage)
- **Wallet authentication** (MetaMask/WalletConnect)
- **Blockchain audit trail** (ready for Hyperledger Fabric)

---

## Start Everything

```powershell
# Start Minikube + deploy all services
.\start-minikube.ps1

# Port forward backend (new terminal)
kubectl port-forward svc/passchain-backend 8080:8080 -n passchain

# Start frontend (new terminal)
cd frontend
npm run dev
```

**Access:** http://localhost:3000

---

## Quick Update

```powershell
# Rebuild and restart backend/frontend
.\update.ps1
```

---

## Architecture

```
Frontend (React/Next.js)
  ├─→ Wallet connection (Web3)
  ├─→ Client-side encryption
  └─→ API calls to backend

Backend (Go)
  ├─→ Vault (Share1 storage)
  ├─→ PostgreSQL (encrypted data)
  ├─→ Fabric SDK (ready for blockchain)
  └─→ REST API

Infrastructure (Kubernetes)
  ├─→ Backend pods
  ├─→ Frontend pods
  ├─→ Vault StatefulSet
  ├─→ PostgreSQL StatefulSet
  └─→ Redis StatefulSet
```

---

## Features

- ✅ Store credentials with wallet signature
- ✅ Reveal passwords (requires wallet signature)
- ✅ Delete credentials  
- ✅ Blockchain explorer UI
- ✅ Audit logs
- ✅ Stats dashboard

---

## Fabric Integration

Backend is **ready** for Hyperledger Fabric. To connect:

1. Install WSL2: `wsl --install`
2. Download Fabric test-network
3. Start peer at `localhost:7051`
4. Backend connects automatically

See `BACKEND_FABRIC_READY.md` for details.

---

**AUUUUFFFF!** 🔥
