package orders

import "context"

type MockOrderRepository struct {
	CreateFunc  func(order *Order) (string, error)
	GetByIDFunc func(ctx context.Context, id string) (*Order, error)
	ListFunc    func(ctx context.Context) ([]*Order, error)
}

func (m *MockOrderRepository) Create(order *Order) (string, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(order)
	}
	return order.ID, nil
}

func (m *MockOrderRepository) GetByID(ctx context.Context, id string) (*Order, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockOrderRepository) List(ctx context.Context) ([]*Order, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx)
	}
	return nil, nil
}
