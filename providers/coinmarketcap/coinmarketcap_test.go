package coinmarketcap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoinmarketcap_GetPrices(t *testing.T) {
	api := New("")

	prices, err := api.GetRates()
	assert.Nil(t, err)

	t.Log(prices)
}
