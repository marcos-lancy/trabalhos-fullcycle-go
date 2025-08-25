.PHONY: up down build run tidy proto graphql

# Subir o banco de dados
up:
	docker compose up -d

# Parar o banco de dados
down:
	docker compose down

# Instalar dependências
tidy:
	go mod tidy

# Gerar código do protobuf
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/order.proto

# Gerar código do GraphQL
graphql:
	go run github.com/99designs/gqlgen generate

# Build da aplicação
build:
	go build -o main .

# Executar a aplicação
run:
	go run main.go

# Setup completo
setup: up tidy proto graphql
