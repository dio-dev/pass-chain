# Quick Start

Get Pass Chain running in 5 minutes!

## Prerequisites

- Docker Desktop
- Kubernetes (Minikube)
- Node.js 18+
- MetaMask browser extension

## Step 1: Start Infrastructure

```powershell
# Clone repository
git clone https://github.com/yourusername/pass-chain
cd pass-chain

# Start Minikube + deploy all services
./start-minikube.ps1
```

This deploys:
- ✅ Backend API (Go)
- ✅ PostgreSQL database
- ✅ HashiCorp Vault
- ✅ Redis cache

**Wait ~2 minutes** for all pods to be ready.

## Step 2: Port Forward Backend

Open a new terminal:

```powershell
kubectl port-forward svc/passchain-backend 8080:8080 -n passchain
```

Keep this running.

## Step 3: Start Frontend

Open another terminal:

```powershell
cd frontend
npm install
npm run dev
```

## Step 4: Access the App

1. Open http://localhost:3000
2. Click "Connect Wallet"
3. Approve MetaMask connection
4. Click "Get Started"

🎉 **You're in!**

## First Credential

1. Click "Add Credential"
2. Fill in:
   - Name: "GitHub"
   - Username: your username
   - Password: your password
   - URL: https://github.com
3. Click "Save"
4. Sign with MetaMask

Your password is now:
- ✅ Encrypted in browser
- ✅ Split across Vault + Database + localStorage
- ✅ Secured by your wallet

## Reveal Password

1. Click the eye icon on a credential
2. Sign with MetaMask
3. Password decrypted in browser
4. Click "Copy" to clipboard

## View Blockchain Explorer

1. Click "Blockchain Explorer" button
2. See all your activity
3. View transaction IDs (when Fabric is connected)

## Next Steps

- [Architecture Overview](../architecture/overview.md)
- [Deploy to AWS](../deployment/aws.md)
- [Security Model](../security/model.md)

---

**That's it! You're securing passwords like a pro.** 🔥

