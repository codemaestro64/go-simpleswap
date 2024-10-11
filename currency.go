package gosimpleswap

import "net/http"

// Response types
type CurrencyResponse struct {
	Name              string   `json:"name"`
	Symbol            string   `json:"symbol"`
	Network           string   `json:"network"`
	ContractAddress   string   `json:"contract_address"`
	HasExtraID        bool     `json:"has_extra_id"`
	ExtraID           string   `json:"extra_id"`
	Image             string   `json:"image"`
	WarningsFrom      []string `json:"warnings_from"`
	WarningsTo        []string `json:"warnings_to"`
	ValidationAddress string   `json:"validation_address"`
	ValidationExtra   string   `json:"validation_extra"`
	AddressExplorer   string   `json:"address_explorer"`
	TxExplorer        string   `json:"tx_explorer"`
	ConfirmationsFrom string   `json:"confirmations_from"`
	IsFiat            bool     `json:"isFiat"`
}

func (c *Client) GetCurrency(symbol string) (CurrencyResponse, *ErrorResponse) {
	params := map[string]string{
		"symbol": symbol,
	}

	var res CurrencyResponse
	err := c.makeRequest(http.MethodGet, getCurrencyEndpoint, params, nil, &res)

	return res, err
}

func (c *Client) GetAllCurrencies() ([]CurrencyResponse, *ErrorResponse) {
	var res []CurrencyResponse
	err := c.makeRequest(http.MethodGet, getAllCurrenciesEndpoint, nil, nil, &res)

	return res, err
}
