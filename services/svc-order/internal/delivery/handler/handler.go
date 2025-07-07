package handler

import (
	"errlib"
	"net/http"
	"ops-monorepo/services/svc-order/internal/delivery/types"
	uc "ops-monorepo/services/svc-order/internal/usecase"
	"ops-monorepo/services/svc-order/validator"
	"ops-monorepo/shared-libs/logger"

	"github.com/gin-gonic/gin"
)

type (
	IOrder interface {
		CreateOrder(c *gin.Context)
	}

	OrderHandler struct {
		validator  validator.IValidator
		logger     logger.Logger
		errHandler errlib.IErrorHandler
		usecase    uc.IOrderUsecase
	}
)

func NewOrderHandler(v validator.IValidator, log logger.Logger, eh errlib.IErrorHandler, uc uc.IOrderUsecase) IOrder {
	return &OrderHandler{
		validator:  v,
		logger:     log,
		errHandler: eh,
		usecase:    uc,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {

	// bind json
	var req types.OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errHandler.HandleAndSendErrorResponse(c.Writer, c.Request, errlib.ErrJSONBinding(err))
		return
	}

	// validate request
	errlist, _ := h.validator.ValidateOrderItems(req.OrderItems)
	if len(errlist) > 0 {
		h.errHandler.HandleAndSendErrorResponse(c.Writer, c.Request, errlib.ErrValidationError(errlist))
		return
	}

	// call usecase
	result, failedReserveStock, err := h.usecase.NewOrder(c.Request.Context(), req)
	if err != nil {
		h.errHandler.HandleAndSendErrorResponse(c.Writer, c.Request, errlib.ErrInternalServer(err))
		return
	}

	data := map[string]interface{}{"order": result}
	if len(failedReserveStock) > 0 {
		data["failed_reserve"] = failedReserveStock
	}

	// log and send success response
	h.logger.Info("order created with pending status")
	c.JSON(http.StatusCreated, types.CreateOrderSuccessResponse{
		Data:       data,
		StatusCode: http.StatusCreated,
		Message:    "order created with pending status",
	})
}
