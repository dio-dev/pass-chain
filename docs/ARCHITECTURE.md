# Pass Chain - Architecture Documentation

## System Overview

Pass Chain is a decentralized password management system that uses a revolutionary split-key architecture to ensure maximum security. The system splits encrypted credentials between HashiCorp Vault and Hyperledger Fabric blockchain, making it virtually impossible for any single component to compromise user data.

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Client Layer                             │
│  ┌──────────────┐         ┌────────────────┐                   │
│  │   Browser    │◄───────►│  Crypto Wallet │                   │
│  │  (Next.js)   │         │  (MetaMask)    │                   │
│  └──────┬───────┘         └────────────────┘                   │
│         │                                                        │
└─────────┼────────────────────────────────────────────────────────┘
          │ HTTPS
          │
┌─────────▼────────────────────────────────────────────────────────┐
│                      Application Layer                            │
│  ┌───────────────────────────────────────────────────┐           │
│  │              Go Backend API (Gin)                  │           │
│  │  ┌──────────┬──────────┬───────────┬──────────┐  │           │
│  │  │   Auth   │Credential│  Payment  │  Crypto  │  │           │
│  │  │  Service │ Service  │  Service  │ Service  │  │           │
│  │  └──────────┴──────────┴───────────┴──────────┘  │           │
│  └───────────────────────────────────────────────────┘           │
│                           │                                       │
└───────────────────────────┼───────────────────────────────────────┘
                            │
            ┌───────────────┼───────────────┐
            │               │               │
┌───────────▼─────┐  ┌──────▼──────┐  ┌────▼────────┐
│  HashiCorp      │  │  Hyperledger │  │  PostgreSQL │
│  Vault          │  │  Fabric      │  │  Database   │
│  (Shard A)      │  │  (Shard B +  │  │  (Metadata) │
│                 │  │   Audit Log) │  │             │
└─────────────────┘  └──────────────┘  └─────────────┘
```

## Component Architecture

### 1. Frontend (Next.js Application)

```
frontend/
├── Pages & Routes
│   ├── Home (/)
│   ├── Dashboard (/dashboard)
│   ├── Credentials (/credentials)
│   └── Analytics (/analytics)
│
├── Web3 Integration
│   ├── WalletConnect
│   ├── MetaMask SDK
│   └── Ethers.js
│
├── State Management
│   ├── Zustand (Global State)
│   ├── React Query (Server State)
│   └── Local Storage (Persistence)
│
└── UI Components
    ├── shadcn/ui
    ├── Tailwind CSS
    └── Framer Motion
```

**Key Features:**
- Client-side encryption/decryption
- Wallet signature generation
- Real-time updates
- Responsive design

### 2. Backend (Go API Server)

```
backend/
├── HTTP Layer (Gin)
│   ├── Middleware
│   │   ├── Authentication
│   │   ├── Rate Limiting
│   │   └── Logging
│   └── Handlers
│       ├── Credentials
│       ├── Auth
│       └── Payments
│
├── Service Layer
│   ├── Credential Service
│   │   ├── Create
│   │   ├── Retrieve
│   │   └── Delete
│   ├── Vault Service
│   │   ├── Shard Storage
│   │   └── Shard Retrieval
│   ├── Blockchain Service
│   │   ├── Fabric SDK
│   │   └── Transaction Submission
│   └── Crypto Service
│       └── Shamir's Secret Sharing
│
├── Repository Layer
│   ├── PostgreSQL
│   └── Redis Cache
│
└── Integration Layer
    ├── Vault Client
    └── Fabric SDK
```

**Key Responsibilities:**
- API request handling
- Business logic execution
- Data sharding
- Payment verification
- Access logging

### 3. HashiCorp Vault

```
Vault
├── Secrets Engines
│   ├── KV v2 (Credential Shards)
│   ├── Transit (Encryption)
│   └── PKI (Certificates)
│
├── Auth Methods
│   ├── Kubernetes (Pod Auth)
│   └── AppRole (Service Auth)
│
├── Policies
│   ├── Credential Read/Write
│   └── Admin Access
│
└── Storage Backend
    └── Raft (Integrated Storage)
```

**Security Features:**
- Auto-unseal with AWS KMS
- Encryption at rest
- Audit logging
- Access policies

### 4. Hyperledger Fabric

```
Fabric Network
├── Organizations
│   ├── Org1 (Pass Chain)
│   └── OrdererOrg
│
├── Peers
│   ├── peer0.org1
│   └── peer1.org1
│
├── Orderer
│   ├── orderer0 (Raft)
│   ├── orderer1
│   └── orderer2
│
├── Channel
│   └── passchainchannel
│
└── Chaincode
    └── credentials
        ├── CreateCredentialShard
        ├── ReadCredentialShard
        ├── LogAccess
        └── GetAccessLogs
```

**Key Features:**
- Private blockchain
- Immutable ledger
- Smart contracts
- Consensus (Raft)

### 5. Database Layer

```
PostgreSQL
├── Tables
│   ├── credentials
│   ├── credential_accesses
│   ├── users
│   └── payments
│
└── Indexes
    ├── wallet_address
    ├── credential_id
    └── created_at
```

```
Redis
├── Cache
│   ├── User sessions
│   ├── Credential metadata
│   └── Rate limiting
│
└── Pub/Sub
    └── Real-time updates
```

## Data Flow

### Credential Storage Flow

```
1. Client Side
   ├── User enters password
   ├── Encrypt with wallet-derived key (AES-256)
   └── Send encrypted data to backend

2. Backend
   ├── Receive encrypted data
   ├── Split using Shamir's Secret Sharing (2-of-2)
   ├── Store Shard A in Vault
   ├── Store Shard B on Blockchain
   ├── Save metadata in PostgreSQL
   └── Return transaction IDs

3. Vault
   ├── Receive Shard A
   ├── Encrypt with Transit Engine
   └── Store in KV engine

4. Blockchain
   ├── Receive Shard B
   ├── Submit transaction to Fabric
   ├── Consensus validation
   └── Store in world state
```

### Credential Retrieval Flow

```
1. Client Request
   ├── Sign request with wallet
   ├── Submit payment proof
   └── Request credential ID

2. Backend Validation
   ├── Verify wallet signature
   ├── Verify payment transaction
   └── Check access permissions

3. Shard Retrieval
   ├── Retrieve Shard A from Vault
   ├── Retrieve Shard B from Blockchain
   └── Reconstruct encrypted data

4. Blockchain Logging
   ├── Log access to Fabric
   ├── Record timestamp, IP, wallet
   └── Create immutable audit entry

5. Return to Client
   ├── Send reconstructed encrypted data
   └── Client decrypts locally
```

## Security Architecture

### Defense in Depth

```
Layer 1: Client-Side Encryption
         └── AES-256-GCM with wallet-derived key

Layer 2: Data Sharding
         └── Shamir's Secret Sharing (2-of-2)

Layer 3: Distributed Storage
         ├── Vault (Shard A)
         └── Blockchain (Shard B)

Layer 4: Access Control
         ├── Wallet authentication
         ├── JWT tokens
         └── Rate limiting

Layer 5: Audit Logging
         └── Immutable blockchain logs

Layer 6: Network Security
         ├── TLS/HTTPS
         ├── VPC isolation
         └── Security groups
```

### Attack Resistance

| Attack Vector | Mitigation |
|---------------|------------|
| Server Compromise | Zero-knowledge architecture |
| Vault Breach | Only has Shard A |
| Blockchain Breach | Only has Shard B |
| MITM Attack | TLS + Pre-encryption |
| Replay Attack | Nonces + Timestamps |
| Brute Force | Rate limiting + Monitoring |
| Wallet Theft | User responsibility + Alerts |

## Infrastructure Architecture (AWS)

```
AWS Cloud
│
├── VPC (10.0.0.0/16)
│   ├── Public Subnets (3 AZs)
│   │   ├── NAT Gateways
│   │   └── Load Balancers
│   │
│   └── Private Subnets (3 AZs)
│       ├── EKS Cluster
│       │   ├── Backend Pods
│       │   ├── Frontend Pods
│       │   ├── Vault Pods
│       │   └── Monitoring
│       ├── RDS PostgreSQL
│       └── ElastiCache Redis
│
├── S3 Buckets
│   ├── Frontend Assets
│   ├── Backups
│   └── Logs
│
├── CloudFront
│   └── CDN for Frontend
│
└── Route53
    └── DNS Management
```

## Deployment Architecture

```
Kubernetes Cluster
│
├── Namespaces
│   ├── production
│   ├── staging
│   ├── vault
│   └── monitoring
│
├── Deployments
│   ├── backend (3 replicas)
│   ├── frontend (3 replicas)
│   └── vault (3 replicas - HA)
│
├── Services
│   ├── backend-service (ClusterIP)
│   ├── frontend-service (ClusterIP)
│   └── vault-service (ClusterIP)
│
├── Ingress
│   ├── ALB Ingress Controller
│   └── TLS Termination
│
└── Monitoring
    ├── Prometheus
    ├── Grafana
    └── Loki
```

## Scalability

### Horizontal Scaling

- **Frontend**: Auto-scale based on requests (3-10 pods)
- **Backend**: Auto-scale based on CPU (3-20 pods)
- **Database**: Read replicas for queries
- **Blockchain**: Add more peers for throughput

### Performance Targets

| Metric | Target | Current |
|--------|--------|---------|
| API Response Time | < 500ms | TBD |
| Credential Storage | < 2s | TBD |
| Credential Retrieval | < 3s | TBD |
| Concurrent Users | 10,000+ | TBD |
| Transactions/sec | 1,000+ | TBD |

## Monitoring & Observability

```
Observability Stack
├── Metrics (Prometheus)
│   ├── API request rates
│   ├── Error rates
│   ├── Latency percentiles
│   └── Resource usage
│
├── Logs (Loki)
│   ├── Application logs
│   ├── Access logs
│   ├── Error logs
│   └── Audit logs
│
├── Traces (Jaeger)
│   ├── Request tracing
│   ├── Service dependencies
│   └── Performance bottlenecks
│
└── Dashboards (Grafana)
    ├── System health
    ├── Business metrics
    ├── Alerts
    └── SLOs
```

## Disaster Recovery

### Backup Strategy

```
Daily Backups
├── PostgreSQL
│   ├── Automated snapshots
│   └── 30-day retention
│
├── Vault
│   ├── Raft snapshots
│   └── S3 storage
│
└── Blockchain
    ├── Ledger snapshots
    └── S3 storage
```

### Recovery Procedures

1. **Database Failure**: Restore from RDS snapshot
2. **Vault Failure**: Deploy new Vault, restore from snapshot
3. **Blockchain Failure**: Restore from ledger backup
4. **Region Failure**: Multi-region failover (future)

## Future Enhancements

### Phase 2
- Multi-signature wallets
- Mobile applications
- Browser extensions
- Advanced analytics

### Phase 3
- Zero-knowledge proofs
- Cross-chain support
- Decentralized identity
- AI-powered security

---

**This architecture is designed for security, scalability, and reliability.**




