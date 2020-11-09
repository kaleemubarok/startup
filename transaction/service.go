package transaction

import (
	"errors"
	"startup/campaign"
)

type service struct {
	repository Repository
	campaignRepository campaign.Repository
}

type Service interface {
	GetTransactionsByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error)
	GetTransactionByUserID(userID int) ([]Transaction, error)
}

func NewService(repository Repository, campaignRepository campaign.Repository) *service{
	return &service{repository, campaignRepository}
}

func (s *service) GetTransactionsByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error)  {
	campaignDetail, err:=s.campaignRepository.FindByID(input.ID)
	if err != nil {
		return []Transaction{}, err
	}

	if input.User.ID!=campaignDetail.UserID {
		return []Transaction{}, errors.New("you are not authorized to see transaction list of this campaign")
	}

	transactions, err:=s.repository.GetByCampaignID(input.ID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *service) GetTransactionByUserID(userID int) ([]Transaction, error)  {
	transactions, err := s.repository.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	return transactions,nil
}