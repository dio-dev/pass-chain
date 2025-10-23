# Hyperledger Fabric Setup

## Backend is Ready

Your backend has full Fabric SDK integration:
- `backend/internal/services/fabric.go` - Fabric client
- `backend/internal/api/handlers/credentials.go` - Stores Share2 on blockchain
- Graceful fallback if Fabric unavailable

---

## Get Fabric Running

### Option 1: WSL2 (Recommended for Windows)

```powershell
# Install WSL2
wsl --install

# Restart computer

# Open Ubuntu terminal
wsl

# Navigate to project
cd /mnt/c/Users/alex/projects/pass-chain/blockchain/fabric

# Download Fabric
curl -sSL https://bit.ly/2ysbOFE | bash -s -- 2.5.0 1.5.5

# Start test network
cd fabric-samples/test-network
./network.sh up createChannel -c passchain

# Deploy chaincode
./network.sh deployCC -ccn credentials -ccp ../../chaincode/credentials -ccl go
```

Peer runs at `localhost:7051` - backend connects automatically!

---

### Option 2: Docker Desktop (Linux Containers)

Same as above but run in Git Bash or PowerShell with Docker.

---

### Option 3: Production (GKE/EKS)

Use Hyperledger Fabric Operator on real Kubernetes cluster.

---

## Verify Connection

```powershell
# Check backend logs
kubectl logs -n passchain -l app=passchain-backend -f

# Should see: "✅ Fabric client initialized - blockchain integration active!"
```

---

## What Happens

**With Fabric:**
- Share1 → Vault
- Share2 → **Fabric blockchain** (immutable!)
- Audit logs → **Fabric blockchain**
- Real transaction IDs in explorer

**Without Fabric:**
- Share1 → Vault
- Share2 → PostgreSQL (fallback)
- Audit logs → PostgreSQL
- App still works perfectly!

---

**AUUUUFFFF!** 🔥

