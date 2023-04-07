package app

import (
	"context"
	"log"
	lomsClient "route256/checkout/internal/clients/loms"
	productClient "route256/checkout/internal/clients/product"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	cartsRepo "route256/checkout/internal/repository/postgresql/carts"
	txm "route256/libs/transactor/postgresql"
	"route256/libs/workerpool"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DependenciesProvider struct {
	database *pgxpool.Pool

	txManager *txm.TransactionManager

	cartsRepository domain.CartsRepository

	lomsServiceClient domain.LOMSServiceClient

	productServiceClientWorkerPool workerpool.WorkerPool
	productServiceClient           domain.ProductServiceClient

	checkoutService *domain.Service
}

func NewDependenciesProvider() *DependenciesProvider {
	return &DependenciesProvider{}
}

func (dp *DependenciesProvider) GetDatabase(ctx context.Context) *pgxpool.Pool {
	if dp.database == nil {
		var err error
		dp.database, err = pgxpool.Connect(ctx, config.Data.Postgres.DSN)
		if err != nil {
			log.Fatalln("Failed to connect to database:", err)
		}
		if err := dp.database.Ping(ctx); err != nil {
			log.Fatalln("Failed to ping database:", err)
		}
	}
	return dp.database
}

func (dp *DependenciesProvider) GetTransationManager(ctx context.Context) *txm.TransactionManager {
	if dp.txManager == nil {
		dp.txManager = txm.New(dp.GetDatabase(ctx))
	}
	return dp.txManager
}

func (dp *DependenciesProvider) GetCartsRepository(ctx context.Context) domain.CartsRepository {
	if dp.cartsRepository == nil {
		dp.cartsRepository = cartsRepo.New(dp.GetTransationManager(ctx))
	}
	return dp.cartsRepository
}

func (dp *DependenciesProvider) GetLomsServiceClient() domain.LOMSServiceClient {
	if dp.lomsServiceClient == nil {
		var err error
		dp.lomsServiceClient, err = lomsClient.New(config.Data.ExternalServices.Loms.Url)
		if err != nil {
			log.Fatalln("Failed to connect to loms service:", err)
		}
	}
	return dp.lomsServiceClient
}

func (dp *DependenciesProvider) GetProductServiceClientWorkerPool() workerpool.WorkerPool {
	if dp.productServiceClientWorkerPool == nil {
		dp.productServiceClientWorkerPool = workerpool.New(
			config.Data.ExternalServices.Product.MaxConcurrentRequests,
		)
	}
	return dp.productServiceClientWorkerPool
}

func (dp *DependenciesProvider) GetProductServiceClient() domain.ProductServiceClient {
	if dp.productServiceClient == nil {
		var err error
		dp.productServiceClient, err = productClient.New(
			config.Data.ExternalServices.Product.Url,
			config.Data.ExternalServices.Product.AccessToken,
			int(config.Data.ExternalServices.Product.RateLimit),
			dp.GetProductServiceClientWorkerPool(),
		)
		if err != nil {
			log.Fatalln("Failed to connect to product service:", err)
		}
	}
	return dp.productServiceClient
}

func (dp *DependenciesProvider) GetCheckoutService(ctx context.Context) *domain.Service {
	if dp.checkoutService == nil {
		dp.checkoutService = domain.New(
			dp.GetTransationManager(ctx),
			dp.GetCartsRepository(ctx),
			dp.GetLomsServiceClient(),
			dp.GetProductServiceClient(),
		)
	}
	return dp.checkoutService
}

func (dp *DependenciesProvider) Close() {
	dp.productServiceClientWorkerPool.WaitClose()
}
