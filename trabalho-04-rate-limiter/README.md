# Rate Limiter

Um rate limiter robusto e configurável implementado em Go, que permite controlar o tráfego de requisições baseado em endereço IP ou token de acesso.

## 🚀 Execução Rápida

### **Método 1: Docker (Recomendado)**
```bash
# 1. Clone o repositório
git clone <repository-url>
cd rate-limiter

# 2. Execute tudo com Docker
docker-compose up --build

# 3. Teste
curl http://localhost:8080/health
```

### **Método 2: Local**
```bash
# 1. Clone o repositório
git clone <repository-url>
cd rate-limiter

# 2. Instale dependências
go mod download

# 3. Inicie Redis
docker-compose up redis -d

# 4. Execute aplicação
go run main.go

# 5. Teste
curl http://localhost:8080/health
```

## 🧪 Testes

```bash
# Execute todos os testes
go test ./...

# Teste rate limiting por IP
for i in {1..6}; do
  curl http://localhost:8080/api/data
  echo "Request $i"
done

# Teste rate limiting por token
for i in {1..11}; do
  curl -H "API_KEY: test-token" http://localhost:8080/api/data
  echo "Request $i"
done
```

## Funcionalidades

- **Rate Limiting por IP**: Limita requisições baseado no endereço IP do cliente
- **Rate Limiting por Token**: Limita requisições baseado em token de acesso (header `API_KEY`)
- **Prioridade de Token**: Configurações de token sobrescrevem as configurações de IP
- **Middleware HTTP**: Integração fácil com servidores web usando Gin
- **Persistência Redis**: Armazenamento de dados do rate limiter no Redis
- **Strategy Pattern**: Arquitetura flexível que permite trocar Redis por outros mecanismos
- **Configuração via Ambiente**: Configuração através de variáveis de ambiente ou arquivo `.env`
- **Tempo de Bloqueio Configurável**: Tempo personalizável de bloqueio quando limite é excedido
- **Resposta HTTP 429**: Resposta padronizada quando limite é excedido

## Arquitetura

O projeto segue uma arquitetura limpa e modular:

```
├── internal/
│   ├── limiter/          # Lógica do rate limiter
│   ├── middleware/       # Middleware HTTP
│   └── storage/          # Interface e implementações de persistência
├── cmd/server/           # Servidor de exemplo
├── main.go              # Ponto de entrada principal
├── docker-compose.yml   # Configuração Docker
└── Dockerfile           # Imagem Docker
```

### Componentes Principais

1. **Storage Interface**: Define contratos para persistência de dados
2. **Redis Storage**: Implementação usando Redis
3. **Memory Storage**: Implementação em memória (para testes)
4. **Rate Limiter**: Lógica principal de controle de taxa
5. **Middleware**: Integração com servidores HTTP
6. **Config**: Gerenciamento de configurações

## Configuração

### Variáveis de Ambiente

Crie um arquivo `config.env` na raiz do projeto:

```env
# Rate Limiter Configuration
RATE_LIMIT_IP_REQUESTS_PER_SECOND=5
RATE_LIMIT_IP_BLOCK_DURATION_MINUTES=5
RATE_LIMIT_TOKEN_REQUESTS_PER_SECOND=10
RATE_LIMIT_TOKEN_BLOCK_DURATION_MINUTES=5

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Server Configuration
SERVER_PORT=8080
```

### Parâmetros de Configuração

- `RATE_LIMIT_IP_REQUESTS_PER_SECOND`: Número máximo de requisições por segundo por IP
- `RATE_LIMIT_IP_BLOCK_DURATION_MINUTES`: Duração do bloqueio em minutos para IP
- `RATE_LIMIT_TOKEN_REQUESTS_PER_SECOND`: Número máximo de requisições por segundo por token
- `RATE_LIMIT_TOKEN_BLOCK_DURATION_MINUTES`: Duração do bloqueio em minutos para token
- `REDIS_HOST`: Host do Redis
- `REDIS_PORT`: Porta do Redis
- `REDIS_PASSWORD`: Senha do Redis (opcional)
- `REDIS_DB`: Número do banco de dados Redis
- `SERVER_PORT`: Porta do servidor web

## 📋 Pré-requisitos

- Go 1.21 ou superior
- Docker e Docker Compose
- Redis (ou usar Docker Compose)

## ⚙️ Configuração

O arquivo `config.env` já está configurado com valores padrão:

```env
RATE_LIMIT_IP_REQUESTS_PER_SECOND=5
RATE_LIMIT_TOKEN_REQUESTS_PER_SECOND=10
REDIS_HOST=localhost
REDIS_PORT=6379
SERVER_PORT=8080
```

## 📡 Endpoints Disponíveis

- `GET /health` - Health check
- `GET /api/data` - Dados protegidos (exemplo)
- `POST /api/data` - Criar dados (exemplo)
- `GET /api/rate-limit/status` - Status do rate limit
- `POST /api/rate-limit/reset` - Reset do rate limit (para testes)

## 🔍 Verificação de Funcionamento

### **Teste 1: Rate Limiting por IP**
```bash
# Faça 6 requisições rapidamente
for i in {1..6}; do
  curl http://localhost:8080/api/data
  echo "Request $i"
done
# Resultado: Requests 1-5 = HTTP 200, Request 6 = HTTP 429
```

### **Teste 2: Rate Limiting por Token**
```bash
# Faça 11 requisições com token
for i in {1..11}; do
  curl -H "API_KEY: test-token" http://localhost:8080/api/data
  echo "Request $i"
done
# Resultado: Requests 1-10 = HTTP 200, Request 11 = HTTP 429
```

### **Teste 3: Verificar Status**
```bash
# Status do rate limit
curl http://localhost:8080/api/rate-limit/status
```

## 🏗️ Arquitetura

```
├── internal/
│   ├── limiter/          # Lógica do rate limiter
│   ├── middleware/       # Middleware HTTP
│   └── storage/          # Persistência (Redis/Memory)
├── main.go              # Servidor principal
├── docker-compose.yml   # Orquestração
└── config.env           # Configurações
```

## ✅ Requisitos Atendidos

- ✅ Rate limiting por IP
- ✅ Rate limiting por Token  
- ✅ Token sobrescreve configurações de IP
- ✅ Middleware HTTP integrado
- ✅ Configuração via variáveis de ambiente
- ✅ Resposta HTTP 429 quando limite excedido
- ✅ Persistência no Redis
- ✅ Strategy pattern para trocar storage
- ✅ Separação de responsabilidades
- ✅ Docker-compose configurado
- ✅ Testes automatizados
- ✅ Documentação completa

## 🚨 Troubleshooting

### **Problema: Redis não conecta**
```bash
# Verificar se Redis está rodando
docker-compose ps
```

### **Problema: Porta 8080 ocupada**
```bash
# Alterar porta no config.env
SERVER_PORT=8081
```

### **Problema: Testes falham**
```bash
# Verificar dependências
go mod tidy
go test ./...
```

---

**Projeto desenvolvido para trabalho de pós-graduação - Rate Limiter em Go**
