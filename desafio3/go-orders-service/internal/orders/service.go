package orders

import (
	"context"
	"errors"
)

type OrderService struct {
	repo OrderRepository
}

func NewOrderService(repo OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(order *Order) (string, error) {
	if order.Amount <= 0 {
		return "", errors.New("invalid order amount")
	}
	return s.repo.Create(order)
}

func (s *OrderService) GetOrder(ctx context.Context, id string) (*Order, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *OrderService) ListOrders(ctx context.Context) ([]*Order, error) {
	return s.repo.List(ctx)
}
