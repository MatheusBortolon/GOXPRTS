package graphql

import (
	"context"

	"github.com/MatheusBortolon/GOXPRTS/desafio3/go-orders-service/internal/orders"
	"github.com/graphql-go/graphql"
)

type Resolver struct {
	orderService *orders.OrderService
}

func NewResolver(orderService *orders.OrderService) *Resolver {
	return &Resolver{orderService: orderService}
}

func (r *Resolver) ListOrdersResolver(p graphql.ResolveParams) (interface{}, error) {
	ctx := context.Background()
	return r.orderService.ListOrders(ctx)
}

func (r *Resolver) GetOrderResolver(p graphql.ResolveParams) (interface{}, error) {
	ctx := context.Background()
	id := p.Args["id"].(string)
	return r.orderService.GetOrder(ctx, id)
}
