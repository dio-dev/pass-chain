# Hyperledger Fabric Multi-Org Network Configuration

## Network Topology

```
Pass Chain Fabric Network
│
├── Organizations
│   ├── OrgVault (Primary storage org)
│   │   ├── peer0.orgvault.passchain.com
│   │   └── peer1.orgvault.passchain.com
│   │
│   ├── OrgAudit (Independent audit org)
│   │   ├── peer0.orgaudit.passchain.com
│   │   └── peer1.orgaudit.passchain.com
│   │
│   └── OrdererOrg (Consensus org)
│       ├── orderer0.passchain.com (Raft leader)
│       ├── orderer1.passchain.com
│       └── orderer2.passchain.com
│
├── Channel: passchain
│   └── Endorsement Policy: AND('OrgVaultMSP.peer', 'OrgAuditMSP.peer')
│
└── Chaincode: credentials_v2
    ├── Version: 2.0
    ├── Language: Go
    └── Private Data Collections: PDC_PASSCHAIN
```

## Private Data Collections

### PDC_PASSCHAIN

**Purpose**: Store Shard B of encrypted credentials

**Access Policy**:
```
OR('OrgVaultMSP.member', 'OrgAuditMSP.member')
```

**Endorsement Policy**:
```
AND('OrgVaultMSP.peer', 'OrgAuditMSP.peer')
```

This means:
- Both OrgVault AND OrgAudit must approve any transaction
- Data is only accessible to members of these two orgs
- Provides separation of concerns and prevents single-org compromise

**Configuration**:
- `requiredPeerCount`: 1 (minimum peers for private data)
- `maxPeerCount`: 3 (distribute to up to 3 peers)
- `blockToLive`: 0 (never purge - permanent storage)
- `memberOnlyRead`: true (only org members can read)
- `memberOnlyWrite`: true (only org members can write)

## Network Setup Scripts

### 1. Generate Crypto Material

```bash
cd blockchain/network

# Generate certificates for all orgs
cryptogen generate --config=crypto-config.yaml
```

### 2. Generate Genesis Block

```bash
# Create system channel genesis block
configtxgen -profile PassChainOrdererGenesis \
  -channelID system-channel \
  -outputBlock ./genesis.block
```

### 3. Create Channel

```bash
# Create channel transaction
configtxgen -profile PassChainChannel \
  -outputCreateChannelTx ./passchain.tx \
  -channelID passchain

# Create anchor peer updates
configtxgen -profile PassChainChannel \
  -outputAnchorPeersUpdate ./OrgVaultMSPanchors.tx \
  -channelID passchain \
  -asOrg OrgVaultMSP

configtxgen -profile PassChainChannel \
  -outputAnchorPeersUpdate ./OrgAuditMSPanchors.tx \
  -channelID passchain \
  -asOrg OrgAuditMSP
```

### 4. Deploy Chaincode with PDC

```bash
# Package chaincode
peer lifecycle chaincode package credentials_v2.tar.gz \
  --path ../chaincode/credentials/ \
  --lang golang \
  --label credentials_v2_1.0

# Install on OrgVault peers
peer lifecycle chaincode install credentials_v2.tar.gz

# Install on OrgAudit peers  
peer lifecycle chaincode install credentials_v2.tar.gz

# Approve for OrgVault
peer lifecycle chaincode approveformyorg \
  --channelID passchain \
  --name credentials \
  --version 2.0 \
  --package-id <package-id> \
  --sequence 1 \
  --collections-config ../config/collections_config.json \
  --signature-policy "AND('OrgVaultMSP.peer','OrgAuditMSP.peer')"

# Approve for OrgAudit
peer lifecycle chaincode approveformyorg \
  --channelID passchain \
  --name credentials \
  --version 2.0 \
  --package-id <package-id> \
  --sequence 1 \
  --collections-config ../config/collections_config.json \
  --signature-policy "AND('OrgVaultMSP.peer','OrgAuditMSP.peer')"

# Commit chaincode (requires both orgs)
peer lifecycle chaincode commit \
  --channelID passchain \
  --name credentials \
  --version 2.0 \
  --sequence 1 \
  --collections-config ../config/collections_config.json \
  --signature-policy "AND('OrgVaultMSP.peer','OrgAuditMSP.peer')" \
  --peerAddresses peer0.orgvault.passchain.com:7051 \
  --peerAddresses peer0.orgaudit.passchain.com:8051
```

## Endorsement Flow

```
1. Client submits transaction proposal
   ↓
2. Sent to OrgVault.peer0 AND OrgAudit.peer0
   ↓
3. Both peers execute chaincode independently
   ↓
4. Both peers must return same result (endorsement)
   ↓
5. If both endorse, client submits to orderer
   ↓
6. Orderer creates block via Raft consensus
   ↓
7. Block distributed to all peers
   ↓
8. Peers validate and commit to their ledgers
```

## Security Benefits

1. **No Single Point of Trust**
   - OrgVault and OrgAudit must both approve
   - If one org is compromised, data remains secure

2. **Separation of Duties**
   - OrgVault: Manages Vault integration
   - OrgAudit: Independent auditor, ensures compliance

3. **Immutable Audit Trail**
   - All transactions logged on blockchain
   - Cannot be tampered with or deleted

4. **Private Data Protection**
   - Shard B never visible on public ledger
   - Only authorized org members can access PDC

5. **Byzantine Fault Tolerance**
   - Raft consensus tolerates up to (n-1)/2 failures
   - With 3 orderers: can tolerate 1 failure

## Network Ports

| Component | Port | TLS Port |
|-----------|------|----------|
| peer0.orgvault | 7051 | 7052 |
| peer1.orgvault | 7151 | 7152 |
| peer0.orgaudit | 8051 | 8052 |
| peer1.orgaudit | 8151 | 8152 |
| orderer0 | 7050 | 7053 |
| orderer1 | 7150 | 7153 |
| orderer2 | 7250 | 7253 |

## Kubernetes Deployment

See `../../infrastructure/k8s/fabric/` for Kubernetes manifests.

Each peer and orderer runs as a StatefulSet for persistent storage and stable network identities.




