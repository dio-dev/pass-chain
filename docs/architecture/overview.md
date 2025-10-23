# System Architecture

Pass Chain's architecture is designed for maximum security through distribution and zero-knowledge principles.

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────┐
│                  User's Browser                          │
│  ┌────────────────────────────────────────────────────┐ │
│  │  1. Enter password                                 │ │
│  │  2. Generate encryption key (256-bit)              │ │
│  │  3. Encrypt with XChaCha20-Poly1305               │ │
│  │  4. Split key into 3 shares (Shamir)              │ │
│  └────────────────────────────────────────────────────┘ │
└──────────────┬──────────────────────────────────────────┘
               │
    ┌──────────┴──────────┐
    │                     │
    ▼                     ▼
┌────────────┐      ┌─────────────┐
│   Share 1  │      │   Share 2   │
│  → Vault   │      │→ Blockchain │
└────────────┘      └─────────────┘
                          
    Share 3 → localStorage (backup)
```

## Components

### Frontend (React/Next.js)
**Purpose**: Client-side encryption and wallet integration

**Responsibilities**:
- Generate encryption keys
- Encrypt/decrypt passwords locally
- Split keys using Shamir Secret Sharing
- Wallet authentication (MetaMask)
- Zero-knowledge proof generation

**Technology**:
- Next.js 14 (App Router)
- Web3.js, Wagmi, RainbowKit
- @noble/ciphers for encryption
- TailwindCSS for UI

### Backend (Go)
**Purpose**: Orchestration and API coordination

**Responsibilities**:
- Coordinate Vault and blockchain operations
- Verify wallet signatures
- Manage encrypted data in PostgreSQL
- Audit logging
- **CANNOT decrypt passwords** (by design)

**Technology**:
- Go 1.21
- Gin web framework
- GORM (PostgreSQL ORM)
- Hyperledger Fabric SDK
- Zap logger

### HashiCorp Vault
**Purpose**: Secure storage of Share 1

**Configuration**:
- KV v2 secrets engine
- Auto-unseal (AWS KMS in production)
- Kubernetes auth
- Audit logging enabled

**Data Stored**:
```json
{
  "share1": "base64_encoded_share",
  "created_at": "timestamp"
}
```

### Hyperledger Fabric (Optional)
**Purpose**: Immutable storage of Share 2 + audit trail

**Network**:
- Channel: `passchain`
- Chaincode: `credentials`
- Organizations: Org1MSP (expandable)
- Consensus: Raft

**Private Data Collections**:
- Share 2 stored in PDC (private to user)
- Audit logs visible to all

### PostgreSQL
**Purpose**: Encrypted data storage

**Schema**:
```sql
credentials (
  id UUID PRIMARY KEY,
  credential_name TEXT NOT NULL,
  username TEXT NOT NULL,
  url TEXT,
  encrypted_data TEXT NOT NULL, -- XChaCha20 ciphertext
  nonce TEXT NOT NULL,
  wallet_address TEXT NOT NULL,
  vault_path TEXT NOT NULL,
  blockchain_tx_id TEXT,
  share2 TEXT, -- Fallback if Fabric unavailable
  created_at TIMESTAMP,
  last_accessed TIMESTAMP
)

audit_logs (
  id UUID PRIMARY KEY,
  credential_id UUID,
  wallet_address TEXT NOT NULL,
  action TEXT NOT NULL, -- create/read/delete
  timestamp TIMESTAMP,
  tx_hash TEXT, -- Fabric transaction ID
  ip_address TEXT
)
```

## Data Flow: Store Credential

```
1. User enters password in browser
   ↓
2. Browser generates random 256-bit key
   ↓
3. Password encrypted: ciphertext = XChaCha20(password, key)
   ↓
4. Key split: [share1, share2, share3] = Shamir(key, 2-of-3)
   ↓
5. User signs message with wallet
   ↓
6. POST /api/v1/credentials {
     encryptedData: ciphertext,
     nonce: random_nonce,
     share1, share2, share3,
     walletAddress,
     signature
   }
   ↓
7. Backend verifies signature
   ↓
8. Backend stores:
   - Share1 → Vault
   - Share2 → Fabric blockchain
   - Ciphertext → PostgreSQL
   ↓
9. Browser stores share3 in localStorage
   ↓
10. Fabric logs "create" action
```

## Data Flow: Retrieve Credential

```
1. User clicks "Reveal Password"
   ↓
2. User signs message with wallet
   ↓
3. GET /api/v1/credentials/:id
   Headers: X-Wallet-Address, X-Signature
   ↓
4. Backend verifies signature
   ↓
5. Backend retrieves:
   - Share1 ← Vault
   - Share2 ← Fabric
   - Ciphertext ← PostgreSQL
   ↓
6. Backend returns {share1, share2, ciphertext, nonce}
   ↓
7. Browser retrieves share3 from localStorage
   ↓
8. Browser reconstructs key: key = Shamir([share1, share2])
   ↓
9. Browser decrypts: password = XChaCha20_decrypt(ciphertext, key)
   ↓
10. Password displayed (copy to clipboard)
    ↓
11. Fabric logs "read" action
```

## Security Layers

| Layer | Protection |
|-------|------------|
| **Network** | TLS 1.3, AWS VPC, Security Groups |
| **Kubernetes** | RBAC, Network Policies, Pod Security |
| **Vault** | Encryption at rest, Access policies |
| **Database** | Encrypted at rest, SSL connections |
| **Fabric** | MSP auth, Endorsement policies |
| **Application** | Zero-knowledge, No plaintext storage |
| **Client** | CSP headers, XSS protection |

## Failure Scenarios

### Scenario 1: Vault Compromised
- ❌ Attacker gets Share1
- ✅ Still need Share2 (blockchain) + wallet signature
- ✅ Passwords remain encrypted

### Scenario 2: Database Compromised
- ❌ Attacker gets ciphertext
- ✅ Still need key (requires 2 shares + wallet)
- ✅ Passwords remain encrypted

### Scenario 3: Blockchain Compromised
- ❌ Attacker gets Share2
- ✅ Still need Share1 (Vault) + wallet signature
- ✅ Passwords remain encrypted

### Scenario 4: Backend Compromised
- ❌ Attacker controls backend
- ✅ Backend never sees plaintext
- ✅ Backend can't decrypt (no wallet signature)
- ✅ Passwords remain encrypted

### Scenario 5: User Loses Device
- ❌ Lost Share3 (localStorage)
- ✅ Still have Share1 (Vault) + Share2 (blockchain)
- ✅ Login from new device works!

## Scaling

**Horizontal Scaling**:
- Backend: Multiple pods (stateless)
- Vault: HA mode with Raft
- PostgreSQL: Read replicas
- Fabric: Multiple peers per org

**Geographic Distribution**:
- Multi-region Kubernetes clusters
- Vault replication
- Fabric multi-region channels
- CDN for frontend

## Monitoring

- **Backend**: Prometheus + Grafana
- **Vault**: Audit logs + CloudWatch
- **Fabric**: Hyperledger Explorer
- **Database**: pgAdmin + CloudWatch
- **APM**: Distributed tracing

---

**Next**: [Client-Side Encryption](./encryption.md)

