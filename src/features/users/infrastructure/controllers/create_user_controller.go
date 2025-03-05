package controllers

import (
	"event-driven/src/features/users/application"
	"event-driven/src/features/users/domain"
	"github.com/gin-gonic/gin"
)

type CreateUserController struct {
	useCase *application.CreateUserUseCase
}

func NewUserCreateController(useCase *application.CreateUserUseCase) *CreateUserController {
	return &CreateUserController{useCase: useCase}
}

func (u *CreateUserController) Execute(c *gin.Context) {
	var user domain.User
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err = u.useCase.Execute(user)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, user)
}