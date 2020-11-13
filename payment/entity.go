package payment

type Transaction struct {
	OrderID string //`json:"id" binding:"required"`
	Amount  int //`json:"amount" binding:"required"`
}
