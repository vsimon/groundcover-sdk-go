package e2e

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apikeysClient "github.com/groundcover-com/groundcover-sdk-go/pkg/client/apikeys"
	policiesClient "github.com/groundcover-com/groundcover-sdk-go/pkg/client/policies"
	serviceaccountsClient "github.com/groundcover-com/groundcover-sdk-go/pkg/client/serviceaccounts"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
)

const (
	testApiKeyNamePrefix     = "sdk-e2e-test-apikey-"
	testApiKeyDesc           = "Created by SDK E2E test"
	testSACreatorEmail       = "sdk-e2e-test@groundcover.com"
	testServiceAccountNameSA = "sdk-e2e-test-sa-for-apikey"
	testPolicyNamePrefix     = "sdk-e2e-test-policy-for-apikey"
)

func TestAPIKeysE2E(t *testing.T) {
	ctx, apiClient := setupTestClient(t)

	// STEP 1: First create a policy as prerequisite
	policyNamePrefix := testPolicyNamePrefix
	policyName := fmt.Sprintf("%s-%d", policyNamePrefix, time.Now().UnixNano())

	// Create the policy using direct struct initialization
	policyReq := models.CreatePolicyRequest{
		Name:        &policyName,
		Description: "Policy for API Keys E2E testing",
		Role: models.RoleMap{
			"admin": "admin",
		},
		DataScope: &models.DataScope{},
	}

	policyParams := policiesClient.NewCreatePolicyParamsWithContext(ctx).WithBody(&policyReq)
	policyResp, err := apiClient.Policies.CreatePolicy(policyParams, nil)
	require.NoError(t, err, "Failed to create prerequisite Policy")
	require.NotNil(t, policyResp.Payload, "Policy response payload is nil")
	require.NotEmpty(t, policyResp.Payload.UUID, "Created Policy UUID is empty")
	policyUUID := policyResp.Payload.UUID
	t.Logf("Successfully created prerequisite Policy with UUID: %s", policyUUID)

	// Cleanup policy at the end of the test
	defer func() {
		t.Logf("Cleaning up temporary policy %s", policyUUID)
		delParams := policiesClient.NewDeletePolicyParamsWithContext(ctx).WithID(policyUUID)
		_, delErr := apiClient.Policies.DeletePolicy(delParams, nil)
		if delErr != nil {
			t.Logf("Warning: Could not delete policy %s: %v", policyUUID, delErr)
		} else {
			t.Logf("Cleaned up temporary policy %s", policyUUID)
		}
	}()

	// STEP 2: Now create the Service Account with the policy attached
	saNamePrefix := testServiceAccountNameSA
	saName := fmt.Sprintf("%s-%d", saNamePrefix, time.Now().UnixNano())
	saEmail := testSACreatorEmail
	saReq := &models.CreateServiceAccountRequest{
		Name:        &saName,
		Email:       &saEmail,
		PolicyUUIDs: []string{policyUUID}, // Attach the created policy
	}
	saParams := serviceaccountsClient.NewCreateServiceAccountParamsWithContext(ctx).WithBody(saReq)

	saResp, err := apiClient.Serviceaccounts.CreateServiceAccount(saParams, nil)
	require.NoError(t, err, "Failed to create prerequisite Service Account")
	require.NotNil(t, saResp.Payload, "Service Account response payload is nil")
	require.NotEmpty(t, saResp.Payload.ServiceAccountID, "Created Service Account ID is empty")
	serviceAccountID := saResp.Payload.ServiceAccountID
	t.Logf("Successfully created prerequisite Service Account with ID: %s", *serviceAccountID)

	// Important: Create all resources before adding cleanups
	// Cleanup for the prerequisite Service Account
	defer func() {
		deleteSAParams := serviceaccountsClient.NewDeleteServiceAccountParamsWithContext(ctx).WithID(*serviceAccountID)
		_, err := apiClient.Serviceaccounts.DeleteServiceAccount(deleteSAParams, nil)
		assert.NoError(t, err, "Failed to delete prerequisite Service Account during cleanup")
		t.Logf("Successfully cleaned up prerequisite Service Account: %s", *serviceAccountID)
	}()

	apiKeyName := fmt.Sprintf("%s%d", testApiKeyNamePrefix, time.Now().UnixNano())
	apiKeyDesc := testApiKeyDesc
	var createdApiKeyID string
	var createdApiKeyToken string

	// --- Create API Key ---
	t.Run("CreateAPIKey", func(t *testing.T) {
		createReq := &models.CreateAPIKeyRequest{
			Name:             &apiKeyName,
			ServiceAccountID: serviceAccountID,
			Description:      apiKeyDesc,
		}
		createParams := apikeysClient.NewCreateAPIKeyParamsWithContext(ctx).WithBody(createReq)

		createResp, err := apiClient.Apikeys.CreateAPIKey(createParams, nil)
		require.NoError(t, err, "CreateApiKey request failed")
		require.NotNil(t, createResp.Payload, "CreateApiKey response payload is nil")
		assert.NotEmpty(t, createResp.Payload.ID, "Created API Key ID is empty")
		assert.NotEmpty(t, createResp.Payload.APIKey, "Created API Key token is empty")

		createdApiKeyID = createResp.Payload.ID
		createdApiKeyToken = createResp.Payload.APIKey
		t.Logf("Successfully created API Key with ID: %s", createdApiKeyID)
		assert.NotEqual(t, "", createdApiKeyToken, "API Key token should not be empty")
	})

	require.NotEmpty(t, createdApiKeyID, "API Key ID was not set after creation")

	// --- List API Keys (Verify Creation) ---
	t.Run("ListAPIKeys_VerifyCreation", func(t *testing.T) {
		listParams := apikeysClient.NewListAPIKeysParamsWithContext(ctx)
		listResp, err := apiClient.Apikeys.ListAPIKeys(listParams, nil)
		require.NoError(t, err, "ListApiKeys request failed")
		require.NotNil(t, listResp.Payload, "ListApiKeys response payload is nil")

		found := false
		for _, key := range listResp.Payload {
			if key.ID == createdApiKeyID {
				// Key found, verify its properties
				assert.Equal(t, apiKeyName, key.Name)
				assert.Equal(t, *serviceAccountID, key.ServiceAccountID)
				assert.Equal(t, apiKeyDesc, key.Description)
				assert.True(t, key.RevokedAt.IsZero(), "Newly created key should not be revoked")
				found = true
				break
			}
		}

		assert.True(t, found, "Newly created API Key ID %s not found in list", createdApiKeyID)
		t.Logf("Successfully verified API Key %s exists via ListApiKeys", createdApiKeyID)
	})

	// --- Delete API Key ---
	t.Run("DeleteAPIKey", func(t *testing.T) {
		deleteParams := apikeysClient.NewDeleteAPIKeyParamsWithContext(ctx).WithID(createdApiKeyID)
		_, err := apiClient.Apikeys.DeleteAPIKey(deleteParams, nil)
		require.NoError(t, err, "DeleteApiKey request failed")
		t.Logf("Successfully initiated deletion for API Key ID: %s", createdApiKeyID)
	})

	// --- List API Keys (Verify Deletion) ---
	t.Run("ListAPIKeys_VerifyDeletion", func(t *testing.T) {
		// Check with withRevoked=true to see the revoked key
		withRevoked := true
		listParams := apikeysClient.NewListAPIKeysParamsWithContext(ctx).WithWithRevoked(&withRevoked)

		listResp, err := apiClient.Apikeys.ListAPIKeys(listParams, nil)
		require.NoError(t, err, "ListApiKeys request failed after deletion attempt")
		require.NotNil(t, listResp.Payload, "ListApiKeys response payload is nil after deletion attempt")

		foundActive := false
		foundRevoked := false

		for _, key := range listResp.Payload {
			if key.ID == createdApiKeyID {
				if !key.RevokedAt.IsZero() {
					foundRevoked = true
				} else {
					foundActive = true
				}
				break
			}
		}

		assert.False(t, foundActive, "Deleted API Key ID %s was found in list as active", createdApiKeyID)
		assert.True(t, foundRevoked, "Deleted (revoked) API Key ID %s not found in list even with withRevoked=true", createdApiKeyID)

		// Also check the default list (without withRevoked flag) - key should not be there
		listParamsDefault := apikeysClient.NewListAPIKeysParamsWithContext(ctx)
		listRespDefault, err := apiClient.Apikeys.ListAPIKeys(listParamsDefault, nil)
		require.NoError(t, err, "ListApiKeys request failed (default view)")

		foundDefault := false
		for _, key := range listRespDefault.Payload {
			if key.ID == createdApiKeyID {
				foundDefault = true
				break
			}
		}

		assert.False(t, foundDefault, "Deleted (revoked) API Key ID %s was found in default list (should be hidden)", createdApiKeyID)
		t.Logf("Successfully verified API Key %s is revoked/deleted via ListApiKeys", createdApiKeyID)
	})
}
