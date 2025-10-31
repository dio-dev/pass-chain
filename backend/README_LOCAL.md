# Quick Start Backend Locally

This guide helps you run the backend locally for development.

## Prerequisites

Make sure these are port-forwarded from Kubernetes:

```powershell
# Terminal 1: PostgreSQL
kubectl port-forward svc/passchain-postgresql 5432:5432 -n passchain

# Terminal 2: Redis
kubectl port-forward svc/passchain-redis-master 6379:6379 -n passchain

# Terminal 3: Vault
kubectl port-forward svc/passchain-vault 8200:8200 -n passchain
```

## Setup

1. **Copy environment file:**
```powershell
cd backend
cp .env.example .env
```

The `.env` file is already configured for local development with port-forwarding.

2. **Run backend:**
```powershell
go run cmd/server/main.go
```

Or build and run:
```powershell
go build -o passchain-backend.exe cmd/server/main.go
.\passchain-backend.exe
```

## Test It

```powershell
# Health check
curl.exe http://localhost:8080/health

# Get credentials (replace with your wallet address)
curl.exe http://localhost:8080/api/v1/credentials?wallet=0xYOUR_WALLET_ADDRESS
```

## Troubleshooting

### "Failed to retrieve encryption key"

**Cause:** Credential was created before Vault was accessible.

**Fix:** Delete and recreate the credential, OR manually add share1 to Vault:

```powershell
.\fix-vault-credential.ps1
```

### "Connection refused" errors

Make sure all services are port-forwarded:
```powershell
# Check what's listening on each port
netstat -an | findstr "5432 6379 8200"
```

You should see:
```
TCP    0.0.0.0:5432    LISTENING   (PostgreSQL)
TCP    0.0.0.0:6379    LISTENING   (Redis)
TCP    0.0.0.0:8200    LISTENING   (Vault)
```

If not, restart the port-forwards.

### Database connection issues

Get the correct password:
```powershell
kubectl get secret passchain-postgresql -n passchain -o jsonpath='{.data.postgres-password}' | ForEach-Object { [System.Text.Encoding]::UTF8.GetString([System.Convert]::FromBase64String($_)) }
```

Update `backend/.env` with the correct password.

---

**AUUUUFFFF!** 🔥

