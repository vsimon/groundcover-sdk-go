# Groundcover Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/groundcover-com/groundcover-sdk-go.svg)](https://pkg.go.dev/github.com/groundcover-com/groundcover-sdk-go)

The official Go SDK for interacting with the Groundcover API.

## Overview

This SDK provides convenient Go interfaces for various Groundcover API endpoints, including RBAC (API Keys, Policies, Service Accounts), Metrics, Monitors, and Kubernetes cluster/workload information.

## Installation

To use the SDK in your Go project, install it using `go get`:

```bash
go get github.com/groundcover-com/groundcover-sdk-go
```

## Development

To contribute or work on the SDK locally:

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/groundcover-com/groundcover-sdk-go.git
    cd groundcover-sdk-go
    ```
2.  **Make your changes.**
3.  **Ensure dependencies are tidy:**
    ```bash
    go mod tidy
    ```

## Usage

### Client Initialization

To use the SDK, you need to initialize an API client. The client requires your Groundcover Base URL, API Key, and Backend ID. It's recommended to provide these via environment variables:

*   `GC_BASE_URL`: The base URL of your Groundcover API endpoint.
*   `GC_API_KEY`: Your Groundcover API key.
*   `GC_BACKEND_ID`: Your Groundcover Backend ID.

```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/groundcover-com/groundcover-sdk-go/sdk/api"
	"github.com/groundcover-com/groundcover-sdk-go/sdk/api/rbac/apikeys"
	"github.com/groundcover-com/groundcover-sdk-go/sdk/httpclient"
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

	// Optional client options (e.g., for tracing)
	clientOpts := []httpclient.ClientOption{}
	traceparent := os.Getenv("GC_TRACEPARENT")
	if traceparent != "" {
		clientOpts = append(clientOpts, httpclient.WithTraceparent(traceparent))
	}

	// Initialize the API client
	client := api.NewClient(baseURL, apiKey, backendID, clientOpts...)

	// Now you can use the client to interact with different services
	// Example: Use the RBAC API Key service
	listApiKeysExample(client.Rbac.Apikeys)
}

func listApiKeysExample(apiKeyService *apikeys.Service) {
    // Example usage is shown in the Examples section below
}

```

### Services

The client provides access to different API services:

*   `client.K8s`: Interact with Kubernetes cluster and workload endpoints.
*   `client.Metrics`: Query metrics data.
*   `client.Monitors`: Manage monitors.
*   `client.Rbac`: Manage Role-Based Access Control:
    *   `client.Rbac.Apikeys`: Manage API Keys.
    *   `client.Rbac.Policies`: Manage RBAC Policies.
    *   `client.Rbac.Serviceaccounts`: Manage Service Accounts.

## Examples

Here are a few examples demonstrating how to use the SDK:

### List API Keys

List all API keys, optionally including revoked or expired ones.

```go
func listApiKeys(apiKeyService *apikeys.Service) {
	ctx := context.Background()

	// List only active keys
	activeKeys, err := apiKeyService.ListApiKeys(ctx, nil, nil)
	if err != nil {
		log.Fatalf("Error listing active API keys: %v", err)
	}
	log.Println("Active API Keys:")
	for _, key := range activeKeys {
		log.Printf("  ID: %s, Name: %s, ServiceAccountID: %s\n", key.Id, key.Name, key.ServiceAccountId)
	}

	// List keys including revoked and expired ones
	withRevoked := true
	withExpired := true
	allKeys, err := apiKeyService.ListApiKeys(ctx, &withRevoked, &withExpired)
	if err != nil {
		log.Fatalf("Error listing all API keys: %v", err)
	}
	log.Println("\nAll API Keys (including revoked/expired):")
	for _, key := range allKeys {
		log.Printf("  ID: %s, Name: %s, Revoked: %v, Expired: %v\n", key.Id, key.Name, key.RevokedAt != nil, key.ExpiredAt != nil)
	}
}

// In your main function or setup:
// listApiKeys(client.Rbac.Apikeys)
```

### Create an API Key

Create a new API key associated with a specific Service Account.

```go
func createApiKey(apiKeyService *apikeys.Service, serviceAccountID string) {
	ctx := context.Background()
	createReq := &apikeys.CreateApiKeyRequest{
		Name:             "my-sdk-generated-key",
		ServiceAccountId: serviceAccountID, // Replace with a valid Service Account ID
		Description:      "API Key created via Go SDK example",
	}

	response, err := apiKeyService.CreateApiKey(ctx, createReq)
	if err != nil {
		log.Fatalf("Error creating API key: %v", err)
	}

	log.Println("Successfully created API Key:")
	log.Printf("Key ID: %s", response.Id)
	log.Println("API Key (only shown once):", response.ApiKey)
}

// In your main function or setup:
// serviceAccountID := "YOUR_SERVICE_ACCOUNT_ID" // Replace with an actual ID
// createApiKey(client.Rbac.Apikeys, serviceAccountID)
```

### Delete an API Key

Delete an API key by its ID.

```go
func deleteApiKey(apiKeyService *apikeys.Service, keyID string) {
	ctx := context.Background()

	_, err := apiKeyService.DeleteApiKey(ctx, keyID) // Replace with the ID of the key to delete
	if err != nil {
		log.Fatalf("Error deleting API key %s: %v", keyID, err)
	}

	log.Printf("Successfully deleted API Key with ID: %s", keyID)
}

// In your main function or setup:
// keyIDToDelete := "YOUR_API_KEY_ID_TO_DELETE" // Replace with an actual ID
// deleteApiKey(client.Rbac.Apikeys, keyIDToDelete)

```
