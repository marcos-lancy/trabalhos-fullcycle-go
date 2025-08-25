package repository

import (
	"trabalho-03/internal/domain"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *domain.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) List() ([]domain.Order, error) {
	var orders []domain.Order
	err := r.db.Find(&orders).Error
	return orders, err
}
