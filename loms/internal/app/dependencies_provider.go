package app

import (
	"context"
	"log"
	"route256/libs/scheduler"
	txm "route256/libs/transactor/postgresql"
	"route256/libs/workerpool"
	notificationsClient "route256/loms/internal/clients/notifications"
	"route256/loms/internal/config"
	"route256/loms/internal/domain"
	ordersRepo "route256/loms/internal/repository/postgresql/orders"
	stocksRepo "route256/loms/internal/repository/postgresql/stocks"

	"github.com/Shopify/sarama"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DependenciesProvider struct {
	database *pgxpool.Pool

	txManager *txm.TransactionManager

	ordersRepository domain.OrdersRespository
	stocksRepository domain.StocksRespository

	ordersCancellingWorkerPool workerpool.WorkerPool
	ordersCancellingScheduler  scheduler.Scheduler

	notificationsSyncProducer sarama.SyncProducer
	notificationsWorkerPool   workerpool.WorkerPool
	notificationsClient       domain.NotificationsClient

	lomsService *domain.Service
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
			log.Fatalln("Failed to connect to notifications kafka cluster:", err)
		}
	}
	return dp.notificationsSyncProducer
}

func (dp *DependenciesProvider) GetNotificationsWorkerPool() workerpool.WorkerPool {
	if dp.notificationsWorkerPool == nil {
		dp.notificationsWorkerPool = workerpool.New(config.Data.ExternalServices.NotificationsService.MaxWorkers)
	}
	return dp.notificationsWorkerPool
}

func (dp *DependenciesProvider) GetNotificationsClient() domain.NotificationsClient {
	if dp.notificationsClient == nil {
		dp.notificationsClient = notificationsClient.New(
			dp.GetNotificationsSyncProducer(),
			config.Data.ExternalServices.NotificationsService.OrderStatusChangeNotificationsTopicName,
			dp.GetNotificationsWorkerPool(),
		)
	}
	return dp.notificationsClient
}

func (dp *DependenciesProvider) GetLOMSService(ctx context.Context) *domain.Service {
	if dp.lomsService == nil {
		dp.lomsService = domain.New(
			dp.GetTransationManager(ctx),
			dp.GetOrdersRepository(ctx),
			dp.GetStocksRepository(ctx),
			config.Data.Service.UnpaidOrderTtl,
			dp.GetOrdersCancellingScheduler(),
			dp.GetNotificationsClient(),
		)
	}
	return dp.lomsService
}

func (dp *DependenciesProvider) Close() {
	dp.GetOrdersCancellingScheduler().WaitClose()
	dp.GetOrdersCancellingWorkerPool().WaitClose()
	dp.GetNotificationsWorkerPool().Close()
}
