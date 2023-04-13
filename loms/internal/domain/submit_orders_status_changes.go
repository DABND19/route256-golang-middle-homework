package domain

import (
	"context"
	"route256/libs/logger"
	"route256/loms/internal/models"
	"time"

	"go.uber.org/zap"
)

func (s *Service) runOrdersStatusChangesSubmission(ctx context.Context) {
	ticker := time.NewTicker(s.ordersStatusChangesSumbitInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				s.submitOrdersStatusChanges(ctx)
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

func (s *Service) submitOrdersStatusChanges(ctx context.Context) {
	var records []models.OrderStatusChange
	err := s.RunReadCommited(ctx, func(ctx context.Context) error {
		var err error
		records, err = s.OrderStatusChangeRepository.GetUnsubmittedChanges(ctx)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logger.Error("Failed to query unsubmitted orders status changes.", zap.Error(err))
		return
	}

	for _, item := range records {
		err := s.NotificationsClient.NotifyAboutOrderStatusChange(ctx, item)
		if err != nil {
			logger.Error("Failed to send notification.", zap.Error(err))
			continue
		}

		err = s.RunReadCommited(ctx, func(ctx context.Context) error {
			return s.OrderStatusChangeRepository.MarkChangeAsSubmitted(ctx, item)
		})
		if err != nil {
			logger.Error("Failed to mark notification as submitted.", zap.Error(err))
		}
	}
}
