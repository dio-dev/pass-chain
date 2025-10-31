# Cloudflare Tunnel Setup for Pass Chain

This guide helps you access Pass Chain from anywhere using Cloudflare Tunnel.

## What You Get

- 🌐 **app.pass-chain.com** → Frontend (localhost:3000)
- 🔌 **api.pass-chain.com** → Backend (localhost:8080)
- 🔐 **vault.pass-chain.com** → Vault UI (localhost:8200)

## Prerequisites

✅ You already have:
- Cloudflare Tunnel installed
- Tunnel ID: `0d558760-2d23-445f-a10e-436a56e18307`
- Config at: `C:\Users\alex\.cloudflared\config.yaml`

## Setup Steps

### 1. Configure DNS in Cloudflare Dashboard

Go to [Cloudflare Dashboard](https://dash.cloudflare.com) → Your domain → DNS

Add CNAME records:

| Type | Name | Content | Proxy |
|------|------|---------|-------|
| CNAME | app | `0d558760-2d23-445f-a10e-436a56e18307.cfargotunnel.com` | ✅ Proxied |
| CNAME | api | `0d558760-2d23-445f-a10e-436a56e18307.cfargotunnel.com` | ✅ Proxied |
| CNAME | vault | `0d558760-2d23-445f-a10e-436a56e18307.cfargotunnel.com` | ✅ Proxied |

### 2. Start Your Services

**Terminal 1 - Frontend:**
```powershell
cd frontend
npm run dev
```

**Terminal 2 - Backend (via Kubernetes port-forward):**
```powershell
kubectl port-forward svc/passchain-backend 8080:8080 -n passchain
```

**Terminal 3 - Vault (optional):**
```powershell
kubectl port-forward svc/passchain-vault 8200:8200 -n passchain
```

### 3. Start Cloudflare Tunnel

**Terminal 4:**
```powershell
cloudflared tunnel run pass-chain.com
```

Or as a service:
```powershell
cloudflared service install
```

### 4. Update Frontend API URL

Edit `frontend/next.config.js`:

```js
module.exports = {
  env: {
    NEXT_PUBLIC_API_URL: 'https://api.pass-chain.com'  // Changed from localhost
  }
}
```

Restart frontend (Ctrl+C and `npm run dev` again).

### 5. Test It!

Visit:
- **Frontend:** https://app.pass-chain.com
- **Backend:** https://api.pass-chain.com/health
- **Vault:** https://vault.pass-chain.com/ui

## Alternative: ngrok (Quick & Easy)

If you don't want to configure DNS:

### Install ngrok
```powershell
choco install ngrok
```

### Expose Frontend
```powershell
ngrok http 3000
```
→ Get URL like: `https://abc123.ngrok.io`

### Expose Backend
```powershell
ngrok http 8080
```
→ Get URL like: `https://def456.ngrok.io`

### Update Frontend
Edit `frontend/.env.local`:
```
NEXT_PUBLIC_API_URL=https://def456.ngrok.io
```

Restart frontend.

## Troubleshooting

### "Cannot reach backend from frontend"

**Check CORS** in `backend/internal/middleware/middleware.go`:

```go
config := cors.Config{
    AllowOrigins:     []string{"https://app.pass-chain.com", "http://localhost:3000"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Wallet-Address", "X-Signature"},
    AllowCredentials: true,
}
```

Restart backend:
```powershell
kubectl delete pod -l app=passchain-backend -n passchain
kubectl port-forward svc/passchain-backend 8080:8080 -n passchain
```

### "Cloudflare tunnel not connecting"

Check tunnel status:
```powershell
cloudflared tunnel list
cloudflared tunnel info pass-chain.com
```

View logs:
```powershell
cloudflared tunnel run pass-chain.com
```

### "Cannot save credentials"

Make sure backend can reach Vault:
```powershell
# In another terminal
kubectl port-forward svc/passchain-vault 8200:8200 -n passchain
```

Backend needs to connect to `http://localhost:8200`.

## Security Notes

### For Development
- ✅ Cloudflare Tunnel is secure (encrypted)
- ✅ No ports opened on your router
- ✅ Traffic goes through Cloudflare

### For Production
- 🔐 Enable Cloudflare Access (authentication)
- 🔐 Use Cloudflare WAF
- 🔐 Set up rate limiting
- 🔐 Deploy to real GKE cluster (not localhost)

## Quick Reference

### Start Everything (PowerShell)

```powershell
# Terminal 1: Frontend
cd C:\Users\alex\projects\pass-chain\frontend
npm run dev

# Terminal 2: Backend
kubectl port-forward svc/passchain-backend 8080:8080 -n passchain

# Terminal 3: Vault
kubectl port-forward svc/passchain-vault 8200:8200 -n passchain

# Terminal 4: Cloudflare Tunnel
cloudflared tunnel run pass-chain.com
```

### Stop Everything

```powershell
# Ctrl+C in each terminal

# Or kill all node processes
taskkill /F /IM node.exe

# Stop port-forwards
kubectl delete pod --field-selector=status.phase==Running -n passchain
```

## Access from Mobile/Other Devices

Once Cloudflare Tunnel is running:
- 📱 **Phone:** Open https://app.pass-chain.com
- 💻 **Other laptop:** Open https://app.pass-chain.com
- 🌍 **Anywhere:** Works globally!

**Just make sure:**
1. Your dev machine is running
2. All services are port-forwarded
3. Cloudflare Tunnel is active

---

**AUUUUFFFF!** 🔥

