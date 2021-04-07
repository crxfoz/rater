package coinlayer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/crxfoz/webclient"
	"github.com/pkg/errors"
)

const (
	BaseURL = "http://api.coinlayer.com/api/live"
)

type rates struct {
	Success   bool               `json:"success"`
	Terms     string             `json:"terms"`
	Privacy   string             `json:"privacy"`
	Timestamp int64              `json:"timestamp"`
	Target    string             `json:"target"`
	Rates     map[string]float64 `json:"rates"`
}

type Coinlayer struct {
	APIKey string
	cleint *webclient.Webclient
}

func (c *Coinlayer) GetAllRates() (map[string]float64, error) {
	return c.getRates()
}

func (c *Coinlayer) GetRate(currency string) (float64, error) {
	data, err := c.getRates(currency)
	if err != nil {
		return 0, err
	}

	if rate, ok := data[currency]; ok {
		return rate, nil
	}

	return 0, fmt.Errorf("requested symbol isn't found in the body")
}

func (c *Coinlayer) getRates(currencies ...string) (map[string]float64, error) {
	request := c.cleint.Get(BaseURL).QueryParam("access_key", c.APIKey)

	if len(currencies) > 0 {
		currencyList := strings.Join(currencies, ",")
		request.QueryParam("symbols", currencyList)
	}

	resp, body, err := request.Do()
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

	if !data.Success {
		return nil, fmt.Errorf("api error")
	}

	return data.Rates, nil
}

func (c *Coinlayer) GetRates(currencies ...string) (map[string]float64, error) {
	return c.getRates(currencies...)
}

func New(apiKey string) *Coinlayer {
	return &Coinlayer{
		APIKey: apiKey,
		cleint: webclient.Config{
			Timeout:        time.Second * 10,
			UseKeepAlive:   false,
			FollowRedirect: false,
		}.New(),
	}
}
