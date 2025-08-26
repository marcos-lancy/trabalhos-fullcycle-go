package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"trabalho-03/graphql"
	"trabalho-03/internal/domain"
	grpcserver "trabalho-03/internal/grpc"
	"trabalho-03/internal/handler"
	"trabalho-03/internal/repository"
	"trabalho-03/internal/usecase"

	gqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Database connection with environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "orders_db")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate
	err = db.AutoMigrate(&domain.Order{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize repositories
	orderRepo := repository.NewOrderRepository(db)

	// Initialize use cases
	orderUseCase := usecase.NewOrderUseCase(orderRepo)

	// Initialize handlers
	restHandler := handler.NewOrderHandler(orderUseCase)
	grpcServer := grpcserver.NewGRPCServer(orderUseCase)

	// Initialize GraphQL resolver
	resolver := &graphql.Resolver{
		OrderUseCase: orderUseCase,
	}

	// Start servers
	go startRESTServer(restHandler)
	go grpcServer.Start()
	go startGraphQLServer(resolver)

	// Keep the main thread alive
	select {}
}

func startRESTServer(orderHandler *handler.OrderHandler) {
	r := gin.Default()

	r.POST("/order", orderHandler.CreateOrder)
	r.GET("/order", orderHandler.ListOrders)

	fmt.Println("REST server starting on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start REST server:", err)
	}
}

func startGraphQLServer(resolver *graphql.Resolver) {
	srv := gqlhandler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: resolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("GraphQL server starting on port 8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("Failed to start GraphQL server:", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
