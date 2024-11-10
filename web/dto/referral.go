package dto

import "github.com/carp-cobain/tracker-pg/domain"

// ReferralRequest is the request type for adding campaign referrals.
type ReferralRequest struct {
	Account domain.Account `json:"account" binding:"required"`
}
