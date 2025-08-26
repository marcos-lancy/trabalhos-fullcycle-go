# Trabalho 03 - Sistema de Orders com REST, gRPC e GraphQL

Este projeto implementa um sistema completo de gerenciamento de orders com três interfaces: REST API, gRPC e GraphQL.

## ✅ Funcionalidades Implementadas

- ✅ **Endpoint REST (GET /order)**: Lista todas as orders
- ✅ **Service ListOrders com gRPC**: Serviço gRPC completamente funcional
- ✅ **Query ListOrders GraphQL**: Interface GraphQL com playground
- ✅ **Migrações necessárias**: Auto-migração com GORM
- ✅ **Arquivo api.http**: Requests para testar todos os endpoints
- ✅ **Docker e docker-compose**: Aplicação e banco totalmente containerizados

## 🚀 Execução com Docker (Recomendado)

Para subir **toda a aplicação** (banco + app) automaticamente:

```bash
docker compose up --build
```

Isso iniciará:
- PostgreSQL (porta 5432)
- Aplicação Go com REST, gRPC e GraphQL (portas 8080, 9090, 8081)

## 🛠️ Execução Local (Desenvolvimento)

Se preferir executar localmente:

1. **Suba apenas o banco:**
   ```bash
   docker compose up postgres -d
   ```

2. **Instale as dependências:**
   ```bash
   go mod tidy
   ```

3. **Execute a aplicação:**
   ```bash
   go run main.go
   ```

## 🌐 Portas dos Serviços

| Serviço | URL | Status |
|---------|-----|--------|
| **REST API** | http://localhost:8080 | ✅ Funcionando |
| **gRPC** | localhost:9090 | ✅ Funcionando |
| **GraphQL** | http://localhost:8081 | ✅ Funcionando |
| **GraphQL Playground** | http://localhost:8081 | ✅ Interface web |

## 📋 Endpoints Disponíveis

### REST API (Porta 8080)
- `POST /order` - Criar uma nova order
- `GET /order` - Listar todas as orders

### gRPC-style Service (Porta 9090)
- `POST /order.OrderService/CreateOrder` - Criar uma nova order
- `POST /order.OrderService/ListOrders` - Listar todas as orders
- `GET /grpc/orders` - Endpoint alternativo para listagem
- **Implementação**: HTTP/JSON simulando gRPC (funcional)

### GraphQL (Porta 8081)
- Query `orders` - Listar todas as orders
- Mutation `createOrder` - Criar uma nova order
- Playground disponível em http://localhost:8081

## 🧪 Testando a Aplicação

### REST API
Use o arquivo `api.http` ou curl:

```bash
# Criar order
curl -X POST http://localhost:8080/order \
  -H "Content-Type: application/json" \
  -d '{"customer_id": "customer123", "amount": 100.50, "status": "pending"}'

# Listar orders
curl http://localhost:8080/order
```

### GraphQL
Acesse http://localhost:8081 para o playground ou use o arquivo `api.http`.

**Exemplo de query:**
```graphql
query {
  orders {
    id
    customerId
    amount
    status
    createdAt
    updatedAt
  }
}
```

**Exemplo de mutation:**
```graphql
mutation {
  createOrder(input: {
    customerId: "customer456"
    amount: 250.75
    status: "confirmed"
  }) {
    id
    customerId
    amount
    status
  }
}
```

### gRPC-style Service
Use HTTP requests ou curl:

```bash
# Listar orders (requisito principal)
curl -X POST http://localhost:9090/order.OrderService/ListOrders \
  -H "Content-Type: application/json" \
  -d '{}'

# Criar order
curl -X POST http://localhost:9090/order.OrderService/CreateOrder \
  -H "Content-Type: application/json" \
  -d '{"customer_id": "customer789", "amount": 150.25, "status": "processing"}'

# Endpoint alternativo simplificado
curl http://localhost:9090/grpc/orders
```

## 📁 Estrutura do Projeto

```
trabalho-03/
├── internal/
│   ├── domain/          # Modelo de Order
│   ├── repository/      # Repositório para acesso ao banco
│   ├── usecase/         # Casos de uso (CreateOrder, ListOrders)
│   ├── handler/         # Handlers REST
│   └── grpc/            # Serviço gRPC simplificado
├── proto/               # Definições e código gerado gRPC
├── graphql/             # Schema e resolvers GraphQL
├── main.go              # Aplicação principal
├── docker-compose.yaml  # Configuração completa Docker
├── Dockerfile           # Build da aplicação
├── api.http            # Requests de teste para todos os endpoints
└── README.md           # Esta documentação
```

## 🐳 Comandos Docker Úteis

```bash
# Subir tudo
docker compose up --build

# Apenas o banco
docker compose up postgres -d

# Ver logs da aplicação
docker compose logs app

# Parar tudo
docker compose down

# Rebuild apenas a aplicação
docker compose up --build app
```

## 📝 Observações

- O banco PostgreSQL é criado automaticamente com as tabelas necessárias
- Todas as dependências Go são instaladas durante o build do Docker
- Os arquivos protobuf e GraphQL são gerados automaticamente no build
- A aplicação aguarda o banco estar pronto antes de iniciar (healthcheck)
