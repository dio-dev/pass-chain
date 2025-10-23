# 🧪 Pass Chain Backend Tests

## Test Coverage

### Credentials API Tests (`credentials_test.go`)

✅ **TestCreateCredential_Success**
- Verifies successful credential creation
- Mocks Vault storage
- Checks database insertion
- Validates response format

✅ **TestCreateCredential_MissingFields**
- Tests validation for missing required fields
- Expects 400 Bad Request

✅ **TestCreateCredential_MissingAuthentication**
- Tests authentication requirements
- Expects 401 Unauthorized when wallet/signature missing

✅ **TestGetCredentials_Success**
- Tests listing credentials for a wallet
- Verifies correct filtering by wallet address

✅ **TestGetCredentialByID_Success**
- Tests retrieving a single credential
- Mocks Vault read for Share1 retrieval
- Verifies Share2 from database

✅ **TestDeleteCredential_Success**
- Tests credential deletion
- Mocks Vault delete
- Verifies database soft delete

---

## Running Tests

### Run all tests:
```bash
cd backend
go test ./... -v
```

### Run with coverage:
```bash
go test ./... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Run specific test:
```bash
go test ./internal/api/handlers -run TestCreateCredential_Success -v
```

### Run tests in Docker (Kubernetes):
```bash
kubectl exec -it deployment/passchain-backend -n passchain -- go test ./... -v
```

---

## Test Structure

```
backend/
├── internal/
│   ├── api/
│   │   └── handlers/
│   │       ├── credentials.go
│   │       └── credentials_test.go  ✅ API tests
│   └── database/
│       ├── database.go
│       └── test_db.go               ✅ Test DB helper
```

---

## Mock Services

### MockVaultService
Mocks HashiCorp Vault operations:
- `WriteSecret()` - Store secrets
- `ReadSecret()` - Retrieve secrets
- `DeleteSecret()` - Delete secrets
- `ListSecrets()` - List secret paths

Uses `github.com/stretchr/testify/mock` for assertions.

---

## Test Database

Uses **SQLite in-memory** for fast, isolated tests:
- No external dependencies
- Fresh database per test
- Automatic cleanup
- Same GORM models as PostgreSQL

---

## CI/CD Integration

Add to your CI pipeline:

```yaml
# .github/workflows/test.yml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run tests
        run: |
          cd backend
          go test ./... -v -cover
```

---

## Next Steps

### Additional Tests to Add:

1. **Integration Tests**
   - Full API flow with real Vault (test mode)
   - End-to-end encryption/decryption

2. **Middleware Tests**
   - CORS validation
   - Wallet signature verification

3. **Crypto Tests**
   - Shamir Secret Sharing
   - XChaCha20-Poly1305 encryption

4. **Load Tests**
   - Concurrent credential creation
   - Rate limiting

5. **Security Tests**
   - SQL injection prevention
   - XSS prevention
   - Authorization bypass attempts

---

**Run tests now!** 🧪✅

