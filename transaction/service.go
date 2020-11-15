package transaction

import (
	"errors"
	"log"
	"math/rand"
	"startup/campaign"
	"startup/payment"
	"strconv"
	"time"
)

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
	paymentService     payment.Service
}

type Service interface {
	GetTransactionsByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error)
	GetTransactionByUserID(userID int) ([]Transaction, error)
	SaveTransaction(input CreateTransactionInput) (Transaction, error)
	ProcessNotification(input NotificationInput) error
}

func NewService(repository Repository, campaignRepository campaign.Repository, paymentService payment.Service) *service {
	return &service{repository, campaignRepository, paymentService}
}

func (s *service) GetTransactionsByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error) {
	campaignDetail, err := s.campaignRepository.FindByID(input.ID)
	if err != nil {
		return []Transaction{}, err
	}

	if input.User.ID != campaignDetail.UserID {
		return []Transaction{}, errors.New("you are not authorized to see transaction list of this campaign")
	}

	transactions, err := s.repository.GetByCampaignID(input.ID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *service) GetTransactionByUserID(userID int) ([]Transaction, error) {
	transactions, err := s.repository.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (s *service) SaveTransaction(input CreateTransactionInput) (Transaction, error) {
	transaction := Transaction{}
	transaction.Amount = input.Amount
	transaction.CampaignID = input.CampaignID
	transaction.UserID = input.User.ID
	transaction.Status = "pending"
	transaction.Code = generateTransactionCode(6, input.User.ID)

	paymentTransaction := payment.Transaction{
		OrderID: transaction.Code,
		Amount:  input.Amount,
	}
	paymentURLResponse, err := s.paymentService.GetPaymentURL(paymentTransaction, input.User)
	if err != nil {
		return transaction, err
	}

	transaction.PaymentURL = paymentURLResponse

	newTransaction, err := s.repository.Save(transaction)
	if err != nil {
		return transaction, err
	}

	return newTransaction, nil
}

func generateTransactionCode(length int, userID int) string {
	userKey := strconv.Itoa(userID)
	var letters = []rune(userKey + "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(int64(userID) + time.Now().UnixNano())
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func (s *service) ProcessNotification(input NotificationInput) error {
	transaction, err := s.repository.GetByCode(input.OrderID)
	if err != nil {
		return err
	}
	log.Println("[ProcessNotification] incoming payment notification: ", transaction)

	//TODO validate signature key
	/*sha:= sha512.New()
	form:=transaction.Code+""
	sha.Write([]byte("Postman-157856885120010000.00VT-server-HJMpl9HLr_ntOKt5mRONdmKj"))
	fmt.Printf("[SHA]: %x", sha.Sum(nil))*/

	var status string = "Canceled"
	if input.PaymentType == "credit_card" {
		if input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
			status = "Paid"
		}
		if input.FraudStatus == "challenge" {
			status = "Pending"
		}
	} else {
		if input.TransactionStatus == "settlement" {
			status = "Paid"
		}
		if input.TransactionStatus == "pending" {
			status = "Pending"
		}
	}

	transaction.Status = status
	updatedTransaction, err := s.repository.Update(transaction)
	if err != nil {
		return err
	}
	log.Println("[ProcessNotification] updated payment status: ", updatedTransaction)

	if updatedTransaction.Status == "Paid" {
		var campaign campaign.Campaign
		if updatedTransaction.ID != 0 {
			campaign, err = s.campaignRepository.FindByID(updatedTransaction.ID)
			if err != nil {
				return err
			}

			campaign.BakerCount += 1
			campaign.CurrentAmount += updatedTransaction.Amount
			campaign.UpdatedAt = time.Now()
		}

		updatedCampaign, err := s.campaignRepository.Update(campaign)
		if err != nil {
			return err
		}
		log.Println("[ProcessNotification] updated campaign details: ", updatedCampaign)

	}

	log.Println("[ProcessNotification] process completed [END]")
	return nil
}
