package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	k8sClient "github.com/groundcover-com/groundcover-sdk-go/pkg/client/k8s"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
)

// Constants for GetEventsOverTime test
const (
	ReasonOOMKilled    = "OOMKilled"
	TypeContainerCrash = "container_crash"
	ReasonFilterKey    = "reason"
	TypeFilterKey      = "type"
)

// Helper to create models.Condition for string equality for GetEventsOverTimeRequest
func newEventsOverTimeEqualStringCondition(key, value string) *models.Condition {
	return &models.Condition{
		Key:    key,
		Origin: "root",
		Type:   "string",
		Filters: []*models.Filter{
			{
				Op:    models.Op("eq"),
				Value: value,
			},
		},
	}
}

func TestK8sAPI(t *testing.T) {
	ctx := context.Background()
	_, client := setupTestClient(t)

	t.Run("list clusters", func(t *testing.T) {
		listParams := k8sClient.NewClustersListParams().
			WithBody(&models.ClustersListRequest{
				// Sources: optional, leave empty for now
			})
		listParams.SetContext(ctx)

		listResp, err := client.K8s.ClustersList(listParams, nil)
		require.NoError(t, err, "Error listing clusters")
		require.NotNil(t, listResp, "Clusters list response is nil")
		require.NotNil(t, listResp.Payload, "Clusters list payload is nil")

		assert.NotNil(t, listResp.Payload.Clusters, "Clusters list payload clusters is nil")
		// TotalCount can be 0 if no clusters are found, which is valid
	})

	t.Run("list workloads", func(t *testing.T) {
		listParams := k8sClient.NewWorkloadsListParams().
			WithBody(&models.WorkloadsListRequest{
				// Sources: optional
				// Conditions: optional
				// Set SortBy and Order explicitly to satisfy binding validation
				SortBy: "rps",  // Default value used by handler
				Order:  "desc", // Default value used by handler
				// Limit/Skip will use backend defaults if not set (0)
			})
		listParams.SetContext(ctx)

		listResp, err := client.K8s.WorkloadsList(listParams, nil)
		require.NoError(t, err, "Error listing workloads")
		require.NotNil(t, listResp, "Workloads list response is nil")
		require.NotNil(t, listResp.Payload, "Workloads list payload is nil")

		assert.NotNil(t, listResp.Payload.Workloads, "Workloads list payload workloads is nil")
		// Total can be 0 if no workloads are found, which is valid
	})

	t.Run("get OOM events over time", func(t *testing.T) {
		now := time.Now()
		startTime := now.Add(-15 * time.Minute)
		endTime := now

		// Required fields from model definition (ensure they match generated enum constants)
		sortBy := models.GetEventsOverTimeRequestSortByTimestamp
		sortOrder := models.GetEventsOverTimeRequestSortOrderDesc

		// Construct the request body
		reqBody := &models.GetEventsOverTimeRequest{
			Start:      (*strfmt.DateTime)(&startTime),
			End:        (*strfmt.DateTime)(&endTime),
			SortBy:     swag.String(sortBy),
			SortOrder:  swag.String(sortOrder),
			Conditions: []*models.Condition{
				// newEventsOverTimeEqualStringCondition(ReasonFilterKey, ReasonOOMKilled), // Uncomment to filter by OOMKilled reason
				// newEventsOverTimeEqualStringCondition(TypeFilterKey, TypeContainerCrash),
			},
			WithRawEvents: true,
		}

		params := k8sClient.NewGetEventsOverTimeParams().
			WithBody(reqBody)
		params.SetContext(ctx)

		resp, err := client.K8s.GetEventsOverTime(params, nil)

		require.NoError(t, err, "Error getting events over time")
		require.NotNil(t, resp, "GetEventsOverTime response is nil")
		require.NotNil(t, resp.Payload, "GetEventsOverTime payload is nil")
		require.NotNil(t, resp.Payload.Events, "GetEventsOverTime payload events is nil")
		require.Greater(t, len(resp.Payload.Events), 0, "No events found in the last 15 minutes")
		require.NotNil(t, resp.Payload.Events[0].Raw, "GetEventsOverTime payload events[0].Raw is nil")

		assert.NotNil(t, resp.Payload.Events, "Payload events slice is nil")
	})
}
