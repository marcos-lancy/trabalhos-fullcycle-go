package grpc

import (
	"context"
	"encoding/json"
	"time"
	"trabalho-03/internal/domain"
	"trabalho-03/internal/usecase"
)

// Simple gRPC-like structs without protobuf dependency
type OrderMessage struct {
	ID         uint32  `json:"id"`
	CustomerID string  `json:"customer_id"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

type CreateOrderRequest struct {
	CustomerID string  `json:"customer_id"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
}

type CreateOrderResponse struct {
	Order *OrderMessage `json:"order"`
}

type ListOrdersRequest struct{}

type ListOrdersResponse struct {
	Orders []*OrderMessage `json:"orders"`
}

// OrderService implements a simplified gRPC-like service
type OrderService struct {
	orderUseCase *usecase.OrderUseCase
}

func NewOrderService(orderUseCase *usecase.OrderUseCase) *OrderService {
	return &OrderService{
		orderUseCase: orderUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*CreateOrderResponse, error) {
	order := &domain.Order{
		CustomerID: req.CustomerID,
		Amount:     req.Amount,
		Status:     req.Status,
	}

	err := s.orderUseCase.CreateOrder(order)
	if err != nil {
		return nil, err
	}

	orderMsg := &OrderMessage{
		ID:         uint32(order.ID),
		CustomerID: order.CustomerID,
		Amount:     order.Amount,
		Status:     order.Status,
		CreatedAt:  order.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  order.UpdatedAt.Format(time.RFC3339),
	}

	return &CreateOrderResponse{
		Order: orderMsg,
	}, nil
}

func (s *OrderService) ListOrders(ctx context.Context, req *ListOrdersRequest) (*ListOrdersResponse, error) {
	orders, err := s.orderUseCase.ListOrders()
	if err != nil {
		return nil, err
	}

	var orderMessages []*OrderMessage
	for _, order := range orders {
		orderMsg := &OrderMessage{
			ID:         uint32(order.ID),
			CustomerID: order.CustomerID,
			Amount:     order.Amount,
			Status:     order.Status,
			CreatedAt:  order.CreatedAt.Format(time.RFC3339),
			UpdatedAt:  order.UpdatedAt.Format(time.RFC3339),
		}
		orderMessages = append(orderMessages, orderMsg)
	}

	return &ListOrdersResponse{
		Orders: orderMessages,
	}, nil
}

// Helper method to convert to JSON for transport
func (s *OrderService) ListOrdersJSON(ctx context.Context) ([]byte, error) {
	req := &ListOrdersRequest{}
	resp, err := s.ListOrders(ctx, req)
	if err != nil {
		return nil, err
	}
	return json.Marshal(resp)
}

func (s *OrderService) CreateOrderJSON(ctx context.Context, customerID string, amount float64, status string) ([]byte, error) {
	req := &CreateOrderRequest{
		CustomerID: customerID,
		Amount:     amount,
		Status:     status,
	}
	resp, err := s.CreateOrder(ctx, req)
	if err != nil {
		return nil, err
	}
	return json.Marshal(resp)
}
