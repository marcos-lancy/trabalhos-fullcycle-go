# Trabalho 02 - Consulta de CEP com Multithreading

Este projeto implementa um sistema de consulta de CEP que utiliza multithreading para buscar dados de duas APIs simultaneamente e retorna o resultado da API mais rápida.

## Requisitos

- ✅ Fazer requisições simultâneas para duas APIs distintas
- ✅ Aceitar a resposta mais rápida e descartar a mais lenta
- ✅ Exibir dados do endereço e qual API respondeu
- ✅ Limitar tempo de resposta em 1 segundo
- ✅ Exibir erro de timeout quando necessário

## APIs Utilizadas

1. **BrasilAPI**: `https://brasilapi.com.br/api/cep/v1/{cep}`
2. **ViaCEP**: `http://viacep.com.br/ws/{cep}/json/`

## Como Executar

```bash
# Executar o programa
go run main.go <CEP>

# Exemplo
go run main.go 01153000
```

## Funcionalidades

- **Multithreading**: Utiliza goroutines para fazer requisições simultâneas
- **Timeout**: Limita o tempo de resposta em 1 segundo
- **Race Condition**: Aceita a primeira resposta que chegar
- **Tratamento de Erros**: Exibe mensagens de erro apropriadas
- **Interface Unificada**: Normaliza os dados das duas APIs

## Estrutura do Projeto

- `main.go`: Arquivo principal com toda a lógica
- `go.mod`: Módulo Go
- `README.md`: Documentação do projeto

## Exemplo de Saída

```
=== Resultado da Consulta de CEP ===
API: BrasilAPI
CEP: 01153-000
Logradouro: Rua Vitorino Carmilo
Bairro: Barra Funda
Cidade: São Paulo
Estado: SP
```
