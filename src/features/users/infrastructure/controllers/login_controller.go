package controllers

import (
	"citasAPI/src/features/users/application"
	"citasAPI/src/features/users/domain"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginController struct {
	loginUseCase *application.LoginUserUseCase
}

func NewLoginController(useCase *application.LoginUserUseCase) *LoginController {
	return &LoginController{loginUseCase: useCase}
}

func (controller *LoginController) Login(c *gin.Context) {
	var request domain.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := controller.loginUseCase.Execute(request.Email, request.Password)
	if err != nil {
		fmt.Printf("Login error: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    user.Id,
		"name":  user.Name,
		"email": user.Email,
	})
}