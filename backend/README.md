# Backend - Pass Chain API

Go-based backend service for Pass Chain password management system.

## Architecture

The backend follows a clean architecture pattern with clear separation of concerns:

```
backend/
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── api/             # HTTP handlers and routes
│   ├── config/          # Configuration management
│   ├── models/          # Data models
│   ├── services/        # Business logic
│   ├── vault/           # HashiCorp Vault integration
│   ├── blockchain/      # Hyperledger Fabric integration
│   ├── crypto/          # Cryptographic operations (Shamir's Secret Sharing)
│   ├── middleware/      # HTTP middleware
│   └── database/        # Database connection and queries
├── pkg/
│   └── logger/          # Logging utilities
└── tests/               # Integration tests
```

## Features

- **Wallet Authentication**: Verify Ethereum wallet signatures
- **Split-Key Storage**: Shard credentials between Vault and blockchain
- **Payment Processing**: Track storage and usage payments
- **Audit Logging**: Immutable access logs on blockchain
- **Rate Limiting**: Prevent abuse
- **API Documentation**: Swagger/OpenAPI specs

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 15+
- Redis 7+
- HashiCorp Vault
- Hyperledger Fabric network
- Docker & Docker Compose

## Getting Started

### 1. Install Dependencies

```bash
cd backend
go mod download
```

### 2. Configure Environment

Create a `.env` file:

```bash
# Server
ENVIRONMENT=development
PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=passchain
DB_SSL_MODE=disable

# Vault
VAULT_ADDR=http://localhost:8200
VAULT_TOKEN=your-vault-token
VAULT_MOUNT_PATH=secret

# Blockchain
FABRIC_CONFIG_PATH=./config/network.yaml
FABRIC_CHANNEL_ID=passchainchannel
FABRIC_CHAINCODE_NAME=credentials
FABRIC_ORG_NAME=Org1
FABRIC_USER_NAME=Admin

# Redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_SECRET=your-secret-key
JWT_TTL=3600
```

### 3. Start Dependencies

```bash
docker-compose up -d postgres redis vault
```

### 4. Run Database Migrations

```bash
go run cmd/server/main.go
```

Migrations run automatically on startup.

### 5. Run the Server

```bash
go run cmd/server/main.go
```

The API will be available at http://localhost:8080

## API Endpoints

### Public Endpoints

- `GET /health` - Health check
- `GET /api/v1/info` - API information

### Protected Endpoints (Require Wallet Authentication)

#### Credentials
- `POST /api/v1/credentials` - Create new credential
- `GET /api/v1/credentials` - List all credentials
- `GET /api/v1/credentials/:id` - Get credential details
- `DELETE /api/v1/credentials/:id` - Delete credential
- `POST /api/v1/credentials/:id/reveal` - Reveal credential (requires payment)

#### Access Logs
- `GET /api/v1/access-logs` - Get all access logs
- `GET /api/v1/access-logs/:credentialId` - Get logs for specific credential

#### User
- `GET /api/v1/user/profile` - Get user profile
- `PUT /api/v1/user/profile` - Update user profile

#### Payments
- `GET /api/v1/payments` - List payments
- `GET /api/v1/payments/:id` - Get payment details

## Development

### Run Tests

```bash
go test ./...
```

### Run with Hot Reload

```bash
# Install air
go install github.com/cosmtrek/air@latest

# Run
air
```

### Linting

```bash
golangci-lint run
```

### Code Coverage

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Security

### Wallet Authentication

Users authenticate by signing a message with their Ethereum wallet. The signature is verified on the server.

### Credential Storage Flow

1. Client encrypts credential with wallet-derived key
2. Server receives encrypted data
3. Data is split using Shamir's Secret Sharing (threshold scheme)
4. Shard A → Stored in Vault
5. Shard B → Stored on Hyperledger Fabric
6. Only combination of both shards can reconstruct the data

### Credential Retrieval Flow

1. Client requests credential with wallet signature
2. Server verifies payment
3. Retrieves Shard A from Vault
4. Retrieves Shard B from blockchain
5. Reconstructs encrypted data
6. Returns to client
7. Client decrypts locally
8. Access logged on blockchain

## Deployment

### Build

```bash
go build -o bin/server cmd/server/main.go
```

### Docker

```bash
docker build -t pass-chain-api .
docker run -p 8080:8080 --env-file .env pass-chain-api
```

### Kubernetes

See `../infrastructure/k8s/backend/` for Kubernetes manifests.

## Monitoring

- **Metrics**: Prometheus endpoint at `/metrics`
- **Health**: Health check at `/health`
- **Logs**: JSON structured logs with Zap

## Troubleshooting

### Connection Issues

- Check Vault is running and unsealed
- Verify Fabric network is up
- Ensure PostgreSQL and Redis are accessible

### Performance

- Enable query caching in Redis
- Optimize Fabric chaincode queries
- Use connection pooling for database

## Contributing

1. Create a feature branch
2. Write tests
3. Ensure code passes linting
4. Submit a pull request

## License

MIT




