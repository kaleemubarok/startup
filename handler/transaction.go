package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"startup/helper"
	"startup/transaction"
	"startup/user"
)

type transactionHandler struct {
	service transaction.Service
}

func NewTransactionHandle(service transaction.Service) *transactionHandler  {
	return &transactionHandler{service}
}

func (h *transactionHandler) GetCampaignTransactions(c *gin.Context)  {
	var input transaction.GetCampaignTransactionsInput
	err:=c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIRespose("Error to get campaign transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	input.User = currentUser

	transactions, err := h.service.GetTransactionsByCampaignID(input)
	if err != nil {
		response := helper.APIRespose("Error to get campaign transactions", http.StatusBadRequest, "error", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIRespose("Campaign transactions lists", http.StatusOK, "success", transaction.FormatCampaignTransactions(transactions))
	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) GetUserTransaction(c *gin.Context)  {
	currentUser:=c.MustGet("currentUser").(user.User)
	userID:=currentUser.ID

	transactions, err:=h.service.GetTransactionByUserID(userID)
	if err != nil {
		response := helper.APIRespose("Error to get user transactions", http.StatusBadRequest, "error", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIRespose("User transactions lists", http.StatusOK, "success", transaction.FormatUserTransactions(transactions))
	c.JSON(http.StatusOK, response)
}