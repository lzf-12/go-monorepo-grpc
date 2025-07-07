// internal/usecase/inventory.go
package usecase

import (
	"context"
	"ops-monorepo/services/svc-inventory/internal/model"
	"ops-monorepo/services/svc-inventory/internal/repository"
	grpcErr "ops-monorepo/shared-libs/grpc/errors"
	"ops-monorepo/shared-libs/logger"
	"strings"
)

type IInventoryUsecase interface {
	CheckStock(ctx context.Context, skus []string) ([]model.StockStatus, error)
	ReserveStock(ctx context.Context, orderId string, skusQuantityMap map[string]float64) (reservationHistory []model.ReservationHistory, failedToReserve []model.StockStatus, err error)
}

type inventoryUsecase struct {
	logger  logger.Logger
	repoSQL repository.IInventorySQLRepository
}

func NewInventoryUsecase(log logger.Logger, repo repository.IInventorySQLRepository) IInventoryUsecase {
	return &inventoryUsecase{
		logger:  log.WithComponent("usecase"),
		repoSQL: repo,
	}
}

func (uc *inventoryUsecase) CheckStock(ctx context.Context, skus []string) ([]model.StockStatus, error) {

	data, _, err := uc.repoSQL.CheckStockWithMultipleSkus(ctx, skus)
	if err != nil {
		uc.logger.Errorf("failed in CheckStockWithMultipleSkus", "error", err.Error())
		return nil, grpcErr.NewAppError(grpcErr.DbError, "something wrong with database: failed in CheckStockWithMultipleSkus", map[string]interface{}{"error": err.Error()})
	}

	return data, nil
}

func (uc *inventoryUsecase) ReserveStock(ctx context.Context, orderId string, skusQuantityMap map[string]float64) (stockStatus []model.ReservationHistory, failedToReserve []model.StockStatus, err error) {

	var successReservedSkus []string

	var skusArr []string
	for sku, _ := range skusQuantityMap {
		skusArr = append(skusArr, sku)
	}

	// begin db transaction
	tx, _ := uc.repoSQL.BeginTransaction(ctx)

	// loop reserve each sku
	for sku, qty := range skusQuantityMap {

		err := uc.repoSQL.ReserveStock(ctx, orderId, sku, qty)
		if err != nil {
			errmsg := err.Error()
			// rollback transaction
			uc.repoSQL.RollbackTransaction(ctx, tx)

			// handle insufficient business logic
			if strings.Contains(errmsg, "insufficient available quantity") {

				// get failed stock current status
				failedToReserve, _, err := uc.repoSQL.CheckStockWithMultipleSkus(ctx, skusArr)
				if err != nil {
					return nil, nil, grpcErr.NewAppError(grpcErr.DbError, "something wrong with db: failed in GetStockStatus", map[string]interface{}{"error": err.Error()})
				}

				// return failed stock status without app error
				uc.logger.Infof("insufficient quantity to reserve stock", "failed_to_reserve", failedToReserve)
				return nil, failedToReserve, nil
			}

			uc.logger.Errorf("something wrong with db: failed in ReserveStock", "error", errmsg)
			return nil, nil, grpcErr.NewAppError(grpcErr.DbError, "something wrong with db: failed in ReserveStock", map[string]interface{}{"error": err.Error()})
		}

		successReservedSkus = append(successReservedSkus, sku)
	}

	// commit transaction
	err = uc.repoSQL.CommitTransaction(ctx, tx)
	if err != nil {
		return nil, nil, grpcErr.NewAppError(grpcErr.DbTransactionError, "something wrong with db transaction: failed in CommitTransaction", map[string]interface{}{"error": err.Error()})
	}

	// get reservation history
	reserveHistory, err := uc.repoSQL.GetReservationHistoryByOrderIdAndstatus(ctx, orderId, model.ReservedStatus)
	if err != nil {
		uc.logger.Errorf("failed in GetReservationHistoryByOrderId", "error", err.Error())
		return nil, nil, grpcErr.NewAppError(grpcErr.DbError, "something wrong with database: failed in CheckStockWithMultipleSkus", map[string]interface{}{"error": err.Error()})
	}

	return reserveHistory, nil, nil
}

func (uc *inventoryUsecase) ReleaseStock(ctx context.Context, orderId string) (reservationHistory []model.ReservationHistory, failedToRelease []model.StockStatus, err error) {

	var successReleasedSkus []string

	// get reservation history by orderId and reserved status
	reservationHistory, err = uc.repoSQL.GetReservationHistoryByOrderIdAndstatus(ctx, orderId, model.ReservedStatus)
	if err != nil {
		uc.logger.Errorf("something wrong with db: failed in GetReservationHistoryByOrderId", "error", err.Error())
		return nil, nil, grpcErr.NewAppError(grpcErr.DbError, "something wrong with db: failed in GetReservationHistoryByOrderId", map[string]interface{}{"error": err.Error()})
	}

	// allskus
	var reservedSkusArr []string
	for _, r := range reservationHistory {
		reservedSkusArr = append(reservedSkusArr, r.Sku)
	}

	// begin db transaction
	tx, _ := uc.repoSQL.BeginTransaction(ctx)

	for _, reservedSku := range reservationHistory {

		err := uc.repoSQL.ReleaseStock(ctx, reservedSku.Sku, reservedSku.Quantity)
		if err != nil {
			// rollback transaction
			uc.repoSQL.RollbackTransaction(ctx, tx)
			errmsg := err.Error()

			// handle insufficient business logic
			if strings.Contains(errmsg, "insufficient available quantity") {

				// get failed stock current status
				failedToRelease, _, err := uc.repoSQL.CheckStockWithMultipleSkus(ctx, reservedSkusArr)
				if err != nil {
					return nil, nil, grpcErr.NewAppError(grpcErr.DbError, "something wrong with db: failed in GetStockStatus", map[string]interface{}{"error": err.Error()})
				}

				// return failed stock status without app error
				uc.logger.Infof("insufficient quantity to release stock", "failed_to_release", failedToRelease)
				return nil, failedToRelease, nil
			}

			// other than insufficient return app error
			uc.logger.Errorf("something wrong with db: failed in ReleaseStock", "error", errmsg)
			return nil, nil, grpcErr.NewAppError(grpcErr.DbError, "something wrong with db: failed in ReserveStock", map[string]interface{}{"error": err.Error()})
		}

		successReleasedSkus = append(successReleasedSkus, reservedSku.Sku)
	}

	// commit transaction
	err = uc.repoSQL.CommitTransaction(ctx, tx)
	if err != nil {
		return nil, nil, grpcErr.NewAppError(grpcErr.DbTransactionError, "something wrong with db transaction: failed in CommitTransaction", map[string]interface{}{"error": err.Error()})
	}

	// get skus inventory status
	releasedReserveHistory, err := uc.repoSQL.GetReservationHistoryByOrderIdAndstatus(ctx, orderId, model.ReleasedStatus)
	if err != nil {
		uc.logger.Errorf("failed in GetReservationHistoryByOrderId", "error", err.Error())
		return nil, nil, grpcErr.NewAppError(grpcErr.DbError, "something wrong with database: failed in CheckStockWithMultipleSkus", map[string]interface{}{"error": err.Error()})
	}

	return releasedReserveHistory, nil, nil
}
