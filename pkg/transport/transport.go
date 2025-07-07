package transport

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/rehttp"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	client "github.com/groundcover-com/groundcover-sdk-go/pkg/client"
)

type contextKey int

const (
	traceparentOverrideKey contextKey = iota
)

const (
	headerAuthorization = "Authorization"
	headerBackendID     = "X-Backend-Id"
	headerUserAgent     = "User-Agent"
	headerTraceparent   = "traceparent"
	userAgent           = "groundcover-go-sdk"
	yamlContentType     = "application/x-yaml"
)

const (
	defaultRetryCount = 3
	minRetryWait      = 1 * time.Second
	maxRetryWait      = 30 * time.Second
)

var getMonitorPathRegex = regexp.MustCompile(`^/api/monitors/[^/]+/?$`) // Matches /api/monitors/{id} but not /api/monitors/silences

// yamlByteConsumer consumes application/x-yaml as raw bytes
type yamlByteConsumer struct{}

// Consume reads the response body directly into data without YAML parsing.
// It expects 'data' to be a pointer to a []byte or similar slice type.
func (c *yamlByteConsumer) Consume(reader io.Reader, data interface{}) error {
	buf, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	// Check if data is a pointer to []byte or *[]byte
	byteSlicePtr, ok := data.(*[]byte)
	if !ok {
		// Also handle *strfmt.Base64, which is essentially *[]byte
		base64Ptr, ok := data.(*strfmt.Base64)
		if !ok {
			// Fallback: Try assigning to []uint8 if that's the underlying type
			uint8SlicePtr, ok := data.(*[]uint8)
			if !ok {
				return fmt.Errorf("yamlByteConsumer requires data to be *[]byte, *[]uint8, or *strfmt.Base64, got %T for content type %s", data, yamlContentType)
			}
			*uint8SlicePtr = buf
			return nil
		}
		*base64Ptr = buf
		return nil
	}

	*byteSlicePtr = buf
	return nil
}

// NewYamlByteConsumer returns a new instance of the YAML byte consumer
// for use with go-openapi runtime transport consumers
func NewYamlByteConsumer() *yamlByteConsumer {
	return &yamlByteConsumer{}
}

// ConfigureRuntimeTransport configures the provided runtime transport with
// the SDK's standard consumers and settings
func ConfigureRuntimeTransport(rt *httptransport.Runtime) {
	// Register the YAML byte consumer for application/x-yaml content type
	rt.Consumers[yamlContentType] = NewYamlByteConsumer()
}

// NewConfiguredRuntimeTransport creates a new runtime transport with
// the SDK's standard configuration applied
func NewConfiguredRuntimeTransport(host, basePath string, schemes []string) *httptransport.Runtime {
	rt := httptransport.New(host, basePath, schemes)
	ConfigureRuntimeTransport(rt)
	return rt
}

// ClientOption allows customization of the SDK client
type ClientOption func(*clientConfig)

type clientConfig struct {
	httpTransport    http.RoundTripper
	retryCount       int
	minWait          time.Duration
	maxWait          time.Duration
	retryStatuses    []int
	transportWrapper func(http.RoundTripper) http.RoundTripper
}

// WithHTTPTransport sets a custom HTTP transport
func WithHTTPTransport(transport http.RoundTripper) ClientOption {
	return func(c *clientConfig) {
		c.httpTransport = transport
	}
}

// WithRetryConfig sets custom retry configuration
func WithRetryConfig(retryCount int, minWait, maxWait time.Duration, retryStatuses []int) ClientOption {
	return func(c *clientConfig) {
		c.retryCount = retryCount
		c.minWait = minWait
		c.maxWait = maxWait
		c.retryStatuses = retryStatuses
	}
}

// WithTransportWrapper allows wrapping the transport (e.g., for debugging)
func WithTransportWrapper(wrapper func(http.RoundTripper) http.RoundTripper) ClientOption {
	return func(c *clientConfig) {
		c.transportWrapper = wrapper
	}
}

// NewSDKClient creates a fully configured groundcover SDK client with all
// standard configurations applied automatically. Use options to customize behavior.
func NewSDKClient(apiKey, backendID, baseURL string, options ...ClientOption) (*client.GroundcoverAPI, error) {
	// Default configuration
	config := &clientConfig{
		httpTransport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
		retryCount:    defaultRetryCount,
		minWait:       minRetryWait,
		maxWait:       maxRetryWait,
		retryStatuses: []int{http.StatusServiceUnavailable, http.StatusTooManyRequests, http.StatusGatewayTimeout, http.StatusBadGateway},
	}

	// Apply options
	for _, option := range options {
		option(config)
	}

	// Parse baseURL
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing base URL: %v", err)
	}

	host := parsedURL.Host
	basePath := parsedURL.Path
	if basePath == "" {
		basePath = client.DefaultBasePath
	}
	if !strings.HasPrefix(basePath, "/") && basePath != "" {
		basePath = "/" + basePath
	}

	schemes := []string{parsedURL.Scheme}
	if len(schemes) == 0 || schemes[0] == "" {
		schemes = client.DefaultSchemes
	}

	// Create transport with SDK functionality
	sdkTransport := NewTransport(
		apiKey,
		backendID,
		config.httpTransport,
		config.retryCount,
		config.minWait,
		config.maxWait,
		config.retryStatuses,
	)

	// Apply custom transport wrapper if provided
	finalTransport := http.RoundTripper(sdkTransport)
	if config.transportWrapper != nil {
		finalTransport = config.transportWrapper(sdkTransport)
	}

	// Create runtime transport with SDK configurations
	runtimeTransport := NewConfiguredRuntimeTransport(host, basePath, schemes)
	runtimeTransport.Transport = finalTransport

	// Create and return client
	return client.New(runtimeTransport, strfmt.Default), nil
}

// WithRequestTraceparent returns a new context with the Traceparent override.
func WithRequestTraceparent(ctx context.Context, traceparent string) context.Context {
	return context.WithValue(ctx, traceparentOverrideKey, traceparent)
}

// transport wraps an existing http.RoundTripper to add custom headers.
type transport struct {
	apiKey         string
	backendID      string
	retryTransport http.RoundTripper
}

// NewTransport creates a new transport.
// traceparent is optional and can be an empty string.
// retryCount, minWait, maxWait configure the retry mechanism.
// If retryCount, minWait, or maxWait are provided as 0, package defaults will be used.
func NewTransport(
	apiKey, backendID string,
	baseHttpTransport http.RoundTripper, // This is the transport *before* retries
	retryCount int,
	minWait, maxWait time.Duration,
	retryStatuses []int,
) *transport {
	if baseHttpTransport == nil {
		baseHttpTransport = http.DefaultTransport
	}

	// Default retry statuses if not provided or empty
	if len(retryStatuses) == 0 {
		retryStatuses = []int{http.StatusServiceUnavailable, http.StatusTooManyRequests, http.StatusGatewayTimeout, http.StatusBadGateway}
	}

	// Apply package defaults if parameters are zero
	if retryCount <= 0 {
		retryCount = defaultRetryCount
	}

	if minWait <= 0 {
		minWait = minRetryWait
	}

	if maxWait <= 0 {
		maxWait = maxRetryWait
	}

	// Configure retry transport
	rt := rehttp.NewTransport(
		baseHttpTransport,
		rehttp.RetryAll(
			rehttp.RetryMaxRetries(retryCount),
			rehttp.RetryStatuses(retryStatuses...),
		),
		rehttp.ExpJitterDelay(minWait, maxWait),
	)

	return &transport{
		apiKey:         apiKey,
		backendID:      backendID,
		retryTransport: rt,
	}
}

// RoundTrip executes a single HTTP transaction, checking context for overrides.
func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()

	var effectiveTraceparent string
	if traceVal, ok := ctx.Value(traceparentOverrideKey).(string); ok {
		effectiveTraceparent = traceVal
	}

	// Clone the request to avoid modifying the original passed to the base transport
	newReq := req.Clone(ctx)

	// --- Add Custom Headers ---
	newReq.Header.Set(headerAuthorization, fmt.Sprintf("Bearer %s", t.apiKey))
	newReq.Header.Set(headerBackendID, t.backendID)
	newReq.Header.Set(headerUserAgent, userAgent)

	if effectiveTraceparent != "" {
		newReq.Header.Set(headerTraceparent, effectiveTraceparent)
	}

	// --- Fix Content-Type for specific endpoints ---
	// Fix request Content-Type for workflow create endpoint
	if newReq.Method == http.MethodPost && newReq.URL.Path == "/api/workflows/create" {
		newReq.Header.Set("Content-Type", "text/plain")
	}

	// Execute the request
	resp, err := t.retryTransport.RoundTrip(newReq)
	if err != nil {
		return nil, err
	}

	// Fix response Content-Type for monitor GET endpoints
	if newReq.Method == http.MethodGet && resp.StatusCode == http.StatusOK &&
		getMonitorPathRegex.MatchString(newReq.URL.Path) &&
		!strings.Contains(newReq.URL.Path, "silences") {
		contentType := resp.Header.Get("Content-Type")
		if contentType == "" || !strings.HasPrefix(contentType, yamlContentType) {
			resp.Header.Set("Content-Type", yamlContentType)
		}
	}

	return resp, nil
}
