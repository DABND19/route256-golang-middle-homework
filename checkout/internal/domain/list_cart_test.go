package domain

import (
	"context"
	"errors"
	lomsClientMock "route256/checkout/internal/clients/mocks/loms"
	productClientMock "route256/checkout/internal/clients/mocks/product"
	"route256/checkout/internal/models"
	cartsRepoMock "route256/checkout/internal/repository/mock/carts"
	trMock "route256/libs/transactor/mock"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestListCart(t *testing.T) {
	// data for success scenario
	expectedCart := make([]models.CartProduct, 10)
	gofakeit.Slice(&expectedCart)
	itemsMock := make([]models.CartItem, 0, len(expectedCart))
	for _, item := range expectedCart {
		itemsMock = append(itemsMock, item.CartItem)
	}
	productsMock := make(map[models.SKU]*models.Product, len(expectedCart))
	for _, item := range expectedCart {
		product := item.Product
		productsMock[item.SKU] = &product
	}

	// data for product not found scenario
	notFoundSKU := itemsMock[0].SKU
	missedProductsMock := make(map[models.SKU]*models.Product)
	for k, v := range productsMock {
		missedProductsMock[k] = v
	}
	delete(missedProductsMock, notFoundSKU)
	expectedCartWithMissedProduct := expectedCart[1:]

	var (
		getCartItemsError      = errors.New("Get cart items error")
		getProductsError       = errors.New("Get products error")
		transactionRunnerError = errors.New("Some transaction runner error")
	)

	defaultTrMock := func(mc *minimock.Controller, expectedTxErr error) TransactionRunner {
		return trMock.NewTransactionRunnerMock(mc).RunReadCommitedMock.Set(
			func(ctx context.Context, txFn func(ctx context.Context) error) error {
				err := txFn(ctx)
				if expectedTxErr != nil {
					require.ErrorIs(mc, expectedTxErr, err)
				} else {
					require.Equal(mc, nil, err)
				}
				return err
			},
		)
	}
	defaultCartsRepoMock := func(mc *minimock.Controller) CartsRepository {
		return cartsRepoMock.NewCartsRepositoryMock(mc).GetCartItemsMock.Return(itemsMock, nil)
	}
	defaultProductClientMock := func(mc *minimock.Controller) ProductServiceClient {
		return productClientMock.NewProductServiceClientMock(mc).GetProductsMock.Return(productsMock, nil)
	}

	tests := []struct {
		Name           string
		ExpectedResult []models.CartProduct
		ExpectedErr    error
		Tr             func(mc *minimock.Controller, expectedTxErr error) TransactionRunner
		ExpectedTxErr  error
		CartsRepo      func(mc *minimock.Controller) CartsRepository
		ProductClient  func(mc *minimock.Controller) ProductServiceClient
	}{
		{
			Name:           "success",
			ExpectedResult: expectedCart,
			ExpectedErr:    nil,
			Tr:             defaultTrMock,
			ExpectedTxErr:  nil,
			CartsRepo:      defaultCartsRepoMock,
			ProductClient:  defaultProductClientMock,
		},
		{
			Name:           "product not found",
			ExpectedResult: expectedCartWithMissedProduct,
			ExpectedErr:    nil,
			Tr:             defaultTrMock,
			ExpectedTxErr:  nil,
			CartsRepo:      defaultCartsRepoMock,
			ProductClient: func(mc *minimock.Controller) ProductServiceClient {
				return productClientMock.NewProductServiceClientMock(mc).GetProductsMock.Return(missedProductsMock, nil)
			},
		},
		{
			Name:           "cart repo error",
			ExpectedResult: nil,
			ExpectedErr:    getCartItemsError,
			Tr:             defaultTrMock,
			ExpectedTxErr:  getCartItemsError,
			CartsRepo: func(mc *minimock.Controller) CartsRepository {
				return cartsRepoMock.NewCartsRepositoryMock(mc).GetCartItemsMock.Return(nil, getCartItemsError)
			},
			ProductClient: defaultProductClientMock,
		},
		{
			Name:           "product service error",
			ExpectedResult: nil,
			ExpectedErr:    getProductsError,
			Tr:             defaultTrMock,
			ExpectedTxErr:  nil,
			CartsRepo:      defaultCartsRepoMock,
			ProductClient: func(mc *minimock.Controller) ProductServiceClient {
				return productClientMock.NewProductServiceClientMock(mc).GetProductsMock.Return(nil, getProductsError)
			},
		},
		{
			Name:           "transaction runner error",
			ExpectedResult: nil,
			ExpectedErr:    transactionRunnerError,
			Tr: func(mc *minimock.Controller, expectedTxError error) TransactionRunner {
				return trMock.NewTransactionRunnerMock(mc).RunReadCommitedMock.Inspect(
					func(ctx context.Context, txFn func(ctx context.Context) error) {
						err := txFn(ctx)
						if expectedTxError != nil {
							require.ErrorIs(mc, err, expectedTxError)
						} else {
							require.Equal(mc, nil, err)
						}
					},
				).Return(transactionRunnerError)
			},
			ExpectedTxErr: nil,
			CartsRepo:     defaultCartsRepoMock,
			ProductClient: defaultProductClientMock,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.Name, func(t *testing.T) {
			mc := minimock.NewController(t)
			s := New(
				testCase.Tr(mc, testCase.ExpectedTxErr),
				testCase.CartsRepo(mc),
				lomsClientMock.NewLOMSServiceClientMock(mc),
				testCase.ProductClient(mc),
			)
			cart, err := s.ListCart(context.Background(), 1)
			if testCase.ExpectedErr != nil {
				require.ErrorIs(t, err, testCase.ExpectedErr)
			} else {
				require.Equal(t, cart, testCase.ExpectedResult)
				require.Equal(t, err, testCase.ExpectedErr)
			}
		})
	}
}

func TestCalculateTotalPrice(t *testing.T) {
	tests := []struct {
		Name           string
		CartItems      []models.CartProduct
		ExpectedResult uint32
	}{
		{
			Name:           "empty cart",
			CartItems:      []models.CartProduct{},
			ExpectedResult: 0,
		},
		{
			Name: "basic",
			CartItems: []models.CartProduct{
				{
					CartItem: models.CartItem{
						SKU:   models.SKU(gofakeit.Uint32()),
						Count: 3,
					},
					Product: models.Product{
						Name:  gofakeit.BeerName(),
						Price: 100,
					},
				},
				{
					CartItem: models.CartItem{
						SKU:   models.SKU(gofakeit.Uint32()),
						Count: 2,
					},
					Product: models.Product{
						Name:  gofakeit.BeerName(),
						Price: 200,
					},
				},
				{
					CartItem: models.CartItem{
						SKU:   models.SKU(gofakeit.Uint32()),
						Count: 1,
					},
					Product: models.Product{
						Name:  gofakeit.BeerName(),
						Price: 300,
					},
				},
			},
			ExpectedResult: 1000,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.Name, func(t *testing.T) {
			mc := minimock.NewController(t)
			s := New(
				trMock.NewTransactionRunnerMock(mc),
				cartsRepoMock.NewCartsRepositoryMock(mc),
				lomsClientMock.NewLOMSServiceClientMock(mc),
				productClientMock.NewProductServiceClientMock(mc),
			)
			total := s.CalculateTotalPrice(testCase.CartItems)
			require.Equal(t, total, testCase.ExpectedResult)
		})
	}
}
