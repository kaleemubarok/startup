package campaign

import (
	"fmt"
	"github.com/gosimple/slug"
)

type Service interface {
	GetCampaigns(userID int) ([]Campaign, error)
	GetCampaign(ID GetCampaignDetailInput) (Campaign, error)
	CreateCampaign(input CreateCampaignInput) (Campaign, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s service) GetCampaigns(userID int) ([]Campaign, error) {
	if userID != 0 {
		campaigns, err := s.repository.FindByUserID(userID)
		if err != nil {
			return campaigns, err
		}
		return campaigns, nil
	}

	campaigns, err := s.repository.FindAll()
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

func (s service) GetCampaign(input GetCampaignDetailInput) (Campaign, error) {
	campaign, err := s.repository.FindByID(input.ID)
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (s service) CreateCampaign(input CreateCampaignInput) (Campaign, error) {
	campaign := Campaign{}
	campaign.Name = input.Name
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.GoalAmount = input.GoalAmount
	campaign.UserID = input.User.ID
	campaign.Perks = input.Perks

	slug := s.generateSlug(slug.Make(input.Name),0)
	campaign.Slug = slug

	newCampaign, err := s.repository.Save(campaign)
	if err != nil {
		return campaign, err
	}

	return newCampaign, nil
}

func (s *service) generateSlug(slug string, slugIndentifier int) string {
	slugString := slug
	if slugIndentifier != 0 {
		slugString = fmt.Sprint(slugString,"-",slugIndentifier)
	}

	findSlug, _ := s.repository.FindBySlug(slugString)
	if findSlug.ID != 0 {
		return s.generateSlug(slug,slugIndentifier+1)
	}

	return slugString
}
