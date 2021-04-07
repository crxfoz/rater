package rater

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Logger interface {
	Error(args ...interface{})
}

type Rater struct {
	providers []Provider
	rates     map[string]float64
	ratesMu   sync.RWMutex
	period    time.Duration
	logger    Logger
}

func (r *Rater) GetRate(currency string) (float64, error) {
	r.ratesMu.RLock()
	defer r.ratesMu.RUnlock()

	if rate, ok := r.rates[currency]; ok {
		return rate, nil
	}

	return 0, fmt.Errorf("requested symbol isn't found in the body")
}

func copyMap(data map[string]float64) map[string]float64 {
	copiedMap := make(map[string]float64)
	for k, v := range data {
		copiedMap[k] = v
	}

	return copiedMap
}

func (r *Rater) GetAllRates() (map[string]float64, error) {
	r.ratesMu.RLock()
	defer r.ratesMu.RUnlock()

	copiedMap := copyMap(r.rates)

	return copiedMap, nil
}

func (r *Rater) GetRates(currencies ...string) (map[string]float64, error) {
	r.ratesMu.RLock()
	defer r.ratesMu.RUnlock()

	out := make(map[string]float64)

	for _, item := range currencies {
		if rate, ok := r.rates[item]; ok {
			out[item] = rate
		}
	}

	return out, nil
}

func New(providers []Provider, period time.Duration, logger Logger) (*Rater, error) {
	if len(providers) == 0 {
		return nil, fmt.Errorf("could not setup Rater with empty list")
	}

	return &Rater{
		providers: providers,
		rates:     make(map[string]float64),
		period:    period,
		logger:    logger,
	}, nil
}

// TODO: retry next provider if some symbols haven't been found
func (r *Rater) updateRates() error {
	for _, provider := range r.providers {
		rates, err := provider.GetAllRates()
		if err == nil {

			r.ratesMu.Lock()

			for symbol, rate := range rates {
				r.rates[symbol] = rate
			}

			r.ratesMu.Unlock()

			return nil
		}
	}

	return fmt.Errorf("all providers have failed")
}

func (r *Rater) Start(ctx context.Context, onUpdateCallback func(map[string]float64, error)) error {
	tk := time.NewTicker(r.period)

	if err := r.updateRates(); err != nil {
		r.logger.Error("could not perform initial updateRates")
	}

	for {
		select {
		case <-tk.C:

			err := r.updateRates()
			if err != nil {
				r.logger.Error("could not update rates: ", err)
			}

			if onUpdateCallback != nil {
				r.ratesMu.RLock()

				data := copyMap(r.rates)
				onUpdateCallback(data, err)

				r.ratesMu.RUnlock()
			}

		case <-ctx.Done():
			tk.Stop()
			return ctx.Err()
		}
	}
}
