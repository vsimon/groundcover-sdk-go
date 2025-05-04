package api

import (
	"github.com/groundcover-com/groundcover-sdk-go/sdk/api/k8s"
	"github.com/groundcover-com/groundcover-sdk-go/sdk/api/metrics"
	"github.com/groundcover-com/groundcover-sdk-go/sdk/api/monitors"
	"github.com/groundcover-com/groundcover-sdk-go/sdk/api/rbac"
	"github.com/groundcover-com/groundcover-sdk-go/sdk/httpclient"
)

type Client struct {
	K8s      *k8s.Service
	Metrics  *metrics.Service
	Monitors *monitors.Service
	Rbac     *rbac.Service
}

func NewClient(baseURL string, apiKey string, backendID string, opts ...httpclient.ClientOption) *Client {
	httpClient := httpclient.NewClient(baseURL, apiKey, backendID, opts...)

	return &Client{
		K8s:      k8s.New(httpClient),
		Metrics:  metrics.New(httpClient),
		Monitors: monitors.New(httpClient),
		Rbac:     rbac.New(httpClient),
	}
}
