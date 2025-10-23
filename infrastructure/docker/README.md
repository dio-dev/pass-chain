# Pass Chain - Docker Deployment Guide

## Quick Start

Start all services with Docker Compose:

```bash
cd infrastructure/docker
docker-compose up -d
```

## Services

The Docker Compose setup includes:

| Service | Port | Description |
|---------|------|-------------|
| Frontend | 3000 | Next.js web application |
| Backend | 8080 | Go API server |
| PostgreSQL | 5432 | Database |
| Redis | 6379 | Cache |
| Vault | 8200 | HashiCorp Vault |

## Access Services

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Vault UI**: http://localhost:8200
- **Health Check**: http://localhost:8080/health

## Development Credentials

### Vault
- **Token**: `dev-only-token`
- **Address**: http://localhost:8200

### PostgreSQL
- **Host**: localhost
- **Port**: 5432
- **User**: postgres
- **Password**: postgres
- **Database**: passchain

### Redis
- **Host**: localhost
- **Port**: 6379
- **Password**: (none)

## Commands

### Start Services
```bash
docker-compose up -d
```

### View Logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend
docker-compose logs -f frontend
```

### Stop Services
```bash
docker-compose down
```

### Rebuild and Restart
```bash
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

### Remove Everything (including volumes)
```bash
docker-compose down -v
```

## Building Images

### Backend
```bash
cd ../../backend
docker build -t passchain/backend:latest .
```

### Frontend
```bash
cd ../../frontend
docker build -t passchain/frontend:latest .
```

## Troubleshooting

### Port Already in Use

**Windows PowerShell:**
```powershell
# Find process using port
netstat -ano | findstr :8080

# Kill process
taskkill /PID <process_id> /F
```

### Container Won't Start

```bash
# Check logs
docker-compose logs <service-name>

# Restart specific service
docker-compose restart <service-name>
```

### Database Connection Error

```bash
# Check PostgreSQL is running
docker-compose ps postgres

# Check database logs
docker-compose logs postgres

# Restart database
docker-compose restart postgres
```

### Vault Not Initialized

```bash
# Check Vault status
docker exec -it passchain-vault vault status

# In dev mode, Vault auto-initializes
```

### Frontend Build Issues

```bash
# Rebuild frontend
cd ../../frontend
npm install
cd ../infrastructure/docker
docker-compose build frontend --no-cache
docker-compose up -d frontend
```

### Backend Build Issues

```bash
# Rebuild backend
cd ../../backend
go mod download
cd ../infrastructure/docker
docker-compose build backend --no-cache
docker-compose up -d backend
```

## Health Checks

### Backend Health
```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "healthy",
  "service": "pass-chain-api"
}
```

### Database Connection
```bash
docker exec -it passchain-postgres psql -U postgres -d passchain -c "SELECT 1;"
```

### Redis Connection
```bash
docker exec -it passchain-redis redis-cli ping
```

Expected: `PONG`

### Vault Status
```bash
curl http://localhost:8200/v1/sys/health
```

## Networking

All services are connected via the `passchain-network` bridge network.

Services can communicate using their service names:
- `backend` → `postgres:5432`
- `backend` → `redis:6379`
- `backend` → `vault:8200`
- `frontend` → `backend:8080`

## Volumes

Persistent data is stored in Docker volumes:

- `postgres_data` - Database data
- `redis_data` - Redis persistence

### Backup Volumes

```bash
# Backup PostgreSQL
docker exec passchain-postgres pg_dump -U postgres passchain > backup.sql

# Restore PostgreSQL
docker exec -i passchain-postgres psql -U postgres passchain < backup.sql
```

## Environment Variables

### Backend Environment

All backend environment variables are configured in `docker-compose.yml`:

- `ENVIRONMENT=development`
- `PORT=8080`
- `DB_HOST=postgres`
- `VAULT_ADDR=http://vault:8200`
- etc.

### Frontend Environment

- `NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1`
- `NEXT_PUBLIC_WALLET_CONNECT_PROJECT_ID=demo-project-id`

To use a real WalletConnect project:
1. Get Project ID from https://cloud.walletconnect.com
2. Update in `docker-compose.yml`
3. Restart: `docker-compose up -d frontend`

## Production Considerations

⚠️ **This setup is for DEVELOPMENT only!**

For production:
- Use secrets management (not hardcoded passwords)
- Enable TLS/HTTPS
- Use production Vault (not dev mode)
- Set up proper backup strategy
- Use health checks in orchestrator
- Configure resource limits
- Set up monitoring and alerting

## Next Steps

1. ✅ Start services: `docker-compose up -d`
2. ✅ Check logs: `docker-compose logs -f`
3. ✅ Access frontend: http://localhost:3000
4. ✅ Test backend: http://localhost:8080/health
5. 🔧 Start developing!

## Support

If you encounter issues:
1. Check the logs: `docker-compose logs -f`
2. Restart services: `docker-compose restart`
3. Rebuild if needed: `docker-compose build --no-cache`
4. Check GitHub issues or documentation




