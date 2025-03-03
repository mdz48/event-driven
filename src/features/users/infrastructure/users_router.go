package infrastructure

import (
	"citasAPI/src/features/users/infrastructure/controllers"
	"github.com/gin-gonic/gin"
)

type UserRouter struct {
	engine                *gin.Engine
	login                 *controllers.LoginController
	getAllUsersController *controllers.GetAllUsersController
	createUserController  *controllers.CreateUserController
}

func NewUserRouter(engine *gin.Engine, getAllUsersController *controllers.GetAllUsersController, login *controllers.LoginController, createUserController *controllers.CreateUserController) *UserRouter {
	return &UserRouter{
		engine:                engine,
		getAllUsersController: getAllUsersController,
		login:                 login,
		createUserController:  createUserController,
	}
}

func (r *UserRouter) SetUpRoutes() {
	users := r.engine.Group("/users")
	{
		users.GET("/", r.getAllUsersController.Execute)
		users.POST("/login", r.login.Login)
		users.POST("/", r.createUserController.Execute)
	}
}

func (r *UserRouter) Run() error {
	return r.engine.Run()
}