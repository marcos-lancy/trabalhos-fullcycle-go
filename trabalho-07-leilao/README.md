# Sistema de Leilão com Fechamento Automático

Este projeto implementa um sistema de leilão com funcionalidade de fechamento automático baseado em tempo configurável.

## Funcionalidades

- Criação de leilões
- Sistema de lances (bids)
- Fechamento automático de leilões baseado em tempo configurável
- API REST para gerenciamento de leilões, lances e usuários
- Goroutine para monitoramento e fechamento automático de leilões vencidos

## Como Executar

### Pré-requisitos

- Docker
- Docker Compose
- Go 1.20+ (para desenvolvimento local)

### Executando com Docker

1. Clone o repositório:
```bash
git clone <repository-url>
cd trabalho-07-leilao
```

2. Configure as variáveis de ambiente:
```bash
cp env.example cmd/auction/.env
```

3. Execute o projeto com Docker Compose:
```bash
docker-compose up --build
```

O sistema estará disponível em `http://localhost:8080`

### Executando Localmente

1. Configure as variáveis de ambiente:
```bash
cp env.example cmd/auction/.env
```

2. Execute o MongoDB:
```bash
docker run -d -p 27017:27017 --name mongodb mongo:latest
```

3. Execute a aplicação:
```bash
go run cmd/auction/main.go
```

## Variáveis de Ambiente

- `MONGODB_URI`: URI de conexão com o MongoDB
- `MONGODB_DATABASE`: Nome do banco de dados
- `AUCTION_DURATION`: Duração do leilão (ex: "5m", "1h", "30s")
- `AUCTION_INTERVAL`: Intervalo para verificação de leilões vencidos

## API Endpoints

### Leilões
- `GET /auction` - Listar leilões
- `GET /auction/:auctionId` - Buscar leilão por ID
- `POST /auction` - Criar novo leilão
- `GET /auction/winner/:auctionId` - Buscar lance vencedor

### Lances
- `POST /bid` - Criar novo lance
- `GET /bid/:auctionId` - Buscar lances por leilão

### Usuários
- `GET /user/:userId` - Buscar usuário por ID

## Funcionalidade de Fechamento Automático

O sistema implementa uma goroutine que:

1. Verifica a cada 10 segundos se existem leilões vencidos
2. Fecha automaticamente leilões que excederam o tempo configurado em `AUCTION_DURATION`
3. Utiliza concorrência segura com mutex para evitar condições de corrida
4. Registra logs das operações de fechamento

## Testes

Para executar os testes:

```bash
go test ./internal/infra/database/auction/...
```

Os testes validam:
- Fechamento automático de leilões após o tempo configurado
- Leilões permanecem ativos antes do tempo de expiração
- Funcionamento correto da goroutine de fechamento

## Estrutura do Projeto

```
├── cmd/auction/           # Ponto de entrada da aplicação
├── internal/
│   ├── entity/           # Entidades de domínio
│   ├── infra/            # Infraestrutura (API, banco de dados)
│   ├── usecase/          # Casos de uso
│   └── internal_error/   # Tratamento de erros
├── configuration/        # Configurações (logger, banco)
└── docker-compose.yml    # Configuração do Docker
```

## Tecnologias Utilizadas

- Go 1.20
- MongoDB
- Gin (Framework Web)
- Docker/Docker Compose
- Goroutines para concorrência
