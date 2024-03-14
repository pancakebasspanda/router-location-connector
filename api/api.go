package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

// API is an interface to be implemented by the client that connects to it to interact with router data location API
type API interface {
	GetRouterLocationData(ctx context.Context) (*RouterLocationData, error)
}

// Option specifies a builder function for configuring an APIs client
type Option func(API)

// Client is an Implementation of the API interface
type Client struct {
	options    options
	httpClient *retryablehttp.Client
}

// this is a check to confirm the implementation is compatible with dependent interfaces
var _ API = (*Client)(nil)

// New initializes the api's client
func New(opts ...Option) API {
	client := &Client{}

	for _, opt := range opts {
		opt(client)
	}

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = client.options.maxRetries
	retryClient.HTTPClient.Timeout = client.options.timeout
	retryClient.Backoff = func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		// example of too many requests
		// not relevant for this task as we ourselves are making the API calls
		// but if or app implements its own REST APU this would be relevant
		if resp != nil {
			if resp.StatusCode == http.StatusTooManyRequests {
				return 1 * time.Minute // 5 request quota per minute
			}
		}

		// any other error we perform an exponential backoff
		backOff := math.Pow(2, float64(attemptNum)) * float64(min)
		sleep := time.Duration(backOff)
		if float64(sleep) != backOff || sleep > max {
			sleep = max
		}

		return sleep
	}

	client.httpClient = retryClient

	return client
}

// GetRouterLocationData makes the call to the API and retrieves JSON data and transforms it to our models structs
func (c *Client) GetRouterLocationData(ctx context.Context) (*RouterLocationData, error) {
	req, err := c.prepareGetRequest(ctx, c.options.baseURL)
	if err != nil {
		return nil, err
	}

	req.Header = map[string][]string{"content-type": {"application/json"}}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("performing get router location data request: %s", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusTooManyRequests {
		return nil, fmt.Errorf("unexpected response status code: %d", resp.StatusCode)
	}

	// transform response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body %s", err)
	}

	var rLocData RouterLocationData
	if err := json.Unmarshal(body, &rLocData); err != nil {
		return nil, fmt.Errorf("unmarshal response body %s", err)
	}

	return &rLocData, nil
}

// prepareGetRequest helper function to define the get http request
func (c *Client) prepareGetRequest(ctx context.Context, requestURL string) (*retryablehttp.Request, error) {
	req, err := retryablehttp.NewRequest("GET", requestURL, nil)

	req = req.WithContext(ctx)

	if err != nil {
		return nil, err
	}

	return req, nil
}
