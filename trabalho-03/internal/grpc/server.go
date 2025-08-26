package grpc

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"trabalho-03/internal/usecase"

	"github.com/gin-gonic/gin"
)

// GRPCServer simulates a gRPC server using HTTP/JSON
type GRPCServer struct {
	orderService *OrderService
}

func NewGRPCServer(orderUseCase *usecase.OrderUseCase) *GRPCServer {
	return &GRPCServer{
		orderService: NewOrderService(orderUseCase),
	}
}

func (s *GRPCServer) Start() {
	r := gin.New()
	r.Use(gin.Logger())
	
	// gRPC-style endpoints
	r.POST("/order.OrderService/CreateOrder", s.handleCreateOrder)
	r.POST("/order.OrderService/ListOrders", s.handleListOrders)
	
	// Alternative REST-like endpoints for easier testing
	r.GET("/grpc/orders", s.handleListOrdersGET)
	r.POST("/grpc/orders", s.handleCreateOrderREST)
	
	// Health check
	r.GET("/grpc/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "OrderService"})
	})

	fmt.Println("gRPC-like server starting on port 9090")
	if err := r.Run(":9090"); err != nil {
		log.Fatal("Failed to start gRPC server:", err)
	}
}

func (s *GRPCServer) handleCreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	ctx := context.Background()
	resp, err := s.orderService.CreateOrder(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (s *GRPCServer) handleListOrders(c *gin.Context) {
	ctx := context.Background()
	req := &ListOrdersRequest{}
	
	resp, err := s.orderService.ListOrders(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list orders", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (s *GRPCServer) handleListOrdersGET(c *gin.Context) {
	ctx := context.Background()
	req := &ListOrdersRequest{}
	
	resp, err := s.orderService.ListOrders(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list orders", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (s *GRPCServer) handleCreateOrderREST(c *gin.Context) {
	customerID := c.PostForm("customer_id")
	amountStr := c.PostForm("amount")
	status := c.PostForm("status")

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount format"})
		return
	}

	req := &CreateOrderRequest{
		CustomerID: customerID,
		Amount:     amount,
		Status:     status,
	}

	ctx := context.Background()
	resp, err := s.orderService.CreateOrder(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
