package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"startup/helper"
	"startup/user"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(service user.Service) *userHandler {
	return &userHandler{service}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		errorMessage := helper.FormatValidationError(err)
		response := helper.APIRespose("Register account failed.", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		errorMessage:=helper.FormatValidationError(err)
		response := helper.APIRespose("Register account failed.", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(newUser, "token")
	response := helper.APIRespose("Account has been created.", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}


