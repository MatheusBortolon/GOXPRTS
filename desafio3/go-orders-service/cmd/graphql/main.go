package main

import (
	"log"
	"net/http"

	"github.com/MatheusBortolon/GOXPRTS/desafio3/go-orders-service/internal/orders"
	"github.com/MatheusBortolon/GOXPRTS/desafio3/go-orders-service/internal/transport/graphql"
	gql "github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

func main() {
	repo := orders.NewPostgresOrderRepository()
	ordersService := orders.NewOrderService(repo)
	resolver := graphql.NewResolver(ordersService)

	schema, err := gql.NewSchema(gql.SchemaConfig{
		Query: gql.NewObject(gql.ObjectConfig{
			Name: "Query",
			Fields: gql.Fields{
				"listOrders": &gql.Field{
					Type:    gql.NewList(gql.String),
					Resolve: resolver.ListOrdersResolver,
				},
				"getOrder": &gql.Field{
					Type: gql.String,
					Args: gql.FieldConfigArgument{
						"id": &gql.ArgumentConfig{
							Type: gql.NewNonNull(gql.String),
						},
					},
					Resolve: resolver.GetOrderResolver,
				},
			},
		}),
	})
	if err != nil {
		log.Fatalf("failed to create schema: %v", err)
	}

	graphqlHandler := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	http.Handle("/graphql", graphqlHandler)

	log.Println("Starting GraphQL server on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
