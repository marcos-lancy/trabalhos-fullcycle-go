#!/bin/bash

# Script para testar o rate limiter
# Execute este script após iniciar o servidor

echo "Testing Rate Limiter..."
echo "======================="

# Teste 1: Health check
echo "1. Testing health check..."
curl -s http://localhost:8080/health | jq .
echo ""

# Teste 2: Rate limiting por IP
echo "2. Testing IP rate limiting (limit: 5 req/s)..."
for i in {1..7}; do
    echo "Request $i:"
    response=$(curl -s -w "HTTP Status: %{http_code}\n" http://localhost:8080/api/data)
    echo "$response"
    echo ""
    sleep 0.1
done

echo "Waiting 2 seconds before next test..."
sleep 2

# Teste 3: Rate limiting por token
echo "3. Testing token rate limiting (limit: 10 req/s)..."
for i in {1..12}; do
    echo "Request $i:"
    response=$(curl -s -w "HTTP Status: %{http_code}\n" -H "API_KEY: test-token" http://localhost:8080/api/data)
    echo "$response"
    echo ""
    sleep 0.1
done

echo "Waiting 2 seconds before next test..."
sleep 2

# Teste 4: Verificar status do rate limit
echo "4. Checking rate limit status..."
echo "IP status:"
curl -s http://localhost:8080/api/rate-limit/status | jq .
echo ""
echo "Token status:"
curl -s -H "API_KEY: test-token" http://localhost:8080/api/rate-limit/status | jq .
echo ""

# Teste 5: Reset rate limit
echo "5. Resetting rate limit..."
curl -s -X POST http://localhost:8080/api/rate-limit/reset | jq .
echo ""

# Teste 6: Verificar se reset funcionou
echo "6. Testing after reset..."
curl -s http://localhost:8080/api/data | jq .
echo ""

echo "Rate limiter test completed!"
