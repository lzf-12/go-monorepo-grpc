package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errlib"
	em "errlib/mocks"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/robaho/fixed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/mock/gomock"

	"ops-monorepo/services/svc-order/internal/delivery/types"
	"ops-monorepo/services/svc-order/internal/model"
	"ops-monorepo/services/svc-order/mocks"
	ml "ops-monorepo/shared-libs/logger/mocks"
)

type handlerDeps struct {
	validator mocks.MockIValidator
	usecase   mocks.MockIOrderUsecase
	logger    ml.MockLogger
	errLib    em.MockIErrorHandler

	// http
	ginWriterRsp *gin.ResponseWriter
	req          *http.Request
}

var mockOrderId = "9680e493-843d-4069-9b38-7495e70d7621"
var mockUserId = "9ae74d58-7cb4-408d-bac0-8c5471a23062"
var mockUserEmail = "user@email.com"
var mockResultUsecase = model.OrderWithItems{
	Order: model.Order{
		Id:          uuid.MustParse(mockOrderId),
		UserId:      mockUserId,
		UserEmail:   mockUserEmail,
		Status:      model.ORDER_STATUS_PENDING,
		TotalAmount: fixed.NewS("100"),
		Currency:    "USD",
	},
}
var mockSuccessResponse = map[string]interface{}{
	"data": mockResultUsecase,
}
var noValidationError = []map[string]interface{}{}

func TestOrderHandler_CreateOrder(t *testing.T) {

	gin.SetMode(gin.TestMode)

	type args struct {
		context context.Context
	}

	testCases := []struct {
		Name       string
		Payload    interface{}
		Mock       func(dep *handlerDeps, w http.ResponseWriter, r *http.Request)
		StatusCode int
	}{
		{
			Name: "valid order creation",
			Payload: types.PostOrdersJSONRequestBody{
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
			Mock: func(dep *handlerDeps, w http.ResponseWriter, r *http.Request) {

				dep.validator.EXPECT().ValidateOrderItems(mock.Anything).Return(noValidationError, nil)
				dep.usecase.EXPECT().NewOrder(mock.Anything, mock.Anything).Return(&mockResultUsecase, nil, nil)
				dep.logger.EXPECT().Info("order created with pending status")
			},
			StatusCode: http.StatusCreated,
		},
		{
			Name:    "invalid json binding",
			Payload: `{"invalid": "json"`,
			Mock: func(dep *handlerDeps, w http.ResponseWriter, r *http.Request) {
				dep.errLib.EXPECT().HandleAndSendErrorResponse(
					mock.Anything,
					mock.AnythingOfType("*http.Request"),
					mock.MatchedBy(func(err *errlib.AppError) bool {
						return err != nil &&
							err.Status == http.StatusBadRequest &&
							strings.Contains(err.Message, "Invalid JSON request unable to bind")
					}),
				).Times(1).Run(func(args mock.Arguments) {
					if w, ok := args.Get(0).(http.ResponseWriter); ok {
						if err, ok := args.Get(2).(*errlib.AppError); ok {
							w.WriteHeader(err.Status)
						}
					}
				})
			},
			StatusCode: http.StatusBadRequest,
		},
		{
			Name: "validation error",
			Payload: types.PostOrdersJSONRequestBody{
				OrderItems: []types.StockItemRequest{
					{
						Sku:            "TSHIRT-M-WHITE",
						QuantityPerUom: 2,
						Uom:            "EA",
					},
					{
						Sku:            "TSHIRT-M-WHITE",
						QuantityPerUom: 2,
						Uom:            "EA",
					},
				},
			},
			Mock: func(dep *handlerDeps, w http.ResponseWriter, r *http.Request) {
				dep.validator.EXPECT().ValidateOrderItems(mock.Anything).Return([]map[string]interface{}{
					{
						"sku": "this row should be unique",
						"row": "2",
					},
				}, nil)
				dep.errLib.EXPECT().HandleAndSendErrorResponse(
					mock.Anything,
					mock.AnythingOfType("*http.Request"),
					mock.MatchedBy(func(err *errlib.AppError) bool {
						return err != nil &&
							err.Status == http.StatusBadRequest &&
							strings.Contains(err.Message, "Validation error")
					}),
				).Times(1).Run(func(args mock.Arguments) {
					if w, ok := args.Get(0).(http.ResponseWriter); ok {
						if err, ok := args.Get(2).(*errlib.AppError); ok {
							w.WriteHeader(err.Status)
						}
					}
				})
			},
			StatusCode: http.StatusBadRequest,
		},
		{
			Name: "failed usecase",
			Payload: types.PostOrdersJSONRequestBody{
				OrderItems: []types.StockItemRequest{
					{
						Sku:            "TSHIRT-M-WHITE",
						QuantityPerUom: 2,
						Uom:            "EA",
					},
				},
			},
			Mock: func(dep *handlerDeps, w http.ResponseWriter, r *http.Request) {
				dep.validator.EXPECT().ValidateOrderItems(mock.Anything).Return(noValidationError, nil)
				dep.usecase.EXPECT().NewOrder(mock.Anything, mock.Anything).Return(nil, nil, errors.New("error"))
				dep.errLib.EXPECT().HandleAndSendErrorResponse(
					mock.Anything,
					mock.AnythingOfType("*http.Request"),
					mock.MatchedBy(func(err *errlib.AppError) bool {
						return err != nil &&
							err.Status == http.StatusInternalServerError &&
							strings.Contains(err.Message, "An internal server error occurred")
					}),
				).Times(1).Run(func(args mock.Arguments) {
					if w, ok := args.Get(0).(http.ResponseWriter); ok {
						if err, ok := args.Get(2).(*errlib.AppError); ok {
							w.WriteHeader(err.Status)
						}
					}
				})
			},
			StatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			path := "/v1/api/order"

			// Setup mocks
			mockValidator := mocks.NewMockIValidator(t)
			mockUsecase := mocks.NewMockIOrderUsecase(t)
			mockLogger := ml.NewMockLogger(t)
			mockerrlib := em.NewMockIErrorHandler(t)

			deps := handlerDeps{
				validator: *mockValidator,
				usecase:   *mockUsecase,
				logger:    *mockLogger,
				errLib:    *mockerrlib,
			}

			// Prepare request payload BEFORE calling Mock
			var payloadBytes []byte
			if strPayload, ok := tc.Payload.(string); ok {
				// Handle invalid JSON string
				payloadBytes = []byte(strPayload)
			} else {
				payloadBytes, _ = json.Marshal(tc.Payload)
			}

			req, _ := http.NewRequest(http.MethodPost, path, bytes.NewBuffer(payloadBytes))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			tc.Mock(&deps, resp, req)

			// Create handler with mocked dependencies
			handler := NewOrderHandler(&deps.validator, &deps.logger, &deps.errLib, &deps.usecase)

			// Setup Gin router
			r := gin.Default()
			r.POST(path, handler.CreateOrder)

			// Execute
			r.ServeHTTP(resp, req)

			// Verify response
			assert.Equal(t, tc.StatusCode, resp.Code)
		})
	}
}
