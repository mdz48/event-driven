package controllers

import (
	"event-driven/src/features/users/application"
	"event-driven/src/features/users/domain"
	"github.com/gin-gonic/gin"
)

type GetAllUsersController struct {
	useCase *application.GetAllUsersUseCase
}

func NewGetAllUsersController(useCase *application.GetAllUsersUseCase) *GetAllUsersController {
	return &GetAllUsersController{useCase: useCase}
}

func (g *GetAllUsersController) Execute(c *gin.Context) {
	users, err := g.useCase.Execute()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	var usersResponse []domain.User
	for _, user := range users {
		usersResponse = append(usersResponse, domain.User{
			Id:    user.Id,
			Name:  user.Name,
			Email: user.Email,
		})
	}

	c.JSON(200, usersResponse)
}