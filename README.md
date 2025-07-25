# groundcover Go SDK

This is the official Go SDK for interacting with the groundcover API. It provides convenient access to groundcover's services, including metrics queries and policy management.

## Prerequisites

*   Go 1.24 or higher.

## Installation

To use the SDK in your Go project, you can install it using `go get`:

```bash
go get github.com/groundcover-com/groundcover-sdk-go
```

## Configuration

### Environment Variables

The SDK requires the following environment variables to be set for authentication and endpoint configuration:

*   `GC_BASE_URL`: The base URL of the groundcover API (e.g., `https://api.groundcover.com`).
*   `GC_API_KEY`: Your groundcover API key.
*   `GC_BACKEND_ID`: Your groundcover Backend ID.

Optionally, you can set:

*   `GC_TRACEPARENT`: A default traceparent header value for distributed tracing.

### Client Initialization

The SDK provides simple functions to create a fully configured client with all necessary features (authentication, retries and other options) automatically included.

#### Simple Client Creation (Recommended)

For most use cases, you can create a client with minimal code:

```go
package main

import (
	"log"
	"os"

	"github.com/groundcover-com/groundcover-sdk-go/pkg/transport"
)

func main() {
	baseURL := os.Getenv("GC_BASE_URL")
	if baseURL == "" {
		log.Fatal("GC_BASE_URL environment variable is required")
	}

	apiKey := os.Getenv("GC_API_KEY")
	if apiKey == "" {
		log.Fatal("GC_API_KEY environment variable is required")
	}

	backendID := os.Getenv("GC_BACKEND_ID")
	if backendID == "" {
		log.Fatal("GC_BACKEND_ID environment variable is required")
	}

	// Create a fully configured client - handles auth, retries, content-type fixes, etc.
	sdkClient, err := transport.NewSDKClient(apiKey, backendID, baseURL)
	if err != nil {
		log.Fatalf("Failed to create SDK client: %v", err)
	}

	// Now you can use sdkClient to make API calls
	// Example: sdkClient.Metrics.MetricsQuery(...)
}
```

#### Client Creation with Custom Settings

If you need custom retry settings or HTTP transport:

```go
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/groundcover-com/groundcover-sdk-go/pkg/transport"
)

func main() {
	baseURL := os.Getenv("GC_BASE_URL")
	apiKey := os.Getenv("GC_API_KEY")
	backendID := os.Getenv("GC_BACKEND_ID")

	// Custom HTTP transport (optional)
	customTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		// Add any custom transport settings here
	}

	// Create client with custom settings using options
	sdkClient, err := transport.NewSDKClient(
		apiKey,
		backendID,
		baseURL,
		transport.WithHTTPTransport(customTransport),
		transport.WithRetryConfig(
			3,                    // retry count
			1*time.Second,        // min wait
			10*time.Second,       // max wait
			[]int{http.StatusServiceUnavailable, http.StatusTooManyRequests},
		),
	)
	if err != nil {
		log.Fatalf("Failed to create SDK client: %v", err)
	}

	// Client is ready to use
}
```

## Usage

### Making an API Call

Here's an example of how to make a metrics query:

```go
	// (Continued from Client Initialization above)
	// --- API Call: Metrics Query ---
	// import models "github.com/groundcover-com/groundcover-sdk-go/pkg/models"
	// import metrics "github.com/groundcover-com/groundcover-sdk-go/pkg/client/metrics"
	// import "github.com/sirupsen/logrus"
	// import "github.com/davecgh/go-spew/spew"

	baseCtx := context.Background()

	logrus.Info("--- Calling Metrics Query ---")

	// Prepare the request body for the metrics query
	startTime := strfmt.DateTime(time.Now().Add(-time.Hour))
	endTime := strfmt.DateTime(time.Now())
	step := "30s"
	queryType := "instant"
	promqlQuery := "avg(groundcover_container_cpu_limit_m_cpu)"

	queryRequestBody := &models.QueryRequest{
		Start:     startTime,
		End:       endTime,
		Step:      step,
		QueryType: queryType,
		Promql:    promqlQuery,
	}

	// Prepare the parameters for metrics query
	metricsParams := metrics.NewMetricsQueryParams().
		WithContext(baseCtx).
		WithTimeout(defaultTimeout). // Overall request timeout
		WithBody(queryRequestBody)

	// Execute the metrics query
	// The second argument (nil) is for AuthInfoWriter, as authentication is handled by our custom transport.
	queryResponse, err := sdkClient.Metrics.MetricsQuery(metricsParams, nil)
	if err != nil {
		// Handle errors (see Error Handling section)
		logrus.Errorf("Error executing metrics query: %v", err)
		return
	}

	// Handle the successful metrics response payload
	logrus.Info("Metrics Query Response:")
	spew.Dump(queryResponse.Payload) // queryResponse.Payload contains the data
```

### Building Conditions for Queries

When making API calls that accept a list of conditions (e.g., for filtering events or certain types of metrics), the SDK provides a convenient way to build these conditions using the `ConditionSet` helper located in the `pkg/utils` package. This builder simplifies creating the `[]*models.Condition` slice.

The `pkg/types` package (e.g., `github.com/groundcover-com/groundcover-sdk-go/pkg/types`) contains predefined constants for common condition keys, values, and operators.

Here's how to use the `ConditionSet`:

```go
// import (
// 	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
// 	"github.com/groundcover-com/groundcover-sdk-go/pkg/types"
// 	"github.com/groundcover-com/groundcover-sdk-go/pkg/utils"
// )

func getMyQueryConditions(namespace, podName string) []*models.Condition {
    cs := utils.NewConditionSet() // Initializes with default origin, type, and operator (eq)

    // Add a condition for namespace using default settings
    cs.Add(types.ConditionKeyNamespace, namespace)

    // Add a condition for pod name using default settings
    cs.Add(types.ConditionKeyPodName, podName)

    // Add predefined conditions for OOMKilled events
    cs.AddOOMEventConditions()

    // If you need to specify non-default parameters for a condition:
    // cs.AddFull(
    // 	types.ConditionKeyWorkload,      // Key
    // 	"customOrigin",                // Origin
    // 	"customType",                  // Type
    // 	"myWorkloadName",              // Value
    // 	types.OperatorContains,        // Operator
    // )

    return cs.Build() // Returns []*models.Condition
}

// Later, when preparing your query, you would use these conditions:
// queryRequestBody := &models.QueryRequest{
// 		Conditions: getMyQueryConditions("my-namespace", "my-pod-123"),
// 		// ... other query parameters ...
// }
```

Key methods for `ConditionSet`:

*   `utils.NewConditionSet()`: Creates a new condition set with defaults (Origin: `ConditionOriginRoot`, Type: `ConditionTypeString`, Operator: `OperatorEqual`).
*   `cs.Add(key, value string)`: Adds a condition using the default origin, type, and operator.
*   `cs.AddFull(key, origin, condType, value, opStr string)`: Adds a condition with explicitly specified parameters.
*   `cs.AddOOMEventConditions()`: A helper to add the standard conditions for detecting OOM events (Reason: `OOMKilled` and Type: `container_crash`).
*   `cs.Build()`: Returns the final `[]*models.Condition` slice.

### Context for Request Overrides

The `pkg/transport` module provides functions to set request-specific values, such as a traceparent, using `context.Context`.

*   **Traceparent**: Set a specific `traceparent` header for a request.
    ```go
    // Set a specific traceparent for this request
    metricsCtx := transport.WithRequestTraceparent(baseCtx, "00-customtraceid-customspanid-01")
    // ... then use metricsCtx in NewMetricsQueryParams().WithContext(metricsCtx)
    ```

### Retry Mechanism

The SDK's custom transport has a built-in retry mechanism that automatically retries requests on transient server errors (e.g., `503 Service Unavailable`, `429 Too Many Requests`). This is configured during client initialization via `transport.NewTransport`.

### Error Handling

API calls can return errors. It's important to handle these appropriately. The SDK uses specific error types for different API responses, and also a generic `runtime.APIError`.

```go
	// import "github.com/go-openapi/runtime"
	// import metrics "github.com/groundcover-com/groundcover-sdk-go/pkg/client/metrics"

	// (inside an API call block like the metrics query example)
	// queryResponse, err := sdkClient.Metrics.MetricsQuery(metricsParams, nil)
	if err != nil {
		switch e := err.(type) {
		case *metrics.MetricsQueryBadRequest: // Example specific error
			logrus.Errorf("Metrics API Error (Bad Request): %s, Payload: %v", e.Error(), e.Payload)
		case *metrics.MetricsQueryInternalServerError: // Example specific error
			logrus.Errorf("Metrics API Error (Internal Server Error): %s, Payload: %v", e.Error(), e.Payload)
		default:
			if apiErr, ok := err.(*runtime.APIError); ok {
				// This is a generic error from the go-openapi runtime
				// apiErr.Code gives the HTTP status code
				// apiErr.Response gives the raw response body (needs to be parsed or read)
				logrus.Errorf("Generic API Error: Code %d, Message: %s, Response: %v", apiErr.Code, apiErr.Error(), apiErr.Response)
			} else {
				// Other unexpected errors
				logrus.Errorf("Error executing API call: %v", err)
			}
		}
		return // Or handle as appropriate
	}
	// Process successful response: queryResponse.Payload
```

## Available Services

The SDK is organized by service, available under the `sdkClient` object. For example:

*   `sdkClient.Metrics`: For querying metrics.
*   `sdkClient.Policies`: For managing policies (example was commented out in `main.go` but shows the pattern).

Refer to the generated SDK code in the `pkg/client` directory for a full list of services and their operations.

## License

This SDK is distributed under the Apache License, Version 2.0. See the [LICENSE](LICENSE) file for more information.

