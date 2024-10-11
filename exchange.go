package gosimpleswap

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type ExchangeRequest struct {
	Fixed             bool   `json:"fixed"`
	CurrencyFrom      string `json:"currency_from"`
	CurrencyTo        string `json:"currency_to"`
	Amount            int    `json:"amount"`
	AddressTo         string `json:"address_to"`
	ExtraIDTo         string `json:"extra_id_to"`
	UserRefundAddress string `json:"user_refund_address"`
	UserRefundExtraID string `json:"user_refund_extra_id"`

	ClientIP        string `json:"-"`
	ClientUserAgent string `json:"-"`
	ClientTimezone  string `json:"-"`
	ClientLanguage  string `json:"-"`
}

type ExchangesRequest struct {
	Limit   int       `json:"limit"`
	Offset  int       `json:"offset"`
	MinTime time.Time `json:"gte"`
	MaxTime time.Time `json:"lte"`
}

type ExchangeResponse struct {
	ID                string `json:"id"`
	Type              string `json:"type"`
	Timestamp         string `json:"timestamp"`
	UpdatedAt         string `json:"updated_at"`
	ValidUntil        string `json:"valid_until"`
	CurrencyFrom      string `json:"currency_from"`
	CurrencyTo        string `json:"currency_to"`
	AmountFrom        string `json:"amount_from"`
	ExpectedAmount    string `json:"expected_amount"`
	AmountTo          string `json:"amount_to"`
	AddressFrom       string `json:"address_from"`
	AddressTo         string `json:"address_to"`
	ExtraIDFrom       string `json:"extra_id_from"`
	ExtraIDTo         string `json:"extra_id_to"`
	UserRefundAddress string `json:"user_refund_address"`
	UserRefundExtraID string `json:"user_refund_extra_id"`
	TxFrom            string `json:"tx_from"`
	TxTo              string `json:"tx_to"`
	Status            string `json:"status"`
	RedirectURL       string `json:"redirect_url"`
	Currencies        struct {
		CurrencyFromTicker CurrencyResponse `json:"currency_from_ticker"`
		CurrencyToTicker   CurrencyResponse `json:"currency_to_ticker"`
	} `json:"currencies"`
}

type RangesRequest struct {
	Fixed        bool
	CurrencyFrom string
	CurrencyTo   string
}

type RangesResponse struct {
	Minimum float64 `json:"min"`
	Maximum float64 `json:"max"`
}

func (c *Client) CreateExchange(req ExchangeRequest) (ExchangeResponse, *ErrorResponse) {
	headers := map[string]string{}
	if req.ClientIP != "" {
		headers["x-forwarded-for"] = req.ClientIP
	}

	if req.ClientLanguage != "" {
		headers["x-user-language"] = req.ClientLanguage
	}

	if req.ClientTimezone != "" {
		headers["x-user-timezone"] = req.ClientTimezone
	}

	if req.ClientUserAgent != "" {
		headers["x-user-agent"] = req.ClientUserAgent
	}

	params := map[string]string{
		"fixed":                fmt.Sprintf("%v", req.Fixed),
		"currency_from":        req.CurrencyFrom,
		"currency_to":          req.CurrencyTo,
		"amount":               fmt.Sprintf("%d", req.Amount),
		"extra_id_to":          req.ExtraIDTo,
		"user_refund_address":  req.UserRefundAddress,
		"user_refund_extra_id": req.UserRefundExtraID,
	}

	var res ExchangeResponse
	err := c.makeRequest(http.MethodPost, createExchangeEndpoint, params, headers, &res)

	return res, err
}

func (c *Client) GetExchange(exchangeID string) (ExchangeResponse, *ErrorResponse) {
	params := map[string]string{
		"id": exchangeID,
	}

	var res ExchangeResponse
	err := c.makeRequest(http.MethodGet, getExchangeEndpoint, params, nil, &res)

	return res, err
}

func (c *Client) GetExchanges(req ExchangesRequest) ([]ExchangeResponse, *ErrorResponse) {
	params := map[string]string{}
	if req.Limit != 0 {
		params["limit"] = fmt.Sprintf("%d", req.Limit)
	}

	if req.Offset != 0 {
		params["offset"] = fmt.Sprintf("%d", req.Offset)
	}

	if !req.MinTime.IsZero() {
		params["gte"] = req.MinTime.Format(time.RFC3339)
	}

	if !req.MaxTime.IsZero() {
		params["lte"] = req.MinTime.Format(time.RFC3339)
	}

	var res []ExchangeResponse
	err := c.makeRequest(http.MethodGet, getExchangesEndpoint, params, nil, &res)

	return res, err
}

func (c *Client) GetRanges(req RangesRequest) (*RangesResponse, *ErrorResponse) {
	params := map[string]string{
		"fixed":         fmt.Sprintf("%v", req.Fixed),
		"currency_from": req.CurrencyFrom,
		"currency_to":   req.CurrencyTo,
	}

	var res map[string]string
	err := c.makeRequest(http.MethodGet, getRangesEndpoint, params, nil, &res)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, c.createErrorResponse(errors.New("unknown error"), "error fetching ranges")
	}

	var result RangesResponse
	if min, ok := res["min"]; ok {
		minFl, e := strconv.ParseFloat(min, 64)
		if e != nil {
			return nil, c.createErrorResponse(e, "error marshalling result (min)")
		}
		result.Minimum = minFl
	}

	if max, ok := res["max"]; ok {
		maxFl, e := strconv.ParseFloat(max, 64)
		if e != nil {
			return nil, c.createErrorResponse(e, "error marshalling result (man)")
		}
		result.Maximum = maxFl
	}

	return &result, nil
}
