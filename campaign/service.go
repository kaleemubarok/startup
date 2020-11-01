package campaign

import (
	"errors"
	"fmt"
	"github.com/gosimple/slug"
	"strings"
)

type Service interface {
	GetCampaigns(userID int) ([]Campaign, error)
	GetCampaign(ID GetCampaignDetailInput) (Campaign, error)
	CreateCampaign(input CreateCampaignInput) (Campaign, error)
	UpdateCampaign(ID GetCampaignDetailInput, campaignInput CreateCampaignInput) (Campaign, error)
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

func (s *service) generateSlug(slug string, slugIdentifier int) string {
	slugString := slug
	if slugIdentifier != 0 {
		slugString = fmt.Sprint(slugString,"-", slugIdentifier)
	}

	findSlug, _ := s.repository.FindBySlug(slugString)
	if findSlug.ID != 0 {
		return s.generateSlug(slug, slugIdentifier+1)
	}

	return slugString
}

func (s *service) UpdateCampaign(ID GetCampaignDetailInput, campaignInput CreateCampaignInput) (Campaign, error)  {
	campaign, err := s.repository.FindByID(ID.ID)
	if err != nil {
		return campaign, err
	}

	if campaign.User.ID!=campaignInput.User.ID {
		return campaign, errors.New("unauthorized to update this campaign")
	}

	if strings.Compare(campaign.Name, campaignInput.Name) != 0 {
		campaign.Slug=s.generateSlug(slug.Make(campaignInput.Name), 0)
	}
	campaign.Name=campaignInput.Name
	campaign.Description=campaignInput.Description
	campaign.ShortDescription=campaignInput.ShortDescription
	campaign.Perks= campaignInput.Perks

	updatedCampaign, err := s.repository.Update(campaign)
	if err != nil {
		return updatedCampaign, err
	}

	return updatedCampaign, nil
}