package usecase

import (
	"trabalho-03/internal/domain"
	"trabalho-03/internal/repository"
)

type OrderUseCase struct {
	orderRepo *repository.OrderRepository
}

func NewOrderUseCase(orderRepo *repository.OrderRepository) *OrderUseCase {
	return &OrderUseCase{orderRepo: orderRepo}
}

func (uc *OrderUseCase) CreateOrder(order *domain.Order) error {
	return uc.orderRepo.Create(order)
}

func (uc *OrderUseCase) ListOrders() ([]domain.Order, error) {
	return uc.orderRepo.List()
}
