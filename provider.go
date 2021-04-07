package rater

type Provider interface {
	GetRate(currency string) (float64, error)
	GetRates(currencies ...string) (map[string]float64, error)
	GetAllRates() (map[string]float64, error)
}
