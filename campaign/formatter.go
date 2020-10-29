package campaign

import (
	"strings"
)

type CampaignFormatter struct {
	ID               int    `json:"id"`
	UserID           int    `json:"user_id"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	ImageURL         string `json:"image_url"`
	GoalAmount       int    `json:"goal_amount"`
	CurrentAmount    int    `json:"current_amount"`
	Slug             string `json:"slug"`
}

func FormatCampaign(campaign Campaign) CampaignFormatter {
	campaignFormatter := CampaignFormatter{}
	campaignFormatter.ID = campaign.ID
	campaignFormatter.Name = campaign.Name
	campaignFormatter.ShortDescription = campaign.ShortDescription
	campaignFormatter.GoalAmount = campaign.GoalAmount
	campaignFormatter.CurrentAmount = campaign.CurrentAmount
	campaignFormatter.UserID = campaign.UserID
	campaignFormatter.Slug = campaign.Slug

	if len(campaign.CampaignImages) > 0 {
		campaignFormatter.ImageURL = campaign.CampaignImages[0].FileName
	}

	return campaignFormatter
}

func FormatCampaigns(campaigns []Campaign) []CampaignFormatter {
	var campaignsFormatter []CampaignFormatter

	if len(campaigns) == 0 {
		return []CampaignFormatter{}
	}

	for _, campaign := range campaigns {
		campaignFormatter := FormatCampaign(campaign)
		campaignsFormatter = append(campaignsFormatter, campaignFormatter)
	}

	return campaignsFormatter
}

type CampaignDetailsFormatter struct {
	ID               int                       `json:"id"`
	Name             string                    `json:"name"`
	ShortDescription string                    `json:"short_description"`
	Description      string                    `json:"description"`
	ImageURL         string                    `json:"image_url"`
	GoalAmount       int                       `json:"goal_amount"`
	CurrentAmount    int                       `json:"current_amount"`
	UserID           int                       `json:"user_id"`
	Slug             string                    `json:"slug"`
	Perks            []string                  `json:"perks"`
	User             CampaignUserFormatter     `json:"user"`
	Images           []CampaignImagesFormatter `json:"images"`
}

type CampaignUserFormatter struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

type CampaignImagesFormatter struct {
	ImageURL  string `json:"image_url"`
	IsPrimary bool   `json:"is_primary"`
}

func FormatCampaignDetails(campaign Campaign) CampaignDetailsFormatter {
	campaignDetailsFormatter := CampaignDetailsFormatter{}
	campaignDetailsFormatter.ID = campaign.ID
	campaignDetailsFormatter.Name = campaign.Name
	campaignDetailsFormatter.ShortDescription = campaign.ShortDescription
	campaignDetailsFormatter.Description = campaign.Description
	campaignDetailsFormatter.GoalAmount = campaign.GoalAmount
	campaignDetailsFormatter.CurrentAmount = campaign.CurrentAmount
	campaignDetailsFormatter.UserID = campaign.UserID
	campaignDetailsFormatter.Slug = campaign.Slug

	if len(campaign.CampaignImages) > 0 {
		campaignDetailsFormatter.ImageURL = campaign.CampaignImages[0].FileName
	}

	var perks []string
	for _, perk := range strings.Split(campaign.Perks, ",") {
		perks = append(perks, perk)
	}
	campaignDetailsFormatter.Perks = perks

	userDetails := campaign.User
	campaignUserFormatter := CampaignUserFormatter{}
	campaignUserFormatter.Name = userDetails.Name
	campaignUserFormatter.ImageURL = userDetails.AvatarFileName
	campaignDetailsFormatter.User = campaignUserFormatter

	images := []CampaignImagesFormatter{}
	for _, image := range campaign.CampaignImages {
		campaignImages := CampaignImagesFormatter{}
		campaignImages.ImageURL = image.FileName

		isPrimary := false

		if image.IsPrimary == 1 {
			isPrimary = true
		}

		campaignImages.IsPrimary = isPrimary

		images = append(images, campaignImages)
	}
	campaignDetailsFormatter.Images = images

	return campaignDetailsFormatter
}
