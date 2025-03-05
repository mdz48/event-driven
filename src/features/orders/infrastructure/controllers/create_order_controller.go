package controllers

import (
	"event-driven/src/features/orders/application"
	"event-driven/src/features/orders/domain"
	"event-driven/src/features/orders/infrastructure/producer"
	"github.com/gin-gonic/gin"
)

type CreateOrderController struct {
	createOrderUseCase *application.CreateOrderUseCase
	rabbitMQ           *producer.RabbitMQ
}

func NewCreateOrderController(createOrderUseCase *application.CreateOrderUseCase, rabbitMQ *producer.RabbitMQ) *CreateOrderController {
	return &CreateOrderController{
		createOrderUseCase: createOrderUseCase,
		rabbitMQ:           rabbitMQ,
	}
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

	// Publicar evento si RabbitMQ está disponible
	if c.rabbitMQ != nil {
		err = c.rabbitMQ.NotifyOrderCreated(createdOrder)
		if err != nil {
			// Solo registrar el error, no afecta la respuesta al cliente
			ctx.Error(err)
		}
	}

	ctx.JSON(201, createdOrder)
}