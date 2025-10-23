# Pass Chain Documentation

<div align="center">
  <img src="https://via.placeholder.com/150/8B5CF6/FFFFFF?text=🔐" alt="Pass Chain Logo" width="150"/>
  
  <h3>The Future of Password Management</h3>
  
  <p>
    <a href="https://github.com/yourusername/pass-chain"><img src="https://img.shields.io/github/stars/yourusername/pass-chain?style=social" alt="GitHub stars"/></a>
    <a href="https://github.com/yourusername/pass-chain/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License"/></a>
    <a href="https://github.com/yourusername/pass-chain"><img src="https://img.shields.io/badge/PRs-welcome-brightgreen.svg" alt="PRs Welcome"/></a>
  </p>
</div>

---

## What is Pass Chain?

Pass Chain is a revolutionary password manager that combines:

- 🔐 **Zero-Knowledge Encryption** - We literally can't decrypt your passwords
- 🔑 **Split-Key Security** - Encryption key split across Vault, blockchain, and your device
- 👛 **Wallet Authentication** - Use MetaMask instead of master passwords
- ⛓️ **Blockchain Audit Trail** - Immutable logs on Hyperledger Fabric
- 🚀 **Enterprise Ready** - Kubernetes deployment, SOC2 compliant

## Quick Start

```bash
# Start everything
./start-minikube.ps1

# Port forward backend
kubectl port-forward svc/passchain-backend 8080:8080 -n passchain

# Start frontend
cd frontend && npm run dev
```

Visit http://localhost:3000

## Why Pass Chain?

### Traditional Password Managers
- ❌ Trust the company not to decrypt
- ❌ Single point of failure (master password)
- ❌ Centralized database (honeypot for hackers)
- ❌ No audit trail

### Pass Chain
- ✅ **Mathematically impossible** for us to decrypt
- ✅ **Split-key** - need 2 of 3 parts + wallet signature
- ✅ **Distributed** - Vault + Blockchain + Your device
- ✅ **Immutable audit logs** on blockchain

## Architecture at a Glance

```
┌──────────────┐
│   Browser    │  Encrypt with XChaCha20-Poly1305
│   (Client)   │  Split key into 3 shares
└──────┬───────┘
       │
       ├─────────┐
       │         │
   ┌───▼──┐  ┌───▼────────┐  ┌──────────────┐
   │Vault │  │Blockchain  │  │localStorage  │
   │(S1)  │  │(S2)        │  │(S3 - Backup) │
   └──────┘  └────────────┘  └──────────────┘
   
   🔐 Need any 2 + wallet signature to decrypt
```

## Key Features

| Feature | Description |
|---------|-------------|
| **Client-Side Encryption** | Passwords encrypted in browser using XChaCha20-Poly1305 |
| **Shamir Secret Sharing** | 2-of-3 threshold, distributed across systems |
| **Web3 Authentication** | MetaMask/WalletConnect, no master password |
| **Blockchain Logs** | Every access recorded on Hyperledger Fabric |
| **Zero-Knowledge** | Backend never sees plaintext |
| **Kubernetes Native** | Production-ready infrastructure |

## Documentation Structure

### 📚 Getting Started
- [Quick Start Guide](getting-started/quickstart.md)
- [Installation](getting-started/installation.md)
- [Your First Credential](getting-started/first-credential.md)

### 🏗️ Architecture
- [System Overview](architecture/overview.md)
- [Encryption Deep Dive](architecture/encryption.md)
- [Split-Key Mechanics](architecture/split-key.md)
- [Wallet Authentication](architecture/wallet-auth.md)
- [Blockchain Integration](architecture/blockchain.md)

### 🚀 Deployment
- [Kubernetes Setup](deployment/kubernetes.md)
- [AWS Deployment](deployment/aws.md)
- [GKE Deployment](deployment/gke.md)
- [Vault Configuration](deployment/vault.md)
- [Fabric Setup](deployment/fabric.md)

### 🔒 Security
- [Security Model](security/model.md)
- [Threat Analysis](security/threats.md)
- [Best Practices](security/best-practices.md)
- [Compliance (SOC2/GDPR)](security/compliance.md)

### 📖 API Reference
- [Backend REST API](api/backend.md)
- [Fabric Chaincode](api/chaincode.md)
- [Frontend SDK](api/frontend.md)

### 💻 Development
- [Development Setup](development/setup.md)
- [Contributing Guidelines](development/contributing.md)
- [Testing](development/testing.md)
- [Debugging](development/debugging.md)

## Technology Stack

**Frontend**
- React, Next.js 14
- TailwindCSS
- Web3.js, Wagmi, RainbowKit
- @noble/ciphers (encryption)

**Backend**
- Go 1.21
- Gin (web framework)
- GORM (ORM)
- Fabric SDK Go

**Infrastructure**
- Kubernetes (Minikube, GKE, EKS)
- HashiCorp Vault
- PostgreSQL 15
- Redis
- Hyperledger Fabric 2.5

## Security Guarantees

1. **Zero-Knowledge**: Backend mathematically cannot decrypt your passwords
2. **Split-Key**: Need 2 of 3 key shares + wallet signature
3. **Client-Side**: Encryption happens in browser, never on server
4. **Immutable Logs**: All access logged on blockchain
5. **No Master Password**: Wallet signature is your auth

## Community

- 🌟 [Star us on GitHub](https://github.com/yourusername/pass-chain)
- 💬 [Join Discord](https://discord.gg/passchain)
- 🐦 [Follow on Twitter](https://twitter.com/passchain)
- 📧 [Email Support](mailto:support@passchain.io)

## License

MIT License - see [LICENSE](https://github.com/yourusername/pass-chain/blob/main/LICENSE)

---

<div align="center">
  <p><strong>Built with ❤️ by developers who care about security</strong></p>
  <p><i>AUUUUFFFF! 🔥</i></p>
</div>

