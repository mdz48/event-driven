package main

import (
	"event-driven/src/core"
	ordersUseCase "event-driven/src/features/orders/application"
	ordersDomain "event-driven/src/features/orders/domain" // Corregido: odersDomain → ordersDomain
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

	// Si RabbitMQ está disponible, configurar manejadores de eventos
	if rabbitMQ != nil {
		// Configurar consumidor para órdenes creadas
		err = rabbitMQ.ConsumeSpecificEvents("created", func(order ordersDomain.Order) error { // Corregido: odersDomain → ordersDomain
			log.Printf("Procesando nueva orden creada: ID=%d, Producto=%s", order.ID, order.Product)
			// Lógica específica para órdenes nuevas
			return nil
		})
		if err != nil {
			log.Printf("Error configurando consumidor de órdenes creadas: %v", err)
		}

		// Configurar consumidor para cambios de estado
		err = rabbitMQ.ConsumeSpecificEvents("status_changed", func(order ordersDomain.Order) error { // Corregido: odersDomain → ordersDomain
			log.Printf("Procesando cambio de estado: ID=%d, Nuevo estado=%s", order.ID, order.Status)
			// Lógica específica para cambios de estado
			return nil
		})
		if err != nil {
			log.Printf("Error configurando consumidor de cambios de estado: %v", err)
		}
	}

	ordersRouter := ordersInfrastructure.NewOrderRouter(d.engine, ordersGetAllController, ordersCreateController)
	ordersRouter.SetUpRoutes()

	return d.engine.Run(":8080")
}