#!/bin/bash

# Setup Ollama Service for Secrets Manager
# This script pulls the recommended model and verifies the setup

set -e

echo "🚀 Setting up Ollama service..."

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Default model (can be overridden)
MODEL=${1:-mistral}

echo -e "${YELLOW}📦 Starting Ollama service...${NC}"
docker-compose up -d ollama

echo -e "${YELLOW}⏳ Waiting for Ollama to be ready...${NC}"
sleep 10

# Wait for health check
MAX_RETRIES=30
RETRY_COUNT=0
while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if docker exec secrets_manager_ollama curl -f http://localhost:11434/api/tags > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Ollama service is ready${NC}"
        break
    fi
    RETRY_COUNT=$((RETRY_COUNT + 1))
    echo -e "${YELLOW}⏳ Waiting for Ollama... ($RETRY_COUNT/$MAX_RETRIES)${NC}"
    sleep 2
done

if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
    echo -e "${RED}❌ Ollama service failed to start${NC}"
    exit 1
fi

echo -e "${YELLOW}📥 Pulling model: $MODEL${NC}"
echo -e "${YELLOW}This may take several minutes depending on your internet connection...${NC}"
docker exec secrets_manager_ollama ollama pull $MODEL

echo -e "${GREEN}✅ Model $MODEL pulled successfully${NC}"

echo -e "${YELLOW}🧪 Testing model with a simple query...${NC}"
TEST_RESPONSE=$(docker exec secrets_manager_ollama curl -s http://localhost:11434/api/generate -d '{
  "model": "'$MODEL'",
  "prompt": "Say hello in one word",
  "stream": false
}')

if echo "$TEST_RESPONSE" | grep -q "response"; then
    echo -e "${GREEN}✅ Model is responding correctly${NC}"
    echo -e "${GREEN}Response: $(echo $TEST_RESPONSE | grep -o '"response":"[^"]*"' | cut -d'"' -f4)${NC}"
else
    echo -e "${RED}❌ Model test failed${NC}"
    echo "$TEST_RESPONSE"
    exit 1
fi

echo ""
echo -e "${GREEN}🎉 Ollama setup complete!${NC}"
echo ""
echo "📊 Service Information:"
echo "  - Ollama API: http://localhost:11434"
echo "  - Model: $MODEL"
echo "  - Health Check: http://localhost:11434/api/tags"
echo ""
echo "🔧 Useful Commands:"
echo "  - List models: docker exec secrets_manager_ollama ollama list"
echo "  - Pull another model: docker exec secrets_manager_ollama ollama pull <model-name>"
echo "  - Test query: curl http://localhost:11434/api/generate -d '{\"model\":\"$MODEL\",\"prompt\":\"test\"}'"
echo ""
echo "📚 Recommended Models:"
echo "  - mistral (default, 7B, balanced)"
echo "  - llama2 (7B, good general purpose)"
echo "  - neural-chat (7B, optimized for chat)"
echo "  - codellama (7B, code-focused)"
echo ""
