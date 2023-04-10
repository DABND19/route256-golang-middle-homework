package app

import (
	"context"
	"route256/libs/logger"
	"route256/libs/scheduler"
	txm "route256/libs/transactor/postgresql"
	"route256/libs/workerpool"
	notificationsClient "route256/loms/internal/clients/notifications"
	"route256/loms/internal/config"
	"route256/loms/internal/domain"
	ordersRepo "route256/loms/internal/repository/postgresql/orders"
	ordersStatusChangesRepo "route256/loms/internal/repository/postgresql/ordersstatuschanges"
	stocksRepo "route256/loms/internal/repository/postgresql/stocks"

	"github.com/Shopify/sarama"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type DependenciesProvider struct {
	database *pgxpool.Pool

	txManager *txm.TransactionManager

	ordersRepository              domain.OrdersRespository
	stocksRepository              domain.StocksRespository
	ordersStatusChangesRepository domain.OrderStatusChangeRepository

	ordersCancellingWorkerPool workerpool.WorkerPool
	ordersCancellingScheduler  scheduler.Scheduler

	notificationsSyncProducer sarama.SyncProducer
	notificationsClient       domain.NotificationsClient

	lomsService *domain.Service
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

func (dp *DependenciesProvider) GetOrdersRepository(ctx context.Context) domain.OrdersRespository {
	if dp.ordersRepository == nil {
		dp.ordersRepository = ordersRepo.New(dp.GetTransationManager(ctx))
	}
	return dp.ordersRepository
}

func (dp *DependenciesProvider) GetStocksRepository(ctx context.Context) domain.StocksRespository {
	if dp.stocksRepository == nil {
		dp.stocksRepository = stocksRepo.New(dp.GetTransationManager(ctx))
	}
	return dp.stocksRepository
}

func (dp *DependenciesProvider) GetOrdersStatusChangesRepository(ctx context.Context) domain.OrderStatusChangeRepository {
	if dp.ordersStatusChangesRepository == nil {
		dp.ordersStatusChangesRepository = ordersStatusChangesRepo.New(dp.GetTransationManager(ctx))
	}
	return dp.ordersStatusChangesRepository
}

func (dp *DependenciesProvider) GetOrdersCancellingWorkerPool() workerpool.WorkerPool {
	if dp.ordersCancellingWorkerPool == nil {
		dp.ordersCancellingWorkerPool = workerpool.New(config.Data.Service.UnpaidOrdersCancellingWorkersCount)
	}
	return dp.ordersCancellingWorkerPool
}

func (dp *DependenciesProvider) GetOrdersCancellingScheduler() scheduler.Scheduler {
	if dp.ordersCancellingScheduler == nil {
		dp.ordersCancellingScheduler = scheduler.New(dp.GetOrdersCancellingWorkerPool())
	}
	return dp.ordersCancellingScheduler
}

func (dp *DependenciesProvider) GetNotificationsSyncProducer() sarama.SyncProducer {
	if dp.notificationsSyncProducer == nil {
		producerConfig := sarama.NewConfig()
		producerConfig.Producer.Idempotent = true
		producerConfig.Producer.RequiredAcks = sarama.WaitForAll
		producerConfig.Producer.Partitioner = sarama.NewHashPartitioner
		producerConfig.Producer.Return.Successes = true
		producerConfig.Net.MaxOpenRequests = 1

		var err error
		dp.notificationsSyncProducer, err = sarama.NewSyncProducer(config.Data.ExternalServices.NotificationsService.KafkaBrokers, producerConfig)
		if err != nil {
			logger.Fatal("Failed to connect to notifications kafka cluster.", zap.Error(err))
		}
	}
	return dp.notificationsSyncProducer
}

func (dp *DependenciesProvider) GetNotificationsClient() domain.NotificationsClient {
	if dp.notificationsClient == nil {
		dp.notificationsClient = notificationsClient.New(
			dp.GetNotificationsSyncProducer(),
			config.Data.ExternalServices.NotificationsService.OrderStatusChangeNotificationsTopicName,
		)
	}
	return dp.notificationsClient
}

func (dp *DependenciesProvider) GetLOMSService(ctx context.Context) *domain.Service {
	if dp.lomsService == nil {
		dp.lomsService = domain.New(
			ctx,
			dp.GetTransationManager(ctx),
			dp.GetOrdersRepository(ctx),
			dp.GetStocksRepository(ctx),
			config.Data.Service.UnpaidOrderTtl,
			dp.GetOrdersCancellingScheduler(),
			dp.GetNotificationsClient(),
			dp.GetOrdersStatusChangesRepository(ctx),
			config.Data.Service.OrdersStatuschangesSubmissionInterval,
		)
	}
	return dp.lomsService
}

func (dp *DependenciesProvider) Close() {
	dp.GetOrdersCancellingScheduler().WaitClose()
	dp.GetOrdersCancellingWorkerPool().WaitClose()
}
