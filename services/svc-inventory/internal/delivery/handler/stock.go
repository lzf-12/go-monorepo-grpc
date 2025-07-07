package handler

import (
	"context"
	"ops-monorepo/services/svc-inventory/internal/model"
	"ops-monorepo/services/svc-inventory/internal/usecase"
	grpcErr "ops-monorepo/shared-libs/grpc/errors"
	"ops-monorepo/shared-libs/logger"
	inventoryv1 "pb_schemas/inventory/v1"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type (
	IInventoryHandler interface {
		inventoryv1.InventoryServiceServer
	}

	inventoryHandler struct {
		inventoryv1.UnimplementedInventoryServiceServer // embed the unimplemented server
		logger                                          logger.Logger
		grpcErr                                         *grpcErr.GRPCErrorHandler
		usecase                                         usecase.IInventoryUsecase
	}
)

func NewInventoryHandler(
	log logger.Logger,
	uc usecase.IInventoryUsecase,
	grpcErr *grpcErr.GRPCErrorHandler,

) IInventoryHandler {
	return &inventoryHandler{
		logger:  log,
		usecase: uc,
		grpcErr: grpcErr,
	}
}

func (h *inventoryHandler) CheckStock(ctx context.Context, req *inventoryv1.StandardInventoryRequest) (*inventoryv1.InventoryStatusResponse, error) {
	if req.Items == nil {
		return nil, grpcErr.NewValidationError("validation error", map[string]string{
			"items": "null",
		})
	}

	skus := []string{}
	mapSkuRequestedQuantityPerUom := map[string]float64{}
	if req.Items != nil {
		for _, item := range req.Items {
			skus = append(skus, item.Sku)
			mapSkuRequestedQuantityPerUom[item.Sku] = item.ReqQtyPerUom
		}
	}

	result, err := h.usecase.CheckStock(ctx, skus)
	if err != nil {
		return nil, h.grpcErr.HandleError(err)
	}
	return toProtoSuccessInventoryStatusResp(result, mapSkuRequestedQuantityPerUom), nil
}

func toProtoSuccessInventoryStatusResp(stocks []model.StockStatus, mapSkuReqQuantityPerUom map[string]float64) *inventoryv1.InventoryStatusResponse {

	var items []*inventoryv1.InventoryStatus
	for _, stock := range stocks {

		requestedQuantity, _ := mapSkuReqQuantityPerUom[stock.SKU]
		pStock := &inventoryv1.InventoryStatus{
			Sku:               stock.SKU,
			RequestedQuantity: requestedQuantity,
			AvailableQuantity: stock.AvailableQuantity,
			ReservedQuantity:  stock.ReservedQuantity,
			TotalQuantity:     stock.TotalQuantity,
			SkuUom:            stock.SKU_UOM,
			SkuPrice:          stock.SKUPrice,
			SkuCurrency:       stock.SKUCurrency,
		}

		items = append(items, pStock)
	}

	return &inventoryv1.InventoryStatusResponse{
		Items:     items,
		Timestamp: timestamppb.New(time.Now()),
	}
}

func (h *inventoryHandler) ReserveStock(ctx context.Context, req *inventoryv1.StandardInventoryRequest) (*inventoryv1.InventoryReservationResponse, error) {
	if req.Items == nil {
		return nil, grpcErr.NewValidationError("validation error", map[string]string{
			"items": "this properties cannot empty",
		})
	}

	mapSkuRequestedQuantityPerUom := map[string]float64{}
	if req.Items != nil {
		for _, item := range req.Items {
			mapSkuRequestedQuantityPerUom[item.Sku] = item.ReqQtyPerUom
		}
	}

	reservationHistory, failedReserve, err := h.usecase.ReserveStock(ctx, req.OrderId, mapSkuRequestedQuantityPerUom)
	if err == nil && failedReserve != nil {
		// give insufficient error response
		return toProtoSuccessInventoryReservationResp(nil, failedReserve, req.OrderId), nil
	}
	if err != nil && failedReserve == nil {
		// give general error reponse
		return nil, h.grpcErr.HandleError(err)
	}

	return toProtoSuccessInventoryReservationResp(reservationHistory, nil, req.OrderId), nil
}

func toProtoSuccessInventoryReservationResp(reservationHistory []model.ReservationHistory, stockStatus []model.StockStatus, orderId string) *inventoryv1.InventoryReservationResponse {

	var (
		resp             *inventoryv1.InventoryReservationResponse
		unprocessedStock inventoryv1.FailedProcessedItems
		processedStock   inventoryv1.SuccessProcessedItems
	)

	resp.OrderId = orderId

	if stockStatus != nil && len(stockStatus) > 0 {
		for _, ss := range stockStatus {

			item := &inventoryv1.InventoryStatus{
				Sku:               ss.SKU,
				AvailableQuantity: ss.AvailableQuantity,
				ReservedQuantity:  ss.ReservedQuantity,
				TotalQuantity:     ss.TotalQuantity,
				SkuUom:            ss.SKU_UOM,
				SkuPrice:          ss.SKUPrice,
				SkuCurrency:       ss.SKUCurrency,
			}
			unprocessedStock.Items = append(unprocessedStock.Items, item)
		}
	}

	if reservationHistory != nil && len(reservationHistory) > 0 {
		for _, r := range reservationHistory {

			item := &inventoryv1.ReservationHistory{
				Id:         r.Id,
				OrderId:    r.OrderId,
				Sku:        r.Sku,
				Quantity:   r.Quantity,
				Uom:        r.Uom,
				Status:     r.Status,
				ReservedAt: timestamppb.New(r.ReservedAt),
				ReleasedAt: timestamppb.New(*r.ReleasedAt),
			}

			processedStock.Items = append(processedStock.Items, item)
		}
	}

	resp.SuccessProcessedItems = &processedStock
	resp.FailedProcessedItems = &unprocessedStock

	return resp
}
