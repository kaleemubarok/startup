package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"startup/helper"
	"startup/payment"
	"startup/transaction"
	"startup/user"
)

type transactionHandler struct {
	service        transaction.Service
	paymentService payment.Service
}

func NewTransactionHandle(service transaction.Service, paymentService payment.Service) *transactionHandler {
	return &transactionHandler{service, paymentService}
}

func (h *transactionHandler) GetCampaignTransactions(c *gin.Context) {
	var input transaction.GetCampaignTransactionsInput
	err := c.ShouldBindUri(&input)
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

func (h *transactionHandler) GetUserTransaction(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	transactions, err := h.service.GetTransactionByUserID(userID)
	if err != nil {
		response := helper.APIRespose("Error to get user transactions", http.StatusBadRequest, "error", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIRespose("User transactions lists", http.StatusOK, "success", transaction.FormatUserTransactions(transactions))
	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) CreateTransaction(c *gin.Context) {
	var input transaction.CreateTransactionInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		response := helper.APIRespose("Error to create transaction", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	input.User = currentUser

	newTransaction, err := h.service.SaveTransaction(input)
	if err != nil {
		response := helper.APIRespose("Error to create transaction", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIRespose("Success to create new transaction", http.StatusOK, "success", transaction.FormatTransaction(newTransaction))

	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) ReceiveNotification(c *gin.Context) {
	var input transaction.NotificationInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		response := helper.APIRespose("Error to update transactions", http.StatusBadRequest, "error", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err = h.service.ProcessNotification(input)
	if err != nil {
		response := helper.APIRespose("Error to process notification transactions", http.StatusBadRequest, "error", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIRespose("Transactions notification status updated", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}
