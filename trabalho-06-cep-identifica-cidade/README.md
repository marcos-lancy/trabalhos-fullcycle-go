# CEP Weather Service

Sistema em Go que recebe um CEP, identifica a cidade e retorna o clima atual em Celsius, Fahrenheit e Kelvin.

## Como executar

### Usando Docker Compose
```bash
docker-compose up --build
```

### Usando Go diretamente
```bash
go run main.go
```

## Como testar

### Executar testes
```bash
go test
```

### Testar API
```bash
# CEP válido
curl http://localhost:8080/weather/01310100

# CEP inválido (formato incorreto)
curl http://localhost:8080/weather/1234567

# CEP não encontrado
curl http://localhost:8080/weather/99999999
```

## Respostas da API

### Sucesso (200)
```json
{
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.5
}
```

### CEP inválido (422)
```
invalid zipcode
```

### CEP não encontrado (404)
```
can not find zipcode
```