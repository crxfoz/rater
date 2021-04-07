package coinlayer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoinlayer_GetPrices(t *testing.T) {
	api := New("")

	_, err := api.GetRates("FOOBAR")
	assert.NotNil(t, err)

	prices, err := api.GetRates("BTC", "ETH")
	assert.Nil(t, err)

	t.Log(prices)
}

func TestCoinlayer_GetPrice(t *testing.T) {
	api := New("")

	_, err := api.GetRate("FOOBAR")
	assert.NotNil(t, err)

	prices, err := api.GetRate("BTC")
	assert.Nil(t, err)

	t.Log(prices)
}

func TestCoinlayer_GetAllPrices(t *testing.T) {
	api := New("")

	data, err := api.GetAllRates()
	assert.Nil(t, err)

	t.Log(data)
}