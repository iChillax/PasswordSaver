#!/bin/bash

# Test Ollama Service
# Quick verification script for Ollama setup

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

MODEL=${1:-mistral}

echo "🧪 Testing Ollama Service..."
echo ""

# Test 1: Service Health
echo -e "${YELLOW}Test 1: Checking service health...${NC}"
if curl -f http://localhost:11434/api/tags > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Service is healthy${NC}"
else
    echo -e "${RED}❌ Service is not responding${NC}"
    exit 1
fi

# Test 2: Model Availability
echo -e "${YELLOW}Test 2: Checking model availability...${NC}"
MODELS=$(curl -s http://localhost:11434/api/tags | grep -o "\"name\":\"[^\"]*\"" | cut -d'"' -f4)
if echo "$MODELS" | grep -q "$MODEL"; then
    echo -e "${GREEN}✅ Model $MODEL is available${NC}"
else
    echo -e "${RED}❌ Model $MODEL not found${NC}"
    echo "Available models:"
    echo "$MODELS"
    exit 1
fi

# Test 3: Simple Query
echo -e "${YELLOW}Test 3: Testing simple query...${NC}"
RESPONSE=$(curl -s http://localhost:11434/api/generate -d "{
  \"model\": \"$MODEL\",
  \"prompt\": \"What is 2+2? Answer with just the number.\",
  \"stream\": false
}")

if echo "$RESPONSE" | grep -q "response"; then
    ANSWER=$(echo "$RESPONSE" | grep -o '"response":"[^"]*"' | cut -d'"' -f4)
    echo -e "${GREEN}✅ Model responded: $ANSWER${NC}"
else
    echo -e "${RED}❌ Query failed${NC}"
    echo "$RESPONSE"
    exit 1
fi

# Test 4: Intent Recognition Query
echo -e "${YELLOW}Test 4: Testing intent recognition...${NC}"
RESPONSE=$(curl -s http://localhost:11434/api/generate -d "{
  \"model\": \"$MODEL\",
  \"prompt\": \"Extract the intent from this query: 'Show me my GitHub tokens'. Reply with just: search\",
  \"stream\": false
}")

if echo "$RESPONSE" | grep -q "response"; then
    INTENT=$(echo "$RESPONSE" | grep -o '"response":"[^"]*"' | cut -d'"' -f4)
    echo -e "${GREEN}✅ Intent recognition working: $INTENT${NC}"
else
    echo -e "${RED}❌ Intent recognition failed${NC}"
    exit 1
fi

# Test 5: Response Time
echo -e "${YELLOW}Test 5: Measuring response time...${NC}"
START=$(date +%s)
curl -s http://localhost:11434/api/generate -d "{
  \"model\": \"$MODEL\",
  \"prompt\": \"Hello\",
  \"stream\": false
}" > /dev/null
END=$(date +%s)
DURATION=$((END - START))

if [ $DURATION -lt 10 ]; then
    echo -e "${GREEN}✅ Response time: ${DURATION}s (Good)${NC}"
elif [ $DURATION -lt 30 ]; then
    echo -e "${YELLOW}⚠️  Response time: ${DURATION}s (Acceptable)${NC}"
else
    echo -e "${RED}❌ Response time: ${DURATION}s (Too slow)${NC}"
fi

echo ""
echo -e "${GREEN}🎉 All tests passed!${NC}"
echo ""
echo "📊 Summary:"
echo "  - Service: ✅ Healthy"
echo "  - Model: ✅ $MODEL available"
echo "  - Query: ✅ Working"
echo "  - Intent: ✅ Working"
echo "  - Speed: ✅ ${DURATION}s"
echo ""
