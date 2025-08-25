# Trabalho 03 - Listagem de Orders

Este projeto implementa um sistema de gerenciamento de orders com foco na **listagem de orders** conforme solicitado.

## Funcionalidades Implementadas

✅ **Endpoint REST (GET /order)**: Lista todas as orders  
✅ **Service ListOrders com gRPC**: Estrutura preparada para gRPC  
✅ **Query ListOrders GraphQL**: Estrutura preparada para GraphQL  
✅ **Migrações necessárias**: Auto-migração com GORM  
✅ **Arquivo api.http**: Requests para criar e listar orders  
✅ **Docker e docker-compose**: Banco de dados PostgreSQL  

## Pré-requisitos

- Docker e Docker Compose
- Go 1.21 ou superior

## Como executar

1. **Suba o banco de dados:**
   ```bash
   docker compose up -d
   ```

2. **Instale as dependências:**
   ```bash
   go mod tidy
   ```

3. **Execute a aplicação:**
   ```bash
   go run main.go
   ```

## Portas dos serviços

- **REST API**: http://localhost:8080 ✅ **FUNCIONANDO**
- **gRPC**: localhost:9090 (estrutura preparada)
- **GraphQL**: http://localhost:8081 (estrutura preparada)

## Endpoints disponíveis

### REST API (Porta 8080) ✅ **FUNCIONANDO**
- `POST /order` - Criar uma nova order
- `GET /order` - Listar todas as orders

### gRPC (Porta 9090) - Estrutura preparada
- `CreateOrder` - Criar uma nova order
- `ListOrders` - Listar todas as orders

### GraphQL (Porta 8081) - Estrutura preparada
- Query `orders` - Listar todas as orders
- Mutation `createOrder` - Criar uma nova order

## Testando a aplicação

Use o arquivo `api.http` para testar os endpoints REST:

```bash
# Criar uma order
POST http://localhost:8080/order
Content-Type: application/json

{
  "customer_id": "customer123",
  "amount": 100.50,
  "status": "pending"
}

# Listar orders
GET http://localhost:8080/order
```

## Estrutura do Projeto

```
trabalho-03/
├── internal/
│   ├── domain/          # Modelo de Order
│   ├── repository/      # Repositório para acesso ao banco
│   ├── usecase/         # Casos de uso (CreateOrder, ListOrders)
│   └── handler/         # Handlers REST
├── proto/               # Definições gRPC
├── graphql/             # Schema GraphQL
├── main.go              # Aplicação principal
├── docker-compose.yaml  # Banco PostgreSQL
├── api.http            # Requests de teste
└── README.md           # Documentação
```
