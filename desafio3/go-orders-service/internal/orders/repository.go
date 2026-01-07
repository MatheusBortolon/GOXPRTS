package orders

import "context"

type PostgresOrderRepository struct{}

func NewPostgresOrderRepository() OrderRepository {
	return &PostgresOrderRepository{}
}

func (r *PostgresOrderRepository) Create(order *Order) (string, error) {
	return order.ID, nil
}

func (r *PostgresOrderRepository) GetByID(ctx context.Context, id string) (*Order, error) {
	return nil, nil
}

func (r *PostgresOrderRepository) List(ctx context.Context) ([]*Order, error) {
	return nil, nil
}
