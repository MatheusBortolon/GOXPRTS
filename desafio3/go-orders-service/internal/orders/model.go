package orders

import "context"

type Order struct {
	ID         string  `json:"id"`
	CustomerID string  `json:"customer_id"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}
type OrderRepository interface {
	Create(order *Order) (string, error)
	GetByID(ctx context.Context, id string) (*Order, error)
	List(ctx context.Context) ([]*Order, error)
}
