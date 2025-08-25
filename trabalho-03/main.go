package main

import (
	"fmt"
	"log"
	"trabalho-03/internal/domain"
	"trabalho-03/internal/handler"
	"trabalho-03/internal/repository"
	"trabalho-03/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Database connection
	dsn := "host=localhost user=postgres password=postgres dbname=orders_db port=5432 sslmode=disable"
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
	orderHandler := handler.NewOrderHandler(orderUseCase)

	// Start REST server
	go startRESTServer(orderHandler)

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
