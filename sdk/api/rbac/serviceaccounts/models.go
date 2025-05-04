package serviceaccounts

import (
	"time"

	"github.com/groundcover-com/groundcover-sdk-go/sdk/models"
)

type CreateServiceAccountRequest struct {
	Name        string   `json:"name" binding:"required"`
	Email       string   `json:"email" binding:"required,email"`
	PolicyUUIDs []string `json:"policyUUIDs" binding:"required"`
}

type CreateServiceAccountResponse struct {
	ServiceAccountId string `json:"serviceAccountId"`
}

type ListServiceAccountsResponseItem struct {
	ServiceAccountId string                  `json:"serviceAccountId"`
	Name             string                  `json:"name"`
	Email            string                  `json:"email"`
	Deleted          bool                    `json:"deleted"`
	Policies         []models.PolicyMetadata `json:"policies"`
	LastActive       *time.Time              `json:"lastActive"`
}

type UpdateServiceAccountRequest struct {
	ServiceAccountId string   `json:"serviceAccountId" binding:"required"`
	Email            string   `json:"email"`
	PolicyUUIDs      []string `json:"policyUUIDs"`
}

type UpdateServiceAccountResponse struct {
	ServiceAccountId string `json:"serviceAccountId"`
}
