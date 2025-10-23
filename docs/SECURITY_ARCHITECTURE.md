# Pass Chain - Security Architecture

## Overview

Pass Chain implements a multi-layered security architecture that ensures credentials are never stored in a single location and can only be accessed by the rightful owner.

## Core Security Principles

### 1. Split-Key Architecture

Credentials are split using **Shamir's Secret Sharing** into two shards:

```
Original Credential
       |
       v
  [Encryption]
       |
       v
   [Sharding]
      / \
     /   \
    v     v
Shard A  Shard B
   |       |
   v       v
 Vault  Blockchain
```

**Key Properties**:
- Neither shard alone can decrypt the credential
- Both shards required to reconstruct
- Threshold scheme: k-of-n (currently 2-of-2)
- Computationally infeasible to break

### 2. Zero-Knowledge Architecture

The backend server **never sees** plaintext credentials:

```
Client Side           Server Side
-----------           -----------
Generate Password
     |
Encrypt with
Wallet Key
     |
  Shard Data  ------>  Receive
     |                 Encrypted
     |                 Shards Only
     |                    |
     |                    v
     |               Store Shard A
     |               in Vault
     |                    |
     |                    v
     |               Store Shard B
     |               on Blockchain
     |
     v
Send to Server
```

**Benefits**:
- Server compromise doesn't expose passwords
- No plaintext ever transmitted
- Client-side encryption/decryption only

### 3. Wallet-Based Authentication

Authentication uses cryptographic signatures from user wallets:

```
1. Server sends challenge nonce
2. Client signs nonce with wallet private key
3. Server verifies signature against wallet address
4. JWT token issued for session
```

**Advantages**:
- No password to steal
- Cryptographically secure
- User controls private key
- Can't be phished easily

## Security Layers

### Layer 1: Client-Side Encryption

```typescript
// Client encrypts before sending to server
const walletKey = await wallet.signMessage("Pass Chain Encryption Key");
const encryptedData = encrypt(password, walletKey);
```

**Encryption Algorithm**: AES-256-GCM
**Key Derivation**: PBKDF2 from wallet signature
**IV**: Unique per credential

### Layer 2: Data Sharding

```go
// Server splits encrypted data
shardA, shardB := shamir.Split(encryptedData, 2, 2)
```

**Algorithm**: Shamir's Secret Sharing
**Threshold**: 2-of-2 (both shards required)
**Shard Storage**: Separate systems

### Layer 3: Vault Storage

```
HashiCorp Vault
├── Transit Engine (encryption)
├── KV Engine (shard storage)
├── Auto-unseal (AWS KMS)
└── Audit Logging
```

**Features**:
- Encryption at rest
- Access policies
- Audit logs
- Auto-unseal with KMS

### Layer 4: Blockchain Immutability

```
Hyperledger Fabric
├── Private Channel
├── Encrypted Shard Storage
├── Immutable Audit Trail
└── Distributed Consensus
```

**Benefits**:
- Tamper-proof logs
- Distributed storage
- Consensus validation
- Complete audit trail

## Attack Vectors & Mitigations

### 1. Server Compromise

**Risk**: Attacker gains access to backend server

**Mitigation**:
- ✅ Zero-knowledge architecture
- ✅ Credentials encrypted client-side
- ✅ Server never sees plaintext
- ✅ Both Vault AND blockchain needed

**Result**: Attacker gets encrypted shards only

### 2. Vault Compromise

**Risk**: Attacker gains access to Vault

**Mitigation**:
- ✅ Only has Shard A
- ✅ Shard B stored on blockchain
- ✅ Both shards required
- ✅ Vault audit logs all access

**Result**: Cannot reconstruct credentials

### 3. Blockchain Compromise

**Risk**: Attacker gains access to blockchain node

**Mitigation**:
- ✅ Only has Shard B
- ✅ Shard A stored in Vault
- ✅ Both shards required
- ✅ Blockchain immutable logs

**Result**: Cannot reconstruct credentials

### 4. Man-in-the-Middle Attack

**Risk**: Attacker intercepts traffic

**Mitigation**:
- ✅ TLS/HTTPS enforced
- ✅ Certificate pinning
- ✅ Data pre-encrypted client-side
- ✅ Wallet signature verification

**Result**: Attacker sees encrypted data only

### 5. Wallet Compromise

**Risk**: User's private key stolen

**Mitigation**:
- ⚠️ This is the critical security boundary
- ✅ Multi-signature option (future)
- ✅ Usage alerts
- ✅ IP-based anomaly detection
- ✅ Rate limiting

**Result**: User must protect their wallet

### 6. Replay Attack

**Risk**: Attacker replays captured requests

**Mitigation**:
- ✅ Nonce in authentication
- ✅ Timestamp validation
- ✅ Short-lived JWT tokens
- ✅ One-time signature challenges

**Result**: Replay blocked

### 7. Brute Force Attack

**Risk**: Attacker tries many requests

**Mitigation**:
- ✅ Rate limiting (per IP)
- ✅ Rate limiting (per wallet)
- ✅ CAPTCHA for suspicious activity
- ✅ Exponential backoff

**Result**: Attack infeasible

## Payment Security

### Storage Payment

```
1. User approves transaction in wallet
2. Smart contract validates payment
3. Backend verifies on-chain transaction
4. Only then stores credential
```

### Usage Payment

```
1. User requests credential reveal
2. Payment required before retrieval
3. Microtransaction via Layer 2
4. Shard retrieval only after payment
```

## Audit Trail

Every action is logged on blockchain:

```json
{
  "action": "credential_access",
  "credential_id": "uuid",
  "wallet_address": "0x...",
  "ip_hash": "sha256(...)",
  "timestamp": "2024-01-01T00:00:00Z",
  "success": true,
  "tx_hash": "0x..."
}
```

**Properties**:
- Immutable
- Timestamped
- Traceable
- Verifiable

## Compliance

### GDPR Compliance

- ✅ Right to access (view credentials)
- ✅ Right to deletion (delete credentials)
- ⚠️ Right to erasure (blockchain consideration)
- ✅ Data portability (export feature)
- ✅ Consent required

**Note**: Blockchain data is immutable. We hash/encrypt PII.

### SOC 2 Targets

- ✅ Security
- ✅ Availability
- ✅ Processing Integrity
- ✅ Confidentiality
- ⚠️ Privacy (in progress)

## Security Best Practices

### For Users

1. **Protect Your Wallet**
   - Use hardware wallet
   - Never share private key
   - Enable 2FA on exchanges

2. **Verify Transactions**
   - Check contract address
   - Verify payment amounts
   - Review transaction before signing

3. **Monitor Access**
   - Check audit logs regularly
   - Set up alerts
   - Report suspicious activity

### For Developers

1. **Code Security**
   - Regular security audits
   - Dependency scanning
   - Static code analysis
   - Penetration testing

2. **Deployment Security**
   - Secrets in Vault only
   - No hardcoded credentials
   - TLS everywhere
   - Principle of least privilege

3. **Monitoring**
   - Real-time alerts
   - Anomaly detection
   - Audit log review
   - Incident response plan

## Security Roadmap

### Phase 1 (Current)
- [x] Split-key architecture
- [x] Wallet authentication
- [x] Blockchain audit logs
- [x] Vault integration

### Phase 2
- [ ] Multi-signature support
- [ ] Hardware wallet integration
- [ ] Biometric authentication
- [ ] Anomaly detection ML

### Phase 3
- [ ] Zero-knowledge proofs
- [ ] Quantum-resistant encryption
- [ ] Decentralized identity
- [ ] Cross-chain support

## Security Audits

- **Code Audit**: Pending
- **Smart Contract Audit**: Pending
- **Penetration Testing**: Pending
- **Bug Bounty**: Launch planned

## Reporting Security Issues

**DO NOT** open public issues for security vulnerabilities.

Email: security@passchain.io
PGP Key: [Link to PGP key]

Expected response time: 24 hours

---

**Security is our top priority. If you have concerns or suggestions, please reach out.**




