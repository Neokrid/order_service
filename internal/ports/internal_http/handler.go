package internalhttp

import (
	"net/http"

	orders "github.com/Neokrid/order_service/internal/application"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InternalOrderHandler struct {
	service *orders.OrderService
}

func NewPublicOrderHandler(s *orders.OrderService) *InternalOrderHandler {
	return &InternalOrderHandler{
		service: s,
	}
}

func (h *InternalOrderHandler) Init(r *gin.RouterGroup) {

	orders := r.Group("/internal/orders")
	{
		orders.PATCH("/:id/status", h.UpdateStatus)
	}
}


func (h *InternalOrderHandler) UpdateStatus(c *gin.Context) {
	orderID, _ := uuid.Parse(c.Param("id"))

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.UpdateOrderStatus(c.Request.Context(), orderID, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
