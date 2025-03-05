package controllers

import (
	"event-driven/src/features/orders/application"
	"github.com/gin-gonic/gin"
)

type GetAllController struct {
	getAllUseCase *application.GetAllUseCase
}

func NewGetAllController(getAllUseCase *application.GetAllUseCase) *GetAllController {
	return &GetAllController{getAllUseCase: getAllUseCase}
}

func (c *GetAllController) GetAllOrders(ctx *gin.Context) {
	orders, err := c.getAllUseCase.Execute()
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to retrieve orders"})
		return
	}

	ctx.JSON(200, orders)
}
