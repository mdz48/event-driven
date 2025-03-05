package infrastructure

import (
	"event-driven/src/features/orders/infrastructure/controllers"
	"github.com/gin-gonic/gin"
)

type OrderRouter struct {
	engine                 *gin.Engine
	getAllOrdersController *controllers.GetAllController
	createOrderController  *controllers.CreateOrderController
}

func NewOrderRouter(engine *gin.Engine, getAllOrdersController *controllers.GetAllController, createOrderController *controllers.CreateOrderController) *OrderRouter {
	return &OrderRouter{
		engine:                 engine,
		getAllOrdersController: getAllOrdersController,
		createOrderController:  createOrderController,
	}
}

func (r *OrderRouter) SetUpRoutes() {
	orders := r.engine.Group("/orders")
	{
		orders.GET("/", r.getAllOrdersController.GetAllOrders)
		orders.POST("/", r.createOrderController.CreateOrder)
	}

}

func (r *OrderRouter) Run() error {
	return r.engine.Run()
}