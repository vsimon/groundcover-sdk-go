package e2e

import (
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	metricsClient "github.com/groundcover-com/groundcover-sdk-go/pkg/client/metrics"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
)

func TestMetricsQueryE2E(t *testing.T) {
	ctx, apiClient := setupTestClient(t)

	t.Run("Execute Metrics Query", func(t *testing.T) {
		// Define query parameters
		endTime := time.Now()
		startTime := endTime.Add(-15 * time.Minute) // Query last 15 minutes
		step := "1m"
		queryType := "range"
		promqlQuery := "up" // Simple query to check if targets are up

		// Convert times to strfmt.DateTime
		startDateTime := strfmt.DateTime(startTime)
		endDateTime := strfmt.DateTime(endTime)

		// Construct the request body
		body := &models.QueryRequest{
			Start:     startDateTime,
			End:       endDateTime,
			Step:      step,
			QueryType: queryType,
			Promql:    promqlQuery, // Use the direct promql field
			// Pipeline, Filters, Conditions, SubPipelines are omitted as we provide direct promql
		}

		// Create parameters
		params := metricsClient.NewMetricsQueryParamsWithContext(ctx).WithBody(body)

		// Execute the query
		resp, err := apiClient.Metrics.MetricsQuery(params, nil)

		// Assertions
		require.NoError(t, err, "Metrics query request failed")
		require.NotNil(t, resp, "Metrics query response should not be nil")
		require.NotNil(t, resp.Payload, "Metrics query response payload should not be nil")

		// Basic check on the payload structure (Prometheus response format)
		payloadMap, ok := resp.Payload.(map[string]interface{})
		require.True(t, ok, "Payload should be a map[string]interface{}")

		status, statusOk := payloadMap["status"].(string)
		assert.True(t, statusOk, "Payload should contain a 'status' field")
		assert.Equal(t, "success", status, "Query status should be 'success'")

		data, dataOk := payloadMap["data"].(map[string]interface{})
		assert.True(t, dataOk, "Payload should contain a 'data' field")
		assert.NotNil(t, data, "Data field should not be nil")

		t.Logf("Successfully executed metrics query '%s' with status '%s'", promqlQuery, status)
		// Further checks on data.resultType, data.result etc. could be added if needed
	})
}
