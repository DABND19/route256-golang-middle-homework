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

func TestPurchase(t *testing.T) {
	t.Parallel()

	userIDMock := models.User(gofakeit.Int64())
	cartItemsMock := make([]models.CartItem, 10)
	gofakeit.Slice(&cartItemsMock)
	orderIDMock := new(models.OrderID)
	*orderIDMock = models.OrderID(gofakeit.Int64())

	defaultTrMock := func(mc *minimock.Controller, expectedTxErr error) *trMock.TransactionRunnerMock {
		return trMock.NewTransactionRunnerMock(mc).RunSerializableMock.Set(
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
	defaultLomsClientMock := func(mc *minimock.Controller) *lomsClientMock.LOMSServiceClientMock {
		lomsClient := lomsClientMock.NewLOMSServiceClientMock(mc)
		return lomsClient.CreateOrderMock.Return(orderIDMock, nil)
	}
	defaultCartRepoMock := func(mc *minimock.Controller) *cartsRepoMock.CartsRepositoryMock {
		cartsRepo := cartsRepoMock.NewCartsRepositoryMock(mc)
		cartsRepo = cartsRepo.GetCartItemsMock.Return(cartItemsMock, nil)
		i := 0
		cartsRepo = cartsRepo.DeleteCartItemMock.Inspect(func(ctx context.Context, user models.User, sku models.SKU) {
			require.Equal(mc, cartItemsMock[i].SKU, sku)
			i++
		}).Return(nil)
		return cartsRepo
	}

	var (
		getCartItemsError      = errors.New("Get cart items error")
		deleteCartItemError    = errors.New("Delete cart item error")
		transactionRunnerError = errors.New("Some transaction runner error")
	)

	tests := []struct {
		Name                             string
		ExpectedResult                   *models.OrderID
		ExpectedErr                      error
		Tr                               func(mc *minimock.Controller, expectedTxErr error) *trMock.TransactionRunnerMock
		ExpectedTxErr                    error
		LOMSClient                       func(mc *minimock.Controller) *lomsClientMock.LOMSServiceClientMock
		CartsRepo                        func(mc *minimock.Controller) *cartsRepoMock.CartsRepositoryMock
		ExpectedDeleteCartItemCallsCount uint64
	}{
		{
			Name:                             "success",
			ExpectedResult:                   orderIDMock,
			ExpectedErr:                      nil,
			Tr:                               defaultTrMock,
			ExpectedTxErr:                    nil,
			LOMSClient:                       defaultLomsClientMock,
			CartsRepo:                        defaultCartRepoMock,
			ExpectedDeleteCartItemCallsCount: uint64(len(cartItemsMock)),
		},
		{
			Name:           "get cart item repo error",
			ExpectedResult: orderIDMock,
			ExpectedErr:    getCartItemsError,
			Tr:             defaultTrMock,
			ExpectedTxErr:  getCartItemsError,
			LOMSClient:     defaultLomsClientMock,
			CartsRepo: func(mc *minimock.Controller) *cartsRepoMock.CartsRepositoryMock {
				cartsRepo := cartsRepoMock.NewCartsRepositoryMock(mc)
				cartsRepo = cartsRepo.GetCartItemsMock.Return(nil, getCartItemsError)
				return cartsRepo
			},
			ExpectedDeleteCartItemCallsCount: 0,
		},
		{
			Name:           "order creation error",
			ExpectedResult: orderIDMock,
			ExpectedErr:    OrderCreationError,
			Tr:             defaultTrMock,
			ExpectedTxErr:  OrderCreationError,
			LOMSClient: func(mc *minimock.Controller) *lomsClientMock.LOMSServiceClientMock {
				lomsClient := lomsClientMock.NewLOMSServiceClientMock(mc)
				lomsClient.CreateOrderMock.Return(nil, OrderCreationError)
				return lomsClient
			},
			CartsRepo:                        defaultCartRepoMock,
			ExpectedDeleteCartItemCallsCount: 0,
		},
		{
			Name:           "delete cart item repo error",
			ExpectedResult: orderIDMock,
			ExpectedErr:    deleteCartItemError,
			Tr:             defaultTrMock,
			ExpectedTxErr:  deleteCartItemError,
			LOMSClient:     defaultLomsClientMock,
			CartsRepo: func(mc *minimock.Controller) *cartsRepoMock.CartsRepositoryMock {
				cartsRepo := cartsRepoMock.NewCartsRepositoryMock(mc)
				cartsRepo = cartsRepo.GetCartItemsMock.Return(cartItemsMock, nil)
				cartsRepo = cartsRepo.DeleteCartItemMock.Return(deleteCartItemError)
				return cartsRepo
			},
			ExpectedDeleteCartItemCallsCount: 1,
		},
		{
			Name:           "transaction error",
			ExpectedResult: orderIDMock,
			ExpectedErr:    transactionRunnerError,
			Tr: func(mc *minimock.Controller, expectedTxErr error) *trMock.TransactionRunnerMock {
				tr := trMock.NewTransactionRunnerMock(mc)
				tr = tr.RunSerializableMock.Inspect(func(ctx context.Context, txFn func(ctx context.Context) error) {
					err := txFn(ctx)
					if expectedTxErr != nil {
						require.ErrorIs(mc, expectedTxErr, err)
					} else {
						require.Equal(mc, nil, err)
					}
				}).Return(transactionRunnerError)
				return tr
			},
			ExpectedTxErr:                    nil,
			LOMSClient:                       defaultLomsClientMock,
			CartsRepo:                        defaultCartRepoMock,
			ExpectedDeleteCartItemCallsCount: uint64(len(cartItemsMock)),
		},
	}
	for _, testCase := range tests {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			mc := minimock.NewController(t)
			cartsRepo := testCase.CartsRepo(mc)
			s := New(
				testCase.Tr(mc, testCase.ExpectedTxErr),
				cartsRepo,
				testCase.LOMSClient(mc),
				productClientMock.NewProductServiceClientMock(mc),
			)
			orderID, err := s.MakePurchase(context.Background(), userIDMock)
			if testCase.ExpectedErr != nil {
				require.ErrorIs(t, err, testCase.ExpectedErr)
			} else {
				require.Equal(t, *orderID, *testCase.ExpectedResult)
				require.Equal(t, err, nil)
			}
			require.Equal(t, testCase.ExpectedDeleteCartItemCallsCount, cartsRepo.DeleteCartItemAfterCounter())
		})
	}
}
