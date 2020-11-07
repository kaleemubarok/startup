package transaction

import "time"

type CampaignTransactionsFormatter struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Amount int `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

func FormatCampaignTransaction(transaction Transaction) CampaignTransactionsFormatter  {
	format:=CampaignTransactionsFormatter{}
	format.ID=transaction.ID
	format.Name=transaction.User.Name
	format.Amount=transaction.Amount
	format.CreatedAt=transaction.CreatedAt
	return format
}

func FormatCampaignTransactions(transactions []Transaction) []CampaignTransactionsFormatter  {
	if len(transactions) == 0 {
		return []CampaignTransactionsFormatter{}
	}

	var formatter []CampaignTransactionsFormatter
	for _, transaction := range transactions {
		formatter = append(formatter, FormatCampaignTransaction(transaction))
	}

	return formatter
}
