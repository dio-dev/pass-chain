---
hidden: true
---

# Blockchain - Hyperledger Fabric Network

Hyperledger Fabric blockchain network for Pass Chain.

## Overview

This directory contains the Hyperledger Fabric network configuration, chaincode (smart contracts), and deployment scripts for Pass Chain.

## Architecture

```
Blockchain Network
├── Organizations
│   ├── Org1 (Pass Chain Org)
│   ├── Org2 (Optional backup org)
│   └── Orderer Org
├── Channel: passchainchannel
└── Chaincode: credentials
```

## Components

### Chaincode (Smart Contracts)

Located in `chaincode/credentials/`:

* **CredentialShard** - Stores encrypted credential shards
* **AccessLog** - Immutable access logs
* **Functions**:
  * `CreateCredentialShard` - Store a credential shard
  * `ReadCredentialShard` - Retrieve a credential shard
  * `DeleteCredentialShard` - Delete a credential
  * `LogAccess` - Log credential access
  * `GetAccessLogs` - Query access logs

### Network Configuration

* **Organizations**: 1-2 peer organizations + orderer organization
* **Peers**: 2 per organization for redundancy
* **Orderer**: Raft consensus (3 orderers)
* **Channel**: Private channel for Pass Chain
* **State Database**: CouchDB for rich queries

## Prerequisites

* Docker 24+
* Docker Compose
* Hyperledger Fabric binaries 2.5+
* Go 1.21+ (for chaincode development)

## Getting Started

### 1. Install Fabric Binaries

```bash
curl -sSL https://bit.ly/2ysbOFE | bash -s -- 2.5.0 1.5.5
```

### 2. Start the Network

```bash
cd blockchain
./network.sh up createChannel -ca
```

This will:

* Generate crypto material
* Start all containers
* Create the channel
* Join peers to channel

### 3. Deploy Chaincode

```bash
./network.sh deployCC -ccn credentials -ccp ./chaincode/credentials -ccl go
```

### 4. Test Chaincode

```bash
# Initialize ledger
./scripts/test-chaincode.sh init

# Create a credential shard
./scripts/test-chaincode.sh create

# Read credential shard
./scripts/test-chaincode.sh read <id>

# Query all shards
./scripts/test-chaincode.sh query
```

## Network Scripts

### network.sh

Main script for network management:

```bash
# Start network
./network.sh up

# Create channel
./network.sh createChannel

# Deploy chaincode
./network.sh deployCC -ccn credentials -ccp ./chaincode/credentials -ccl go

# Upgrade chaincode
./network.sh upgradeCC -ccn credentials -ccv 2.0

# Stop network
./network.sh down

# Restart network
./network.sh restart
```

## Chaincode Development

### Project Structure

```
chaincode/credentials/
├── chaincode.go      # Main chaincode logic
├── go.mod           # Go dependencies
└── go.sum           # Dependency checksums
```

### Local Testing

```bash
cd chaincode/credentials
go test -v
```

### Upgrading Chaincode

1. Update version in chaincode
2. Update `go.mod` if needed
3. Deploy with new version:

```bash
./network.sh upgradeCC -ccn credentials -ccv 2.0
```

## Data Models

### CredentialShard

```json
{
  "id": "uuid",
  "wallet_address": "0x...",
  "shard_data": "encrypted_shard_b",
  "vault_path": "secret/data/credentials/uuid",
  "storage_fee": "0.001",
  "created_at": "2024-01-01T00:00:00Z"
}
```

### AccessLog

```json
{
  "id": "log_uuid",
  "credential_id": "uuid",
  "wallet_address": "0x...",
  "ip_address": "1.2.3.4",
  "usage_fee": "0.0001",
  "accessed_at": "2024-01-01T00:00:00Z",
  "success": true
}
```

## API Integration

The Go backend connects to Fabric using the Fabric SDK:

```go
import (
    "github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

// Connect to network
gw, err := gateway.Connect(
    gateway.WithConfig(config.FromFile(configPath)),
    gateway.WithIdentity(wallet, "appUser"),
)

// Get contract
network, _ := gw.GetNetwork("passchainchannel")
contract := network.GetContract("credentials")

// Invoke chaincode
result, err := contract.SubmitTransaction("CreateCredentialShard", ...)
```

## Monitoring

### View Logs

```bash
# All containers
docker-compose logs -f

# Specific peer
docker logs peer0.org1.example.com -f

# Orderer
docker logs orderer.example.com -f
```

### Query Chaincode Logs

```bash
docker logs dev-peer0.org1.example.com-credentials_1.0
```

## Security

* **Private Channel**: Only authorized organizations can join
* **Identity Management**: Fabric CA for certificate issuance
* **Access Control**: Chaincode-level ACLs
* **Encryption**: Data encrypted before storing on chain
* **Audit Trail**: Immutable transaction history

## Backup & Recovery

### Backup Ledger

```bash
docker exec peer0.org1.example.com peer node pause
docker cp peer0.org1.example.com:/var/hyperledger/production ./ledger-backup
docker exec peer0.org1.example.com peer node resume
```

### Restore Ledger

```bash
docker cp ./ledger-backup peer0.org1.example.com:/var/hyperledger/production
docker restart peer0.org1.example.com
```

## Troubleshooting

### Network won't start

```bash
# Clean everything
./network.sh down
docker system prune -a --volumes
./network.sh up createChannel -ca
```

### Chaincode errors

```bash
# Check chaincode container logs
docker logs $(docker ps -q -f name=dev-peer0.org1)

# Reinstall chaincode
./network.sh deployCC -ccn credentials -ccp ./chaincode/credentials -ccl go
```

### Connection issues

* Check Docker network: `docker network inspect fabric_test`
* Verify all containers running: `docker ps`
* Check peer connectivity: `docker exec peer0.org1.example.com peer channel list`

## Performance Tuning

* **Block Size**: Adjust in `configtx.yaml`
* **Batch Timeout**: Configure orderer timeout
* **CouchDB Indexes**: Create indexes for common queries
* **Connection Pooling**: Use SDK connection pools

## Production Deployment

See `../infrastructure/k8s/blockchain/` for Kubernetes deployment.

## License

MIT
