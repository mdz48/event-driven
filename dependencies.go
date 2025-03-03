package main

import (
	"citasAPI/src/core"
	usersUseCase "citasAPI/src/features/users/application"
	usersInfrastructure "citasAPI/src/features/users/infrastructure"
	usersControllers "citasAPI/src/features/users/infrastructure/controllers"
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

	usersDataBase := usersInfrastructure.NewMySQL(database.Conn)
	usersLogin := usersUseCase.NewLoginUserUseCase(usersDataBase)
	usersLoginController := usersControllers.NewLoginController(usersLogin)
	usersGetAllUseCase := usersUseCase.NewGetAllUsersUseCase(usersDataBase)
	usersGetAllController := usersControllers.NewGetAllUsersController(usersGetAllUseCase)
	usersCreateUseCase := usersUseCase.NewCreateUserUseCase(usersDataBase)
	usersCreateController := usersControllers.NewUserCreateController(usersCreateUseCase)
	usersRouter := usersInfrastructure.NewUserRouter(d.engine, usersGetAllController, usersLoginController, usersCreateController)
	usersRouter.SetUpRoutes()

	return d.engine.Run(":8080")
}