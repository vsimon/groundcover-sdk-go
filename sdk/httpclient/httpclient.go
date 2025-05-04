package httpclient

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/rehttp"
	"gopkg.in/yaml.v2"
)

const (
	headerAuthorization = "Authorization"
	headerBackendID     = "X-Backend-Id"
	headerContentType   = "Content-Type"
	headerContentLength = "Content-Length"
	headerUserAgent     = "User-Agent"
	headerTraceparent   = "traceparent"

	ContentTypeJSON = "application/json"
	ContentTypeYAML = "application/x-yaml"
	encodingGzip    = "gzip"
	userAgent       = "groundcover-go-sdk"
)

type CallRequest struct {
	Method      string
	Path        string
	QueryParams url.Values
	Body        interface{}
	Response    interface{}
	ContentType string
}

type Doer interface {
	Do(ctx context.Context, callRequest *CallRequest, opts ...RequestOption) error
}

type Config struct {
	Traceparent        string
	RetryCount         int
	GzipRequestEnabled bool
	Transport          http.RoundTripper
}

func (c *Config) Clone() *Config {
	return &Config{
		Traceparent:        c.Traceparent,
		RetryCount:         c.RetryCount,
		GzipRequestEnabled: c.GzipRequestEnabled,
	}
}

type Client struct {
	*Config
	BaseURL    string
	APIKey     string
	BackendID  string
	HTTPClient *http.Client
}

func NewClient(baseURL string, apiKey string, backendID string, opts ...ClientOption) *Client {
	client := &Client{
		BaseURL:   baseURL,
		APIKey:    apiKey,
		BackendID: backendID,
		HTTPClient: &http.Client{
			Timeout: time.Second * 30,
		},
		Config: &Config{
			RetryCount:         5,
			GzipRequestEnabled: false,
		},
	}

	client.applyOptions(opts)
	return client
}

func (c *Client) applyOptions(opts []ClientOption) {
	for _, opt := range opts {
		opt(c)
	}

	c.HTTPClient.Transport = createTransportFromConfig(c.Config)
}

func createTransportFromConfig(config *Config) http.RoundTripper {
	return rehttp.NewTransport(
		&http.Transport{
			DisableCompression: false,
		},
		rehttp.RetryAll(
			rehttp.RetryMaxRetries(config.RetryCount),
			rehttp.RetryStatuses(http.StatusServiceUnavailable, http.StatusTooManyRequests),
		),
		rehttp.ExpJitterDelay(time.Second, 5*time.Second),
	)
}

func (c *Client) applyRequestOptions(opts []RequestOption) *Config {
	config := c.Config.Clone()
	for _, opt := range opts {
		opt(config)
	}
	return config
}

func (c *Client) Do(ctx context.Context, callRequest *CallRequest, options ...RequestOption) error {
	httpReq, err := c.buildHttpRequest(ctx, callRequest.Method, callRequest.Path, callRequest.Body, callRequest.ContentType, options)
	if err != nil {
		return fmt.Errorf("error building request: %v", err)
	}

	// Handle query parameters if provided
	if len(callRequest.QueryParams) > 0 {
		// Get existing query values and merge with provided ones
		q := httpReq.URL.Query()
		for k, v := range callRequest.QueryParams {
			for _, vItem := range v {
				q.Add(k, vItem)
			}
		}
		// Set the combined query string
		httpReq.URL.RawQuery = q.Encode()
	}

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}

	if err := c.decodeResponse(resp, callRequest.Response); err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}
	return nil
}

func (c *Client) buildUrl(requestUrl string) (*url.URL, error) {
	joinedPath, err := url.JoinPath(c.BaseURL, requestUrl)
	if err != nil {
		return nil, fmt.Errorf("error joining path: %v", err)
	}
	parsedUrl, err := url.Parse(joinedPath)
	if err != nil {
		return nil, fmt.Errorf("error parsing url: %v", err)
	}

	return parsedUrl, nil
}

func (c *Client) buildHttpRequest(ctx context.Context, method string, urlPath string, request interface{}, contentType string, opts []RequestOption) (*http.Request, error) {
	config := c.applyRequestOptions(opts)
	var bodyReader io.Reader
	var contentLength int
	var err error

	if contentType == "" {
		contentType = ContentTypeJSON
	}

	if request != nil {
		bodyReader, contentLength, err = c.marshalBody(contentType, request)
		if err != nil {
			return nil, fmt.Errorf("error marshaling body: %v", err)
		}
	}

	parsedUrl, err := c.buildUrl(urlPath)
	if err != nil {
		return nil, fmt.Errorf("error parsing url: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, parsedUrl.String(), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	c.setHeaders(method, contentLength, req, config, contentType)
	return req, nil
}

func (c *Client) decodeResponse(resp *http.Response, response interface{}) error {
	if resp == nil {
		return fmt.Errorf("response is nil")
	}

	if resp.Body != nil && response != nil {
		if err := c.unmarshalResponseBody(resp, response); err != nil {
			return fmt.Errorf("error unmarshaling response body: %v", err)
		}
	}

	return nil
}

func (c *Client) setHeaders(method string, length int, httpReq *http.Request, conf *Config, contentType string) {
	httpReq.Header.Set(headerAuthorization, fmt.Sprintf("Bearer %s", c.APIKey))
	httpReq.Header.Set(headerBackendID, c.BackendID)
	httpReq.Header.Set(headerUserAgent, userAgent)

	if conf.Traceparent != "" {
		httpReq.Header.Set(headerTraceparent, conf.Traceparent)
	}

	if length > 0 {
		httpReq.Header.Set(headerContentType, contentType)
		httpReq.Header.Set(headerContentLength, fmt.Sprintf("%d", length))
	}
}

func (c *Client) marshalBody(contentType string, body interface{}) (io.Reader, int, error) {
	buf := bytes.NewBuffer(nil)
	if c.GzipRequestEnabled {
		gzw := gzip.NewWriter(buf)
		if err := json.NewEncoder(gzw).Encode(body); err != nil {
			return nil, 0, fmt.Errorf("error marshaling the body: %s", err)
		}
		if err := gzw.Close(); err != nil {
			return nil, 0, fmt.Errorf("error closing gzip writer: %s", err)
		}
	} else if contentType == ContentTypeYAML {
		switch body.(type) {
		case string:
			buf.WriteString(body.(string))
		case []byte:
			buf.Write(body.([]byte))
		default:
			if err := yaml.NewEncoder(buf).Encode(body); err != nil {
				return nil, 0, fmt.Errorf("error marshaling the body: %s", err)
			}
		}
	} else if contentType == ContentTypeJSON {
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, 0, fmt.Errorf("error marshaling the body: %s", err)
		}
	}

	return buf, buf.Len(), nil
}

func (c *Client) unmarshalResponseBody(resp *http.Response, target interface{}) error {
	if resp == nil {
		return errors.New("response is nil")
	}

	if resp.Body == nil {
		return errors.New("response body is nil")
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}
	if statusCode := resp.StatusCode; statusCode < 200 || statusCode >= 300 {
		return fmt.Errorf("error response: status code %d, body: %s length %d", statusCode, string(body), len(body))
	}

	switch target.(type) {
	case *[]byte:
		*target.(*[]byte) = body
	default:
		if err := json.Unmarshal(body, target); err != nil {
			return fmt.Errorf("error decoding response body: %v", err)
		}
	}

	return nil
}
