# Detailed Application Plan: AI-Powered Secrets Manager

## Project Overview
A web application that securely stores credentials (usernames, passwords, tokens, URLs) with an integrated AI chatbot for natural language search and retrieval of stored secrets.

---

## 1. Architecture Overview

### Tech Stack
- **Backend**: Go with Gin framework
- **Frontend**: Vue.js with Tailwind CSS
- **Database**: MongoDB (NoSQL for flexible schema)
- **AI/Chatbot**: Self-hosted Ollama/LLaMA (local LLM inference)
- **Authentication**: JWT-based auth with optional 2FA
- **Deployment**: Docker containerization (M1 Mac compatible)
- **Development Environment**: macOS M1 with Docker Desktop

### System Components
1. Backend API (Go/Gin)
2. Frontend UI (Vue.js)
3. Database Layer
4. AI Chatbot Service
5. Authentication Service
6. Encryption Service

---

## 2. Core Features

### 2.1 Secrets Management
- **Create**: Add new secrets (username/password, tokens, URLs, API keys)
- **Read**: Retrieve stored secrets
- **Update**: Modify existing secrets
- **Delete**: Remove secrets with confirmation
- **Categorize**: Organize secrets by type/tags (e.g., "social", "work", "banking")
- **Search**: Full-text search by name, category, or metadata

### 2.2 AI Chatbot Integration
- **Natural Language Search**: "Show me my GitHub token" → retrieves GitHub secret
- **Smart Filtering**: "Find all my work passwords" → filters by category
- **Context Understanding**: Understand user intent and return relevant secrets
- **Conversation History**: Maintain chat context within a session
- **Security Prompts**: Confirm sensitive operations before execution

### 2.3 Security Features
- **Encryption**: AES-256 encryption for stored secrets
- **Authentication**: User login with JWT tokens
- **Authorization**: Role-based access control (RBAC)
- **Audit Logging**: Track all access and modifications
- **Session Management**: Auto-logout after inactivity
- **Password Hashing**: bcrypt for user passwords
- **HTTPS Only**: Enforce secure communication

### 2.4 User Management
- **Registration**: Create new user accounts
- **Login/Logout**: Secure authentication
- **Profile Management**: Update user settings
- **Password Reset**: Secure password recovery flow

---

## 3. Backend (Go/Gin)

### API Endpoints

#### Authentication
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `POST /api/auth/logout` - User logout
- `POST /api/auth/refresh` - Refresh JWT token

#### Secrets Management
- `GET /api/secrets` - List all user secrets
- `GET /api/secrets/:id` - Get specific secret
- `POST /api/secrets` - Create new secret
- `PUT /api/secrets/:id` - Update secret
- `DELETE /api/secrets/:id` - Delete secret
- `GET /api/secrets/search?q=query` - Search secrets

#### Chatbot
- `POST /api/chat` - Send message to chatbot
- `GET /api/chat/history` - Get conversation history
- `DELETE /api/chat/history` - Clear chat history

#### User
- `GET /api/user/profile` - Get user profile
- `PUT /api/user/profile` - Update profile
- `POST /api/user/change-password` - Change password

### Database Schema (MongoDB Collections)

#### Users Collection
```json
{
  "_id": ObjectId,
  "email": "user@example.com",
  "password_hash": "bcrypt_hash",
  "created_at": ISODate,
  "updated_at": ISODate,
  "last_login": ISODate
}
```

#### Secrets Collection
```json
{
  "_id": ObjectId,
  "user_id": ObjectId,
  "name": "GitHub Token",
  "type": "token",
  "encrypted_value": "encrypted_string",
  "category": "development",
  "tags": ["github", "api"],
  "metadata": {
    "url": "https://github.com",
    "username": "optional_username"
  },
  "created_at": ISODate,
  "updated_at": ISODate
}
```

#### Chat History Collection
```json
{
  "_id": ObjectId,
  "user_id": ObjectId,
  "message": "Show me my GitHub token",
  "response": "Found your GitHub token...",
  "timestamp": ISODate,
  "session_id": "session_uuid"
}
```

#### Audit Log Collection
```json
{
  "_id": ObjectId,
  "user_id": ObjectId,
  "action": "create|read|update|delete",
  "resource_type": "secret",
  "resource_id": ObjectId,
  "timestamp": ISODate,
  "ip_address": "user_ip"
}
```

### MongoDB Indexes
- Users: unique index on email
- Secrets: compound index on (user_id, category), text index on name/tags
- Chat History: index on (user_id, timestamp)
- Audit Log: index on (user_id, timestamp)

### Middleware
- Authentication middleware (JWT validation)
- Authorization middleware (permission checks)
- Logging middleware
- Error handling middleware
- CORS middleware

---

## 4. Frontend (Vue.js)

### Pages/Views

#### Authentication Pages
- Login page
- Registration page
- Password reset page

#### Main Dashboard
- Sidebar navigation
- Secrets list/grid view
- Quick stats (total secrets, recent access)

#### Secrets Management
- Create/Edit secret modal
- Secret detail view
- Bulk operations (delete, export)
- Category/tag management

#### Chatbot Interface
- Chat window (sidebar or modal)
- Message input
- Conversation history
- Clear history button

#### User Settings
- Profile management
- Security settings
- Session management
- Logout

### Components
- Header/Navigation
- Sidebar
- Secret Card
- Secret Form
- Chat Widget
- Modal dialogs
- Toast notifications
- Loading spinners

### State Management
- Vue Pinia or Vuex for global state
- User authentication state
- Secrets list state
- Chat history state

---

## 5. AI Chatbot Integration

### Capabilities
- Parse natural language queries
- Map user intent to secret retrieval
- Filter secrets by category, type, or metadata
- Provide contextual responses
- Handle follow-up questions

### Self-Hosted LLM Setup (Ollama/LLaMA)

#### Ollama Installation & Configuration
- **Installation**: Download Ollama from ollama.ai (M1 Mac native support)
- **Model Selection**: Use lightweight models for M1 Mac
  - `ollama pull mistral` (7B, ~4GB VRAM)
  - `ollama pull neural-chat` (7B, optimized for chat)
  - `ollama pull llama2` (7B, ~4GB VRAM)
- **API Endpoint**: Ollama runs on `http://localhost:11434`
- **Memory Requirements**: M1 Mac with 8GB+ RAM recommended

#### Backend Integration
- Call Ollama API from Go backend via HTTP
- Endpoint: `POST http://localhost:11434/api/generate`
- Response streaming for real-time chat
- Error handling for model unavailability

### Prompt Engineering
- System prompt defining chatbot behavior and security constraints
- Context injection with user's secret metadata (names, categories, tags)
- Safety guardrails to prevent unauthorized access
- Confirmation prompts for sensitive operations
- Few-shot examples for intent recognition

### Example Ollama Integration
```go
// Backend calls Ollama for chat inference
POST http://localhost:11434/api/generate
{
  "model": "mistral",
  "prompt": "User query with context",
  "stream": true
}
```

---

## 6. Security Considerations

### Data Protection
- Encrypt secrets at rest (AES-256)
- Encrypt secrets in transit (HTTPS/TLS)
- Never log sensitive data
- Secure key management (environment variables, vaults)

### Access Control
- JWT tokens with expiration
- Refresh token rotation
- Rate limiting on API endpoints
- CORS configuration

### Compliance
- GDPR compliance (data deletion, export)
- SOC 2 considerations
- Regular security audits
- Dependency vulnerability scanning

---

## 7. Development Phases

### Phase 1: MVP (Weeks 1-2)
- Set up Docker Compose with MongoDB and Ollama
- User authentication (login/register) with JWT
- Basic secrets CRUD operations
- MongoDB integration in backend
- Basic UI with Tailwind CSS
- Local testing on M1 Mac

### Phase 2: Chatbot Integration (Weeks 3-4)
- Integrate Ollama API in Go backend
- Implement chatbot endpoint (`POST /api/chat`)
- Natural language query parsing with LLM
- Chat UI component in Vue.js
- Conversation history storage in MongoDB
- Test with local Ollama model

### Phase 3: Enhancement (Weeks 5-6)
- Advanced filtering and categorization
- Audit logging to MongoDB
- User settings and profile management
- Export/import functionality
- Optimize Ollama model selection for M1 Mac
- Performance tuning for local development

### Phase 4: Polish & Deployment (Weeks 7-8)
- Security hardening (encryption, CORS, rate limiting)
- Performance optimization
- Docker Compose finalization
- Comprehensive testing (unit, integration, E2E)
- Documentation for M1 Mac setup
- Bug fixes and refinements

---

## 8. Testing Strategy

### Backend Testing
- Unit tests for business logic
- Integration tests for API endpoints
- Security tests (SQL injection, XSS prevention)
- Load testing

### Frontend Testing
- Component unit tests (Vue Test Utils)
- Integration tests
- E2E tests (Cypress/Playwright)
- Accessibility testing

---

## 9. Deployment

### Docker Setup (M1 Mac Compatible)

#### Docker Compose Services
- **Backend**: Go/Gin API (arm64 compatible)
- **Frontend**: Vue.js dev server or nginx
- **MongoDB**: Official MongoDB image (arm64 support)
- **Ollama**: Self-hosted LLM service (arm64 native)

#### M1 Mac Specific Configuration
- Use `arm64` base images for all services
- Docker Desktop for Mac with native M1 support
- Resource allocation: 4GB+ RAM for Ollama + MongoDB
- Volume mounts for persistent data

#### Docker Compose File Structure
```yaml
services:
  backend:
    image: golang:1.21-alpine (arm64)
    ports: 8080
  frontend:
    image: node:20-alpine (arm64)
    ports: 5173
  mongodb:
    image: mongo:7.0 (arm64)
    ports: 27017
  ollama:
    image: ollama/ollama (arm64)
    ports: 11434
    volumes: model cache
```

#### Environment Configuration
- `.env` files for each service
- MongoDB connection string
- Ollama API endpoint
- JWT secret keys
- Encryption keys

### Local Development on M1 Mac
- Clone repository
- Run `docker-compose up` to start all services
- Backend API: `http://localhost:8080`
- Frontend: `http://localhost:5173`
- MongoDB: `localhost:27017`
- Ollama: `http://localhost:11434`
- First run: Ollama downloads model (~4GB, one-time)

### Production Considerations
- MongoDB Atlas or self-hosted MongoDB with backups
- SSL/TLS certificates (Let's Encrypt)
- Monitoring and logging (ELK stack or similar)
- Error tracking (Sentry)
- Performance monitoring
- Ollama model caching and optimization
- Resource limits for containers

---

## 10. Future Enhancements

- Two-factor authentication (2FA)
- Biometric authentication
- Secret sharing with other users
- Mobile app (React Native/Flutter)
- Browser extension for auto-fill
- Advanced analytics dashboard
- Webhook integrations
- API key management for third-party apps
- Backup and recovery features
- Multi-language support

---

## 11. Success Metrics

- User registration and retention
- Chatbot query success rate
- API response time (<200ms)
- System uptime (>99.9%)
- User satisfaction score
- Security audit results
