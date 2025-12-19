# Quick Start Guide - Secrets Manager MVP

## Prerequisites

- Docker Desktop for Mac (M1 native support)
- 8GB+ RAM available
- Terminal/Command line access

## Setup (5 minutes)

### 1. Generate Encryption Key

```bash
openssl rand -hex 32
```

Copy the output (64 hex characters).

### 2. Create Backend Environment File

```bash
cp backend/.env.example backend/.env
```

Edit `backend/.env` and replace:
- `ENCRYPTION_KEY` with the generated key from step 1
- `JWT_SECRET` with a random string (e.g., `your-super-secret-key-change-this`)

Example:
```
ENCRYPTION_KEY=a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6
JWT_SECRET=my-secret-jwt-key-12345
```

### 3. Start Docker Compose

```bash
cd dockerfiles
docker-compose up -d mongodb
```

Wait for MongoDB to be ready.

### 4. Setup Ollama (Optional - for AI Chat)

```bash
./dockerfiles/setup-ollama.sh
```

This will:
- Start Ollama service
- Pull the Mistral model (~4GB)
- Verify setup with test queries

Skip this step if you don't need the AI chatbot feature yet.

### 5. Access the Application

- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:8080/api/v1
- **MongoDB**: localhost:27017
- **Ollama**: http://localhost:11434 (if enabled)

## First Steps

### 1. Register Account

1. Go to http://localhost:5173
2. Click "Register"
3. Enter email and password
   - Password must have: 8+ chars, uppercase, number, special char
4. Click "Create account"

### 2. Create Your First Secret

1. Click "View All" or navigate to Secrets
2. Click "+ New Secret"
3. Fill in:
   - **Name**: e.g., "GitHub Token"
   - **Type**: Select from dropdown (password, token, url, api_key)
   - **Value**: Your actual secret
   - **Category**: e.g., "development"
   - **Tags**: e.g., "github", "api"
4. Click "Create"

### 3. View Secret

1. Click on the secret card
2. Click "Show" to reveal the value
3. Click "Copy" to copy to clipboard
4. Click "Delete" to remove

## API Testing

### Register User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123!"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123!"
  }'
```

Response includes `token` - use this for authenticated requests.

### Create Secret

```bash
curl -X POST http://localhost:8080/api/v1/secrets \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "name": "GitHub Token",
    "type": "token",
    "value": "ghp_xxxxxxxxxxxx",
    "category": "development",
    "tags": ["github", "api"]
  }'
```

### List Secrets

```bash
curl -X GET http://localhost:8080/api/v1/secrets \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### Get Secret (Decrypted)

```bash
curl -X GET http://localhost:8080/api/v1/secrets/SECRET_ID \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### Search Secrets

```bash
curl -X GET "http://localhost:8080/api/v1/secrets/search?q=github" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## Troubleshooting

### Port Already in Use

If ports 5173, 8080, 27017, or 11434 are already in use:

```bash
# Find process using port
lsof -i :8080

# Kill process
kill -9 <PID>
```

Or modify `docker-compose.yaml` to use different ports.

### MongoDB Connection Error

```bash
# Check MongoDB logs
docker logs secrets_manager_mongodb

# Restart MongoDB
docker-compose restart mongodb
```

### Frontend Not Loading

```bash
# Check frontend logs
docker logs secrets_manager_frontend

# Rebuild frontend
docker-compose build frontend
docker-compose up frontend
```

### Encryption Key Error

Make sure `ENCRYPTION_KEY` in `.env` is exactly 64 hex characters (32 bytes).

```bash
# Verify key length
echo -n "YOUR_KEY" | wc -c  # Should output 64
```

## Stopping Services

```bash
cd dockerfiles
docker-compose down
```

To also remove volumes (delete all data):

```bash
docker-compose down -v
```

## Development

### Local Backend Development

```bash
cd backend
go mod download
go run main.go
```

### Local Frontend Development

```bash
cd frontend
npm install
npm run dev
```

## Next Steps

1. **Phase 2**: Chatbot integration with Ollama
2. **Phase 3**: Advanced features (2FA, audit logging, etc.)
3. **Production**: Deploy with proper key management

## Documentation

- [Full Plan](./plan.md)
- [Encryption Strategy](./docs/specs/encryption-strategy.md)
- [Phase 1 Tasks](./docs/tasks/phase-1-tasks.md)
- [Phase 2 Tasks](./docs/tasks/phase-2-tasks.md)
- [Ollama Setup Guide](./docs/OLLAMA_SETUP.md)
- [Implementation Summary](./docs/tasks/phase-1-implementation-summary.md)

## Support

For issues or questions, check:
1. Docker logs: `docker-compose logs -f`
2. Browser console: F12 → Console tab
3. Backend logs: `docker logs secrets_manager_backend`

---

**Happy secret managing! 🔐**
