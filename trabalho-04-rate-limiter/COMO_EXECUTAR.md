# 🚀 Como Executar o Rate Limiter

## ⚡ Execução Rápida (2 comandos)

```bash
# 1. Clone e entre na pasta
git clone <repository-url>
cd rate-limiter

# 2. Execute com Docker
docker-compose up --build
```

**Pronto!** A aplicação estará rodando em `http://localhost:8080`

## 🧪 Testes Básicos

### **Teste 1: Health Check**
```bash
curl http://localhost:8080/health
```

### **Teste 2: Rate Limiting por IP**
```bash
# Faça 6 requisições rapidamente
for i in {1..6}; do
  curl http://localhost:8080/api/data
  echo "Request $i"
done
# Resultado: Requests 1-5 = HTTP 200, Request 6 = HTTP 429
```

### **Teste 3: Rate Limiting por Token**
```bash
# Faça 11 requisições com token
for i in {1..11}; do
  curl -H "API_KEY: test-token" http://localhost:8080/api/data
  echo "Request $i"
done
# Resultado: Requests 1-10 = HTTP 200, Request 11 = HTTP 429
```

### **Teste 4: Testes Automatizados**
```bash
go test ./...
```

## ✅ Verificação de Requisitos

- ✅ Rate limiting por IP (5 req/s)
- ✅ Rate limiting por token (10 req/s)
- ✅ Token sobrescreve IP
- ✅ Resposta HTTP 429
- ✅ Persistência Redis
- ✅ Strategy pattern
- ✅ Middleware HTTP
- ✅ Configuração via ambiente

## 🚨 Problemas Comuns

**Redis não conecta:**
```bash
docker-compose ps
```

**Porta ocupada:**
```bash
# Alterar SERVER_PORT no config.env
```

**Testes falham:**
```bash
go mod tidy
go test ./...
```

---

**Tempo estimado para execução: 2-5 minutos**
