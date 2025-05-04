package apikeys

import (
	"time"

	"github.com/groundcover-com/groundcover-sdk-go/sdk/models"
)

type CreateApiKeyRequest struct {
	Name             string     `json:"name" binding:"required"`
	ServiceAccountId string     `json:"serviceAccountId" binding:"required"`
	Description      string     `json:"description"`
	ExpirationDate   *time.Time `json:"expirationDate"`
}

type CreateApiKeyResponse struct {
	ApiKey string `json:"apiKey"`
	Id     string `json:"id"`
}

type ListApiKeysResponseItem struct {
	Id                 string                  `json:"id"`
	Name               string                  `json:"name"`
	ServiceAccountId   string                  `json:"serviceAccountId"`
	ServiceAccountName string                  `json:"serviceAccountName"`
	CreatedBy          string                  `json:"createdBy"`
	CreationDate       time.Time               `json:"creationDate"`
	Policies           []models.PolicyMetadata `json:"policies"`
	LastActive         *time.Time              `json:"lastActive"`
	Description        string                  `json:"description"`
	RevokedAt          *time.Time              `json:"revokedAt"`
	ExpiredAt          *time.Time              `json:"expiredAt"`
}
