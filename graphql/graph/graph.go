package graph

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/Dzhodddi/EcommerceAPI/gateway/generated"

	account "github.com/Dzhodddi/EcommerceAPI/account/client"
	order "github.com/Dzhodddi/EcommerceAPI/order/client"
	payment "github.com/Dzhodddi/EcommerceAPI/payment/client"
	product "github.com/Dzhodddi/EcommerceAPI/product/client"
)

type Server struct {
	accountClient *account.Client
	productClient *product.Client
	orderClient   *order.Client
	paymentClient *payment.Client
}

func NewGraphQLServer(accountUrl, productUrl, orderUrl, paymentUrl, _ string) (*Server, error) {
	accClient, err := account.NewClient(accountUrl)
	if err != nil {
		return nil, err
	}

	prodClient, err := product.NewClient(productUrl)
	if err != nil {
		accClient.Close()
		return nil, err
	}

	ordClient, err := order.NewClient(orderUrl)
	if err != nil {
		accClient.Close()
		prodClient.Close()
		return nil, err
	}

	paymentClient, err := payment.NewClient(paymentUrl)
	if err != nil {
		accClient.Close()
		prodClient.Close()
		ordClient.Close()
	}

	return &Server{
		accountClient: accClient,
		productClient: prodClient,
		orderClient:   ordClient,
		paymentClient: paymentClient,
	}, nil
}

func (server *Server) Mutation() generated.MutationResolver {
	return &mutationResolver{
		server: server,
	}
}

func (server *Server) Query() generated.QueryResolver {
	return &queryResolver{
		server: server,
	}
}

func (server *Server) Account() generated.AccountResolver {
	return &accountResolver{
		server: server,
	}
}

func (server *Server) ToExecutableSchema() graphql.ExecutableSchema {
	return generated.NewExecutableSchema(generated.Config{
		Resolvers: server,
	})
}
