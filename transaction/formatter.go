package transaction

import (
	"time"
)

type CampaignTransactionsFormatter struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

func FormatCampaignTransaction(transaction Transaction) CampaignTransactionsFormatter {
	format := CampaignTransactionsFormatter{}
	format.ID = transaction.ID
	format.Name = transaction.User.Name
	format.Amount = transaction.Amount
	format.CreatedAt = transaction.CreatedAt
	return format
}

func FormatCampaignTransactions(transactions []Transaction) []CampaignTransactionsFormatter {
	if len(transactions) == 0 {
		return []CampaignTransactionsFormatter{}
	}

	var formatter []CampaignTransactionsFormatter
	for _, transaction := range transactions {
		formatter = append(formatter, FormatCampaignTransaction(transaction))
	}

	return formatter
}

type UserTransactionFormatter struct {
	ID        int       `json:"id"`
	Amount    int       `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	Campaign  UserTransactionCampaignImages
}

type UserTransactionCampaignImages struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

func FormatUserTransaction(transaction Transaction) UserTransactionFormatter {
	format := UserTransactionFormatter{}
	format.ID = transaction.ID
	format.Amount = transaction.Amount
	format.Status = transaction.Status
	format.CreatedAt = transaction.CreatedAt

	campaignImageFormat := UserTransactionCampaignImages{}
	campaignImageFormat.Name = transaction.Campaign.Name

	if len(transaction.Campaign.CampaignImages) > 0 {
		campaignImageFormat.ImageURL = transaction.Campaign.CampaignImages[0].FileName
	}

	format.Campaign = campaignImageFormat

	return format
}

func FormatUserTransactions(transactions []Transaction) []UserTransactionFormatter {
	if len(transactions) == 0 {
		return []UserTransactionFormatter{}
	}

	var formatter []UserTransactionFormatter
	for _, transaction := range transactions {
		formatter = append(formatter, FormatUserTransaction(transaction))
	}

	return formatter
}

type TransactionFormatter struct {
	ID         int    `json:"id"`
	CampaignID int    `json:"campaign_id"`
	UserID     int    `json:"user_id"`
	Amount     int    `json:"amount"`
	Status     string `json:"status"`
	Code       string `json:"code"`
	PaymentURL string `json:"payment_url"`
}

func FormatTransaction(transaction Transaction) TransactionFormatter {
	format := TransactionFormatter{
		ID:         transaction.ID,
		CampaignID: transaction.CampaignID,
		Amount:     transaction.Amount,
		UserID:     transaction.UserID,
		Status:     transaction.Status,
		Code:       transaction.Code,
		PaymentURL: transaction.PaymentURL,
	}
	return format
}
