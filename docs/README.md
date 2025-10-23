# Pass Chain Documentation

<div align="center">

![Pass Chain](https://img.shields.io/badge/Pass%20Chain-v1.0-purple?style=for-the-badge)
[![License](https://img.shields.io/badge/license-MIT-blue?style=for-the-badge)](../LICENSE)
[![GitHub](https://img.shields.io/github/stars/yourusername/pass-chain?style=for-the-badge&logo=github)](https://github.com/yourusername/pass-chain)

**Secure, decentralized password management with blockchain audit trail**

[Quick Start](#quick-start) • [Architecture](#architecture) • [Security](#security) • [API](#api) • [Deploy](#deployment)

</div>

---

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

**[Full Quick Start Guide →](getting-started/quickstart.md)**

---

## 📚 Documentation

### Getting Started
- 📖 [Quick Start](getting-started/quickstart.md)
- 🔧 [Installation](getting-started/installation.md)
- 🎯 [First Credential](getting-started/first-credential.md)

### Architecture
- 🏗️ [System Overview](architecture/overview.md)
- 🔐 [Encryption Deep Dive](architecture/encryption.md)
- 🔑 [Split-Key Security](architecture/split-key.md)
- 👛 [Wallet Authentication](architecture/wallet-auth.md)
- ⛓️ [Blockchain Integration](architecture/blockchain.md)

### Deployment
- ☸️ [Kubernetes Setup](deployment/kubernetes.md)
- ☁️ [AWS Deployment](deployment/aws.md)
- 🌐 [GKE Deployment](deployment/gke.md)
- 🔐 [Vault Configuration](deployment/vault.md)
- ⛓️ [Fabric Setup](deployment/fabric.md)

### Security
- 🛡️ [Security Model](security/model.md)
- ⚠️ [Threat Analysis](security/threats.md)
- ✅ [Best Practices](security/best-practices.md)
- 📋 [Compliance (SOC2/GDPR)](security/compliance.md)

### API Reference
- 🔌 [Backend REST API](api/backend.md)
- ⛓️ [Fabric Chaincode](api/chaincode.md)
- 💻 [Frontend SDK](api/frontend.md)

### Development
- 💻 [Dev Setup](development/setup.md)
- 🤝 [Contributing](development/contributing.md)
- 🧪 [Testing](development/testing.md)
- 🐛 [Debugging](development/debugging.md)

### FAQ
- ❓ [General Questions](faq/general.md)
- 🔒 [Security Questions](faq/security.md)
- 🔧 [Troubleshooting](faq/troubleshooting.md)

---

## 💡 What is Pass Chain?

Pass Chain is a revolutionary password manager that uses:

- 🔐 **Zero-Knowledge Encryption** - We literally can't decrypt your passwords
- 🔑 **Split-Key Security** - Key split across Vault, blockchain, and your device
- 👛 **Wallet Authentication** - Use MetaMask instead of master passwords
- ⛓️ **Blockchain Audit Trail** - Immutable logs on Hyperledger Fabric
- ☸️ **Enterprise Ready** - Kubernetes deployment, SOC2 compliant

### How It Works

```
1. Enter password → 2. Encrypt in browser → 3. Split key into 3 parts
                    ↓
    Part 1 → Vault | Part 2 → Blockchain | Part 3 → Your device
                    ↓
4. To decrypt: Wallet signature + any 2 of 3 parts
```

**[Deep Dive into Architecture →](architecture/overview.md)**

---

## 🛡️ Security Guarantees

| What We **CAN'T** Do | What You **GET** |
|----------------------|------------------|
| ❌ Decrypt your passwords | ✅ Client-side encryption |
| ❌ Access your credentials | ✅ Split-key architecture |
| ❌ See plaintext data | ✅ Blockchain audit trail |
| ❌ Recover without wallet | ✅ Complete control |

**[Security Model →](security/model.md)**

---

## 🏗️ Tech Stack

**Frontend**: React, Next.js 14, TailwindCSS, Web3.js, Wagmi  
**Backend**: Go 1.21, Gin, GORM, Fabric SDK  
**Infrastructure**: Kubernetes, Vault, PostgreSQL, Redis, Fabric  
**Blockchain**: Hyperledger Fabric 2.5

**[Full Stack Details →](../ProjectStack.md)**

---

## 📦 Project Structure

```
pass-chain/
├── frontend/          # React/Next.js app
├── backend/           # Go API server
├── blockchain/        # Fabric chaincode
├── infrastructure/    # K8s manifests
└── docs/             # This documentation
```

---

## 🤝 Contributing

We welcome contributions! See [Contributing Guide](development/contributing.md)

---

## 📄 License

MIT License - see [LICENSE](../LICENSE)

---

## 🔗 Links

- 🏠 [Main Site](https://passchain.io)
- 💻 [GitHub](https://github.com/yourusername/pass-chain)
- 📖 [Documentation](https://docs.passchain.io)
- 💬 [Discord](https://discord.gg/passchain)
- 🐦 [Twitter](https://twitter.com/passchain)

---

<div align="center">

**Built with ❤️ by developers who care about security**

*AUUUUFFFF!* 🔥

</div>
