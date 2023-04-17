package app

import (
	"context"
	lomsClient "route256/checkout/internal/clients/loms"
	productClient "route256/checkout/internal/clients/product"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/models"
	cartsRepo "route256/checkout/internal/repository/postgresql/carts"
	"route256/libs/cachemetrics"
	"route256/libs/logger"
	"route256/libs/lrucache"
	txm "route256/libs/transactor/postgresql"
	"route256/libs/workerpool"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type DependenciesProvider struct {
	database *pgxpool.Pool

	txManager *txm.TransactionManager

	cartsRepository domain.CartsRepository

	lomsServiceClient domain.LOMSServiceClient

	productServiceClientWorkerPool workerpool.WorkerPool
	productServiceClientCache      *cachemetrics.CacheMetrics[models.SKU, *models.Product]
	productServiceClient           domain.ProductServiceClient

	checkoutService *domain.Service
}

func NewDependenciesProvider() *DependenciesProvider {
	return &DependenciesProvider{}
}

func (dp *DependenciesProvider) GetDatabase(ctx context.Context) *pgxpool.Pool {
	if dp.database == nil {
		dbConfig, err := pgxpool.ParseConfig(config.Data.Postgres.DSN)
		if err != nil {
			logger.Fatal("Failed to parse database config.", zap.Error(err))
		}
		dbConfig.ConnConfig.PreferSimpleProtocol = true

		dp.database, err = pgxpool.ConnectConfig(ctx, dbConfig)
		if err != nil {
			logger.Fatal("Failed to connect to database.", zap.Error(err))
		}
		if err := dp.database.Ping(ctx); err != nil {
			logger.Fatal("Failed to ping database.", zap.Error(err))
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
			logger.Fatal("Failed to connect to loms service.", zap.Error(err))
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

func (dp *DependenciesProvider) GetProductServiceClientCache() *cachemetrics.CacheMetrics[models.SKU, *models.Product] {
	if dp.productServiceClientCache == nil {
		cache := lrucache.New[models.SKU, *models.Product](
			config.Data.ExternalServices.Product.CacheMaxSize,
		)
		dp.productServiceClientCache = cachemetrics.New[models.SKU, *models.Product](cache, "product_service")
	}
	return dp.productServiceClientCache
}

func (dp *DependenciesProvider) GetProductServiceClient() domain.ProductServiceClient {
	if dp.productServiceClient == nil {
		var err error
		dp.productServiceClient, err = productClient.New(
			config.Data.ExternalServices.Product.Url,
			config.Data.ExternalServices.Product.AccessToken,
			int(config.Data.ExternalServices.Product.RateLimit),
			dp.GetProductServiceClientWorkerPool(),
			dp.GetProductServiceClientCache(),
		)
		if err != nil {
			logger.Fatal("Failed to connect to product service.", zap.Error(err))
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
