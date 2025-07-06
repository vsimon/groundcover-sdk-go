package transport

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/rehttp"
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
