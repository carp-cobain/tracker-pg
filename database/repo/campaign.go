package repo

import (
	"time"

	"github.com/carp-cobain/tracker-pg/database/model"
	"github.com/carp-cobain/tracker-pg/database/query"
	"github.com/carp-cobain/tracker-pg/domain"

	"gorm.io/gorm"
)

// CampaignRepo manages campaigns in a database.
type CampaignRepo struct {
	db *gorm.DB
}

// NewCampaignRepo creates a new repository for managing campaigns.
func NewCampaignRepo(db *gorm.DB) CampaignRepo {
	return CampaignRepo{db}
}

// GetCampaign gets a campaign by ID
func (self CampaignRepo) GetCampaign(campaignID domain.CampaignID) (campaign domain.Campaign, err error) {
	var result model.Campaign
	if result, err = query.SelectCampaign(self.db, campaignID.String()); err == nil {
		campaign = result.ToDomain()
	}
	return
}

// GetCampaigns gets a page of campaigns for a blockchain account
func (self CampaignRepo) GetCampaigns(
	account domain.Account, pageParams domain.PageParams) domain.Page[domain.Campaign] {

	var nextCursor uint64
	results := query.SelectCampaigns(self.db, account.String(), pageParams.Cursor, pageParams.Limit)
	campaigns := make([]domain.Campaign, len(results))
	for i, result := range results {
		campaigns[i] = result.ToDomain()
		nextCursor = max(nextCursor, uint64(result.CreatedAt))
	}
	return domain.NewPage(nextCursor, pageParams.Limit, campaigns)
}

// CreateCampaign creates a new named campaign
func (self CampaignRepo) CreateCampaign(account domain.Account, name string) (campaign domain.Campaign, err error) {
	var result model.Campaign
	if result, err = query.InsertCampaign(self.db, account.String(), name); err == nil {
		campaign = result.ToDomain()
	}
	return
}

// UpdateCampaign updates campaign fields.
func (self CampaignRepo) UpdateCampaign(
	campaignID domain.CampaignID, name string, expiresAt time.Time) (campaign domain.Campaign, err error) {

	// Ensure campaign exists
	var existing model.Campaign
	existing, err = query.SelectCampaign(self.db, campaignID.String())
	if err != nil {
		return
	}

	// Only apply non-zero updates, keeping existing values
	var expires model.DateTime
	if expiresAt.IsZero() {
		expires = existing.ExpiresAt
	} else {
		expires = model.DateTime(expiresAt.Unix())
	}
	if name == "" {
		name = existing.Name
	}

	// Apply any updates
	var updated model.Campaign
	if updated, err = query.UpdateCampaign(self.db, existing, name, expires); err == nil {
		campaign = updated.ToDomain()
	}
	return
}