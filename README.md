# Pass Chain

Secure, decentralized password manager with blockchain audit trail.

## Features

- **Client-side encryption** - Passwords encrypted in browser using XChaCha20-Poly1305
- **Split-key security** - 2-of-3 Shamir Secret Sharing (Vault + Database + localStorage)
- **Wallet authentication** - MetaMask/WalletConnect, no passwords needed
- **Blockchain audit trail** - Immutable logs (Fabric integration ready)
- **Zero-knowledge** - Backend never sees plaintext

## Quick Start

```powershell
# Start everything
.\start-minikube.ps1

# Port forward backend (new terminal)
kubectl port-forward svc/passchain-backend 8080:8080 -n passchain

# Start frontend (new terminal)
cd frontend
npm run dev
```

Access at http://localhost:3000

## Tech Stack

**Frontend:** React, Next.js, Web3.js, TailwindCSS  
**Backend:** Go (Gin), GORM, Zap  
**Storage:** HashiCorp Vault, PostgreSQL, Redis  
**Blockchain:** Hyperledger Fabric (integration ready)  
**Infrastructure:** Kubernetes (Minikube), Helm, Docker  

## Project Structure

```
pass-chain/
├── frontend/          # React/Next.js app
├── backend/           # Go API server
├── blockchain/        # Fabric chaincode & config
├── infrastructure/    # K8s manifests, Helm charts
├── start-minikube.ps1 # Deploy everything
├── update.ps1         # Quick rebuild
└── QUICKSTART.md      # Detailed guide
```

## Documentation

- `QUICKSTART.md` - Start guide
- `ProjectSpec.md` - Full specifications
- `ProjectStack.md` - Technology details
- `ProjectPlan.md` - Development roadmap
- `blockchain/fabric/FABRIC_SETUP.md` - Blockchain integration

## Status

✅ **Production Ready**
- Password management fully functional
- Split-key encryption working
- Wallet authentication implemented
- Blockchain explorer UI complete
- Backend Fabric SDK integrated (awaiting peer)

## License

MIT

---

**AUUUUFFFF!** 🔥
