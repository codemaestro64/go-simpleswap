package gosimpleswap

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	apiKey     string
	httpClient *resty.Client
}

type ErrorResponse struct {
	IsAPIError  bool   `json:"is_api_error"`
	Code        int    `json:"code"`
	Error       string `json:"error"`
	Description string `json:"description"`
	TraceID     string `json:"trace_id"`
}

const baseURL = "https://api.simpleswap.io"

// Endpoints
const (
	getCurrencyEndpoint      = "/get_currency"
	getAllCurrenciesEndpoint = "/get_all_currencies"
	createExchangeEndpoint   = "/create_exchange"
	getExchangeEndpoint      = "/get_exchange"
	getExchangesEndpoint     = "/get_exchanges"
	getRangesEndpoint        = "/get_ranges"
	getEstimatedEndpoint     = "/get_estimated"
	getPairsEndpoint         = "/get_pairs"
	getPairEndpoint          = "/get_pair"
)

// New creates a new Client with the provided API key.
func New(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: resty.New().SetBaseURL(baseURL),
	}
}

// makeRequest handles the HTTP request and response marshaling.
func (c *Client) makeRequest(method, endpoint string, params, headers map[string]string, res any) *ErrorResponse {
	r := c.httpClient.R().SetQueryParam("api_key", c.apiKey)

	// Set headers and params
	c.setRequestHeadersAndParams(r, method, params, headers)

	// Perform the request
	resp, err := c.executeRequest(r, method, endpoint)
	if err != nil {
		return c.createErrorResponse(err, "error making API request")
	}

	// Handle the response
	return c.handleResponse(resp, res)
}

// setRequestHeadersAndParams sets headers and parameters for the request.
func (c *Client) setRequestHeadersAndParams(r *resty.Request, method string, params, headers map[string]string) {
	for header, value := range headers {
		r.SetHeader(header, value)
	}

	if params != nil {
		if method == http.MethodPost {
			r.SetBody(params)
		} else {
			r.SetQueryParams(params)
		}
	}
}

// executeRequest performs the actual HTTP request.
func (c *Client) executeRequest(r *resty.Request, method, endpoint string) (*resty.Response, error) {
	if method == http.MethodPost {
		return r.Post(endpoint)
	}
	return r.Get(endpoint)
}

// handleResponse processes the HTTP response and unmarshals it into the destination.
func (c *Client) handleResponse(resp *resty.Response, res any) *ErrorResponse {
	if resp.StatusCode() == http.StatusOK {
		return c.unmarshalResponse(resp.Body(), res)
	}
	return c.unmarshalErrorResponse(resp.Body())
}

// unmarshalResponse unmarshals a successful response.
func (c *Client) unmarshalResponse(body []byte, dest any) *ErrorResponse {
	if err := json.Unmarshal(body, dest); err != nil {
		return c.createErrorResponse(err, "error unmarshalling result")
	}
	return nil
}

// unmarshalErrorResponse handles error responses from the API.
func (c *Client) unmarshalErrorResponse(body []byte) *ErrorResponse {
	var errResp ErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		return c.createErrorResponse(err, "error unmarshalling error response")
	}
	return &errResp
}

// createErrorResponse generates a structured error response.
func (c *Client) createErrorResponse(err error, msg string) *ErrorResponse {
	return &ErrorResponse{
		IsAPIError:  false,
		Description: fmt.Sprintf("%s: %v", msg, err),
	}
}
