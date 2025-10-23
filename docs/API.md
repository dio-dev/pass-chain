# API Documentation

## Base URL

```
http://localhost:8080/api/v1
```

Production: `https://api.passchain.io/api/v1`

## Authentication

All protected endpoints require a JWT token obtained through wallet authentication.

### Wallet Authentication Flow

1. Request challenge nonce
2. Sign nonce with wallet
3. Submit signature to get JWT

```
POST /api/v1/auth/challenge
GET /api/v1/auth/verify
```

### Using JWT Token

Include in `Authorization` header:

```
Authorization: Bearer <your-jwt-token>
```

## Endpoints

### Health & Info

#### GET /health

Health check endpoint.

**Response:**
```json
{
  "status": "healthy",
  "service": "pass-chain-api"
}
```

#### GET /api/v1/info

API information.

**Response:**
```json
{
  "name": "Pass Chain API",
  "version": "1.0.0",
  "environment": "development"
}
```

### Credentials

#### POST /api/v1/credentials

Create a new credential.

**Headers:**
- `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "credential_name": "Gmail Password",
  "description": "My personal Gmail account",
  "encrypted_data": "base64_encrypted_data..."
}
```

**Response: 201 Created**
```json
{
  "id": "uuid",
  "wallet_address": "0x...",
  "credential_name": "Gmail Password",
  "description": "My personal Gmail account",
  "vault_path": "secret/data/credentials/uuid",
  "blockchain_tx_id": "fabric_tx_id",
  "storage_fee": "0.001",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### GET /api/v1/credentials

List all credentials for authenticated wallet.

**Headers:**
- `Authorization: Bearer <token>`

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 20)

**Response: 200 OK**
```json
{
  "credentials": [
    {
      "id": "uuid",
      "wallet_address": "0x...",
      "credential_name": "Gmail Password",
      "description": "My personal Gmail account",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 10,
  "page": 1,
  "limit": 20
}
```

#### GET /api/v1/credentials/:id

Get credential details (without revealing password).

**Headers:**
- `Authorization: Bearer <token>`

**Response: 200 OK**
```json
{
  "id": "uuid",
  "wallet_address": "0x...",
  "credential_name": "Gmail Password",
  "description": "My personal Gmail account",
  "vault_path": "secret/data/credentials/uuid",
  "blockchain_tx_id": "fabric_tx_id",
  "storage_fee": "0.001",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### POST /api/v1/credentials/:id/reveal

Reveal credential (requires payment).

**Headers:**
- `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "payment_tx_hash": "0x...",
  "payment_amount": "0.0001"
}
```

**Response: 200 OK**
```json
{
  "encrypted_data": "base64_encrypted_data...",
  "blockchain_tx_id": "fabric_access_log_tx_id"
}
```

**Note**: Client must decrypt `encrypted_data` locally using wallet key.

#### DELETE /api/v1/credentials/:id

Delete a credential.

**Headers:**
- `Authorization: Bearer <token>`

**Response: 204 No Content**

### Access Logs

#### GET /api/v1/access-logs

Get all access logs for authenticated wallet.

**Headers:**
- `Authorization: Bearer <token>`

**Query Parameters:**
- `page` (optional): Page number
- `limit` (optional): Items per page

**Response: 200 OK**
```json
{
  "logs": [
    {
      "id": "uuid",
      "credential_id": "uuid",
      "wallet_address": "0x...",
      "ip_address": "1.2.3.4",
      "usage_fee": "0.0001",
      "blockchain_tx_id": "fabric_tx_id",
      "success": true,
      "accessed_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 50,
  "page": 1,
  "limit": 20
}
```

#### GET /api/v1/access-logs/:credentialId

Get access logs for a specific credential.

**Headers:**
- `Authorization: Bearer <token>`

**Response: 200 OK**
```json
{
  "logs": [
    {
      "id": "uuid",
      "credential_id": "uuid",
      "wallet_address": "0x...",
      "ip_address": "1.2.3.4",
      "usage_fee": "0.0001",
      "blockchain_tx_id": "fabric_tx_id",
      "success": true,
      "accessed_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### User Profile

#### GET /api/v1/user/profile

Get user profile.

**Headers:**
- `Authorization: Bearer <token>`

**Response: 200 OK**
```json
{
  "id": "uuid",
  "wallet_address": "0x...",
  "email": "user@example.com",
  "settings": {
    "notifications": true,
    "theme": "dark"
  },
  "created_at": "2024-01-01T00:00:00Z"
}
```

#### PUT /api/v1/user/profile

Update user profile.

**Headers:**
- `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "email": "newemail@example.com",
  "settings": {
    "notifications": false,
    "theme": "light"
  }
}
```

**Response: 200 OK**
```json
{
  "id": "uuid",
  "wallet_address": "0x...",
  "email": "newemail@example.com",
  "settings": {
    "notifications": false,
    "theme": "light"
  },
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Payments

#### GET /api/v1/payments

List payment transactions.

**Headers:**
- `Authorization: Bearer <token>`

**Query Parameters:**
- `page` (optional): Page number
- `limit` (optional): Items per page
- `type` (optional): Filter by type (`storage` or `usage`)

**Response: 200 OK**
```json
{
  "payments": [
    {
      "id": "uuid",
      "wallet_address": "0x...",
      "credential_id": "uuid",
      "amount": "0.001",
      "currency": "ETH",
      "payment_type": "storage",
      "tx_hash": "0x...",
      "status": "confirmed",
      "blockchain_tx_id": "fabric_tx_id",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 25,
  "page": 1,
  "limit": 20
}
```

#### GET /api/v1/payments/:id

Get payment details.

**Headers:**
- `Authorization: Bearer <token>`

**Response: 200 OK**
```json
{
  "id": "uuid",
  "wallet_address": "0x...",
  "credential_id": "uuid",
  "amount": "0.001",
  "currency": "ETH",
  "payment_type": "storage",
  "tx_hash": "0x...",
  "status": "confirmed",
  "blockchain_tx_id": "fabric_tx_id",
  "created_at": "2024-01-01T00:00:00Z"
}
```

## Error Responses

### 400 Bad Request

```json
{
  "error": "Invalid request",
  "details": "credential_name is required"
}
```

### 401 Unauthorized

```json
{
  "error": "Unauthorized",
  "details": "Invalid or missing authentication token"
}
```

### 403 Forbidden

```json
{
  "error": "Forbidden",
  "details": "You don't have permission to access this resource"
}
```

### 404 Not Found

```json
{
  "error": "Not found",
  "details": "Credential not found"
}
```

### 500 Internal Server Error

```json
{
  "error": "Internal server error",
  "details": "An unexpected error occurred"
}
```

## Rate Limiting

- **Anonymous**: 10 requests/minute
- **Authenticated**: 100 requests/minute
- **Premium**: 1000 requests/minute

Rate limit headers:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1640000000
```

## Webhook Events

Subscribe to webhook events for:
- `credential.created`
- `credential.accessed`
- `credential.deleted`
- `payment.confirmed`

Configure webhooks in user settings.

## SDKs

Coming soon:
- JavaScript/TypeScript SDK
- Python SDK
- Go SDK

## Support

- API Issues: api@passchain.io
- Documentation: docs@passchain.io
- Discord: [Join our community]




