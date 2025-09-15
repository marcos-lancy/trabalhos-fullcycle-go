# Build stage
FROM golang:1.21-alpine AS builder

# Instalar dependências necessárias
RUN apk add --no-cache git

# Definir diretório de trabalho
WORKDIR /app

# Copiar arquivos de dependências
COPY go.mod go.sum ./

# Baixar dependências
RUN go mod download

# Copiar código fonte
COPY . .

# Compilar a aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o stress-test .

# Runtime stage
FROM alpine:latest

# Instalar ca-certificates para requisições HTTPS
RUN apk --no-cache add ca-certificates

# Criar usuário não-root
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Definir diretório de trabalho
WORKDIR /app

# Copiar o binário compilado
COPY --from=builder /app/stress-test .

# Mudar para usuário não-root
USER appuser

# Expor porta (opcional, para casos onde a aplicação precise de uma porta)
EXPOSE 8080

# Comando padrão
ENTRYPOINT ["./stress-test"]
