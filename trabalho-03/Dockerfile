FROM golang:1.23-alpine

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code (protobuf and graphql files are pre-generated)
COPY . .

# Build the application
RUN go build -o main .

EXPOSE 8080 9090 8081

CMD ["./main"]
