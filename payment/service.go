package payment

import (
	"github.com/veritrans/go-midtrans"
	"log"
	"startup/campaign"
	"startup/user"
)

type service struct {
	campaignRepository campaign.Repository
}

type Service interface {
	GetPaymentURL(transaction Transaction, user user.User) (string, error)
}

func NewService(campaignRepository campaign.Repository) *service {
	return &service{campaignRepository}
}

func (s *service) GetPaymentURL(transaction Transaction, user user.User) (string, error) {
	midclient := midtrans.NewClient()
	midclient.ServerKey = ""
	midclient.ClientKey = ""
	midclient.APIEnvType = midtrans.Sandbox

	snapGateway := midtrans.SnapGateway{
		Client: midclient,
	}

	snapReq := &midtrans.SnapReq{
		CustomerDetail: &midtrans.CustDetail{
			FName: user.Name,
			Email: user.Email,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  transaction.OrderID,
			GrossAmt: int64(transaction.Amount),
		},
		UserId: "",
	}

	log.Println("GetPaymentURL:")
	snapTokenResp, err := snapGateway.GetToken(snapReq)

	if err != nil {
		return "", err
	}
	return snapTokenResp.RedirectURL, err
}
