# Ollama Quick Reference

## 🚀 Quick Start

```bash
# Setup (one-time)
./dockerfiles/setup-ollama.sh

# Test
./dockerfiles/test-ollama.sh

# Verify
curl http://localhost:11434/api/tags
```

## 📋 Common Commands

### Service Management
```bash
# Start
docker-compose up -d ollama

# Stop
docker-compose stop ollama

# Restart
docker-compose restart ollama

# Logs
docker logs -f secrets_manager_ollama

# Status
docker-compose ps ollama
```

### Model Management
```bash
# List models
docker exec secrets_manager_ollama ollama list

# Pull model
docker exec secrets_manager_ollama ollama pull mistral

# Remove model
docker exec secrets_manager_ollama ollama rm mistral

# Show model info
docker exec secrets_manager_ollama ollama show mistral
```

### Testing
```bash
# Health check
curl http://localhost:11434/api/tags

# Simple query
curl http://localhost:11434/api/generate -d '{
  "model": "mistral",
  "prompt": "Hello",
  "stream": false
}'

# Streaming query
curl http://localhost:11434/api/generate -d '{
  "model": "mistral",
  "prompt": "Count to 5",
  "stream": true
}'
```

## 🎯 Recommended Models

| Model | Size | Best For |
|-------|------|----------|
| mistral | 7B | General purpose (recommended) |
| llama2 | 7B | Reliable, well-tested |
| neural-chat | 7B | Fast chat responses |
| codellama | 7B | Code understanding |

## 🔧 Configuration

### Environment Variables (backend/.env)
```env
OLLAMA_API_URL=http://localhost:11434
OLLAMA_MODEL=mistral
OLLAMA_TIMEOUT=30
```

### Model Parameters
```json
{
  "temperature": 0.7,
  "top_p": 0.9,
  "num_predict": 256
}
```

## 🐛 Troubleshooting

### Service won't start
```bash
docker logs secrets_manager_ollama
docker-compose restart ollama
```

### Model not found
```bash
docker exec secrets_manager_ollama ollama list
docker exec secrets_manager_ollama ollama pull mistral
```

### Slow responses
- Check: `docker stats secrets_manager_ollama`
- Reduce `num_predict` parameter
- Use smaller model

### Port conflict
```bash
lsof -i :11434
kill -9 <PID>
```

## 📊 Performance Targets

- Simple queries: < 2s
- Complex queries: < 5s
- Cold start: < 10s
- Memory: ~4GB
- Disk: ~4GB per model

## 🔗 Endpoints

- Health: `GET http://localhost:11434/api/tags`
- Generate: `POST http://localhost:11434/api/generate`
- Models: `GET http://localhost:11434/api/tags`

## 📚 Documentation

- Full guide: [docs/OLLAMA_SETUP.md](../docs/OLLAMA_SETUP.md)
- Phase 2 tasks: [docs/tasks/phase-2-tasks.md](../docs/tasks/phase-2-tasks.md)
- Implementation: [docs/tasks/task-2.1-implementation.md](../docs/tasks/task-2.1-implementation.md)

---

**Need help?** Check the full documentation or run `./dockerfiles/test-ollama.sh`
