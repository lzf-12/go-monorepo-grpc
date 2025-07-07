package usecase

import (
	"context"
	"errlib"
	inventoryv1 "pb_schemas/inventory/v1"

	// "internal/runtime/math"
	"time"

	// "ops-monorepo/services/svc-order/internal"

	"ops-monorepo/services/svc-order/internal/delivery/types"
	"ops-monorepo/services/svc-order/internal/model"
	"ops-monorepo/services/svc-order/internal/repository"
	grpc "ops-monorepo/shared-libs/grpc/client"
	"ops-monorepo/shared-libs/logger"

	"github.com/google/uuid"
	"github.com/robaho/fixed"
)

type (
	IOrderUsecase interface {
		NewOrder(ctx context.Context, request types.OrderRequest) (*model.OrderWithItems, []*model.OrderedItemStockStatus, error)
	}

	OrderUsecase struct {
		logger              logger.Logger
		repoSQL             repository.IOrderSQLRepository
		inventoryGrpcClient grpc.InvClient
	}
)

func NewOrderUsecase(sql repository.IOrderSQLRepository, log logger.Logger, invClient grpc.InvClient) IOrderUsecase {
	return &OrderUsecase{
		logger:              log,
		repoSQL:             sql,
		inventoryGrpcClient: invClient,
	}
}

func (u *OrderUsecase) NewOrder(ctx context.Context, request types.OrderRequest) (*model.OrderWithItems, []*model.OrderedItemStockStatus, error) {

	// check stock
	var InventoryItems []*inventoryv1.InventoryItem
	for _, item := range request.OrderItems {
		invItem := &inventoryv1.InventoryItem{
			Sku:          item.Sku,
			ReqQtyPerUom: item.QuantityPerUom,
			Uom:          item.Uom,
		}
		InventoryItems = append(InventoryItems, invItem)
	}

	stockStatus, err := u.inventoryGrpcClient.CheckStock(ctx, &inventoryv1.StandardInventoryRequest{
		Items: InventoryItems,
	})
	// handle error sku not found here

	if err != nil {

		// handle error

		u.logger.Errorf("failed check stock to inventory service", map[string]interface{}{"error": err})
		return nil, nil, errlib.ErrInternalServer(err)
	}

	// makesure quantity available

	// prepare order data
	useriddummy := "80212e88-5d6d-446b-981c-5dfdc426e867"
	orderId := uuid.New()
	order := model.Order{
		Id:        orderId,
		Status:    model.ORDER_STATUS_PENDING,
		CreatedAt: time.Now(),
		UserId:    useriddummy,
		UserEmail: "dummy@email.com",
		Currency:  "USD",
	}
	var items []model.ItemOrder
	total := fixed.NewF(0)

	for _, i := range stockStatus.Items {
		item := model.ItemOrder{
			Id:             uuid.New(),
			OrderId:        orderId,
			Sku:            i.Sku,
			QuantityPerUom: fixed.NewF(i.RequestedQuantity),
			PricePerUom:    fixed.NewF(i.SkuPrice),
			UomCode:        i.SkuUom,
		}
		items = append(items, item)

		amount := item.QuantityPerUom.Mul(item.PricePerUom)
		total = total.Add(amount)
	}
	order.TotalAmount = total

	// insert order with pending status
	err = u.repoSQL.InsertOrderWithItems(ctx, &order, items)
	if err != nil {
		u.logger.Errorf("failed in InsertOrderWithItems", "error", err)
		return nil, nil, err
	}

	// reserve stock
	reserveResp, errReserv := u.inventoryGrpcClient.ReserveStock(ctx, &inventoryv1.StandardInventoryRequest{})
	var failedReserveStockStatus []*model.OrderedItemStockStatus
	if errReserv != nil && reserveResp != nil {
		failedReserveStockStatus = reserveResp.FailedProcessedItems.GetItems()
	}

	// handle failed to reserve caused by insufficient, with no app error
	if errReserv == nil && len(reserveResp.FailedProcessedItems.Items) > 0 {

		// update order status to failed reservation
		if err := u.repoSQL.UpdateOrderStatus(ctx, orderId, model.ORDER_STATUS_FAILED_RESERVATION); err != nil {
			u.logger.Errorf("failed update order status to reserved", "error", err.Error())
			return nil, reserveResp.FailedProcessedItems.Items, errlib.ErrDBQuery()
		}

		return &model.OrderWithItems{Order: order, Items: items}, reserveResp.FailedProcessedItems.Items, nil
	}
	if errReserv != nil {
		// update order to canceled
		if err := u.repoSQL.UpdateOrderStatus(ctx, orderId, model.ORDER_STATUS_CANCELLED); err != nil {
			u.logger.Errorf("failed update order status to reserved", "error", err.Error())
			return nil, nil, errlib.ErrDBQuery()
		}
		return nil, nil, errlib.ErrReservationStock(errReserv)
	}

	// update order status to reserved
	if err = u.repoSQL.UpdateOrderStatus(ctx, orderId, model.ORDER_STATUS_CONFIRMED); err != nil {
		u.logger.Errorf("failed update order status to reserved", "error", err.Error())
		return nil, nil, errlib.ErrDBQuery()
	}

	return &model.OrderWithItems{
		Order: order,
		Items: items,
	}, failedReserveStockStatus, nil
}
