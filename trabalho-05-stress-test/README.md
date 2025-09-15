# Stress Test Tool

Uma ferramenta CLI em Go para realizar testes de carga em serviços web.

## Descrição

Esta aplicação permite executar testes de carga em qualquer serviço web, fornecendo métricas detalhadas sobre performance, distribuição de códigos de status HTTP e estatísticas de tempo de resposta.

## Funcionalidades

- ✅ Testes de carga com controle de concorrência
- ✅ Relatório detalhado com métricas de performance
- ✅ Distribuição de códigos de status HTTP
- ✅ Estatísticas de tempo de resposta (mínimo, máximo, médio)
- ✅ Cálculo de requests por segundo
- ✅ Containerização com Docker

## Pré-requisitos

- Go 1.21 ou superior
- Docker (opcional, para execução via container)

## Instalação e Execução

### Método 1: Execução Direta com Go

1. **Clone o repositório:**
   ```bash
   git clone <url-do-repositorio>
   cd trabalho-05-stress-test
   ```

2. **Instale as dependências:**
   ```bash
   go mod download
   ```

3. **Compile a aplicação:**
   ```bash
   go build -o stress-test .
   ```

4. **Execute o teste:**
   ```bash
   ./stress-test --url=http://google.com --requests=1000 --concurrency=10
   ```

### Método 2: Execução via Docker

1. **Construa a imagem Docker:**
   ```bash
   docker build -t stress-test .
   ```

2. **Execute o container:**
   ```bash
   docker run stress-test --url=http://google.com --requests=1000 --concurrency=10
   ```

## Parâmetros de Entrada

| Parâmetro | Descrição | Obrigatório | Exemplo |
|-----------|-----------|-------------|---------|
| `--url` | URL do serviço a ser testado | Sim | `http://google.com` |
| `--requests` | Número total de requests | Sim | `1000` |
| `--concurrency` | Número de chamadas simultâneas | Sim | `10` |

## Exemplo de Uso

```bash
# Teste básico
./stress-test --url=http://httpbin.org/get --requests=100 --concurrency=5

# Teste com mais requests
./stress-test --url=https://api.github.com --requests=1000 --concurrency=20

# Via Docker
docker run stress-test --url=http://google.com --requests=500 --concurrency=15
```

## Relatório de Saída

A aplicação gera um relatório detalhado contendo:

- **Tempo total de execução**
- **Total de requests realizados**
- **Quantidade de requests com status 200**
- **Distribuição de códigos de status HTTP**
- **Estatísticas de tempo de resposta:**
  - Tempo médio
  - Tempo mínimo
  - Tempo máximo
- **Requests por segundo**

### Exemplo de Relatório

```
========================================
RELATÓRIO DO TESTE DE CARGA
========================================
Tempo total de execução: 2.345s
Total de requests realizados: 1000
Requests com status 200: 950

Distribuição de códigos de status:
  Status 200: 950
  Status 404: 30
  Status 500: 20

Estatísticas de tempo de resposta:
  Tempo médio: 45ms
  Tempo mínimo: 12ms
  Tempo máximo: 234ms
  Requests por segundo: 426.44
========================================
```

## Estrutura do Projeto

```
trabalho-05-stress-test/
├── main.go          # Código principal da aplicação
├── go.mod           # Dependências do Go
├── go.sum           # Checksums das dependências
├── Dockerfile       # Configuração do Docker
├── README.md        # Este arquivo
└── required.md      # Requisitos do trabalho
```

## Tecnologias Utilizadas

- **Go 1.21**: Linguagem de programação
- **Cobra**: Biblioteca para CLI
- **Docker**: Containerização
- **Goroutines**: Concorrência para testes de carga

## Considerações Técnicas

- A aplicação utiliza goroutines para controlar a concorrência
- Implementa semáforos para limitar o número de requests simultâneos
- Coleta métricas em tempo real durante a execução
- Trata erros de conexão e timeouts adequadamente
- Otimizada para performance com uso eficiente de memória

## Troubleshooting

### Erro de Conexão
Se você receber erros de conexão, verifique:
- Se a URL está acessível
- Se não há firewall bloqueando as conexões
- Se o serviço de destino está funcionando

### Timeout
Para serviços lentos, considere:
- Reduzir o número de requests simultâneos
- Aumentar o timeout (modificar o código se necessário)

### Performance
Para melhor performance:
- Ajuste o número de goroutines baseado na capacidade do sistema
- Monitore o uso de CPU e memória durante os testes

## Contribuição

Este projeto foi desenvolvido como parte de um trabalho de pós-graduação. Para sugestões ou melhorias, abra uma issue no repositório.

## Licença

Este projeto é de uso educacional e acadêmico.
