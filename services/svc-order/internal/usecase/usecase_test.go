package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/robaho/fixed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"ops-monorepo/services/svc-order/internal/delivery/types"
	"ops-monorepo/services/svc-order/internal/model"
	"ops-monorepo/services/svc-order/mocks"
	grpcMocks "ops-monorepo/shared-libs/grpc/client/mocks"
	loggerMocks "ops-monorepo/shared-libs/logger/mocks"
	inventoryv1 "pb_schemas/inventory/v1"
)

type usecaseDeps struct {
	logger              loggerMocks.MockLogger
	repoSQL             mocks.MockIOrderSQLRepository
	inventoryGrpcClient grpcMocks.MockInvClient
}

var (
	mockOrderId   = uuid.MustParse("9680e493-843d-4069-9b38-7495e70d7621")
	mockUserId    = "9ae74d58-7cb4-408d-bac0-8c5471a23062"
	mockUserEmail = "user@email.com"
	mockOrder     = model.Order{
		Id:          mockOrderId,
		UserId:      mockUserId,
		UserEmail:   mockUserEmail,
		Status:      model.ORDER_STATUS_PENDING,
		TotalAmount: fixed.NewS("100"),
		Currency:    "USD",
		CreatedAt:   time.Now(),
	}
	mockItems = []model.ItemOrder{
		{
			Id:             uuid.New(),
			OrderId:        mockOrderId,
			Sku:            "OLIVE-OIL-1L",
			QuantityPerUom: fixed.NewS("0.5"),
			PricePerUom:    fixed.NewS("50"),
			UomCode:        "L",
		},
		{
			Id:             uuid.New(),
			OrderId:        mockOrderId,
			Sku:            "TSHIRT-M-WHITE",
			QuantityPerUom: fixed.NewS("2"),
			PricePerUom:    fixed.NewS("25"),
			UomCode:        "EA",
		},
	}
	mockStockResponse = &inventoryv1.InventoryStatusResponse{
		Items: []*inventoryv1.InventoryStatus{
			{
				Sku:               "OLIVE-OIL-1L",
				RequestedQuantity: 0.5,
				SkuPrice:          50,
				SkuUom:            "L",
			},
			{
				Sku:               "TSHIRT-M-WHITE",
				RequestedQuantity: 2,
				SkuPrice:          25,
				SkuUom:            "EA",
			},
		},
	}
	mockReserveSuccessResponse = &inventoryv1.InventoryReservationResponse{
		FailedProcessedItems: &inventoryv1.FailedProcessedItems{
			Items: []*model.OrderedItemStockStatus{},
		},
	}
	mockReserveFailedResponse = &inventoryv1.InventoryReservationResponse{
		FailedProcessedItems: &inventoryv1.FailedProcessedItems{
			Items: []*model.OrderedItemStockStatus{
				{
					Sku:               "OLIVE-OIL-1L",
					AvailableQuantity: 1,
					RequestedQuantity: 3,
					ReservedQuantity:  2,
				},
			},
		},
	}
)

func TestOrderUsecase_NewOrder(t *testing.T) {
	type args struct {
		ctx     context.Context
		request types.OrderRequest
	}

	testCases := []struct {
		Name        string
		Args        args
		Mock        func(dep *usecaseDeps)
		ExpectedErr bool
		Expected    *model.OrderWithItems
	}{
		{
			Name: "successful order creation with stock reservation",
			Args: args{
				ctx: context.Background(),
				request: types.OrderRequest{
					OrderItems: []types.StockItemRequest{
						{
							Sku:            "OLIVE-OIL-1L",
							QuantityPerUom: 0.5,
							Uom:            "L",
						},
						{
							Sku:            "TSHIRT-M-WHITE",
							QuantityPerUom: 2,
							Uom:            "EA",
						},
					},
				},
			},
			Mock: func(dep *usecaseDeps) {
				dep.inventoryGrpcClient.EXPECT().CheckStock(mock.Anything, mock.Anything).
					Return(mockStockResponse, nil)
				dep.repoSQL.EXPECT().InsertOrderWithItems(mock.Anything, mock.AnythingOfType("*model.Order"), mock.AnythingOfType("[]model.ItemOrder")).
					Return(nil)
				dep.inventoryGrpcClient.EXPECT().ReserveStock(mock.Anything, mock.Anything).
					Return(mockReserveSuccessResponse, nil)
				dep.repoSQL.EXPECT().UpdateOrderStatus(mock.Anything, mock.AnythingOfType("uuid.UUID"), model.ORDER_STATUS_CONFIRMED).
					Return(nil)
			},
			ExpectedErr: false,
			Expected: &model.OrderWithItems{
				Order: model.Order{
					Status:      model.ORDER_STATUS_PENDING,
					UserId:      "80212e88-5d6d-446b-981c-5dfdc426e867",
					UserEmail:   "dummy@email.com",
					Currency:    "USD",
					TotalAmount: fixed.NewS("75"),
				},
			},
		},
		{
			Name: "failed stock check",
			Args: args{
				ctx: context.Background(),
				request: types.OrderRequest{
					OrderItems: []types.StockItemRequest{
						{
							Sku:            "INVALID-SKU",
							QuantityPerUom: 1,
							Uom:            "EA",
						},
					},
				},
			},
			Mock: func(dep *usecaseDeps) {
				dep.inventoryGrpcClient.EXPECT().CheckStock(mock.Anything, mock.Anything).
					Return(nil, errors.New("inventory service error"))
				dep.logger.EXPECT().Errorf("failed check stock to inventory service", mock.Anything)
			},
			ExpectedErr: true,
			Expected:    nil,
		},
		{
			Name: "failed to insert order with items",
			Args: args{
				ctx: context.Background(),
				request: types.OrderRequest{
					OrderItems: []types.StockItemRequest{
						{
							Sku:            "OLIVE-OIL-1L",
							QuantityPerUom: 0.5,
							Uom:            "L",
						},
					},
				},
			},
			Mock: func(dep *usecaseDeps) {
				dep.inventoryGrpcClient.EXPECT().CheckStock(mock.Anything, mock.Anything).
					Return(&inventoryv1.InventoryStatusResponse{
						Items: []*inventoryv1.InventoryStatus{
							{
								Sku:               "OLIVE-OIL-1L",
								RequestedQuantity: 0.5,
								SkuPrice:          50,
								SkuUom:            "L",
							},
						},
					}, nil)
				dep.repoSQL.EXPECT().InsertOrderWithItems(mock.Anything, mock.AnythingOfType("*model.Order"), mock.AnythingOfType("[]model.ItemOrder")).
					Return(errors.New("database error"))
				dep.logger.EXPECT().Errorf("failed in InsertOrderWithItems", mock.Anything)
			},
			ExpectedErr: true,
			Expected:    nil,
		},
		{
			Name: "stock reservation failed - insufficient stock",
			Args: args{
				ctx: context.Background(),
				request: types.OrderRequest{
					OrderItems: []types.StockItemRequest{
						{
							Sku:            "OLIVE-OIL-1L",
							QuantityPerUom: 0.5,
							Uom:            "L",
						},
					},
				},
			},
			Mock: func(dep *usecaseDeps) {
				dep.inventoryGrpcClient.EXPECT().CheckStock(mock.Anything, mock.Anything).
					Return(&inventoryv1.InventoryStatusResponse{
						Items: []*inventoryv1.InventoryStatus{
							{
								Sku:               "OLIVE-OIL-1L",
								RequestedQuantity: 0.5,
								SkuPrice:          50,
								SkuUom:            "L",
							},
						},
					}, nil)
				dep.repoSQL.EXPECT().InsertOrderWithItems(mock.Anything, mock.AnythingOfType("*model.Order"), mock.AnythingOfType("[]model.ItemOrder")).
					Return(nil)
				// insufficient stock scenario
				dep.inventoryGrpcClient.EXPECT().ReserveStock(mock.Anything, mock.Anything).
					Return(mockReserveFailedResponse, nil)
				dep.repoSQL.EXPECT().UpdateOrderStatus(mock.Anything, mock.AnythingOfType("uuid.UUID"), model.ORDER_STATUS_FAILED_RESERVATION).
					Return(nil)
			},
			ExpectedErr: false,
			Expected: &model.OrderWithItems{
				Order: model.Order{
					Status:    model.ORDER_STATUS_PENDING,
					UserId:    "80212e88-5d6d-446b-981c-5dfdc426e867",
					UserEmail: "dummy@email.com",
					Currency:  "USD",
				},
			},
		},
		{
			Name: "stock reservation service error",
			Args: args{
				ctx: context.Background(),
				request: types.OrderRequest{
					OrderItems: []types.StockItemRequest{
						{
							Sku:            "OLIVE-OIL-1L",
							QuantityPerUom: 0.5,
							Uom:            "L",
						},
					},
				},
			},
			Mock: func(dep *usecaseDeps) {
				dep.inventoryGrpcClient.EXPECT().CheckStock(mock.Anything, mock.Anything).
					Return(&inventoryv1.InventoryStatusResponse{
						Items: []*inventoryv1.InventoryStatus{
							{
								Sku:               "OLIVE-OIL-1L",
								RequestedQuantity: 0.5,
								SkuPrice:          50,
								SkuUom:            "L",
							},
						},
					}, nil)
				dep.repoSQL.EXPECT().InsertOrderWithItems(mock.Anything, mock.AnythingOfType("*model.Order"), mock.AnythingOfType("[]model.ItemOrder")).
					Return(nil)
				// service error during reservation
				dep.inventoryGrpcClient.EXPECT().ReserveStock(mock.Anything, mock.Anything).
					Return(nil, errors.New("inventory service error"))
				dep.repoSQL.EXPECT().UpdateOrderStatus(mock.Anything, mock.AnythingOfType("uuid.UUID"), model.ORDER_STATUS_CANCELLED).
					Return(nil)
			},
			ExpectedErr: true,
			Expected:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			mockRepo := mocks.NewMockIOrderSQLRepository(t)
			mockLogger := loggerMocks.NewMockLogger(t)
			mockInvClient := grpcMocks.NewMockInvClient(t)

			deps := usecaseDeps{
				logger:              *mockLogger,
				repoSQL:             *mockRepo,
				inventoryGrpcClient: *mockInvClient,
			}

			tc.Mock(&deps)

			usecase := NewOrderUsecase(&deps.repoSQL, &deps.logger, &deps.inventoryGrpcClient)
			result, failedItems, err := usecase.NewOrder(tc.Args.ctx, tc.Args.request)

			if tc.ExpectedErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				if tc.Expected != nil {
					assert.NotNil(t, result)
					assert.Equal(t, tc.Expected.Order.Status, result.Order.Status)
					assert.Equal(t, tc.Expected.Order.UserId, result.Order.UserId)
					assert.Equal(t, tc.Expected.Order.UserEmail, result.Order.UserEmail)
					assert.Equal(t, tc.Expected.Order.Currency, result.Order.Currency)

					// verify failed items for insufficient stock scenario
					if tc.Name == "stock reservation failed - insufficient stock" {
						assert.NotNil(t, failedItems)
						assert.Len(t, failedItems, 1)
					}
				}
			}
		})
	}
}
