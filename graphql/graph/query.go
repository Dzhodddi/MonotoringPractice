package graph

import (
	"context"
	"github.com/Dzhodddi/EcommerceAPI/gateway/generated"
	"github.com/Dzhodddi/EcommerceAPI/gateway/models"
	"github.com/Dzhodddi/EcommerceAPI/gateway/utils"
	"time"
)

type queryResolver struct {
	server *Server
}

func (resolver *queryResolver) Accounts(
	ctx context.Context,
	pagination *generated.PaginationInput,
	id *int,
) ([]*models.Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if id != nil {
		res, err := resolver.server.accountClient.GetAccount(ctx, uint64(*id))
		if err != nil {
			return nil, err
		}
		return []*models.Account{{
			ID:    res.ID,
			Name:  res.Name,
			Email: res.Email,
		}}, nil
	}

	skip, take := uint64(0), uint64(0)
	if pagination != nil {
		skip, take = utils.Bounds(pagination)
	}
	accountList, err := resolver.server.accountClient.GetAccounts(ctx, skip, take)
	if err != nil {
		return nil, err
	}

	var accounts []*models.Account
	for _, account := range accountList {
		accounts = append(accounts, &models.Account{
			ID:    account.ID,
			Name:  account.Name,
			Email: account.Email,
		})
	}

	return accounts, nil
}

func (resolver *queryResolver) Product(
	ctx context.Context,
	pagination *generated.PaginationInput,
	query, id *string,
	_ []*string,
	_ *bool,
) ([]*generated.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if id != nil {
		res, err := resolver.server.productClient.GetProduct(ctx, *id)
		if err != nil {
			return nil, err
		}
		return []*generated.Product{{
			ID:          res.ID,
			Name:        res.Name,
			Description: res.Description,
			Price:       res.Price,
		}}, nil
	}
	skip, take := uint64(0), uint64(0)
	if pagination != nil {
		skip, take = utils.Bounds(pagination)
	}

	q := ""
	if query != nil {
		q = *query
	}
	productList, err := resolver.server.productClient.GetProducts(ctx, skip, take, nil, q)
	if err != nil {
		return nil, err
	}

	var products []*generated.Product
	for _, product := range productList {
		products = append(products,
			&generated.Product{
				ID:          product.ID,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
			},
		)
	}

	return products, nil
}
