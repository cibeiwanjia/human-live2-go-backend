#!/bin/bash

BASE_URL="http://localhost:8881"
VERSION="v0"

echo "=========================================="
echo "Go Backend API Test Suite"
echo "=========================================="

# Test Agent API
echo ""
echo "--- Agent API Tests ---"

echo "1. GET /adh/agent/${VERSION}/engine"
curl -s "${BASE_URL}/adh/agent/${VERSION}/engine" | jq .

echo ""
echo "2. GET /adh/agent/${VERSION}/engine/default"
curl -s "${BASE_URL}/adh/agent/${VERSION}/engine/default" | jq .

echo ""
echo "3. GET /adh/agent/${VERSION}/engine/RepeaterAgent"
curl -s "${BASE_URL}/adh/agent/${VERSION}/engine/RepeaterAgent" | jq .

echo ""
echo "4. POST /adh/agent/${VERSION}/engine (SSE Stream)"
curl -N -s "${BASE_URL}/adh/agent/${VERSION}/engine" \
  -H "Content-Type: application/json" \
  -d '{"engine":"RepeaterAgent","data":"Hello, this is a test!"}'

# Test TTS API
echo ""
echo ""
echo "--- TTS API Tests ---"

echo "1. GET /adh/tts/${VERSION}/engine"
curl -s "${BASE_URL}/adh/tts/${VERSION}/engine" | jq .

echo ""
echo "2. GET /adh/tts/${VERSION}/engine/default"
curl -s "${BASE_URL}/adh/tts/${VERSION}/engine/default" | jq .

echo ""
echo "3. GET /adh/tts/${VERSION}/engine/EdgeTTS/voice"
curl -s "${BASE_URL}/adh/tts/${VERSION}/engine/EdgeTTS/voice" | jq .

# Test ASR API
echo ""
echo ""
echo "--- ASR API Tests ---"

echo "1. GET /adh/asr/${VERSION}/engine"
curl -s "${BASE_URL}/adh/asr/${VERSION}/engine" | jq .

echo ""
echo "2. GET /adh/asr/${VERSION}/engine/default"
curl -s "${BASE_URL}/adh/asr/${VERSION}/engine/default" | jq .

echo ""
echo "=========================================="
echo "Tests Complete"
echo "=========================================="