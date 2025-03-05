package controllers

import (
	"event-driven/src/features/orders/application"
	"event-driven/src/features/orders/domain"
	"github.com/gin-gonic/gin"
)

type CreateOrderController struct {
	createOrderUseCase *application.CreateOrderUseCase
}

func NewCreateOrderController(createOrderUseCase *application.CreateOrderUseCase) *CreateOrderController {
	return &CreateOrderController{createOrderUseCase: createOrderUseCase}
}

func (c *CreateOrderController) CreateOrder(ctx *gin.Context) {
	var order domain.Order
	if err := ctx.ShouldBindJSON(&order); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	createdOrder, err := c.createOrderUseCase.Execute(order)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to create order"})
		return
	}

	ctx.JSON(201, createdOrder)
}
