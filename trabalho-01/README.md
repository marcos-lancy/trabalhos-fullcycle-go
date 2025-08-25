# Sistema de Cotação do Dólar

Este projeto implementa um sistema cliente-servidor em Go para obter e persistir cotações do dólar.

## Estrutura do Projeto

- `server.go`: Servidor HTTP que consome a API de cotação e persiste dados em arquivo JSON
- `client.go`: Cliente que faz requisição ao servidor e salva a cotação em arquivo
- `go.mod`: Gerenciamento de dependências

## Requisitos

- Go 1.21 ou superior
- Conexão com internet para acessar a API de cotação

## Como Executar

### 1. Instalar dependências
```bash
go mod tidy
```

### 2. Executar o servidor
```bash
go run server.go
```
O servidor estará disponível em `http://localhost:8080`

### 3. Executar o cliente (em outro terminal)
```bash
go run client.go
```

## Funcionalidades

### Server.go
- Endpoint: `/cotacao`
- Porta: 8080
- Consome API: https://economia.awesomeapi.com.br/json/last/USD-BRL
- Timeout para API: 200ms
- Timeout para banco: 10ms
- Persiste dados em arquivo JSON (cotacao.json)

### Client.go
- Timeout para requisição: 300ms
- Salva cotação no arquivo `cotacao.txt`
- Formato: "Dólar: {valor}"

## Timeouts e Tratamento de Erros

Todos os contextos implementam timeouts conforme especificado:
- API de cotação: 200ms
- Banco de dados: 10ms  
- Cliente: 300ms

Erros de timeout são logados quando os tempos são excedidos.

## Arquivos Gerados

- `cotacao.json`: Arquivo com histórico de cotações
- `cotacao.txt`: Arquivo com a cotação atual
- `server.exe` e `client.exe`: Executáveis compilados
