package coinmarketcap

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/crxfoz/webclient"
	"github.com/pkg/errors"
)

const (
	BaseURL = "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest"
)

type rates struct {
	Data []struct {
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
		Quote  map[string]struct {
			Price float64 `json:"price"`
		} `json:"quote"`
	} `json:"data"`
}

type Coinmarketcap struct {
	APIKey string
	cleint *webclient.Webclient
}

func (c *Coinmarketcap) GetAllRates() (map[string]float64, error) {
	return c.getRates()
}

func (c *Coinmarketcap) getRates(_ ...string) (map[string]float64, error) {
	resp, body, err := c.cleint.Get(BaseURL).SetHeader("X-CMC_PRO_API_KEY", c.APIKey).Do()
	if err != nil {
		return nil, errors.Wrap(err, "could not perform http-request")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code isn't valie: %d", resp.StatusCode)
	}

	var data rates

	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return nil, errors.Wrap(err, "unexpected body")
	}

	result := make(map[string]float64)

	for _, item := range data.Data {
		if usdRate, ok := item.Quote["USD"]; ok {
			result[item.Symbol] = usdRate.Price
		}
	}

	return result, nil
}

func (c *Coinmarketcap) GetRate(currency string) (float64, error) {
	prices, err := c.getRates(currency)
	if err != nil {
		return 0, err
	}

	if rate, ok := prices[currency]; ok {
		return rate, nil
	}

	return 0, fmt.Errorf("requested symbol isn't found in the body")
}

func (c *Coinmarketcap) GetRates(currencies ...string) (map[string]float64, error) {
	return c.getRates(currencies...)
}

func New(apiKey string) *Coinmarketcap {
	return &Coinmarketcap{
		APIKey: apiKey,
		cleint: webclient.Config{
			Timeout:        time.Second * 10,
			UseKeepAlive:   false,
			FollowRedirect: false,
		}.New(),
	}
}
