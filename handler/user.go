package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"startup/auth"
	"startup/helper"
	"startup/user"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(service user.Service, auth auth.Service) *userHandler {
	return &userHandler{service, auth}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		errorMessage := gin.H{"errors": helper.FormatValidationError(err)}
		response := helper.APIRespose("Register account failed.", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIRespose("Register account failed.", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := h.authService.GenerateToken(newUser.ID)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIRespose("Register account failed.", http.StatusInternalServerError, "error", errorMessage)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	formatter := user.FormatUser(newUser, token)
	response := helper.APIRespose("Account has been created.", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c *gin.Context) {
	var input user.LoginUser

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errorMessage := gin.H{"errors": helper.FormatValidationError(err)}
		response := helper.APIRespose("Login failed.", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	loginUser, err := h.userService.Login(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIRespose("Login failed.", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	token, err := h.authService.GenerateToken(loginUser.ID)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIRespose("Login failed.", http.StatusInternalServerError, "error", errorMessage)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	formatter := user.FormatUser(loginUser, token)
	response := helper.APIRespose("Successfully login.", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) CheckEmailAvailability(c *gin.Context) {
	var input user.EmailCheckInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errorMessage := gin.H{"errors": helper.FormatValidationError(err)}
		response := helper.APIRespose("Email checking failed.", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errorMessage := gin.H{"error": "Server error"}
		response := helper.APIRespose("Email checking failed.", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	metaMessage := "Email has been registered"
	if isEmailAvailable {
		metaMessage = "Email is available"
	}
	data := gin.H{"is_available": isEmailAvailable}
	response := helper.APIRespose(metaMessage, http.StatusOK, "success", data)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		log.Println("FormFile: " + err.Error())
		errorMessage := gin.H{"is_uploaded": false}
		response := helper.APIRespose("Failed to upload avatar file.", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	user := c.MustGet("currentUser").(user.User)
	userID := user.ID

	path := fmt.Sprintf("images/A%d-%s", userID, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		log.Println("SaveUploadedFile: " + err.Error())

		errorMessage := gin.H{"is_uploaded": false}
		response := helper.APIRespose("Failed to upload avatar file.", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	_, err = h.userService.SaveAvatar(userID, path)
	if err != nil {
		log.Println("SaveAvatar: " + err.Error())

		errorMessage := gin.H{"is_uploaded": false}
		response := helper.APIRespose("Failed to upload avatar file.", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	data := gin.H{"is_uploaded": true}
	response := helper.APIRespose("Avatar successfully uploaded.", http.StatusOK, "success", data)

	c.JSON(http.StatusOK, response)
}
