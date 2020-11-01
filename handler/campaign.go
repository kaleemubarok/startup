package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"startup/campaign"
	"startup/helper"
	"startup/user"
	"strconv"
)

type campaignHandler struct {
	campaignService campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))
	campaigns, err := h.campaignService.GetCampaigns(userID)
	if err != nil {
		response := helper.APIRespose("Error to get campaigns", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIRespose("List of campaigns", http.StatusOK, "success", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)

}

func (h *campaignHandler) GetCampaign(c *gin.Context) {
	var input campaign.GetCampaignDetailInput
	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIRespose("Error to get campaignDetails details", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	campaignDetails, errs := h.campaignService.GetCampaign(input)
	if errs != nil {
		response := helper.APIRespose("Error to get campaignDetails details", http.StatusBadRequest, "error", errs)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIRespose("Campaign details", http.StatusOK, "success", campaign.FormatCampaignDetails(campaignDetails))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) CreateCampaign(c *gin.Context) {
	var input campaign.CreateCampaignInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		response := helper.APIRespose("Failed to create campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	input.User = currentUser

	newCampaign, err := h.campaignService.CreateCampaign(input)
	if err != nil {
		response := helper.APIRespose("Failed to create campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIRespose("Success to create campaign", http.StatusOK, "success", campaign.FormatCampaign(newCampaign))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) UpdateCampaign(c *gin.Context) {
	var inputID campaign.GetCampaignDetailInput
	var inputCampaignDetail campaign.CreateCampaignInput

	err := c.ShouldBindUri(&inputID)
	println(err)
	if err != nil {
		response := helper.APIRespose("Failed to update campaign", http.StatusBadRequest, "error", helper.FormatValidationError(err))
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err = c.ShouldBindJSON(&inputCampaignDetail)
	if err != nil {
		response := helper.APIRespose("Failed to update campaign", http.StatusBadRequest, "error", helper.FormatValidationError(err))
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	inputCampaignDetail.User = currentUser

	updatedCampaign, err := h.campaignService.UpdateCampaign(inputID, inputCampaignDetail)
	if err != nil {
		errorMessage := helper.FormatValidationError(err)
		response := helper.APIRespose("Failed to update campaign", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIRespose("Success to update campaign", http.StatusOK, "success", campaign.FormatCampaign(updatedCampaign))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) CreateCampaignImage(c *gin.Context) {
	input := campaign.CreateCampaignImageInput{}
	err := c.ShouldBind(&input)
	if err != nil {
		errorMessage := gin.H{"errors": helper.FormatValidationError(err)}
		response := helper.APIRespose("Failed to upload campaign image", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		log.Println("FormFile: " + err.Error())
		errorMessage := gin.H{"is_uploaded": false}
		response := helper.APIRespose("Failed to upload campaign image", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	user := c.MustGet("currentUser").(user.User)
	input.User=user

	path := fmt.Sprintf("images/CI%d%d-%s", input.CampaignID, user.ID, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		log.Println("SaveCampaignImage: " + err.Error())

		errorMessage := gin.H{"is_uploaded": false}
		response := helper.APIRespose("Failed to upload campaign image", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	_, err = h.campaignService.CreateCampaignImage(input, path)
	if err != nil {
		log.Println("FormFile: " + err.Error())
		errorMessage := gin.H{"is_uploaded": false}
		response := helper.APIRespose("Failed to upload campaign image", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIRespose("Success to update campaign image", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)

}
