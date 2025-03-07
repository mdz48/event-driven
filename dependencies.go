package main

import (
	"event-driven/src/core"
	ordersUseCase "event-driven/src/features/orders/application"
	ordersInfrastructure "event-driven/src/features/orders/infrastructure"
	ordersControllers "event-driven/src/features/orders/infrastructure/controllers"
	"event-driven/src/features/orders/infrastructure/producer"
	usersUseCase "event-driven/src/features/users/application"
	usersInfrastructure "event-driven/src/features/users/infrastructure"
	usersControllers "event-driven/src/features/users/infrastructure/controllers"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	engine *gin.Engine
}

func NewDependencies() *Dependencies {
	return &Dependencies{
		engine: gin.Default(),
	}

}

func (d *Dependencies) Run() error {
	database := core.NewDatabase()

	// Configuración de usuarios
	usersDataBase := usersInfrastructure.NewMySQL(database.Conn)
	usersLogin := usersUseCase.NewLoginUserUseCase(usersDataBase)
	usersLoginController := usersControllers.NewLoginController(usersLogin)
	usersGetAllUseCase := usersUseCase.NewGetAllUsersUseCase(usersDataBase)
	usersGetAllController := usersControllers.NewGetAllUsersController(usersGetAllUseCase)
	usersCreateUseCase := usersUseCase.NewCreateUserUseCase(usersDataBase)
	usersCreateController := usersControllers.NewUserCreateController(usersCreateUseCase)
	usersRouter := usersInfrastructure.NewUserRouter(d.engine, usersGetAllController, usersLoginController, usersCreateController)
	usersRouter.SetUpRoutes()

	// Configuración de órdenes con MySQL y RabbitMQ
	ordersDataBase := ordersInfrastructure.NewMySQL(database.Conn)

	rabbitURL := os.Getenv("RABBITMQ_URL")

	rabbitMQ, err := producer.NewRabbitMQ(rabbitURL)
	if err != nil {
		log.Printf("Error inicializando RabbitMQ: %v. Continuando sin mensajería...", err)
	} else {
		defer rabbitMQ.Close()
	}

	ordersGetAllUseCase := ordersUseCase.NewGetAllUseCase(ordersDataBase)
	ordersGetAllController := ordersControllers.NewGetAllController(ordersGetAllUseCase)
	ordersCreateUseCase := ordersUseCase.NewCreateOrderUseCase(ordersDataBase)
	ordersCreateController := ordersControllers.NewCreateOrderController(ordersCreateUseCase, rabbitMQ)

	ordersRouter := ordersInfrastructure.NewOrderRouter(d.engine, ordersGetAllController, ordersCreateController)
	ordersRouter.SetUpRoutes()

	return d.engine.Run(":8080")
}