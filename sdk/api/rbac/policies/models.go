package policies

import (
	"time"

	"github.com/groundcover-com/groundcover-sdk-go/sdk/models"
)

type ApplyPolicyRequest struct {
	PolicyUUIDs []string `json:"policyUUIDs" binding:"required,min=1,dive,required"`
	Emails      []string `json:"emails" binding:"required,min=1,dive,required"`
	Override    bool     `json:"override"`
}

type CreatePolicyRequest struct {
	Name        string    `json:"name" binding:"required"`
	Description *string   `json:"description"`
	Role        RoleMap   `json:"role" binding:"required"`
	DataScope   DataScope `json:"dataScope"`
	ClaimRole   *string   `json:"claimRole"`
}

type UpdatePolicyRequest struct {
	Name            string    `json:"name" binding:"required"`
	Description     *string   `json:"description"`
	Role            RoleMap   `json:"role" binding:"required"`
	DataScope       DataScope `json:"dataScope"`
	ClaimRole       *string   `json:"claimRole"`
	CurrentRevision int32     `json:"currentRevision" binding:"required_if=Override false,excluded_if=Override true"`
}

type RoleMap map[string]string

type AdvancedDataScope struct {
	Workloads *models.Group `json:"workloads" validate:"required_without_all=Logs Traces Events Metrics"`
	Logs      *models.Group `json:"logs" validate:"required_without_all=Workloads Traces Events Metrics"`
	Traces    *models.Group `json:"traces" validate:"required_without_all=Workloads Logs Events Metrics"`
	Events    *models.Group `json:"events" validate:"required_without_all=Workloads Logs Traces Metrics"`
	Metrics   *models.Group `json:"metrics" validate:"required_without_all=Workloads Logs Traces Events"`
}

type DataScope struct {
	Simple   *models.Group      `json:"simple"`
	Advanced *AdvancedDataScope `json:"advanced"`
}

type Policy struct {
	UUID             string    `json:"uuid"`
	TenantUUID       string    `json:"tenantUuid"`
	Name             string    `json:"name"`
	Description      *string   `json:"description,omitempty"`
	Role             RoleMap   `json:"role"`      // Overrides original string Role
	DataScope        DataScope `json:"dataScope"` // Overrides original string DataScope
	ClaimRole        *string   `json:"claimRole"`
	ReadOnly         bool      `json:"readOnly"`
	CreatedBy        string    `json:"createdBy"`
	UpdatedBy        string    `json:"updatedBy"`
	CreatedTimestamp time.Time `json:"createdTimestamp"`
	UpdatedTimestamp time.Time `json:"updatedTimestamp"`
	RevisionNumber   int32     `json:"revisionNumber"`
}

type PolicyWithEntityCount struct {
	Policy      `json:",inline"`
	EntityCount int64 `json:"entityCount"`
}
