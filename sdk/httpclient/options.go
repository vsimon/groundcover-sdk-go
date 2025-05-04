package httpclient

import "time"

type ClientOption func(*Client)

type RequestOption func(*Config)

func WithTraceparent(traceparent string) ClientOption {
	return func(c *Client) {
		c.Traceparent = traceparent
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.HTTPClient.Timeout = timeout
	}
}

func WithRetry(retryCount int) ClientOption {
	return func(c *Client) {
		c.RetryCount = retryCount
	}
}

func WithGzipAllRequests(enabled bool) ClientOption {
	return func(c *Client) {
		c.GzipRequestEnabled = enabled
	}
}

func WithGzipRequest(enabled bool) RequestOption {
	return func(c *Config) {
		c.GzipRequestEnabled = enabled
	}
}

func WithRequesTraceparent(traceparent string) RequestOption {
	return func(c *Config) {
		c.Traceparent = traceparent
	}
}
